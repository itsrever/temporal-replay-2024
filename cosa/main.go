package main

import (
	"context"
	"log"

	"go.temporal.io/sdk/client"
)

func main() {
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create temporal client", err)
	}
	defer c.Close()

	opt := client.StartWorkflowOptions{
		TaskQueue: "greetings",
	}
	_, _ = c.SignalWithStartWorkflow(context.Background(), "workflow-id", "signal-name", nil, opt, "Greet")
}
