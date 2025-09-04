package provider

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"google.golang.org/api/sheets/v4"
)

var _ ephemeral.EphemeralResourceWithConfigure = &SheetEphemeralResource{}

func NewSheetEphemeralResource() ephemeral.EphemeralResource {
	return &SheetEphemeralResource{}
}

type SheetEphemeralResource struct {
	service *sheets.Service
}

type SheetEphemeralResourceModel struct {
	SheetId types.String `tfsdk:"sheet_id"`
	Range   types.String `tfsdk:"range"`
	Json    types.String `tfsdk:"json"`
}

func (d *SheetEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sheet"
}

func (d *SheetEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"sheet_id": schema.StringAttribute{
				Required: true,
			},
			"range": schema.StringAttribute{
				Required: true,
			},
			"json": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
		},
	}
}

func (d *SheetEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (d *SheetEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data SheetEphemeralResourceModel
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
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}
