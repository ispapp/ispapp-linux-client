#!/bin/bash

# Configuration
REMOTE_HOST="$1"  # First argument is the device IP
REMOTE_USER="${2:-root}"  # Second argument is the user, default to root
REMOTE_PORT="${3:-22}"  # Third argument is the SSH port, default to 22
REMOTE_PATH="/tmp/ispapp-agent"  # Remote path to deploy to
LOCAL_BINARY="./dist/agent"  # Local path to the binary
SOURCE_DIR="/mnt/c/Users/lenovo/ispapp-linux-client/ispapp-agent"  # Source directory to watch

# Check if device IP is provided
if [ -z "$REMOTE_HOST" ]; then
    echo "Usage: $0 <device-ip> [user] [port]"
    exit 1
fi
# Ensure we have fswatch or inotifywait for file watching
WATCH_CMD=""
if command -v fswatch >/dev/null 2>&1; then
    WATCH_CMD="fswatch -o $SOURCE_DIR"
elif command -v inotifywait >/dev/null 2>&1; then
    WATCH_CMD="inotifywait -m -r -e modify,create,delete $SOURCE_DIR"
else
    echo "Error: fswatch or inotifywait is required for hot-reloading"
    exit 1
fi

# Function to build and deploy
build_and_deploy() {
    echo "Building agent..."
    GOOS=linux GOARCH=arm go build -o $LOCAL_BINARY cmd/agent/main.go
    
    if [ $? -ne 0 ]; then
        echo "Build failed!"
        return 1
    fi
    
    # Stop any running instance first
    echo "Stopping any running agent on $REMOTE_HOST..."
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "killall -q agent || true"
    
    # Wait a moment to ensure process is completely terminated
    sleep 1
    
    echo "Copying binary to $REMOTE_HOST..."
    # Copy to a temporary file first to avoid "text file busy" error
    scp -P $REMOTE_PORT $LOCAL_BINARY $REMOTE_USER@$REMOTE_HOST:${REMOTE_PATH}.new
    
    if [ $? -ne 0 ]; then
        echo "Failed to copy binary to remote host!"
        return 1
    fi
    
    # Move the new file into place
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "mv ${REMOTE_PATH}.new $REMOTE_PATH && chmod +x $REMOTE_PATH"
    
    # Copy the configuration if it exists
    if [ -f "./dist/config/ispapp" ]; then
        echo "Copying configuration to $REMOTE_HOST..."
        # Create the config directory if it doesn't exist
        ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "mkdir -p /etc/config"
        scp -P $REMOTE_PORT ./dist/config/ispapp $REMOTE_USER@$REMOTE_HOST:/etc/config/ispapp
    fi
    
    # Run the binary in the background and capture output
    echo "Starting agent on $REMOTE_HOST..."
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "chmod +x $REMOTE_PATH && $REMOTE_PATH > /tmp/agent.log 2>&1 &"
    
    # Start showing logs
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "tail -f /tmp/agent.log"
}

# Initial build and deploy
build_and_deploy

# Set up background SSH connection to follow logs
LOG_PID=""
follow_logs() {
    if [ ! -z "$LOG_PID" ]; then
        kill $LOG_PID 2>/dev/null || true
    fi
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "tail -f /tmp/agent.log" &
    LOG_PID=$!
}

# Watch for changes and rebuild/redeploy
echo "Watching for changes. Press Ctrl+C to stop."
$WATCH_CMD | while read event; do
    echo "Change detected: $event"
    # Kill log follower
    if [ ! -z "$LOG_PID" ]; then
        kill $LOG_PID 2>/dev/null || true
        LOG_PID=""
    fi
    
    # Stop remote process with a delay to ensure termination
    ssh -p $REMOTE_PORT $REMOTE_USER@$REMOTE_HOST "killall -q agent || true"
    sleep 1
    
    # Rebuild and redeploy
    build_and_deploy
    
    # Restart log following
    follow_logs
done
