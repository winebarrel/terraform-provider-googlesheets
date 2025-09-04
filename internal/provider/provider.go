package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

var _ provider.ProviderWithEphemeralResources = &GoogleSheetsProvider{}

type GoogleSheetsProvider struct {
	version string
}

type GoogleSheetsProviderModel struct {
	CredentialsJson types.String `tfsdk:"credentials_json"`
	CredentialsEnv  types.String `tfsdk:"credentials_env"`
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
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("credentials_env")),
				},
			},
			"credentials_env": schema.StringAttribute{
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ConflictsWith(path.MatchRelative().AtParent().AtName("credentials_json")),
				},
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

	if !data.CredentialsEnv.IsNull() {
		envName := data.CredentialsEnv.ValueString()
		envValue := os.Getenv(envName)

		if envValue == "" {
			resp.Diagnostics.AddError("Unable to Get Credentials from environment variable", fmt.Sprintf("$%s is empty", envName))
			return
		}

		creds, err = google.CredentialsFromJSON(ctx, []byte(envValue), sheets.SpreadsheetsReadonlyScope)
	} else if !data.CredentialsJson.IsNull() {
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
	resp.EphemeralResourceData = svr
}

func (p *GoogleSheetsProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		// No Resources
	}
}

func (p *GoogleSheetsProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewSheetDataSource,
		NewSensitiveSheetDataSource,
	}
}

func (p *GoogleSheetsProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{
		NewSheetEphemeralResource,
	}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &GoogleSheetsProvider{
			version: version,
		}
	}
}
