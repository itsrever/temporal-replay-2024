install:
	brew install temporal

start:
	./start-temporal.sh

init:
	go run ./cmd/main.go

run-wf:
	temporal workflow start --task-queue greetings --type Greet --input '"Eric"'