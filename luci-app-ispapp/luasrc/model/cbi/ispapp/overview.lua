local dsp = require "luci.dispatcher"
local m, s, o

m = Map("ispapp", translate("ISPApp Overview"), translate("Overview of ISPApp service status and statistics."))

-- Section for displaying service status
s = m:section(TypedSection, "overview", translate("Service Overview"))
s.anonymous = true

-- Service Status
local service_status = s:option(DummyValue, "service_status", translate("Service Status"))
function service_status.va(self, section)
    local result = util.exec("ubus call luci.ispapp get_service_status")
    local json_result = jsonc.parse(result)
    return json_result and json_result.service_status or translate("Unknown")
end

-- Last Edit Time
local last_edit_time = s:option(DummyValue, "last_edit_time", translate("Last Edit Time"))
function last_edit_time.cfgvalue(self, section)
    local result = util.exec("ubus call luci.ispapp get_last_edit_time")
    local json_result = jsonc.parse(result)
    return json_result and json_result.last_edit_time or translate("Unknown")
end

-- Active Time
local active_time = s:option(DummyValue, "active_time", translate("Active Time"))
function active_time.cfgvalue(self, section)
    local result = util.exec("ubus call luci.ispapp get_active_time")
    local json_result = jsonc.parse(result)
    return json_result and json_result.active_time or translate("Unknown")
end

-- CPU Usage
local cpu_usage = s:option(DummyValue, "cpu_usage", translate("CPU Usage"))
function cpu_usage.cfgvalue(self, section)
    local result = util.exec("ubus call luci.ispapp get_cpu_usage")
    local json_result = jsonc.parse(result)
    return json_result and json_result.cpu_usage or translate("Unknown")
end

-- Network Speed
local network_speed = s:option(DummyValue, "network_speed", translate("Network Speed"))
function network_speed.cfgvalue(self, section)
    local result = util.exec("ubus call luci.ispapp get_network_speed")
    local json_result = jsonc.parse(result)
    return json_result and json_result.network_speed or translate("Unknown")
end

-- Device Mode
local device_mode = s:option(DummyValue, "device_mode", translate("Device Mode"))
function device_mode.cfgvalue(self, section)
    local result = util.exec("ubus call luci.ispapp get_device_mode")
    local json_result = jsonc.parse(result)
    return json_result and json_result.device_mode or translate("Unknown")
end

return m