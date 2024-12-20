-- /usr/lib/lua/luci/model/cbi/ispapp/
local util = require("luci.util")
local jsonc = require "luci.jsonc"
local ubus = require "ubus"
local conn = ubus.connect()
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
    luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/overview"))
end

-- Button for stopping service
local stop_button = s:option(Button, "_stop", translate("Stop Service"))
stop_button.inputstyle = "remove"
stop_button.write = function(self, section)
    luci.sys.exec("/etc/init.d/ispapp stop")
    luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/overview"))
end

-- Button for restarting service
local restart_button =
    s:option(Button, "_restart", translate("Restart Service"))
restart_button.inputstyle = "reset"
restart_button.write = function(self, section)
    luci.sys.exec("/etc/init.d/ispapp restart")
    luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/overview"))
end

-- Last Edit Time
local last_edit_time = s:option(DummyValue, "last_edit_time",
                                translate("Last Edit Time"))
last_edit_time.rawhtml = true
last_edit_time.cfgvalue = function(self, section)
    local json_result = conn:call("ispapp", "get_last_edit_time", {})
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
    local json_result = conn:call("ispapp", "get_active_time", {})
    return "<span style='font-style:italic;'>" ..
               (json_result and json_result.active_time or "N/A") .. "</span>"
end

-- CPU Usage
local cpu_usage = s:option(DummyValue, "cpu_usage", translate("CPU Usage"))
cpu_usage.rmempty = false
cpu_usage.rawhtml = true
cpu_usage.cfgvalue = function(self, section)
    local result = conn:call("ispapp", "get_cpu_usage", {})
    return "<span style='font-style:italic;'>" ..
               (result and result.cpu_usage or "N/A") .. "</span>"
end
local process_stats = conn:call("ispapp", "process_stats", {}) or
                          {
        VmRSS = "N/A",
        Threads = "N/A",
        Cpus_allowed = "N/A",
        State = "N/A"
    }
-- Memory Usage (VmRSS)
local memory_usage = s:option(DummyValue, "memory_usage",
                              translate("Memory Usage (VmRSS)"))
memory_usage.rmempty = false
memory_usage.rawhtml = true
memory_usage.cfgvalue = function(self, section)
    return "<span style='font-style:italic;'>" ..
               (process_stats and process_stats.VmRSS or "N/A") .. "</span>"
end

-- Threads
local threads = s:option(DummyValue, "threads", translate("Threads"))
threads.rmempty = false
threads.rawhtml = true
threads.cfgvalue = function(self, section)
    return "<span style='font-style:italic;'>" ..
               (process_stats and process_stats.Threads or "N/A") .. "</span>"
end

-- CPU Affinity
local cpu_affinity = s:option(DummyValue, "cpu_affinity",
                              translate("CPU Affinity"))
cpu_affinity.rmempty = false
cpu_affinity.rawhtml = true
cpu_affinity.cfgvalue = function(self, section)
    return "<span style='font-style:italic;'>" ..
               (process_stats and process_stats.Cpus_allowed or "N/A") ..
               "</span>"
end

-- Process State
local process_state = s:option(DummyValue, "process_state",
                               translate("Process State"))
process_state.rmempty = false
process_state.rawhtml = true
process_state.cfgvalue = function(self, section)
    return "<span style='font-style:italic;'>" ..
               (process_stats and process_stats.State or "N/A") .. "</span>"
end

-- Device Mode
local device_mode =
    s:option(DummyValue, "device_mode", translate("Device Mode"))
device_mode.rmempty = false
device_mode.rawhtml = true
device_mode.cfgvalue = function(self, section)
    local json_result = conn:call("ispapp", "get_device_mode", {})
    return "<span style='font-style:italic;'>" ..
               (json_result and json_result.device_mode or translate("Unknown")) ..
               "</span>"
end

-- Server Live Status
local server_live_status = s:option(DummyValue, "server_live_status",
                                    translate("Server Live Status"))
server_live_status.rawhtml = true
server_live_status.cfgvalue = function(self, section)
    local json_result = conn:call("ispapp", "check_domain", {})
    if json_result and json_result.status == 200 then
        return "<span style='color:green; font-weight:bold;'>" .. "Live" ..
                   "</span>"
    else
        return "<span style='color:red; font-weight:bold;'>" .. "Down" ..
                   "</span>"
    end
end

return m
