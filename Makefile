install:
	brew install temporal

start:
	temporal server start-dev

init:
	temporal operator search-attribute create --name ReverStatus --type Keyword
	go run ./cmd/main.go