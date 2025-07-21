#!/bin/bash

# Start the backend services for development

echo "Starting Resume Optimizer Services..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "Creating .env file from .env.example..."
    cp .env.example .env
    echo "Please update .env with your actual configuration values"
fi

# Load environment variables from .env file
if [ -f .env ]; then
    echo "Loading environment variables from .env file..."
    export $(grep -v '^#' .env | xargs)
fi

# Build services
echo "Building services..."
go build -o services/auth/auth-service ./services/auth
go build -o services/resume-processor/resume-processor-service ./services/resume-processor

# Start auth service in background
echo "Starting Auth Service on port 8080..."
cd services/auth
PORT=8080 ./auth-service &
AUTH_PID=$!

# Start resume processor service in background
echo "Starting Resume Processor Service on port 8081..."
cd ../resume-processor
PORT=8081 ./resume-processor-service &
RESUME_PID=$!

cd ../..

echo "Services started!"
echo "Auth Service PID: $AUTH_PID"
echo "Resume Processor Service PID: $RESUME_PID"
echo ""
echo "Auth Service: http://localhost:8080"
echo "Resume Processor Service: http://localhost:8081"
echo ""
echo "Press Ctrl+C to stop all services"

# Function to cleanup processes on exit
cleanup() {
    echo "Stopping services..."
    kill $AUTH_PID $RESUME_PID 2>/dev/null
    exit
}

trap cleanup INT

# Wait for processes to finish
wait
