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

func Greet(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 2,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	TryActivity(ctx, name)

	return "Hello, " + name, nil
}

func TryActivity(ctx workflow.Context, name string) {
	err := workflow.ExecuteActivity(ctx, SayGreetingReference, name).Get(ctx, nil)
	if err == nil {
		return
	}

	isRetry := OnActivityFailed(ctx)
	if isRetry {
		TryActivity(ctx, name)
	}
}

func OnActivityFailed(ctx workflow.Context) bool {
	attributes := map[string]interface{}{
		"ReverStatus": "ON_HOLD",
	}
	// This won't persist search attributes on server because commands are not sent to server,
	// but local cache will be updated.
	_ = workflow.UpsertSearchAttributes(ctx, attributes)

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

	return isRetry
}

var SayGreetingReference = (&Activities{}).SayGreeting

type Activities struct{}

func (a *Activities) SayGreeting(name string) (string, error) {
	return "", errors.New("failed to say greeting")
}
