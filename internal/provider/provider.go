package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ provider.Provider = &GoogleSheetsProvider{}

type GoogleSheetsProvider struct {
	version string
}

type GoogleSheetsProviderModel struct {
	CredentialsJson types.String `tfsdk:"credentials_json"`
}

func (p *GoogleSheetsProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "googlesheets"
	resp.Version = p.version
}

func (p *GoogleSheetsProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"credentials_json": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
		},
	}
}

func (p *GoogleSheetsProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data GoogleSheetsProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var creds *google.Credentials
	var err error

	if !data.CredentialsJson.IsNull() {
		creds, err = google.CredentialsFromJSON(ctx, []byte(data.CredentialsJson.ValueString()), sheets.SpreadsheetsReadonlyScope)
	} else {
		creds, err = google.FindDefaultCredentials(ctx, sheets.SpreadsheetsReadonlyScope)
	}

	if err != nil {
		resp.Diagnostics.AddError("Unable to Load Credentials", err.Error())
		return
	}

	svr, err := sheets.NewService(ctx, option.WithCredentials(creds))

	if err != nil {
		resp.Diagnostics.AddError("Unable to Create Service", err.Error())
		return
	}

	resp.DataSourceData = svr
	resp.ResourceData = svr
}

func (p *GoogleSheetsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// No Resources
	}
}

func (p *GoogleSheetsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSheetDataSource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GoogleSheetsProvider{
			version: version,
		}
	}
}
