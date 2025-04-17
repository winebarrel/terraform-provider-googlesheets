package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/sheets/v4"
)

var _ datasource.DataSourceWithConfigure = &SheetDataSource{}

func NewSheetDataSource() datasource.DataSource {
	return &SheetDataSource{}
}

type SheetDataSource struct {
	service *sheets.Service
}

type SheetDataSourceModel struct {
	SheetId       types.String `tfsdk:"sheet_id"`
	Range         types.String `tfsdk:"range"`
	Json          types.String `tfsdk:"json"`
	SensitiveJson types.String `tfsdk:"sensitive_json"`
}

func (d *SheetDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sheet"
}

func (d *SheetDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sheet_id": schema.StringAttribute{
				Required: true,
			},
			"range": schema.StringAttribute{
				Required: true,
			},
			"json": schema.StringAttribute{
				Computed: true,
			},
			"sensitive_json": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (d *SheetDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	svr, ok := req.ProviderData.(*sheets.Service)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *sheets.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	d.service = svr
}

func (d *SheetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data SheetDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	sheetId := data.SheetId.ValueString()
	range_ := data.Range.ValueString()

	valRange, err := d.service.Spreadsheets.Values.Get(sheetId, range_).Do()

	if err != nil {
		resp.Diagnostics.AddError("Error getting values", err.Error())
		return
	}

	rawJson, err := json.Marshal(valRange.Values)

	if err != nil {
		resp.Diagnostics.AddError("Error marshalling values", err.Error())
		return
	}

	jsonStr := types.StringValue(string(rawJson))
	data.Json = jsonStr
	data.SensitiveJson = jsonStr
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
