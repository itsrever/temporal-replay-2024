install:
	brew install temporal

start:
	./start-temporal.sh

init:
	go run ./cmd/main.go

run-wf:
	temporal workflow start --task-queue greetings --type Greet


# temporal workflow describe --workflow-id 15e9c0c6-48ce-42a0-a897-1f27a67cb9b0