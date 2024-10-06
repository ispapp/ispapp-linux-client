#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
#include "types.h"  // Assuming your structs are defined in types.h

// Helper functions to run shell commands and capture their output
char *run_command(const char *cmd) {
    char buffer[128];
    char *result = (char *)malloc(1024); // Allocate 1 KB for the result
    if (!result) return NULL;
    
    FILE *pipe = popen(cmd, "r");
    if (!pipe) return NULL;
    
    result[0] = '\0'; // Initialize empty string
    while (fgets(buffer, sizeof(buffer), pipe) != NULL) {
        strcat(result, buffer);
    }
    
    pclose(pipe);
    return result;
}
// Helper function to get system uptime
int64_t get_system_uptime() {
    struct sysinfo info;
    if (sysinfo(&info) == 0) {
        return (int64_t)info.uptime;
    }
    return -1; // Fallback value if sysinfo fails
}
// Helper function to get wireless configuration using `iw` command
WirelessInterface collect_wireless_info() {
    WirelessInterface wifi;
    memset(&wifi, 0, sizeof(WirelessInterface));

    // Example command to get SSID (you might want to tweak the command based on your requirements)
    char *ssid = run_command("iw dev wlan0 info | grep ssid | awk '{print $2}'");
    wifi.ssid = ssid ? strdup(ssid) : NULL;

    // Get MAC address
    char *mac_address = run_command("ifconfig wlan0 | grep HWaddr | awk '{print $5}'");
    wifi.mac_address = mac_address ? strdup(mac_address) : NULL;

    // Get band (2.4GHz, 5GHz, etc.)
    char *band = run_command("iw dev wlan0 info | grep channel | awk '{print $4}'");
    wifi.band = band ? strdup(band) : NULL;

    // Fill other fields using similar commands, or leave them empty
    wifi.running = true;  // Just an example; you may need to adjust based on real data
    return wifi;
}

// Helper function to get interface information using `ifconfig` or `ip` command
Interface collect_interface_info() {
    Interface iface;
    memset(&iface, 0, sizeof(Interface));

    // Get interface name (eth0 as an example)
    iface.interface_name = strdup("eth0");

    // Get MAC address
    char *mac_address = run_command("ifconfig eth0 | grep HWaddr | awk '{print $5}'");
    iface.mac = mac_address ? strdup(mac_address) : NULL;

    // Get bytes received
    char *rec_bytes = run_command("ifconfig eth0 | grep 'RX bytes' | awk '{print $2}' | cut -d':' -f2");
    iface.rec_bytes = rec_bytes ? atoll(rec_bytes) : 0;

    // Get bytes sent
    char *sent_bytes = run_command("ifconfig eth0 | grep 'TX bytes' | awk '{print $6}' | cut -d':' -f2");
    iface.sent_bytes = sent_bytes ? atoll(sent_bytes) : 0;

    // Fill other fields similarly based on the command outputs
    return iface;
}

// Collect configuration from UCI system for OpenWrt
void collect_uci_info(HostRequest *request) {
    // Example to collect the hostname using `uci`
    char *hostname = run_command("uci get system.@system[0].hostname");
    request->hostname = hostname ? strdup(hostname) : NULL;

    // Collect WAN IP from UCI (if configured)
    char *wan_ip = run_command("uci get network.wan.ipaddr");
    request->wan_ip = wan_ip ? strdup(wan_ip) : NULL;

    // Collect firmware status
    char *fw_status = run_command("uci get system.@system[0].fw_status");
    request->fw_status = fw_status ? strdup(fw_status) : NULL;

    // Add additional UCI-based configurations as necessary
}
// Helper function to get the hardware serial number (example using `cat` and `/proc` file)
char* get_hardware_serial() {
    return run_command("cat /proc/cpuinfo | grep Serial | awk '{print $3}'");
}

// Main function to collect all configurations and fill the HostRequest structure
void CollectConfigs(HostRequest *request) {
    memset(request, 0, sizeof(HostRequest));

    // Collect basic system information using UCI or shell commands
    collect_uci_info(request);

    // Collect wireless information using `iw`
    request->wireless_configured = (WirelessInterface *)malloc(sizeof(WirelessInterface));
    if (request->wireless_configured) {
        request->wireless_configured_count = 1;
        request->wireless_configured[0] = collect_wireless_info();
    }

    // Collect interface information using `ifconfig` or `ip`
    request->interfaces = (Interface *)malloc(sizeof(Interface));
    if (request->interfaces) {
        request->interface_count = 1;
        request->interfaces[0] = collect_interface_info();
    }

    // Collect any other necessary information (e.g., Cmd array, security profiles, etc.)
    // You can add more collection logic here as required
}

