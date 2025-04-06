package controllers

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CollectDeviceInfo gathers all device information and populates the Device struct
func (d *Device) CollectDeviceInfo() error {
	// Set basic device information
	hostname, err := os.Hostname()
	if err == nil {
		d.Hostname = hostname
	}

	// Collect uptime
	if uptime, err := d.GetUptime(); err == nil {
		d.Uptime = uptime
	}

	// Collect system information
	sysInfo, err := d.GetSystemInfo()
	if err == nil && d.Collectors == nil {
		d.Collectors = &Collector{}
		d.Collectors.System = sysInfo
	} else if err == nil {
		d.Collectors.System = sysInfo
	}

	// Collect operating system info
	osInfo, err := d.GetOSInfo()
	if err == nil {
		d.Os = osInfo.ID
		d.OsVersion = osInfo.Version
		d.OsBuildDate = osInfo.BuildDate
	}

	// Collect hardware info
	hwInfo, err := d.GetHardwareInfo()
	if err == nil {
		d.HardwareMake = hwInfo.hardwareMake
		d.HardwareModel = hwInfo.hardwareModel
		d.HardwareModelNumber = hwInfo.hardwareModelNumber
		d.HardwareSerialNumber = hwInfo.hardwareSerialNumber
		d.HardwareCpuInfo = hwInfo.hardwareCpuInfo
	}

	// Collect network interfaces
	interfaces, err := d.GetNetworkInterfaces()
	if err == nil {
		d.Interfaces = interfaces
		if d.Collectors != nil {
			d.Collectors.Interface = interfaces
		}
	}

	// Collect wireless configuration
	wirelessConfig, err := d.GetWirelessConfigured()
	if err == nil {
		d.WirelessConfigured = wirelessConfig
	}

	// Collect security profiles
	securityProfiles, err := d.GetSecurityProfiles()
	if err == nil {
		d.SecurityProfiles = securityProfiles
	}

	// Collect wireless access point info
	wapInfo, err := d.GetWapInfo()
	if err == nil && d.Collectors != nil {
		d.Collectors.Wap = wapInfo
	}

	// Collect TCP info
	tcpInfo, err := d.GetTcpInfo()
	if err == nil && d.Collectors != nil {
		d.Collectors.Tcp = tcpInfo
	}

	// Collect DHCP leases
	dhcpLeases, err := d.GetDhcpLeases()
	if err == nil && d.Collectors != nil {
		d.Collectors.Leases = dhcpLeases
	}

	// Collect ARP table
	arpTable, err := d.GetArpTable()
	if err == nil && d.Collectors != nil {
		d.Collectors.Arptable = arpTable
	}

	// Collect ping stats
	pingStats, err := d.GetPingStats("1.1.1.1")
	if err == nil && d.Collectors != nil {
		pings := []Ping{pingStats}
		d.Collectors.Ping = pings
	}

	// Get external IP information
	ipInfo, err := d.GetExternalIpInfo()
	if err == nil {
		d.OutsideIP = ipInfo.IP
		if ipInfo.Latitude != 0 && ipInfo.Longitude != 0 {
			d.Lat = &ipInfo.Latitude
			d.Lng = &ipInfo.Longitude

			// Set location collector
			if d.Collectors != nil {
				d.Collectors.LocationCollector = LocationCollector{
					Lat: ipInfo.Latitude,
					Lng: ipInfo.Longitude,
				}
			}
		}
	}

	// Collect WAN IP
	wanIP, err := d.GetWanIP()
	if err == nil {
		d.WanIp = wanIP
	}

	// Get device capabilities
	d.WebshellSupport = d.HasWebshellSupport()
	d.BandwidthTestSupport = d.HasBandwidthTestSupport()
	d.FirmwareUpgradeSupport = d.HasFirmwareUpgradeSupport()

	// Initialize empty collections if needed
	if d.Collectors != nil && d.Collectors.Topos == nil {
		d.Collectors.Topos = []Topo{}
	}

	return nil
}

// GetUptime retrieves the system uptime in seconds
func (d *Device) GetUptime() (int64, error) {
	content, err := os.ReadFile("/proc/uptime")
	if err != nil {
		return 0, err
	}

	fields := strings.Fields(string(content))
	if len(fields) < 1 {
		return 0, fmt.Errorf("invalid uptime format")
	}

	uptime, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return 0, err
	}

	return int64(uptime), nil
}

// GetSystemInfo collects system information like CPU load, memory usage, and disk usage
func (d *Device) GetSystemInfo() (System, error) {
	system := System{}

	// Get load average
	load, err := d.GetLoadAvg()
	if err == nil {
		system.Load = load
	}

	// Get memory info
	memory, err := d.GetMemoryInfo()
	if err == nil {
		system.Memory = memory
	}

	// Get disk info
	disks, err := d.GetDiskInfo()
	if err == nil {
		system.Disks = disks
	}

	return system, nil
}

// GetLoadAvg retrieves the system load average
func (d *Device) GetLoadAvg() (Load, error) {
	content, err := os.ReadFile("/proc/loadavg")
	if err != nil {
		return Load{}, err
	}

	fields := strings.Fields(string(content))
	if len(fields) < 4 {
		return Load{}, fmt.Errorf("invalid loadavg format")
	}

	one, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return Load{}, err
	}

	five, err := strconv.ParseFloat(fields[1], 64)
	if err != nil {
		return Load{}, err
	}

	fifteen, err := strconv.ParseFloat(fields[2], 64)
	if err != nil {
		return Load{}, err
	}

	// Process count is usually in the format "1/123"
	procParts := strings.Split(fields[3], "/")
	if len(procParts) < 2 {
		return Load{}, fmt.Errorf("invalid process count format")
	}

	procCount, err := strconv.ParseUint(procParts[1], 10, 64)
	if err != nil {
		return Load{}, err
	}

	return Load{
		One:          one,
		Five:         five,
		Fifteen:      fifteen,
		ProcessCount: procCount,
	}, nil
}

// GetMemoryInfo retrieves memory usage information
func (d *Device) GetMemoryInfo() (Memory, error) {
	content, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return Memory{}, err
	}

	var memory Memory
	memValues := make(map[string]uint64)

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}

		// Remove colon from key
		key := strings.TrimSuffix(parts[0], ":")
		value, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}

		memValues[key] = value
	}

	// Extract relevant memory values
	memory.Total = memValues["MemTotal"]
	memory.Free = memValues["MemFree"]
	memory.Buffers = memValues["Buffers"]
	memory.Cache = memValues["Cached"]

	return memory, nil
}

// GetDiskInfo retrieves information about disk usage
func (d *Device) GetDiskInfo() ([]Disk, error) {
	cmd := exec.Command("df", "-k", "/")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	if len(lines) < 2 {
		return nil, fmt.Errorf("unexpected df output format")
	}

	// Parse df output format
	fields := strings.Fields(lines[1])
	if len(fields) < 6 {
		return nil, fmt.Errorf("invalid df output format")
	}

	used, err := strconv.ParseUint(fields[2], 10, 64)
	if err != nil {
		return nil, err
	}

	avail, err := strconv.ParseUint(fields[3], 10, 64)
	if err != nil {
		return nil, err
	}

	mount := fields[5]

	return []Disk{
		{
			Mount: mount,
			Used:  used,
			Avail: avail,
		},
	}, nil
}

// GetNetworkInterfaces collects information about network interfaces
func (d *Device) GetNetworkInterfaces() ([]Interface, error) {
	interfaces := []Interface{}

	// Get all network interfaces
	netInterfaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}

	// Map of bridges to keep track of bridge members
	bridgeMembers := make(map[string][]string)

	// First pass to identify bridges
	for _, iface := range netInterfaces {
		// Check if this is a bridge
		if isBridge, members := d.isBridge(iface.Name); isBridge {
			bridgeMembers[iface.Name] = members
		}
	}

	// Second pass to collect interface details
	for _, iface := range netInterfaces {
		// Skip loopback interfaces
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		// Find which bridge this interface belongs to
		defaultIf := iface.Name
		for bridge, members := range bridgeMembers {
			for _, member := range members {
				if member == iface.Name {
					defaultIf = bridge
					break
				}
			}
		}

		// Get interface stats
		stats, err := d.getInterfaceStats(iface.Name)
		if err != nil {
			// Just use zeros if we can't get stats
			stats = map[string]uint64{
				"rx_bytes":   0,
				"tx_bytes":   0,
				"rx_packets": 0,
				"tx_packets": 0,
				"rx_errors":  0,
				"tx_errors":  0,
				"rx_dropped": 0,
				"tx_dropped": 0,
			}
		}

		// Get link speed
		linkSpeed := d.getLinkSpeed(iface.Name)

		// Determine interface type
		ifaceType := "unknown"
		if strings.HasPrefix(iface.Name, "eth") {
			ifaceType = "ethernet"
		} else if strings.HasPrefix(iface.Name, "wlan") {
			ifaceType = "wireless"
		} else if strings.HasPrefix(iface.Name, "br") {
			ifaceType = "bridge"
		}

		// Create interface object
		netInterface := Interface{
			If:            iface.Name,
			DefaultIf:     defaultIf,
			Mac:           iface.HardwareAddr.String(),
			RecBytes:      stats["rx_bytes"],
			Up:            iface.Flags&net.FlagUp != 0,
			Type:          ifaceType,
			BridgeMembers: bridgeMembers[iface.Name],
			RecPackets:    stats["rx_packets"],
			RecErrors:     stats["rx_errors"],
			RecDrops:      stats["rx_dropped"],
			SentBytes:     stats["tx_bytes"],
			SentPackets:   stats["tx_packets"],
			LinkSpeed:     linkSpeed,
			SentErrors:    stats["tx_errors"],
			SentDrops:     stats["tx_dropped"],
			Present:       true,
			External:      false, // Determine this based on additional logic if needed
			Ipv6:          d.hasIPv6(iface.Name),
			Mtu:           uint64(iface.MTU),
		}

		// Get carrier changes
		carrierChanges, _ := d.getCarrierChanges(iface.Name)
		netInterface.CarrierChanges = carrierChanges

		interfaces = append(interfaces, netInterface)
	}

	return interfaces, nil
}

// isBridge checks if the given interface is a bridge and returns its members
func (d *Device) isBridge(ifName string) (bool, []string) {
	// Check if this is a bridge by looking at sysfs
	bridgePath := filepath.Join("/sys/class/net", ifName, "bridge")
	if _, err := os.Stat(bridgePath); err != nil {
		return false, nil
	}

	// Get bridge members
	membersPath := filepath.Join(bridgePath, "slaves")
	members, err := os.ReadDir(membersPath)
	if err != nil {
		return true, []string{}
	}

	memberNames := []string{}
	for _, member := range members {
		memberNames = append(memberNames, member.Name())
	}

	return true, memberNames
}

// getInterfaceStats retrieves statistics for a network interface
func (d *Device) getInterfaceStats(ifName string) (map[string]uint64, error) {
	stats := make(map[string]uint64)

	// Read stats from /sys/class/net/<iface>/statistics/
	statsPaths := []string{
		"rx_bytes", "tx_bytes", "rx_packets", "tx_packets",
		"rx_errors", "tx_errors", "rx_dropped", "tx_dropped",
	}

	for _, statName := range statsPaths {
		statPath := filepath.Join("/sys/class/net", ifName, "statistics", statName)
		statData, err := os.ReadFile(statPath)
		if err != nil {
			continue
		}

		statValue, err := strconv.ParseUint(strings.TrimSpace(string(statData)), 10, 64)
		if err != nil {
			continue
		}

		stats[statName] = statValue
	}

	return stats, nil
}

// getLinkSpeed gets the speed of a network interface in Mbps
func (d *Device) getLinkSpeed(ifName string) uint64 {
	speedPath := filepath.Join("/sys/class/net", ifName, "speed")
	speedData, err := os.ReadFile(speedPath)
	if err != nil {
		return 0
	}

	speed, err := strconv.ParseUint(strings.TrimSpace(string(speedData)), 10, 64)
	if err != nil {
		return 0
	}

	return speed
}

// hasIPv6 checks if an interface has IPv6 configured
func (d *Device) hasIPv6(ifName string) bool {
	// Check for IPv6 addresses on the interface
	cmd := exec.Command("ip", "-6", "addr", "show", "dev", ifName)
	output, err := cmd.Output()
	if err != nil {
		return false
	}

	// If the output contains an inet6 line, IPv6 is configured
	return strings.Contains(string(output), "inet6")
}

// getCarrierChanges gets the number of carrier changes for an interface
func (d *Device) getCarrierChanges(ifName string) (uint64, error) {
	carrierChangesPath := filepath.Join("/sys/class/net", ifName, "carrier_changes")
	carrierData, err := os.ReadFile(carrierChangesPath)
	if err != nil {
		return 0, err
	}

	carrierChanges, err := strconv.ParseUint(strings.TrimSpace(string(carrierData)), 10, 64)
	if err != nil {
		return 0, err
	}

	return carrierChanges, nil
}

// GetWapInfo collects information about wireless access points
func (d *Device) GetWapInfo() ([]Wap, error) {
	waps := []Wap{}

	// Get wireless devices
	wirelessDevices, err := d.getWirelessDevices()
	if err != nil {
		return nil, err
	}

	for _, device := range wirelessDevices {
		wapInfo, err := d.getWapForDevice(device)
		if err != nil {
			continue
		}

		waps = append(waps, wapInfo)
	}

	return waps, nil
}

// getWirelessDevices gets a list of wireless devices
func (d *Device) getWirelessDevices() ([]string, error) {
	// Try to use 'iw dev' to list wireless devices
	cmd := exec.Command("iw", "dev")
	output, err := cmd.Output()
	if err != nil {
		// Fallback: Check for wireless interfaces in /sys/class/net
		var devices []string
		entries, err := os.ReadDir("/sys/class/net")
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			name := entry.Name()
			if strings.HasPrefix(name, "wlan") || strings.HasPrefix(name, "ath") {
				devices = append(devices, name)
			}
		}

		return devices, nil
	}

	// Parse iw dev output to get interface names
	var devices []string
	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "Interface") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				devices = append(devices, parts[1])
			}
		}
	}

	return devices, nil
}

// getWapForDevice collects wireless AP info for a specific device
func (d *Device) getWapForDevice(device string) (Wap, error) {
	wap := Wap{
		Interface: device,
		Stations:  []Station{},
	}

	// Get SSID
	ssid, err := d.getWirelessSSID(device)
	if err == nil {
		wap.Ssid = ssid
	}

	// Get wireless parameters
	wifiParams, err := d.getWirelessParams(device)
	if err == nil {
		wap.Signal = wifiParams.Signal
		wap.Channel = wifiParams.Channel
		wap.Frequency = wifiParams.Frequency
		wap.TxPower = wifiParams.TxPower
		wap.Quality = wifiParams.Quality
		wap.Rate = wifiParams.Rate
		wap.Chutil = wifiParams.Chutil
		wap.BandWidth = wifiParams.BandWidth
		wap.Mac = wifiParams.Mac
		wap.Mode = wifiParams.Mode
		wap.Phy = wifiParams.Phy
		wap.Noise = wifiParams.Noise
	}

	// Get connected stations
	stations, err := d.getConnectedStations(device)
	if err == nil {
		wap.Stations = stations
	}

	// Get encryption key and type
	key, keyTypes := d.getWirelessSecurity(device)
	wap.Key = key
	wap.Keytypes = keyTypes

	return wap, nil
}

// getWirelessSSID gets the SSID for a wireless interface
func (d *Device) getWirelessSSID(device string) (string, error) {
	cmd := exec.Command("iw", "dev", device, "info")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "ssid") {
			parts := strings.SplitN(line, "ssid", 2)
			if len(parts) == 2 {
				return strings.TrimSpace(parts[1]), nil
			}
		}
	}

	return "", fmt.Errorf("SSID not found for device %s", device)
}

// getWirelessParams gets wireless parameters for an interface
type wirelessParams struct {
	Signal    float64
	Channel   int
	Frequency int
	TxPower   int
	Quality   int
	Rate      int
	Chutil    int
	BandWidth int
	Mac       string
	Mode      string
	Phy       string
	Noise     float64
}

func (d *Device) getWirelessParams(device string) (wirelessParams, error) {
	params := wirelessParams{}

	// Get signal and noise from /proc/net/wireless
	wirelessFile, err := os.ReadFile("/proc/net/wireless")
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(wirelessFile))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, device) {
				fields := strings.Fields(line)
				if len(fields) >= 4 {
					// Parse signal (field 3)
					signal, err := strconv.ParseFloat(fields[3], 64)
					if err == nil {
						params.Signal = signal
					}

					// Parse noise (field 4)
					if len(fields) >= 5 {
						noise, err := strconv.ParseFloat(fields[4], 64)
						if err == nil {
							params.Noise = noise
						}
					}
				}
			}
		}
	}

	// Get channel, frequency, txpower using iw
	cmd := exec.Command("iw", "dev", device, "info")
	output, err := cmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(output))
		for scanner.Scan() {
			line := scanner.Text()

			// Extract channel and frequency
			if strings.Contains(line, "channel") {
				re := regexp.MustCompile(`channel\s+(\d+)\s+\((\d+)\s+MHz\)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 3 {
					params.Channel, _ = strconv.Atoi(matches[1])
					params.Frequency, _ = strconv.Atoi(matches[2])
				}
			}

			// Extract txpower
			if strings.Contains(line, "txpower") {
				re := regexp.MustCompile(`txpower\s+(\d+)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 2 {
					params.TxPower, _ = strconv.Atoi(matches[1])
				}
			}

			// Extract mode
			if strings.Contains(line, "type") {
				re := regexp.MustCompile(`type\s+(\w+)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 2 {
					params.Mode = matches[1]
				}
			}

			// Extract bandwidth
			if strings.Contains(line, "width") {
				re := regexp.MustCompile(`width:\s+(\d+)`)
				matches := re.FindStringSubmatch(line)
				if len(matches) >= 2 {
					params.BandWidth, _ = strconv.Atoi(matches[1])
				}
			}
		}
	}

	// Get MAC address
	macPath := filepath.Join("/sys/class/net", device, "address")
	macData, err := os.ReadFile(macPath)
	if err == nil {
		params.Mac = strings.TrimSpace(string(macData))
	}

	// Get physical device name
	phyPath := filepath.Join("/sys/class/net", device, "phy80211/name")
	phyData, err := os.ReadFile(phyPath)
	if err == nil {
		params.Phy = strings.TrimSpace(string(phyData))
	}

	return params, nil
}

// getConnectedStations gets information about stations connected to a wireless AP
func (d *Device) getConnectedStations(device string) ([]Station, error) {
	var stations []Station

	// Try using iw station dump
	cmd := exec.Command("iw", "dev", device, "station", "dump")
	output, err := cmd.Output()
	if err != nil {
		return stations, nil // Return empty list if command fails
	}

	var currentStation *Station

	scanner := bufio.NewScanner(bytes.NewReader(output))
	for scanner.Scan() {
		line := scanner.Text()

		// New station starts with MAC address
		if macRegex := regexp.MustCompile(`^Station\s+([0-9a-fA-F:]{17})`); macRegex.MatchString(line) {
			// Save previous station if any
			if currentStation != nil {
				stations = append(stations, *currentStation)
			}

			// Start new station
			matches := macRegex.FindStringSubmatch(line)
			currentStation = &Station{
				Mac:  matches[1],
				Info: "Station " + matches[1],
			}
		} else if currentStation != nil {
			// Parse station parameters
			if signalRegex := regexp.MustCompile(`signal:\s+(-?\d+\.\d+)\s+dBm`); signalRegex.MatchString(line) {
				matches := signalRegex.FindStringSubmatch(line)
				signal, _ := strconv.ParseFloat(matches[1], 64)
				currentStation.Rssi = signal
				currentStation.Signal0 = signal
				currentStation.Signal1 = signal
				currentStation.Signal2 = signal
				currentStation.Signal3 = signal
			} else if bytesRegex := regexp.MustCompile(`rx bytes:\s+(\d+)`); bytesRegex.MatchString(line) {
				matches := bytesRegex.FindStringSubmatch(line)
				currentStation.RecBytes, _ = strconv.ParseUint(matches[1], 10, 64)
			} else if bytesTxRegex := regexp.MustCompile(`tx bytes:\s+(\d+)`); bytesTxRegex.MatchString(line) {
				matches := bytesTxRegex.FindStringSubmatch(line)
				currentStation.SentBytes, _ = strconv.ParseUint(matches[1], 10, 64)
			} else if timeRegex := regexp.MustCompile(`connected time:\s+(\d+)\s+seconds`); timeRegex.MatchString(line) {
				matches := timeRegex.FindStringSubmatch(line)
				currentStation.AssocTime, _ = strconv.ParseUint(matches[1], 10, 64)
			}
		}
	}

	// Add the last station
	if currentStation != nil {
		stations = append(stations, *currentStation)
	}

	// Add DHCP info to each station
	dhcpLeases, _ := d.GetDhcpLeases()
	for i, station := range stations {
		if lease, ok := dhcpLeases[station.Mac]; ok {
			stations[i].Dhcp = lease
		} else {
			stations[i].Dhcp = Lease{
				Mac:       station.Mac,
				Ip:        "N/A",
				Hostname:  "N/A",
				Timestamp: uint64(time.Now().In(time.UTC).Unix()),
				Duid:      "N/A",
			}
		}
	}

	return stations, nil
}

// getWirelessSecurity gets security settings for a wireless interface
func (d *Device) getWirelessSecurity(device string) (string, string) {
	// Try OpenWrt/LEDE uci configuration first
	uciCmd := exec.Command("uci", "show", "wireless")
	uciOutput, err := uciCmd.Output()
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(uciOutput))
		var key, encryption string
		deviceSection := ""

		// Find the section for this device
		for scanner.Scan() {
			line := scanner.Text()
			// Look for lines that define the interface name
			if strings.Contains(line, ".ifname='"+device+"'") ||
				strings.Contains(line, ".device='"+device+"'") {
				parts := strings.Split(line, ".")
				if len(parts) > 0 {
					deviceSection = parts[0]
				}
			}

			// If we found our section, look for encryption and key
			if deviceSection != "" && strings.HasPrefix(line, deviceSection) {
				if strings.Contains(line, ".encryption=") {
					parts := strings.Split(line, "=")
					if len(parts) > 1 {
						encryption = strings.Trim(parts[1], "'")
					}
				} else if strings.Contains(line, ".key=") {
					parts := strings.Split(line, "=")
					if len(parts) > 1 {
						key = strings.Trim(parts[1], "'")
					}
				}
			}
		}

		if encryption != "" {
			return key, encryption
		}
	}

	// Try to read from wpa_supplicant configuration
	wpaPath := "/etc/wpa_supplicant/wpa_supplicant-" + device + ".conf"
	if _, err := os.Stat(wpaPath); os.IsNotExist(err) {
		// Try default config if device-specific doesn't exist
		wpaPath = "/etc/wpa_supplicant/wpa_supplicant.conf"
	}

	wpaContent, err := os.ReadFile(wpaPath)
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(wpaContent))
		var key, encryption string
		inNetwork := false

		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "network={") {
				inNetwork = true
			} else if strings.Contains(line, "}") && inNetwork {
				inNetwork = false
			} else if inNetwork {
				if strings.Contains(line, "psk=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) > 1 {
						key = strings.Trim(parts[1], "\"")
					}
					encryption = "psk2+ccmp"
				} else if strings.Contains(line, "key_mgmt=NONE") {
					encryption = "none"
				} else if strings.Contains(line, "key_mgmt=WPA-PSK") {
					encryption = "psk+ccmp"
				} else if strings.Contains(line, "key_mgmt=WPA-EAP") {
					encryption = "wpa-eap"
				} else if strings.Contains(line, "key_mgmt=SAE") {
					encryption = "sae"
				}
			}
		}

		if encryption != "" {
			return key, encryption
		}
	}

	// Try using iw command to get encryption info (won't get the key)
	iwCmd := exec.Command("iw", device, "info")
	iwOutput, err := iwCmd.Output()
	if err == nil {
		encryption := "none"

		if bytes.Contains(iwOutput, []byte("WPA")) ||
			bytes.Contains(iwOutput, []byte("wpa")) {
			encryption = "psk+ccmp"
		}

		if bytes.Contains(iwOutput, []byte("WPA2")) ||
			bytes.Contains(iwOutput, []byte("wpa2")) {
			encryption = "psk2+ccmp"
		}

		if bytes.Contains(iwOutput, []byte("WPA3")) ||
			bytes.Contains(iwOutput, []byte("wpa3")) {
			encryption = "sae"
		}

		return "********", encryption
	}

	// Default fallback
	return "********", "psk2+ccmp"
}

// GetTcpInfo collects TCP/IP statistics
func (d *Device) GetTcpInfo() (TcpCollector, error) {
	tcp := TcpCollector{
		UniqueIps:         0,
		SlowedPairPackets: 0,
		Cwr:               0,
		Ece:               0,
		Rst:               0,
		Syn:               0,
		Urg:               0,
	}

	// Count unique IPs from connection tracking
	uniqueIps, err := d.countUniqueIps()
	if err == nil {
		tcp.UniqueIps = uniqueIps
	}

	// Get TCP stats
	err = d.parseTcpStats(&tcp)

	return tcp, err
}

// countUniqueIps counts unique IPs in the connection tracking table
func (d *Device) countUniqueIps() (uint64, error) {
	// Try to read connection tracking table
	conntrackPath := "/proc/net/nf_conntrack"
	if _, err := os.Stat(conntrackPath); os.IsNotExist(err) {
		return 0, err
	}

	cmd := exec.Command("sh", "-c", "cat "+conntrackPath+" | grep -o 'src=[0-9.]\\+' | sort -u | wc -l")
	output, err := cmd.Output()
	if err != nil {
		return 0, err
	}

	count, err := strconv.ParseUint(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// parseTcpStats parses TCP statistics from /proc/net/netstat
func (d *Device) parseTcpStats(tcp *TcpCollector) error {
	netstatPath := "/proc/net/netstat"
	content, err := os.ReadFile(netstatPath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	var headers, values []string

	for i, line := range lines {
		if strings.HasPrefix(line, "TcpExt:") {
			headers = strings.Fields(line)
			if i+1 < len(lines) {
				values = strings.Fields(lines[i+1])
				break
			}
		}
	}

	if len(headers) != len(values) {
		return fmt.Errorf("mismatched headers and values in netstat")
	}

	// Map of header to index
	headerMap := make(map[string]int)
	for i, header := range headers {
		headerMap[header] = i
	}

	// Extract relevant stats
	extractStat := func(name string) uint64 {
		if idx, ok := headerMap[name]; ok && idx < len(values) {
			val, err := strconv.ParseUint(values[idx], 10, 64)
			if err == nil {
				return val
			}
		}
		return 0
	}

	tcp.Syn = extractStat("TCPSynRetrans")
	tcp.Urg = extractStat("TCPOrigDataSent")
	tcp.Rst = extractStat("TCPAbortOnData")
	tcp.Cwr = extractStat("TCPWinProbe")

	// Calculate slowed pair packets
	if tcp.Cwr > 0 && extractStat("TCPDSACKOldSent") > 0 {
		tcp.SlowedPairPackets++
	}

	return nil
}

// GetDhcpLeases retrieves DHCP lease information
func (d *Device) GetDhcpLeases() (map[string]Lease, error) {
	leases := make(map[string]Lease)

	// Try to read DHCP leases file
	leasesPath := "/tmp/dhcp.leases"
	if _, err := os.Stat(leasesPath); os.IsNotExist(err) {
		return leases, nil
	}

	content, err := os.ReadFile(leasesPath)
	if err != nil {
		return leases, err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 4 {
			continue
		}

		timestamp, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			continue
		}

		mac := fields[1]
		ip := fields[2]
		hostname := fields[3]
		duid := "N/A"
		if len(fields) >= 5 {
			duid = fields[4]
		}

		// Replace "*" with "Unknown" for hostname
		if hostname == "*" {
			hostname = "Unknown"
		}

		// Replace "*" with MAC for DUID if missing
		if duid == "*" {
			duid = mac
		}

		leases[mac] = Lease{
			Timestamp: timestamp,
			Mac:       mac,
			Ip:        ip,
			Hostname:  hostname,
			Duid:      duid,
		}
	}

	return leases, nil
}

// GetArpTable retrieves the ARP table
func (d *Device) GetArpTable() (map[string]ArpDevice, error) {
	arpTable := make(map[string]ArpDevice)

	// Try to read ARP table
	arpPath := "/proc/net/arp"
	content, err := os.ReadFile(arpPath)
	if err != nil {
		return arpTable, err
	}

	lines := strings.Split(string(content), "\n")
	if len(lines) < 2 {
		return arpTable, nil
	}

	// Skip header line
	for _, line := range lines[1:] {
		fields := strings.Fields(line)
		if len(fields) < 6 {
			continue
		}

		ip := fields[0]
		hwtype := fields[1]
		flags := fields[2]
		mac := fields[3]
		mask := fields[4]
		device := fields[5]

		// Initialize device entry if not exists
		if _, ok := arpTable[device]; !ok {
			arpTable[device] = ArpDevice{
				Name:   device,
				Device: device,
				Clients: []struct {
					Ip     string `json:"ip,omitempty"`
					Mac    string `json:"mac,omitempty"`
					Flags  string `json:"flags,omitempty"`
					Hwtype string `json:"hwtype,omitempty"`
					Mask   string `json:"mask,omitempty"`
				}{},
			}
		}

		// Add client
		dev := arpTable[device]
		dev.Clients = append(dev.Clients, struct {
			Ip     string `json:"ip,omitempty"`
			Mac    string `json:"mac,omitempty"`
			Flags  string `json:"flags,omitempty"`
			Hwtype string `json:"hwtype,omitempty"`
			Mask   string `json:"mask,omitempty"`
		}{
			Ip:     ip,
			Mac:    mac,
			Flags:  flags,
			Hwtype: hwtype,
			Mask:   mask,
		})
		arpTable[device] = dev
	}

	return arpTable, nil
}

// GetPingStats performs a ping test and returns statistics
func (d *Device) GetPingStats(host string) (Ping, error) {
	ping := Ping{
		Host:   host,
		AvgRtt: 0,
		MinRtt: 0,
		MaxRtt: 0,
		Loss:   0,
	}

	// Run ping command
	cmd := exec.Command("ping", "-c", "5", "-q", host)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return ping, err
	}

	// Parse ping statistics
	outputStr := string(output)
	lines := strings.Split(outputStr, "\n")

	for _, line := range lines {
		if strings.Contains(line, "min/avg/max") {
			// Extract min/avg/max values
			re := regexp.MustCompile(`= ([\d.]+)/([\d.]+)/([\d.]+)`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 4 {
				ping.MinRtt, _ = strconv.ParseFloat(matches[1], 64)
				ping.AvgRtt, _ = strconv.ParseFloat(matches[2], 64)
				ping.MaxRtt, _ = strconv.ParseFloat(matches[3], 64)
			}
		} else if strings.Contains(line, "packet loss") {
			// Extract packet loss percentage
			re := regexp.MustCompile(`([\d.]+)% packet loss`)
			matches := re.FindStringSubmatch(line)
			if len(matches) >= 2 {
				ping.Loss, _ = strconv.ParseFloat(matches[1], 64)
			}
		}
	}

	return ping, nil
}

// GetExternalIpInfo retrieves information about the external IP address
type ExternalIpInfo struct {
	IP         string  `json:"ip"`
	Country    string  `json:"country"`
	City       string  `json:"city"`
	ASN        string  `json:"asn"`
	RegionName string  `json:"region_name"`
	RegionCode string  `json:"region_code"`
	TimeZone   string  `json:"time_zone"`
	Latitude   float64 `json:"latitude"`
	Longitude  float64 `json:"longitude"`
	CountryIso string  `json:"country_iso"`
}

func (d *Device) GetExternalIpInfo() (ExternalIpInfo, error) {
	info := ExternalIpInfo{
		IP:         "N/A",
		Country:    "N/A",
		City:       "N/A",
		ASN:        "N/A",
		RegionName: "N/A",
		RegionCode: "N/A",
		TimeZone:   "N/A",
		Latitude:   0,
		Longitude:  0,
		CountryIso: "N/A",
	}

	// Make HTTP request to get IP info
	resp, err := http.Get("http://ip.longshot-router.com/json")
	if err != nil {
		return info, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return info, fmt.Errorf("HTTP request failed with status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return info, err
	}

	// Parse JSON response
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		return info, err
	}

	// Extract information
	if ip, ok := data["realIp"].(string); ok {
		info.IP = ip
	} else if ip, ok := data["ip"].(string); ok {
		info.IP = ip
	}

	if country, ok := data["country"].(string); ok {
		info.Country = country
	}

	if city, ok := data["city"].(string); ok {
		info.City = city
	}

	// Extract ASN information
	if asn, ok := data["asn"].(map[string]interface{}); ok {
		asnNumber, _ := asn["asn"].(float64)
		info.ASN = fmt.Sprintf("%.0f", asnNumber)

		// Extract entity information if available
		if entities, ok := asn["entities"].([]interface{}); ok && len(entities) > 0 {
			entityNames := []string{}
			for _, entity := range entities {
				if e, ok := entity.(map[string]interface{}); ok {
					if fn, ok := e["fn"].(string); ok {
						entityType, _ := e["type"].(string)
						entityNames = append(entityNames, fn+" ("+entityType+")")
					}
				}
			}
			if len(entityNames) > 0 {
				info.ASN += " (" + strings.Join(entityNames, ", ") + ")"
			}
		}
	}

	if region, ok := data["region"].(string); ok {
		info.RegionName = region
	}

	if regionCode, ok := data["regionCode"].(string); ok {
		info.RegionCode = regionCode
	}

	if timezone, ok := data["timezone"].(string); ok {
		info.TimeZone = timezone
	}

	if lat, ok := data["latitude"].(string); ok {
		info.Latitude, _ = strconv.ParseFloat(lat, 64)
	}

	if lng, ok := data["longitude"].(string); ok {
		info.Longitude, _ = strconv.ParseFloat(lng, 64)
	}

	if countryIso, ok := data["country_iso"].(string); ok {
		info.CountryIso = countryIso
	}

	return info, nil
}

// GetWanIP retrieves the WAN IP address
func (d *Device) GetWanIP() (string, error) {
	// Try to get global IPv4 addresses
	cmd := exec.Command("sh", "-c", "ip -4 addr show | awk '/global/ && /inet/ {print $2}' | cut -d'/' -f1 | head -1")
	output, err := cmd.Output()
	if err != nil {
		return "N/A", err
	}

	wanIP := strings.TrimSpace(string(output))
	if wanIP == "" {
		return "N/A", fmt.Errorf("no WAN IP found")
	}

	return wanIP, nil
}

// GetWirelessConfigured retrieves wireless interface configuration
func (d *Device) GetWirelessConfigured() ([]WirelessInterface, error) {
	// In a production system, this would parse UCI configuration
	// or other system configuration files

	var wirelessInterfaces []WirelessInterface

	// Get wireless devices
	devices, err := d.getWirelessDevices()
	if err != nil {
		return wirelessInterfaces, err
	}

	for i, device := range devices {
		disabled := false
		hidden := false
		running := true

		// Check if interface is running
		if _, err := os.Stat(filepath.Join("/sys/class/net", device)); os.IsNotExist(err) {
			running = false
		}

		// Get wireless parameters
		params, _ := d.getWirelessParams(device)

		// Get security information
		key, encryption := d.getWirelessSecurity(device)

		// Get SSID
		ssid, _ := d.getWirelessSSID(device)

		// Create wireless interface object
		wifiInterface := WirelessInterface{
			ID:              strPtr("*" + strconv.Itoa(i)),
			Disabled:        &disabled,
			HideSSID:        &hidden,
			InterfaceType:   strPtr("wireless"),
			Key:             &key,
			MACAddress:      &params.Mac,
			Encryption:      &encryption,
			MasterInterface: &device,
			Name:            &device,
			Running:         &running,
			SecurityProfile: strPtr("*" + strconv.Itoa(i)),
			SSID:            &ssid,
			Band:            strPtr(d.getBandFromFrequency(params.Frequency)),
		}

		wirelessInterfaces = append(wirelessInterfaces, wifiInterface)
	}

	return wirelessInterfaces, nil
}

// Helper function to convert a string to a pointer
func strPtr(s string) *string {
	return &s
}

// getBandFromFrequency determines the wireless band from frequency
func (d *Device) getBandFromFrequency(frequency int) string {
	if frequency >= 2400 && frequency <= 2500 {
		return "2.4ghz-g"
	} else if frequency >= 5000 && frequency <= 6000 {
		return "5ghz-ac"
	} else if frequency >= 6000 {
		return "6ghz-ax"
	}
	return "unknown"
}

// GetSecurityProfiles retrieves security profiles for wireless interfaces
func (d *Device) GetSecurityProfiles() ([]SecurityProfile, error) {
	// In a production system, this would parse UCI configuration
	// or other system configuration files

	var profiles []SecurityProfile

	// Get wireless devices
	devices, err := d.getWirelessDevices()
	if err != nil {
		return profiles, err
	}

	for i, device := range devices {
		// Get security information
		_, encryption := d.getWirelessSecurity(device)

		// Create profile ID
		id := "*" + strconv.Itoa(i)

		// Create default variables
		defaultVal := false
		mode := "ap"
		name := device
		technology := "wireless"

		// Create authentication types based on encryption
		var authTypes []string
		if encryption == "none" {
			authTypes = []string{"none"}
		} else if strings.Contains(encryption, "psk") {
			authTypes = []string{encryption}
		} else {
			authTypes = []string{"psk2+ccmp"}
		}

		// Get key
		key, _ := d.getWirelessSecurity(device)

		// Create security profile
		profile := SecurityProfile{
			ID:                  &id,
			AuthenticationTypes: authTypes,
			Default:             &defaultVal,
			Mode:                &mode,
			Name:                &name,
			Technology:          &technology,
		}

		// Add appropriate key based on encryption type
		if strings.Contains(encryption, "psk2") || strings.Contains(encryption, "wpa2") {
			profile.WPA2PreSharedKey = &key
		}

		if strings.Contains(encryption, "psk") && !strings.Contains(encryption, "psk2") {
			profile.WPAPreSharedKey = &key
		}

		if strings.Contains(encryption, "sae") || strings.Contains(encryption, "wpa3") {
			profile.WPA3PreSharedKey = &key
		}

		profiles = append(profiles, profile)
	}

	return profiles, nil
}

// GetOSInfo retrieves information about the operating system
type OSInfo struct {
	Name      string
	Version   string
	ID        string
	BuildDate string
}

func (d *Device) GetOSInfo() (OSInfo, error) {
	info := OSInfo{
		Name:      "Unknown",
		Version:   "Unknown",
		ID:        "Unknown",
		BuildDate: "",
	}

	// Try to read /etc/os-release
	osRelease, err := os.ReadFile("/etc/os-release")
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(osRelease))
		for scanner.Scan() {
			line := scanner.Text()
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}

			key := parts[0]
			value := strings.Trim(parts[1], "\"'")

			switch key {
			case "NAME":
				info.Name = value
			case "VERSION":
				info.Version = value
			case "ID":
				info.ID = value
			}
		}
	}

	// Try to get build date from /etc directory timestamp
	if stat, err := os.Stat("/etc"); err == nil {
		info.BuildDate = stat.ModTime().Format(time.RFC3339)
	}

	return info, nil
}

// GetHardwareInfo retrieves information about the hardware
type HardwareInfo struct {
	hardwareMake         string
	hardwareModel        string
	hardwareModelNumber  string
	hardwareSerialNumber string
	hardwareCpuInfo      string
}

func (d *Device) GetHardwareInfo() (HardwareInfo, error) {
	info := HardwareInfo{
		hardwareMake:         "Generic",
		hardwareModel:        "Generic",
		hardwareModelNumber:  "N/A",
		hardwareSerialNumber: "N/A",
		hardwareCpuInfo:      "N/A",
	}

	// Try to read /etc/device_info
	deviceInfo, err := os.ReadFile("/etc/device_info")
	if err == nil {
		scanner := bufio.NewScanner(bytes.NewReader(deviceInfo))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "DEVICE_MANUFACTURER") {
				info.hardwareMake = strings.Trim(strings.Split(line, "=")[1], "'")
			} else if strings.Contains(line, "DEVICE_PRODUCT") {
				info.hardwareModel = strings.Trim(strings.Split(line, "=")[1], "'")
			} else if strings.Contains(line, "DEVICE_REVISION") {
				info.hardwareModelNumber = strings.Trim(strings.Split(line, "=")[1], "'")
			}
		}
	}

	// Try to get serial number
	cmd := exec.Command("fw_printenv", "sn")
	output, err := cmd.Output()
	if err == nil && len(output) > 0 {
		re := regexp.MustCompile(`sn=(\S+)`)
		matches := re.FindStringSubmatch(string(output))
		if len(matches) >= 2 {
			info.hardwareSerialNumber = matches[1]
		}
	}

	// If serial number not found, try reading from /proc/cpuinfo
	if info.hardwareSerialNumber == "N/A" {
		cpuInfo, err := os.ReadFile("/proc/cpuinfo")
		if err == nil {
			scanner := bufio.NewScanner(bytes.NewReader(cpuInfo))
			for scanner.Scan() {
				line := scanner.Text()
				if strings.Contains(line, "Serial") {
					parts := strings.SplitN(line, ":", 2)
					if len(parts) >= 2 {
						info.hardwareSerialNumber = strings.TrimSpace(parts[1])
					}
				}
			}
		}
	}

	// Get CPU info
	cpuInfo, err := os.ReadFile("/proc/cpuinfo")
	if err == nil {
		var model string
		var cores int

		scanner := bufio.NewScanner(bytes.NewReader(cpuInfo))
		for scanner.Scan() {
			line := scanner.Text()
			if strings.Contains(line, "model name") {
				parts := strings.SplitN(line, ":", 2)
				if len(parts) >= 2 {
					model = strings.TrimSpace(parts[1])
				}
			} else if strings.Contains(line, "processor") {
				cores++
			}
		}

		info.hardwareCpuInfo = fmt.Sprintf("%s, %d cores", model, cores)
	}

	return info, nil
}

// HasWebshellSupport checks if the device supports webshell
func (d *Device) HasWebshellSupport() bool {
	// Check if the shell is available
	_, err := exec.LookPath("sh")
	return err == nil
}

// HasBandwidthTestSupport checks if the device supports bandwidth testing
func (d *Device) HasBandwidthTestSupport() bool {
	// Check for iperf or other bandwidth testing tools
	_, err1 := exec.LookPath("iperf")
	_, err2 := exec.LookPath("iperf3")
	_, err3 := exec.LookPath("btest")

	return err1 == nil || err2 == nil || err3 == nil
}

// HasFirmwareUpgradeSupport checks if the device supports firmware upgrades
func (d *Device) HasFirmwareUpgradeSupport() bool {
	// Check for firmware upgrade tools
	_, err := exec.LookPath("fwupgrade-tools")
	return err == nil
}

// PerformSpeedTest runs a bandwidth test
func (d *Device) PerformSpeedTest() (map[string]interface{}, error) {
	result := map[string]interface{}{
		"up":   0.0,
		"down": 0.0,
	}

	// Determine which speed test tool to use
	var testTool string
	if _, err := exec.LookPath("iperf3"); err == nil {
		testTool = "iperf3"
	} else if _, err := exec.LookPath("iperf"); err == nil {
		testTool = "iperf"
	} else {
		return result, fmt.Errorf("no speed test tool available")
	}

	// Use a predefined server for testing
	server := "iperf.longshot-router.com"

	// Run the speed test
	var cmd *exec.Cmd
	if testTool == "iperf3" {
		cmd = exec.Command(testTool, "-c", server, "-i", "1", "-t", "1", "-P", "5", "-f", "m", "-J")
	} else {
		cmd = exec.Command(testTool, "-c", server, "-i", "1", "-t", "1", "-P", "5", "-f", "m")
	}

	output, err := cmd.Output()
	if err != nil {
		return result, err
	}

	// Parse the output based on which tool was used
	if testTool == "iperf3" {
		var jsonResult map[string]interface{}
		if err := json.Unmarshal(output, &jsonResult); err != nil {
			return result, err
		}

		if end, ok := jsonResult["end"].(map[string]interface{}); ok {
			if sumSent, ok := end["sum_sent"].(map[string]interface{}); ok {
				if bits, ok := sumSent["bits_per_second"].(float64); ok {
					result["up"] = bits / 1024 / 1024 // Convert to Mbps
				}
			}

			if sumReceived, ok := end["sum_received"].(map[string]interface{}); ok {
				if bits, ok := sumReceived["bits_per_second"].(float64); ok {
					result["down"] = bits / 1024 / 1024 // Convert to Mbps
				}
			}
		}
	} else {
		// Parse iperf output
		outputStr := string(output)
		upRe := regexp.MustCompile(`SUM.*sender.*\s+(\d+\.\d+)\s+Mbits/sec`)
		downRe := regexp.MustCompile(`SUM.*receiver.*\s+(\d+\.\d+)\s+Mbits/sec`)

		upMatches := upRe.FindStringSubmatch(outputStr)
		if len(upMatches) >= 2 {
			result["up"], _ = strconv.ParseFloat(upMatches[1], 64)
		}

		downMatches := downRe.FindStringSubmatch(outputStr)
		if len(downMatches) >= 2 {
			result["down"], _ = strconv.ParseFloat(downMatches[1], 64)
		}
	}

	// Convert values to KB/s for consistency
	result["up"] = math.Floor(result["up"].(float64) * 1024)
	result["down"] = math.Floor(result["down"].(float64) * 1024)

	return result, nil
}
