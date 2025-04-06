# ISP App OpenWrt Agent

A Go agent for OpenWrt devices that collects system and network information.

## Development

### Prerequisites
- Go 1.21 or higher
- For hot reload: fswatch or inotifywait

### Building Locally
```bash
# Standard build for your platform
make -f Makefile.dev build

# Cross-compile for different OpenWrt targets
make -f Makefile.dev openwrt-mipsle  # For MIPS little-endian devices
make -f Makefile.dev openwrt-arm     # For ARM devices
make -f Makefile.dev openwrt-arm_cortex # For ARM devices
make -f Makefile.dev openwrt-targets # Build all supported platforms
```

### Debugging on a Device
```bash
make -f Makefile.dev debug DEVICE_IP=192.168.1.170 DEVICE_USER=root DEVICE_PORT=22
```

This will:
1. Build the agent for the target architecture
2. Deploy it to the device via SSH
3. Run it and display logs
4. Watch for code changes and automatically rebuild/redeploy

## OpenWrt SDK Integration

### Adding to an OpenWrt Feed
1. Add this repository to your feeds:
```bash
echo "src-git ispapp https://github.com/ispapp/ispapp-linux-client.git^websocket" >> feeds.conf.default
./scripts/feeds update ispapp
./scripts/feeds install ispapp-agent
```

2. Select the package in menuconfig:
```bash
make menuconfig
# Go to Network -> ispapp-agent
```

3. Build the package:
```bash
make package/ispapp-agent/compile
```

### Building Directly in OpenWrt
1. Copy this repository to the OpenWrt source tree:
```bash
mkdir -p package/utils/ispapp-agent
cp -r /path/to/ispapp-agent/* package/utils/ispapp-agent/
```

2. Select and build the package as described above.

### Installing on an OpenWrt Device
Transfer the compiled IPK and install:
```bash
opkg install ispapp-agent_1.0.0-1_mipsel_24kc.ipk
```

## Configuration

The agent uses OpenWrt's UCI configuration system. The configuration is stored in `/etc/config/ispapp`.

### UCI Configuration Format
```
package ispapp
config settings
        option enabled '1'           # Enable/disable the agent
        option login '00000000-0000-0000-0000-000000000000'  # Device ID
        option Domain 'prv.cloud.ispapp.co'  # API domain
        option ListenerPort '443'    # API port
        option connected '0'         # Connection status (managed by agent)
        option Key ''                # API key
        option refreshToken ''       # Refresh token
        option accessToken ''        # Access token
        option updateInterval '1'    # Update interval in minutes
        option IperfServer ''        # Iperf server for speed tests
        list pingTargets 'cloud.ispapp.co'  # Ping targets

config overview
        option enabled '1'           # Enable/disable the overview
        option highrefresh '0'       # Enable high refresh rate
        option refreshInterval '0'   # Refresh interval
```

### Editing Configuration
You can edit the configuration using UCI commands:

```bash
# Edit settings
uci set ispapp.settings.Domain=prv.cloud.ispapp.co
uci set ispapp.settings.login=your-device-id

# Add a ping target
uci add_list ispapp.settings.pingTargets=8.8.8.8

# Commit changes
uci commit ispapp

# Restart the agent to apply changes
/etc/init.d/ispapp-agent restart
```

## Service Management
```bash
# Start the service
/etc/init.d/ispapp-agent start

# Enable autostart
/etc/init.d/ispapp-agent enable

# Stop the service
/etc/init.d/ispapp-agent stop
```
