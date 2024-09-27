package main

import (
	"context"
	"fmt"

	"go.temporal.io/api/enums/v1"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
)

func main() {
	clientOptions := client.Options{
		Namespace:         "default",
		HostPort:          "localhost:7233",
		ConnectionOptions: client.ConnectionOptions{},
	}
	c, err := client.Dial(clientOptions)
	if err != nil {
		panic(err)
	}
	defer c.Close()

	resp, err := c.CountWorkflow(context.Background(), &workflowservice.CountWorkflowExecutionsRequest{})
	if err != nil {
		panic(err)
	}

	fmt.Println(resp.GetCount())

	query := "WorkflowStatus='ON_HOLD'"
	list, err := c.ListWorkflow(context.Background(), &workflowservice.ListWorkflowExecutionsRequest{
		Query:    query,
		PageSize: 500,
	})
	if err != nil {
		panic(err)
	}

	for _, w := range list.GetExecutions() {
		exec, err := c.DescribeWorkflowExecution(context.Background(), w.GetExecution().GetWorkflowId(), w.GetExecution().GetRunId())
		if err != nil {
			panic(err)
		}

		if exec.GetWorkflowExecutionInfo().GetStatus() != enums.WORKFLOW_EXECUTION_STATUS_RUNNING {
			continue
		}

		err = c.TerminateWorkflow(context.Background(), w.GetExecution().GetWorkflowId(), w.GetExecution().GetRunId(), "hola")
		if err != nil {
			panic(err)
		}
	}

	// func (client.Client) CountWorkflow(ctx context.Context, request *workflowservice.CountWorkflowExecutionsRequest) (*workflowservice.CountWorkflowExecutionsResponse, error)
}

// lastEvent := hist.Events[len(hist.Events)-1]
// return lastEvent.GetEventType() == enums.EVENT_TYPE_WORKFLOW_EXECUTION_COMPLETED
