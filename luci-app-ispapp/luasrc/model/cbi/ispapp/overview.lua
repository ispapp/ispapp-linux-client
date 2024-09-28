local m, s, o

m = SimpleForm("ispapp_status", translate("ISP App Status"), translate("Displays current status of the ISP App"))

-- Define sections and options
s = m:section(SimpleSection)
s.title = translate("Service Overview")

-- Service Status
o = s:option(DummyValue, "service_status", translate("Service Status"))
o.rawhtml = true
o.cfgvalue = function(self, section)
    local status = luci.sys.call("ubus call luci.ispapp get_service_status | jsonfilter -e '@.response'")
    if status:match("Running") then
        return translate("Running")
    else
        return translate("Not Running")
    end
end

-- Last Edit Time
o = s:option(DummyValue, "last_edit_time", translate("Last Edit Time"))
o.cfgvalue = function(self, section)
    return luci.sys.exec("ubus call luci.ispapp get_last_edit_time | jsonfilter -e '@.response.last_edit_time'") or "Unknown"
end

-- Active Time
o = s:option(DummyValue, "active_time", translate("Active Time"))
o.cfgvalue = function(self, section)
    return luci.sys.exec("ubus call luci.ispapp get_active_time | jsonfilter -e '@.response.active_time'") or "Unknown"
end

-- CPU Usage
o = s:option(DummyValue, "cpu_usage", translate("CPU Usage"))
o.cfgvalue = function(self, section)
    return luci.sys.exec("ubus call luci.ispapp get_cpu_usage | jsonfilter -e '@.response.cpu_usage'") or "Unknown"
end

-- Network Speed
o = s:option(DummyValue, "network_speed", translate("Network Speed"))
o.cfgvalue = function(self, section)
    return luci.sys.exec("ubus call luci.ispapp get_network_speed | jsonfilter -e '@.response.network_speed'") or "Unknown"
end

-- Device Mode
o = s:option(DummyValue, "device_mode", translate("Device Mode"))
o.cfgvalue = function(self, section)
    return luci.sys.exec("ubus call luci.ispapp get_device_mode | jsonfilter -e '@.response.device_mode'") or "Unknown"
end

-- Buttons for actions
edit_btn = s:option(Button, "_edit", translate("Edit Settings"))
edit_btn.inputstyle = "apply"
edit_btn.write = function(self, section)
    luci.http.redirect(luci.dispatcher.build_url("admin/ispapp/settings"))
end

sign_in_btn = s:option(Button, "_sign_in", translate("Sign In to ISPApp Cloud"))
sign_in_btn.inputstyle = "apply"
sign_in_btn.write = function(self, section)
    luci.sys.exec("xdg-open 'https://ispapp.co/signin'")
end

return m
