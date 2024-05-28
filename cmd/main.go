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

	ExecuteWrappedActivity(ctx, func() (Empty, error) {
		return Empty{}, workflow.ExecuteActivity(ctx, (&Activities{}).GreetActivity).Get(ctx, nil)
	})

	SetCompleted(ctx)

	return nil
}

type Activities struct{}

func (a *Activities) GreetActivity() error {
	return errors.New("failed to say greeting")
	// return nil
}

type Empty struct{}
