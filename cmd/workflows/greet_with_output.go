package workflows

import (
	"errors"
	"time"

	"github.com/ericvg97/temporal-replay/cmd/customsearchattributes"
	"github.com/ericvg97/temporal-replay/cmd/onhold"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

func GreetWithOutput(ctx workflow.Context, name string) (string, error) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 2,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	output := onhold.ExecuteWrappedActivity(ctx, func() (Output, error) {
		output := Output{}
		err := workflow.ExecuteActivity(ctx, (&Activities{}).GreetActivityWithOutput).Get(ctx, &output)
		return output, err
	})

	customsearchattributes.SetCompleted(ctx)

	return output.GetGreeting(), nil
}

type Output struct {
	Greeting string `json:"greeting"`
}

func (o *Output) GetGreeting() string {
	if o == nil {
		return ""
	}

	return o.Greeting
}

func (a *Activities) GreetActivityWithOutput() (Output, error) {
	return Output{}, errors.New("failed to say greeting")
	// return Output{"hello"}, nil
}
