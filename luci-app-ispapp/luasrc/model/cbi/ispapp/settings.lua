local dsp = require "luci.dispatcher"
local m, s, o

-- Map object (the form itself)
m = Map("ispapp", translate("ISPApp Configuration"),
    translate("Configure ISPApp settings."))

-- Section for configuration settings
s = m:section(TypedSection, "settings", translate("Settings"))
s.anonymous = true
s.addremove = false

-- Options for each configuration setting
o = s:option(Flag, "enabled", translate("Enabled"), translate("Enable ISPApp."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.enabled'") or '0'
o.rmempty = false

o = s:option(Value, "login", translate("Login"), translate("MAC address for login."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.login'") or ''
o.rmempty = false

o = s:option(Value, "topDomain", translate("Top Domain"), translate("Domain name for top domain."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.topDomain'") or ''
o.rmempty = false

o = s:option(Value, "topListenerPort", translate("Top Listener Port"), translate("Port for listening."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.topListenerPort'") or ''
o.rmempty = false

o = s:option(Value, "topSmtpPort", translate("Top SMTP Port"), translate("SMTP port."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.topSmtpPort'") or ''
o.rmempty = false

o = s:option(Value, "topKey", translate("Top Key"), translate("Key for authentication."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.topKey'") or ''
o.rmempty = true

o = s:option(Value, "ipbandswtestserver", translate("IP Bandwidth Test Server"), translate("IP address for bandwidth test server."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.ipbandswtestserver'") or ''
o.rmempty = false

o = s:option(Value, "btuser", translate("BT User"), translate("Username for BT."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.btuser'") or ''
o.rmempty = false

o = s:option(Value, "btpwd", translate("BT Password"), translate("Password for BT."))
o.default = luci.sys.exec("ubus call luci.ispapp read_ispapp_config | jsonfilter -e '@.response.btpwd'") or ''
o.rmempty = false

-- Button for saving and applying settings
local apply = s:option(Button, "_apply", translate("Save and Apply"))
apply.inputstyle = "apply"
apply.write = function(self, section)
    local data = {
        enabled = m:get(section, "enabled"),
        login = m:get(section, "login"),
        topDomain = m:get(section, "topDomain"),
        topListenerPort = m:get(section, "topListenerPort"),
        topSmtpPort = m:get(section, "topSmtpPort"),
        topKey = m:get(section, "topKey"),
        ipbandswtestserver = m:get(section, "ipbandswtestserver"),
        btuser = m:get(section, "btuser"),
        btpwd = m:get(section, "btpwd")
    }

    -- Submit the form data via RPC call
    local cmd = "ubus call luci.ispapp write_ispapp_config '" .. luci.util.serialize_json({ config = data }) .. "'"
    local result = luci.sys.exec(cmd)
    
    if result:match('"error"') then
        luci.http.redirect(dsp.build_url("admin/ispapp/settings"))
        luci.http.write(translate("Unable to save changes: ") .. result)
    else
        luci.http.redirect(dsp.build_url("admin/ispapp/settings"))
        luci.http.write(translate("Configuration changes have been saved."))
    end
end

return m
