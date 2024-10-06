#ifndef TYPES_H
#define TYPES_H

#include <jansson.h>
#include <stdint.h>
#include <stdbool.h>

// Cmd struct
typedef struct {
    char *cmd;
    char *type;
    char *ws_id;
    char *uuidv4;
    char *stdout;
    char *stderr;
} Cmd;

// Interface struct
typedef struct {
    char *interface_name;
    char *default_if;
    char *mac;
    uint64_t rec_bytes;
    uint64_t rec_packets;
    uint64_t rec_errors;
    uint64_t rec_drops;
    uint64_t sent_bytes;
    uint64_t sent_packets;
    uint64_t sent_errors;
    uint64_t sent_drops;
    uint64_t carrier_changes;
    uint64_t macs;
    char *found_descriptor;
} Interface;

// WirelessInterface struct
typedef struct {
    char *id;
    bool *disabled;
    bool *hide_ssid;
    char *interface_type;
    char *key;
    char *mac_address;
    char *master_interface;
    char *name;
    bool *running;
    char *security_profile;
    char *ssid;
    char *band;
} WirelessInterface;

// SecurityProfile struct
typedef struct {
    char *id;
    char **authentication_types;
    int auth_type_count;
    bool *default_profile;
    char **eap_methods;
    int eap_method_count;
    char **group_ciphers;
    int group_cipher_count;
    char *mode;
    char *name;
    char *radius_called_format;
    char *technology;
    char *wpa_pre_shared_key;
    char *wpa2_pre_shared_key;
    char *wpa3_pre_shared_key;
} SecurityProfile;

// HostRequest struct
typedef struct {
    char *type;
    char *host_id;
    char *login;
    char *key;
    char *simple_rotated_key;
    int64_t uptime;
    char *wan_ip;
    char *fw_status;
    char *client_info;
    char *os;
    char *os_version;
    char *fw;
    char *hardware_make;
    char *hardware_model;
    char *hardware_model_number;
    char *hardware_serial_number;
    char *hardware_cpu_info;
    uint64_t os_build_date;
    int64_t reboot;
    Cmd *cmds;
    int cmd_count;
    char *cmd;
    char *ws_id;
    char *uuidv4;
    char *stdout;
    char *stderr;
    double *lat;
    double *lng;
    Interface *interfaces;
    int interface_count;
    WirelessInterface *wireless_configured;
    int wireless_configured_count;
    SecurityProfile *security_profiles;
    int security_profile_count;
    bool webshell_support;
    bool bandwidth_test_support;
    bool firmware_upgrade_support;
    char *hostname;
    char *outside_ip;
    int64_t last_config_request;
    bool using_websocket;
} HostRequest;

// Function declarations for JSON conversion

json_t *cmd_to_json(Cmd *cmd);
json_t *interface_to_json(Interface *iface);
json_t *wireless_interface_to_json(WirelessInterface *wifi);
json_t *security_profile_to_json(SecurityProfile *profile);
json_t *host_request_to_json(HostRequest *request);

#endif // TYPES_H
