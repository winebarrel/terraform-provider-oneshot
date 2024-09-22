package provider

import (
	"context"
	"fmt"

	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/winebarrel/terraform-provider-oneshot/internal/util"
)

var _ resource.ResourceWithModifyPlan = &RunResource{}

func NewRunResource() resource.Resource {
	return &RunResource{}
}

type RunResource struct {
	defaultShell string
}

type RunResourceModel struct {
	Command       types.String `tfsdk:"command"`
	PlanCommand   types.String `tfsdk:"plan_command"`
	Shell         types.String `tfsdk:"shell"`
	Stdout        types.String `tfsdk:"stdout"`
	Stderr        types.String `tfsdk:"stderr"`
	PlanStdoutLog types.String `tfsdk:"plan_stdout_log"`
	PlanStderrLog types.String `tfsdk:"plan_stderr_log"`
	RunAt         types.String `tfsdk:"run_at"`
}

func (data RunResourceModel) Run(shell string) (string, string, error) {
	if !data.Shell.IsNull() {
		shell = data.Shell.ValueString()
	}

	cmd := util.NewCmd(shell)
	return cmd.Run(data.Command.ValueString())
}

func (data RunResourceModel) Plan(shell string) (string, string, error) {
	if !data.Shell.IsNull() {
		shell = data.Shell.ValueString()
	}

	cmd := util.NewCmdWithLog(shell, data.PlanStdoutLog.ValueString(), data.PlanStderrLog.ValueString())
	return cmd.Run(data.PlanCommand.ValueString(), "ONESHOT_PLAN=1")
}

func (r *RunResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_run"
}

func (r *RunResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				MarkdownDescription: "Command to execute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"plan_command": schema.StringAttribute{
				MarkdownDescription: "Command to plan.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shell": schema.StringAttribute{
				MarkdownDescription: "Shell to execute the command.",
				Optional:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"stdout": schema.StringAttribute{
				MarkdownDescription: "Stdout of the command.",
				Computed:            true,
			},
			"stderr": schema.StringAttribute{
				MarkdownDescription: "Stderr of the command.",
				Computed:            true,
			},
			"plan_stdout_log": schema.StringAttribute{
				MarkdownDescription: "Stdout log file of the plan command.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("plan-stdout.log"),
			},
			"plan_stderr_log": schema.StringAttribute{
				MarkdownDescription: "Stderr log file of the plan command.",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("plan-stderr.log"),
			},
			"run_at": schema.StringAttribute{
				MarkdownDescription: "Command execution time.",
				Computed:            true,
			},
		},
	}
}

func (r *RunResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	providerData, ok := req.ProviderData.(OneshotProviderModel)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected OneshotProviderModel, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
	}

	r.defaultShell = providerData.DefaultShell.ValueString()
}

func (r *RunResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data RunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	stdout, stderr, err := data.Run(r.defaultShell)

	if err != nil {
		resp.Diagnostics.AddError("Run Command Error", fmt.Sprintf("Unable to run command, got error: %s", err))
	}

	data.Stdout = types.StringValue(stdout)
	data.Stderr = types.StringValue(stderr)
	data.RunAt = types.StringValue(time.Now().Local().String())

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *RunResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Nothing to do
}

func (r *RunResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Nothing to do
}

func (r *RunResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.State.RemoveResource(ctx)
}

func (r *RunResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		resp.Diagnostics.AddWarning(
			"Resource Destruction Considerations",
			"Applying this resource destruction will only remove the resource from the Terraform state "+
				"and will not call the deletion API due to API limitations. Manually use the web "+
				"interface to fully destroy this resource.",
		)
		return
	}

	var data RunResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.PlanCommand.IsNull() {
		return
	}

	_, _, err := data.Plan(r.defaultShell)

	if err != nil {
		resp.Diagnostics.AddError("Plan Command Error", fmt.Sprintf("Unable to plan command, got error: %s", err))
	}
}
