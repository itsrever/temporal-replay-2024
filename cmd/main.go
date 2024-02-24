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

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
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
		// This needs to be idempotent
		err := workflow.ExecuteActivity(ctx, FirstActivityReference).Get(ctx, nil)
		if err != nil {
			return err
		}

		return workflow.ExecuteActivity(ctx, SecondActivityReference).Get(ctx, nil)
	})

	// TryThings(ctx, func() error {
	// 	return workflow.ExecuteActivity(ctx, SecondActivityReference).Get(ctx, nil)
	// })

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
	attributes := map[string]interface{}{
		"ReverStatus": "ON_HOLD",
	}
	_ = workflow.UpsertSearchAttributes(ctx, attributes)
}

func SetRunning(ctx workflow.Context) {
	attributes := map[string]interface{}{
		"ReverStatus": "RUNNING",
	}
	_ = workflow.UpsertSearchAttributes(ctx, attributes)
}

func SetCompleted(ctx workflow.Context) {
	attributes := map[string]interface{}{
		"ReverStatus": "COMPLETED",
	}
	_ = workflow.UpsertSearchAttributes(ctx, attributes)
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

var FirstActivityReference = (&Activities{}).FirstActivity
var SecondActivityReference = (&Activities{}).SecondActivity

type Activities struct{}

func (a *Activities) FirstActivity() error {
	// return errors.New("failed to say greeting")
	return nil
}

func (a *Activities) SecondActivity() error {
	return errors.New("failed to say greeting")
}
