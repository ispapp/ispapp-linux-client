#!/bin/sh
#
# OpenWrt Auto-install Script for ISPApp
# This script automates the installation and configuration of ISPApp on OpenWrt
#

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Logging function
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[ERROR] $1${NC}" >&2
}

warning() {
    echo -e "${YELLOW}[WARNING] $1${NC}"
}

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    error "This script must be run as root"
    exit 1
fi


log "Starting ISPApp installation on OpenWrt..."

# Configure opkg repositories and settings
log "Configuring opkg repositories..."

# Backup existing opkg.conf if it exists
if [ -f /etc/opkg.conf ]; then
    cp /etc/opkg.conf /etc/opkg.conf.backup
    log "Backed up existing opkg.conf"
fi

# Create new opkg.conf with proper configuration
cat > /etc/opkg.conf << 'EOF'
dest root /
dest ram /tmp
lists_dir ext /var/opkg-lists
option overlay_root /overlay
option check_signature 1
option force_space
arch aarch64_cortex-a73_neon-vfpv4 10
arch aarch64_generic 1
arch all 10
arch noarch 1
EOF
cat >> /etc/opkg/distfeeds.conf << 'EOF'
src/gz openwrt_base https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/base
src/gz openwrt_packages https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/packages
src/gz openwrt_telephony https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/telephony
src/gz openwrt_luci https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/luci
src/gz openwrt_routing https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/routing
src/gz openwrt_freifunk https://downloads.openwrt.org/releases/19.07.3/packages/aarch64_generic/freifunk
EOF

log "Updated opkg.conf with proper repositories"

# Update package lists
log "Updating package lists..."
opkg update

# Create necessary directories and files
log "Creating package info files..."
mkdir -p /usr/lib/opkg/info
touch /usr/lib/opkg/info/coreutils.list
touch /usr/lib/opkg/info/ca-bundle.list
touch /usr/lib/opkg/info/ca-certificates.list

# Install core dependencies
log "Installing core dependencies..."
opkg install ca-bundle ca-certificates coreutils

# Install Lua dependencies
log "Installing Lua dependencies..."
opkg install luasocket luasec

# Download and install ISPApp package
log "Downloading ISPApp package..."
ISPAPP_URL="https://github.com/ispapp/ispapp-linux-client/releases/download/skynet%2Faarch64_generic-23.05-SNAPSHOT/luci-app-ispapp_1.0.0_all.ipk"
ISPAPP_PKG="/tmp/luci-app-ispapp_1.0.0_all.ipk"

if ! wget "$ISPAPP_URL" -O "$ISPAPP_PKG"; then
    error "Failed to download ISPApp package"
    exit 1
fi

log "Installing ISPApp package..."
if ! opkg install "$ISPAPP_PKG" --nodeps; then
    error "Failed to install ISPApp package"
    exit 1
fi

# Set executable permissions
log "Setting executable permissions..."
chmod +x /usr/libexec/rpcd/ispapp
chmod +x /etc/init.d/ispapp
chmod +x /bin/ispappd

# Configure ISPApp
log "Configuring ISPApp..."

# Configuration variables (modify these as needed)
ISPAPP_KEY="IOZlWIahLacr8hZr"
ISPAPP_DOMAIN="prv.cloud.ispapp.co"
ISPAPP_PORT="443"

# UCI configuration
uci set ispapp.@settings[0].enabled=1
uci set ispapp.@settings[0].Key="$ISPAPP_KEY"
uci set ispapp.@settings[0].Domain="$ISPAPP_DOMAIN"
uci set ispapp.@settings[0].ListenerPort="$ISPAPP_PORT"
uci commit ispapp

# Set environment variable (if fw_setenv is available)
if command -v fw_setenv >/dev/null 2>&1; then
    log "Setting firmware environment variable..."
    fw_setenv Key="$ISPAPP_KEY"
else
    warning "fw_setenv not available, skipping firmware environment variable"
fi

# Enable and start ISPApp service
log "Enabling and starting ISPApp service..."
/etc/init.d/rpcd reload
/etc/init.d/ispapp enable
/etc/init.d/ispapp start

# Wait a moment for service to start
sleep 3

# Sync agent version and signup
log "Syncing agent version..."
if ! /usr/libexec/rpcd/ispapp call sync_agent_version; then
    warning "Failed to sync agent version"
fi

log "Performing signup..."
if ! /usr/libexec/rpcd/ispapp call signup; then
    warning "Failed to perform signup"
fi

# Cleanup
log "Cleaning up temporary files..."
rm -f "$ISPAPP_PKG"

# Verify installation
log "Verifying installation..."
if pgrep -f ispapp >/dev/null; then
    log "ISPApp is running successfully!"
else
    warning "ISPApp may not be running properly"
fi

log "ISPApp installation completed!"
log "Configuration:"
log "  Key: $ISPAPP_KEY"
log "  Domain: $ISPAPP_DOMAIN"
log "  Port: $ISPAPP_PORT"
log ""
log "You can check the service status with: /etc/init.d/ispapp status"
log "View logs with: logread | grep ispapp"
