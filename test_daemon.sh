#!/bin/bash

# Test script for ispappd daemon
DAEMON_PATH="/Users/kimo/Desktop/kmoz000/ispapp/ispapp-linux-client/luci-app-ispapp/root/bin/ispappd"

echo "Testing ispappd daemon..."
echo "========================="

# Make sure the daemon is executable
chmod +x "$DAEMON_PATH"

echo "1. Testing daemon syntax..."
if lua -c "$DAEMON_PATH" 2>/dev/null; then
    echo "✓ Daemon syntax is valid"
else
    echo "✗ Daemon syntax has errors"
    lua -c "$DAEMON_PATH"
    exit 1
fi

echo ""
echo "2. Testing daemon commands..."

echo "Testing stop command (should be safe even if not running)..."
lua "$DAEMON_PATH" stop

echo ""
echo "Testing cleanup command..."
lua "$DAEMON_PATH" cleanup

echo ""
echo "Testing help..."
lua "$DAEMON_PATH"

echo ""
echo "✓ All tests completed!"
echo "Note: To fully test the daemon, you would need to run it in a proper OpenWrt environment"
echo "      with the required dependencies (nixio, uci, etc.)"
