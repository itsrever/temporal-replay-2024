install-temporal-cli:
	brew install temporal

# Starts Temporal Cluster locally and sets up WorkflowStatus search attribute:
# 	- Frontend: localhost:7233
# 	- UI: localhost:8233
start-temporal-cluster:
	./start-temporal-cluster.sh

start-worker:
	go run ./cmd/main.go

run-workflow:
	temporal workflow start --task-queue greetings --type Greet


# temporal workflow describe --workflow-id 15e9c0c6-48ce-42a0-a897-1f27a67cb9b0