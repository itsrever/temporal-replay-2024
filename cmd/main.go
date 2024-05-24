package main

import (
	"errors"
	"log"
	"time"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"go.temporal.io/sdk/workflow"
)

var WorkflowStatusSearchAttribute = temporal.NewSearchAttributeKeyKeyword("WorkflowStatus")

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create temporal client", err)
	}
	defer c.Close()

	w := worker.New(c, "greetings", worker.Options{})

	w.RegisterWorkflow(Greet)

	activities := &Activities{}
	w.RegisterActivity(activities)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func Greet(ctx workflow.Context, name string) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 2,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	TryThings(ctx, func() error {
		return workflow.ExecuteActivity(ctx, (&Activities{}).GreetActivity).Get(ctx, nil)
	})

	SetCompleted(ctx)

	return nil
}

func TryThings(ctx workflow.Context, c func() error) {
	err := c()
	if err == nil {
		return
	}

	isRetry := OnActivityFailed(ctx)
	if isRetry {
		TryThings(ctx, c)
	}
}

func SetOnHold(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("ON_HOLD"))
}

func SetRunning(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("RUNNING"))

}

func SetCompleted(ctx workflow.Context) {
	_ = workflow.UpsertTypedSearchAttributes(ctx, WorkflowStatusSearchAttribute.ValueSet("COMPLETED"))
}

func OnActivityFailed(ctx workflow.Context) bool {
	SetOnHold(ctx)
	isRetry := false
	selector := workflow.NewSelector(ctx)
	selector.AddReceive(workflow.GetSignalChannel(ctx, "retry"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		isRetry = true
	})
	selector.AddReceive(workflow.GetSignalChannel(ctx, "skip"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
		isRetry = false
	})
	selector.Select(ctx)
	SetRunning(ctx)

	return isRetry
}

type Activities struct{}

func (a *Activities) GreetActivity() error {
	return errors.New("failed to say greeting")
	// return nil
}
