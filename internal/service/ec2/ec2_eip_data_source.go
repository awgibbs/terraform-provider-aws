// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package ec2

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
	tftags "github.com/hashicorp/terraform-provider-aws/internal/tags"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

// @SDKDataSource("aws_eip", name="EIP)
// @Tags
func dataSourceEIP() *schema.Resource {
	return &schema.Resource{
		ReadWithoutTimeout: dataSourceEIPRead,

		Timeouts: &schema.ResourceTimeout{
			Read: schema.DefaultTimeout(20 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"association_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"carrier_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"customer_owned_ipv4_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"domain": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"filter": customFiltersSchema(),
			"id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"instance_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"network_interface_owner_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_ip": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"private_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ip": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"public_dns": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"public_ipv4_pool": {
				Type:     schema.TypeString,
				Computed: true,
			},
			names.AttrTags: tftags.TagsSchemaComputed(),
		},
	}
}

func dataSourceEIPRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	conn := meta.(*conns.AWSClient).EC2Client(ctx)

	input := &ec2.DescribeAddressesInput{}

	if v, ok := d.GetOk("id"); ok {
		input.AllocationIds = []string{v.(string)}
	}

	if v, ok := d.GetOk("public_ip"); ok {
		input.PublicIps = []string{v.(string)}
	}

	input.Filters = append(input.Filters, newTagFilterListV2(
		TagsV2(tftags.New(ctx, d.Get("tags").(map[string]interface{}))),
	)...)

	input.Filters = append(input.Filters, newCustomFilterListV2(
		d.Get("filter").(*schema.Set),
	)...)

	if len(input.Filters) == 0 {
		// Don't send an empty filters list; the EC2 API won't accept it.
		input.Filters = nil
	}

	eip, err := findEIP(ctx, conn, input)

	if err != nil {
		return sdkdiag.AppendFromErr(diags, tfresource.SingularDataSourceFindError("EC2 EIP", err))
	}

	if eip.Domain == types.DomainTypeVpc {
		d.SetId(aws.ToString(eip.AllocationId))
	} else {
		d.SetId(aws.ToString(eip.PublicIp))
	}
	d.Set("association_id", eip.AssociationId)
	d.Set("carrier_ip", eip.CarrierIp)
	d.Set("customer_owned_ip", eip.CustomerOwnedIp)
	d.Set("customer_owned_ipv4_pool", eip.CustomerOwnedIpv4Pool)
	d.Set("domain", eip.Domain)
	d.Set("instance_id", eip.InstanceId)
	d.Set("network_interface_id", eip.NetworkInterfaceId)
	d.Set("network_interface_owner_id", eip.NetworkInterfaceOwnerId)
	d.Set("public_ipv4_pool", eip.PublicIpv4Pool)
	d.Set("private_ip", eip.PrivateIpAddress)
	if v := aws.ToString(eip.PrivateIpAddress); v != "" {
		d.Set("private_dns", PrivateDNSNameForIP(ctx, meta.(*conns.AWSClient), v))
	}
	d.Set("public_ip", eip.PublicIp)
	if v := aws.ToString(eip.PublicIp); v != "" {
		d.Set("public_dns", PublicDNSNameForIP(ctx, meta.(*conns.AWSClient), v))
	}

	setTagsOutV2(ctx, eip.Tags)

	return diags
}
