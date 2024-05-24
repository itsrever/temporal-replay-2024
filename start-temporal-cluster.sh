#!/bin/bash

# Function to handle the SIGINT signal (CTRL+C)
# and stop the Temporal server
cleanup() {
    echo "Stopping Temporal server..."
    kill $TEMPORAL_PID
    exit 0
}

# Start Temporal server in the background
temporal server start-dev &
TEMPORAL_PID=$!

# Setup a trap to catch SIGINT (CTRL+C) and call the cleanup function
trap cleanup SIGINT

# Wait a bit for the Temporal server to initialize
sleep 5

# Run your command
temporal operator search-attribute create --name WorkflowStatus --type Keyword

echo "Temporal server started successfully, you can now start the worker"

# Wait for user to press CTRL+C
wait $TEMPORAL_PID