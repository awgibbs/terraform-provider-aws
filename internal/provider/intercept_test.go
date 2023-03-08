package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-provider-aws/internal/errs/sdkdiag"
)

func TestInterceptorsWhy(t *testing.T) {
	var interceptors interceptorItems

	interceptors = append(interceptors, interceptorItem{
		When: Before,
		Why:  Create,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})
	interceptors = append(interceptors, interceptorItem{
		When: After,
		Why:  Delete,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})
	interceptors = append(interceptors, interceptorItem{
		When: Before,
		Why:  Create,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})

	if got, want := len(interceptors.Why(Create)), 2; got != want {
		t.Errorf("length of interceptors.Why(Create) = %v, want %v", got, want)
	}
	if got, want := len(interceptors.Why(Read)), 0; got != want {
		t.Errorf("length of interceptors.Why(Read) = %v, want %v", got, want)
	}
	if got, want := len(interceptors.Why(Update)), 0; got != want {
		t.Errorf("length of interceptors.Why(Update) = %v, want %v", got, want)
	}
	if got, want := len(interceptors.Why(Delete)), 1; got != want {
		t.Errorf("length of interceptors.Why(Delete) = %v, want %v", got, want)
	}
}

func TestInterceptedHandler(t *testing.T) {
	var interceptors interceptorItems

	interceptors = append(interceptors, interceptorItem{
		When: Before,
		Why:  Create,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})
	interceptors = append(interceptors, interceptorItem{
		When: After,
		Why:  Delete,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})
	interceptors = append(interceptors, interceptorItem{
		When: Before,
		Why:  Create,
		Interceptor: interceptorFunc(func(ctx context.Context, d *schema.ResourceData, meta any, when When, why Why, diags diag.Diagnostics) (context.Context, diag.Diagnostics) {
			return ctx, diags
		}),
	})

	var read schema.ReadContextFunc = func(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
		var diags diag.Diagnostics
		return sdkdiag.AppendErrorf(diags, "read error")
	}

	diags := interceptedHandler(interceptors, read, Read)(context.Background(), nil, 42)
	if got, want := len(diags), 1; got != want {
		t.Errorf("length of diags = %v, want %v", got, want)
	}
}
