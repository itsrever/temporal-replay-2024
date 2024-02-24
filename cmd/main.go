package main

import (
	"log"

	"go.temporal.io/sdk/client"
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

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}

func Greet(ctx workflow.Context, name string) (string, error) {
	attributes := map[string]interface{}{
		"ReverStatus": "ON_HOLD",
	}
	// This won't persist search attributes on server because commands are not sent to server,
	// but local cache will be updated.
	err := workflow.UpsertSearchAttributes(ctx, attributes)
	if err != nil {
		return "", err
	}

	selector := workflow.NewSelector(ctx)
	selector.AddReceive(workflow.GetSignalChannel(ctx, "retry"), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, nil)
	})

	selector.Select(ctx)

	return "Hello, " + name, nil
}
