install-temporal-cli:
	brew install temporal

start-temporal-cluster:
	./start-temporal-cluster.sh

start-worker:
	go run ./cmd/main.go

run-greet-workflow:
	temporal workflow start --task-queue greetings --type Greet

run-greet-workflow-with-output:
	temporal workflow start --task-queue greetings --type GreetWithOutput
