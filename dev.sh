#!/bin/bash

# NexusTerm One-Click Dev Script
# Starts Backend (Go) and Frontend (Vite) concurrently

# Function to kill child processes on exit
cleanup() {
    echo ""
    echo "Shutting down servers..."
    kill $(jobs -p) 2>/dev/null
    wait
    echo "Done."
}

# Trap SIGINT (Ctrl+C)
trap cleanup SIGINT

echo ">>> Starting NexusTerm Development Environment <<<"

# Start Backend
echo "[1/2] Starting Backend..."
cd backend
if [ ! -f "go.mod" ]; then
    echo "Error: go.mod not found in backend/"
    exit 1
fi
# Run in background
go run main.go &
BACKEND_PID=$!
cd ..

# Start Frontend
echo "[2/2] Starting Frontend..."
cd frontend
if [ ! -d "node_modules" ]; then
    echo "Installing frontend dependencies..."
    npm install
fi
# Run in background
npm run dev &
FRONTEND_PID=$!
cd ..

echo ">>> Services Started <<<"
echo "Backend PID: $BACKEND_PID"
echo "Frontend PID: $FRONTEND_PID"
echo ""
echo "Press Ctrl+C to stop both services."

# Wait for both processes
wait
