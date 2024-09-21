package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const (
	DefaultShell = "/bin/bash -c"
)

var _ provider.Provider = &OneshotProvider{}

type OneshotProvider struct {
	version string
}

type OneshotProviderModel struct {
	DefaultShell types.String `tfsdk:"default_shell"`
}

func (p *OneshotProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "oneshot"
	resp.Version = p.version
}

func (p *OneshotProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"default_shell": schema.StringAttribute{
				MarkdownDescription: "Default shell to execute the command. (default: " + DefaultShell + ")",
				Optional:            true,
			},
		},
	}
}

func (p *OneshotProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data OneshotProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.DefaultShell.IsNull() {
		data.DefaultShell = types.StringValue(DefaultShell)
	}

	resp.DataSourceData = data
	resp.ResourceData = data
}

func (p *OneshotProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewRunResource,
	}
}

func (p *OneshotProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		// No Data Sources
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &OneshotProvider{
			version: version,
		}
	}
}
