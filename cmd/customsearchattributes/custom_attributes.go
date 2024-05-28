package customsearchattributes

import (
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

var WorkflowStatusSearchAttribute = temporal.NewSearchAttributeKeyKeyword("WorkflowStatus")

func SetOnHold(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("ON_HOLD"))
}

func SetRunning(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("RUNNING"))

}

func SetCompleted(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("COMPLETED"))
}
