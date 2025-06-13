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
checkinternet() {
    # Check if the internet is reachable
    if ! ping -c 1 -W 1 8.8.8.8 >/dev/null 2>&1; then
        error "Internet connection is not available"
        sleep 15
        # retry checking internet connection
        log "Retrying to check internet connection..."
        checkinternet
    fi
}

# Check if running as root
if [ "$(id -u)" -ne 0 ]; then
    error "This script must be run as root"
    exit 1
fi

checkinternet

log "Starting ISPApp installation on OpenWrt..."

# Configure opkg repositories and settings
log "Configuring opkg repositories..."

# Backup existing opkg.conf if it exists
if [ -f /etc/opkg.conf ]; then
    cp /etc/opkg.conf /etc/opkg.conf.backup
    log "Backed up existing opkg.conf"
fi
# Clear existing opkg configuration
echo "" > /etc/opkg.conf
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
# Create distfeeds.conf if it doesn't exist
if [ ! -f /etc/opkg/distfeeds.conf ]; then
    touch /etc/opkg/distfeeds.conf
    log "Created distfeeds.conf"
else
    log "distfeeds.conf already exists, will append repositories"
fi
# Clear existing distfeeds.conf
echo "" > /etc/opkg/distfeeds.conf
# Add the default OpenWrt repositories
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
# download and install the latest version of banner
log "Installing banner package..."
BANNER_URL="https://github.com/ispapp/ispapp-linux-client/raw/refs/heads/websocket/banner"
BANNER_PKG="/etc/banner"
if ! wget "$BANNER_URL" -O "$BANNER_PKG"; then
    error "Failed to download banner package"
fi
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
ISPAPP_DOMAIN="prv.cloud.ispapp.co"
ISPAPP_PORT="443"

# UCI configuration
uci set ispapp.@settings[0].enabled=1
uci set ispapp.@settings[0].Domain="$ISPAPP_DOMAIN"
uci set ispapp.@settings[0].ListenerPort="$ISPAPP_PORT"
uci set ispapp.@settings[0].updateInterval=1
uci set ispapp.@settings[0].connected=0
# check if fw_setenv is available
if fw_printenv Key >/dev/null 2>&1; then
    log "Setting firmware environment variable..."
    uci set ispapp.@settings[0].Key="$(fw_printenv Key | cut -d'=' -f2)"
else
    warning "fw_setenv not available, skipping firmware environment variable"
fi
if fw_printenv Domain >/dev/null 2>&1; then
    log "Setting firmware environment variable for Domain..."
    uci set ispapp.@settings[0].Domain="$(fw_printenv Domain | cut -d'=' -f2)"
else
    warning "fw_setenv not available, skipping firmware environment variable for Domain"
fi
if fw_printenv AccessToken >/dev/null 2>&1; then
    log "Setting firmware environment variable for AccessToken..."
    uci set ispapp.@settings[0].accessToken="$(fw_printenv AccessToken | cut -d'=' -f2)"
else
    warning "fw_setenv not available, skipping firmware environment variable for AccessToken"
fi
if fw_printenv RefreshToken >/dev/null 2>&1; then
    log "Setting firmware environment variable for RefreshToken..."
    uci set ispapp.@settings[0].refreshToken="$(fw_printenv RefreshToken | cut -d'=' -f2)"
else
    warning "fw_setenv not available, skipping firmware environment variable for RefreshToken"
fi
uci commit ispapp

# Set environment variable (if fw_setenv is available)
if command -v fw_setenv >/dev/null 2>&1; then
    log "Setting firmware environment variable..."
    fw_setenv Key="$ISPAPP_KEY"
else
    warning "fw_setenv not available, skipping firmware environment variable"
fi
/etc/init.d/rpcd reload || {
    error "Failed to reload rpcd service"ß
}
if [ -f /etc/rc.d/S99ispapp ]; then
    log "disable unused service"
    /etc/rc.d/S99ispapp disable
else
    error "Failed to disable ISPApp service"
    exit 1
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
uci get ispapp.@settings[0].Key >/dev/null 2>&1 || {
    error "ISPApp configuration not found, installation may have failed"
    exit 1
}
uci get ispapp.@settings[0].Key && {
     log "ISPApp configuration found successfully"
    ISPAPP_KEY=$(uci get ispapp.@settings[0].Key)
    ISPAPP_DOMAIN=$(uci get ispapp.@settings[0].Domain)
    ISPAPP_AccessToken=$(uci get ispapp.@settings[0].accessToken)
    ISPAPP_RefreshToken=$(uci get ispapp.@settings[0].refreshToken)
    fw_setenv Key="$ISPAPP_KEY" || {
        warning "Failed to set firmware environment variable"
    }
    fw_setenv Domain="$ISPAPP_DOMAIN" || {
        warning "Failed to set firmware environment variable for Domain"
    }
    fw_setenv AccessToken="$ISPAPP_AccessToken" || {
        warning "Failed to set firmware environment variable for AccessToken"
    }
    fw_setenv RefreshToken="$ISPAPP_RefreshToken" || {
        warning "Failed to set firmware environment variable for RefreshToken"
    }
    log "ISPApp configuration set successfully"
}
fi
log "ISPApp installation completed!"
log "Configuration:"
log "  Key: $ISPAPP_KEY"
log "  Domain: $ISPAPP_DOMAIN"
log "  Port: $ISPAPP_PORT"
log ""
log "You can check the service status with: /etc/init.d/ispapp status"
log "View logs with: logread | grep ispapp"


