# ISP App OpenWrt Agent

A Go agent for OpenWrt devices that collects system and network information and reports it to the ISP App platform.

## Features

- WebSocket communication with the ISP App platform
- System statistics collection (CPU, memory, disk)
- Network monitoring (interfaces, traffic, latency)
- Automatic reconnection with exponential backoff
- Remote command execution support
- OpenWrt UCI configuration integration

## Development

### Prerequisites

- Go 1.21 or higher
- For hot reload: fswatch or inotifywait
- OpenWrt SDK (for cross-compilation)

### Project Structure

```
ispapp-agent/
├── cmd/
│   └── agent/           # Entry point for the agent
├── internal/
│   ├── agent/           # Core agent logic
│   ├── config/          # Configuration management
│   ├── handlers/        # Feature handlers
│   │   ├── network/     # Network monitoring
│   │   ├── system/      # System stats collection
│   │   └── websocket/   # Communication with server
│   ├── tools/           # Utility functions
│   └── websocket/       # WebSocket client
├── files/               # OpenWrt-specific files
│   ├── ispapp           # UCI config template
│   └── ispapp-agent.init # Init script
├── Makefile             # OpenWrt build system
├── Makefile.dev         # Development tasks
└── README.md            # This file
```

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
make -f Makefile.dev debug DEVICE_IP=192.168.1.1 DEVICE_USER=root DEVICE_PORT=22
```

This will:
1. Build the agent for the target architecture
2. Deploy it to the device via SSH
3. Run it and display logs
4. Watch for code changes and automatically rebuild/redeploy

### Running Tests

```bash
make -f Makefile.dev test
```

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

To edit the configuration through UCI:

```bash
# Edit configuration
uci set ispapp.settings.updateInterval='5'
uci commit ispapp

# Restart service to apply changes
/etc/init.d/ispapp-agent restart
```

## Service Management

Start/stop/restart the agent:
```bash
/etc/init.d/ispapp-agent start
/etc/init.d/ispapp-agent stop
/etc/init.d/ispapp-agent restart
```

Enable/disable autostart:
```bash
/etc/init.d/ispapp-agent enable
/etc/init.d/ispapp-agent disable
```

Check status:
```bash
/etc/init.d/ispapp-agent status
```

## License

MIT
