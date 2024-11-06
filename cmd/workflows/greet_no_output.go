package workflows

import (
	"errors"
	"time"

	"go.temporal.io/sdk/workflow"
)

func Greet(ctx workflow.Context, name string) error {
	// ao := workflow.ActivityOptions{
	// 	StartToCloseTimeout: 10 * time.Second,
	// 	RetryPolicy: &temporal.RetryPolicy{
	// 		InitialInterval: time.Second,
	// 		MaximumAttempts: 2,
	// 	},
	// }
	// ctx = workflow.WithActivityOptions(ctx, ao)

	// onhold.ExecuteWrappedActivity(ctx, func() (Empty, error) {
	// 	return Empty{}, workflow.ExecuteActivity(ctx, (&Activities{}).GreetActivity).Get(ctx, nil)
	// })

	// customsearchattributes.SetCompleted(ctx)

	workflow.Sleep(ctx, time.Duration(10)*time.Second)

	return nil
}

type Activities struct{}

func (a *Activities) GreetActivity() error {
	return errors.New("failed to say greeting")
	// return nil
}

type Empty struct{}
