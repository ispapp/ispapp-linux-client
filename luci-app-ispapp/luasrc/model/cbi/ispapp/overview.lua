-- /usr/lib/lua/luci/model/cbi/ispapp/
local dsp = require "luci.dispatcher"
local util = require("luci.util")
local jsonc = require "luci.jsonc"
local sys = require "luci.sys"
local uci = require"luci.model.uci".cursor()
local m, s, o

m = Map("ispapp", translate("ISPApp Overview"),
        translate("Overview of ISPApp service status and statistics."))

-- Section for displaying service status
s = m:section(TypedSection, "overview")
s.anonymous = true
s.rmempty = false

-- Service Status
local service_status = s:option(DummyValue, "service_status",
                                translate("Service Status"))
service_status.rmempty = false
service_status.rawhtml = true
service_status.cfgvalue = function(self, section)
    local pid_file_check = luci.sys.exec(
                               "[ -s /var/run/ispappd.pid ] && echo 'running' || echo 'stopped'")
    if pid_file_check:find("running") then
        return "<span style='color:green; font-weight:bold;'>" .. "Running" ..
                   "</span>"
    else
        return "<span style='color:red; font-weight:bold;'>" .. "Stopped" ..
                   "</span>"
    end
end

-- Button for starting service
local start_button = s:option(Button, "_start", translate("Start Service"))
start_button.inputstyle = "apply"
start_button.write = function(self, section)
    luci.sys.exec("/etc/init.d/ispapp start")
    -- luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/logread"))
end

-- Button for stopping service
local stop_button = s:option(Button, "_stop", translate("Stop Service"))
stop_button.inputstyle = "remove"
stop_button.write = function(self, section)
    luci.sys.exec("/etc/init.d/ispapp stop")
    -- luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/logread"))
end

-- Button for restarting service
local restart_button =
    s:option(Button, "_restart", translate("Restart Service"))
restart_button.inputstyle = "reset"
restart_button.write = function(self, section)
    luci.sys.exec("/etc/init.d/ispapp restart")
    -- luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/logread"))
end

-- Last Edit Time
local last_edit_time = s:option(DummyValue, "last_edit_time",
                                translate("Last Edit Time"))
last_edit_time.rawhtml = true
last_edit_time.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_last_edit_time")
    local json_result = jsonc.parse(result)
    return "<span style='font-style:italic;'>" ..
               (json_result and json_result.last_edit_time or
                   translate("Unknown")) .. "</span>"
end

-- Active Time
local active_time =
    s:option(DummyValue, "active_time", translate("Active Time"))
active_time.rmempty = false
active_time.rawhtml = true
active_time.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_active_time")
    local json_result = jsonc.parse(result)
    return "<span style='font-style:italic;'>" ..
               (json_result and json_result.active_time or "N/A") .. "</span>"
end

-- CPU Usage
local cpu_usage = s:option(DummyValue, "cpu_usage", translate("CPU Usage"))
cpu_usage.rmempty = false
cpu_usage.rawhtml = true
cpu_usage.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_cpu_usage")
    local ok, json_result = pcall(jsonc.parse, result)
    return "<span style='font-style:italic;'>" ..
               (ok and json_result.cpu_usage or "N/A") .. "</span>"
end

-- Device Mode
local device_mode =
    s:option(DummyValue, "device_mode", translate("Device Mode"))
device_mode.rmempty = false
device_mode.rawhtml = true
device_mode.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp get_device_mode")
    local json_result = jsonc.parse(result)
    return "<span style='font-style:italic;'>" ..
               (json_result and json_result.device_mode or translate("Unknown")) ..
               "</span>"
end

-- Server Live Status
local server_live_status = s:option(DummyValue, "server_live_status",
                                    translate("Server Live Status"))
server_live_status.rawhtml = true
server_live_status.cfgvalue = function(self, section)
    local result = util.exec("ubus call ispapp check_domain")
    local json_result = jsonc.parse(result)
    if json_result and json_result.status == 200 then
        return "<span style='color:green; font-weight:bold;'>" .. "Live" ..
                   "</span>"
    else
        return "<span style='color:red; font-weight:bold;'>" .. "Down" ..
                   "</span>"
    end
end

return m
