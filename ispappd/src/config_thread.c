#include "ispappd.h"
#include "types.h"
#include <jansson.h>
#include <stdlib.h>

void *configs_thread(void *arg)
{
    // Configs thread logic
    while (1)
    {
        // Perform periodic config checks
        sleep(60); // Placeholder: Adjust as needed
    }
    return NULL;
}

json_t *cmd_to_json(Cmd *cmd) {
    json_t *json = json_object();
    json_object_set_new(json, "cmd", cmd->cmd ? json_string(cmd->cmd) : json_null());
    json_object_set_new(json, "type", cmd->type ? json_string(cmd->type) : json_null());
    json_object_set_new(json, "ws_id", cmd->ws_id ? json_string(cmd->ws_id) : json_null());
    json_object_set_new(json, "uuidv4", cmd->uuidv4 ? json_string(cmd->uuidv4) : json_null());
    json_object_set_new(json, "stdout", cmd->stdout ? json_string(cmd->stdout) : json_null());
    json_object_set_new(json, "stderr", cmd->stderr ? json_string(cmd->stderr) : json_null());
    return json;
}

json_t *interface_to_json(Interface *iface) {
    json_t *json = json_object();
    json_object_set_new(json, "if", iface->interface_name ? json_string(iface->interface_name) : json_null());
    json_object_set_new(json, "defaultIf", iface->default_if ? json_string(iface->default_if) : json_null());
    json_object_set_new(json, "mac", iface->mac ? json_string(iface->mac) : json_null());
    json_object_set_new(json, "recBytes", json_integer(iface->rec_bytes));
    json_object_set_new(json, "recPackets", json_integer(iface->rec_packets));
    json_object_set_new(json, "recErrors", json_integer(iface->rec_errors));
    json_object_set_new(json, "recDrops", json_integer(iface->rec_drops));
    json_object_set_new(json, "sentBytes", json_integer(iface->sent_bytes));
    json_object_set_new(json, "sentPackets", json_integer(iface->sent_packets));
    json_object_set_new(json, "sentErrors", json_integer(iface->sent_errors));
    json_object_set_new(json, "sentDrops", json_integer(iface->sent_drops));
    json_object_set_new(json, "carrierChanges", json_integer(iface->carrier_changes));
    json_object_set_new(json, "macs", json_integer(iface->macs));
    json_object_set_new(json, "foundDescriptor", iface->found_descriptor ? json_string(iface->found_descriptor) : json_null());
    return json;
}

json_t *wireless_interface_to_json(WirelessInterface *wifi) {
    json_t *json = json_object();
    json_object_set_new(json, ".id", wifi->id ? json_string(wifi->id) : json_null());
    json_object_set_new(json, "disabled", wifi->disabled ? json_boolean(*wifi->disabled) : json_null());
    json_object_set_new(json, "hide-ssid", wifi->hide_ssid ? json_boolean(*wifi->hide_ssid) : json_null());
    json_object_set_new(json, "interface-type", wifi->interface_type ? json_string(wifi->interface_type) : json_null());
    json_object_set_new(json, "key", wifi->key ? json_string(wifi->key) : json_null());
    json_object_set_new(json, "mac-address", wifi->mac_address ? json_string(wifi->mac_address) : json_null());
    json_object_set_new(json, "master-interface", wifi->master_interface ? json_string(wifi->master_interface) : json_null());
    json_object_set_new(json, "name", wifi->name ? json_string(wifi->name) : json_null());
    json_object_set_new(json, "running", wifi->running ? json_boolean(*wifi->running) : json_null());
    json_object_set_new(json, "security-profile", wifi->security_profile ? json_string(wifi->security_profile) : json_null());
    json_object_set_new(json, "ssid", wifi->ssid ? json_string(wifi->ssid) : json_null());
    json_object_set_new(json, "band", wifi->band ? json_string(wifi->band) : json_null());
    return json;
}

json_t *security_profile_to_json(SecurityProfile *profile) {
    json_t *json = json_object();

    json_object_set_new(json, ".id", profile->id ? json_string(profile->id) : json_null());
    json_t *auth_types = json_array();
    for (int i = 0; i < profile->auth_type_count; i++) {
        json_array_append_new(auth_types, json_string(profile->authentication_types[i]));
    }
    json_object_set_new(json, "authentication-types", auth_types);

    json_object_set_new(json, "default", profile->default_profile ? json_boolean(*profile->default_profile) : json_null());

    json_t *eap_methods = json_array();
    for (int i = 0; i < profile->eap_method_count; i++) {
        json_array_append_new(eap_methods, json_string(profile->eap_methods[i]));
    }
    json_object_set_new(json, "eap-methods", eap_methods);

    json_t *group_ciphers = json_array();
    for (int i = 0; i < profile->group_cipher_count; i++) {
        json_array_append_new(group_ciphers, json_string(profile->group_ciphers[i]));
    }
    json_object_set_new(json, "group-ciphers", group_ciphers);

    json_object_set_new(json, "mode", profile->mode ? json_string(profile->mode) : json_null());
    json_object_set_new(json, "name", profile->name ? json_string(profile->name) : json_null());
    json_object_set_new(json, "radius-called-format", profile->radius_called_format ? json_string(profile->radius_called_format) : json_null());
    json_object_set_new(json, "technology", profile->technology ? json_string(profile->technology) : json_null());
    json_object_set_new(json, "wpa-pre-shared-key", profile->wpa_pre_shared_key ? json_string(profile->wpa_pre_shared_key) : json_null());
    json_object_set_new(json, "wpa2-pre-shared-key", profile->wpa2_pre_shared_key ? json_string(profile->wpa2_pre_shared_key) : json_null());
    json_object_set_new(json, "wpa3-pre-shared-key", profile->wpa3_pre_shared_key ? json_string(profile->wpa3_pre_shared_key) : json_null());

    return json;
}

json_t *host_request_to_json(HostRequest *request) {
    json_t *json = json_object();

    // Encode basic fields
    json_object_set_new(json, "type", request->type ? json_string(request->type) : json_null());
    json_object_set_new(json, "hostId", request->host_id ? json_string(request->host_id) : json_null());
    json_object_set_new(json, "login", request->login ? json_string(request->login) : json_null());
    json_object_set_new(json, "key", request->key ? json_string(request->key) : json_null());
    json_object_set_new(json, "simpleRotatedKey", request->simple_rotated_key ? json_string(request->simple_rotated_key) : json_null());
    json_object_set_new(json, "uptime", json_integer(request->uptime));
    json_object_set_new(json, "wanIp", request->wan_ip ? json_string(request->wan_ip) : json_null());
    json_object_set_new(json, "fwStatus", request->fw_status ? json_string(request->fw_status) : json_null());
    json_object_set_new(json, "clientInfo", request->client_info ? json_string(request->client_info) : json_null());
    json_object_set_new(json, "os", request->os ? json_string(request->os) : json_null());
    json_object_set_new(json, "osVersion", request->os_version ? json_string(request->os_version) : json_null());
    json_object_set_new(json, "fw", request->fw ? json_string(request->fw) : json_null());
    json_object_set_new(json, "hardwareMake", request->hardware_make ? json_string(request->hardware_make) : json_null());
    json_object_set_new(json, "hardwareModel", request->hardware_model ? json_string(request->hardware_model) : json_null());
    json_object_set_new(json, "hardwareModelNumber", request->hardware_model_number ? json_string(request->hardware_model_number) : json_null());
    json_object_set_new(json, "hardwareSerialNumber", request->hardware_serial_number ? json_string(request->hardware_serial_number) : json_null());
    json_object_set_new(json, "hardwareCpuInfo", request->hardware_cpu_info ? json_string(request->hardware_cpu_info) : json_null());
    json_object_set_new(json, "osBuildDate", json_integer(request->os_build_date));
    json_object_set_new(json, "reboot", json_integer(request->reboot));

    // Handle Cmd array
    json_t *cmd_array = json_array();
    for (int i = 0; i < request->cmd_count; i++) {
        json_array_append_new(cmd_array, cmd_to_json(&request->cmds[i]));
    }
    json_object_set_new(json, "cmds", cmd_array);

    // Handle Interfaces array
    json_t *interface_array = json_array();
    for (int i = 0; i < request->interface_count; i++) {
        json_array_append_new(interface_array, interface_to_json(&request->interfaces[i]));
    }
    json_object_set_new(json, "interfaces", interface_array);

    // Handle Wireless Configured array
    json_t *wifi_array = json_array();
    for (int i = 0; i < request->wireless_configured_count; i++) {
        json_array_append_new(wifi_array, wireless_interface_to_json(&request->wireless_configured[i]));
    }
    json_object_set_new(json, "wirelessConfigured", wifi_array);

    // Handle Security Profiles array
    json_t *security_profile_array = json_array();
    for (int i = 0; i < request->security_profile_count; i++) {
        json_array_append_new(security_profile_array, security_profile_to_json(&request->security_profiles[i]));
    }
    json_object_set_new(json, "securityProfiles", security_profile_array);

    // Encode other fields
    json_object_set_new(json, "webshellSupport", json_boolean(request->webshell_support));
    json_object_set_new(json, "bandwidthTestSupport", json_boolean(request->bandwidth_test_support));
    json_object_set_new(json, "firmwareUpgradeSupport", json_boolean(request->firmware_upgrade_support));
    json_object_set_new(json, "hostname", request->hostname ? json_string(request->hostname) : json_null());
    json_object_set_new(json, "outsideIp", request->outside_ip ? json_string(request->outside_ip) : json_null());
    json_object_set_new(json, "lastConfigRequest", json_integer(request->last_config_request));
    json_object_set_new(json, "usingWebSocket", json_boolean(request->using_websocket));

    return json;
}
