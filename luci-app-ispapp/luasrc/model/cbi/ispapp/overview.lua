-- /usr/lib/lua/luci/model/cbi/ispapp/overview.lua
local dsp = require "luci.dispatcher"
local util = require("luci.util")
local jsonc = require "luci.jsonc"
local m, s, o

m = Map("ispapp", translate("ISPApp Overview"),
        translate("Overview of ISPApp service status and statistics."))

-- Section for displaying service status
s = m:section(TypedSection, "overview", translate("Service Overview"))
s.anonymous = true
s.rmempty = false
-- s.addremove = false
-- s1.anonymous = true
-- Service Status
-- s.extedit = luci.dispatcher.build_url("admin/services/shadowsocks/servers/%s")

-- function s.create(...)
-- 	local sid = TypedSection.create(...)
-- 	if sid then
-- 		luci.http.redirect(s.main % sid)
-- 		return
-- 	end
-- end
-- s:option(Flag, "enabled", "Enabled")
-- status = s:option(DummyValue, "status", "Status")
-- status.rawhtml = true
-- status.cfgvalue = function(self, section)
--     local pid = luci.sys.exec("pidof ispappd")
--     if pid == "" then
--         return "<span style='color:red'>" .. "Stopped" .. "</span>"
--     else
--         return "<span style='color:green'>" .. "Running" .. "</span>"
--     end
-- end
local service_status = s:option(DummyValue, "service_status", translate("Service Status"))
service_status.rmempty = false
service_status.rawhtml = true
service_status.cfgvalue = function(self, section)
    local pid = luci.sys.exec("pidof ispappd")
    if pid == "" then
        return "<span style='color:red'>" .. "Stopped" .. "</span>"
    else
        return "<span style='color:green'>" .. "Running" .. "</span>"
    end
end

-- Last Edit Time
local last_edit_time = s:option(DummyValue, "last_edit_time",
                                translate("Last Edit Time"))
last_edit_time.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_last_edit_time")
    local json_result = jsonc.parse(result)
    return json_result and json_result.last_edit_time or translate("Unknown")
end

-- Active Time
local active_time =
    s:option(DummyValue, "active_time", translate("Active Time"))
service_status.rmempty = false
active_time.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_active_time")
    local json_result = jsonc.parse(result)
    return json_result and json_result.active_time or "N/A"
end

-- CPU Usage
local cpu_usage = s:option(DummyValue, "cpu_usage", translate("CPU Usage"))
service_status.rmempty = false
cpu_usage.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_cpu_usage")
    local json_result = jsonc.parse(result)
    return json_result and json_result.cpu_usage or translate("Unknown")
end

-- Device Mode
local device_mode =
s:option(DummyValue, "device_mode", translate("Device Mode"))
service_status.rmempty = false
device_mode.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_device_mode")
    local json_result = jsonc.parse(result)
    return json_result and json_result.device_mode or translate("Unknown")
end

return m
