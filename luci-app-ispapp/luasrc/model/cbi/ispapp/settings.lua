-- /usr/lib/lua/luci/model/cbi/ispapp/
local dsp = require "luci.dispatcher"
local util = require("luci.util")
local uci = luci.model.uci.cursor()
local m, s, o

-- Map object (the form itself)
m = Map("ispapp", translate("ISPApp Configuration"),
        translate("Configure ISPApp settings."))

-- Section for configuration settings
s = m:section(TypedSection, "settings", translate("Settings"))
s.anonymous = true
s.rmempty = false
-- s.addremove = true

-- Options for each configuration setting

-- Enable/Disable ISPApp
enabled = s:option(Flag, "enabled", translate("Enabled"),
                   translate("Enable ISPApp."))
enabled.rmempty = false
enabled.datatype = "bool"
enabled.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "@settings[0]", "enabled", m:get(section, "enabled"))
    uci:commit("ispapp")
end
enabled.default = uci:get("ispapp", "@settings[0]", "enabled") or 0

-- MAC Address for Login
login = s:option(Value, "login", translate("Login"),
                 translate("MAC address for login."))
login.rmempty = false
login.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "login", "login", m:get(section, "login"))
    uci:commit("ispapp")
end
login.default = uci:get("ispapp", "@settings[0]", "login") or "N/A"

-- Top Domain
Domain = s:option(Value, "Domain", translate("Domain"),
                  translate("Domain name for top domain."))
Domain.rmempty = false
Domain.datatype = "host"
Domain.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "Domain", m:get(section, "Domain"))
    uci:commit("ispapp")
end
Domain.default = uci:get("ispapp", "@settings[0]", "Domain") or "N/A"

-- Top Listener Port
ListenerPort = s:option(Value, "ListenerPort", translate("Listener Port"),
                        translate("Port for listening."))
ListenerPort.rmempty = false
ListenerPort.datatype = "port"
ListenerPort.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "@settings[0]", "ListenerPort",
            m:get(section, "ListenerPort"))
    uci:commit("ispapp")
end
ListenerPort.default = uci:get("ispapp", "@settings[0]", "ListenerPort") or
                           "N/A"

-- Top Key
Key = s:option(Value, "Key", translate("Key"),
               translate("Key for authentication."))
Key.rmempty = true
Key.datatype = "string"
Key.password = true
Key.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    local key = m:get(section, "Key")
    uci:set("ispapp", "@settings[0]", "Key", key)
    local ok, _, _ = pcall(util.exec, "fw_setenv Key " .. key)
    if ok then
        debug.traceback(
            "\npersisting Key into the fw boot environment:" .. key .. "\n", 2)
    end
    uci:commit("ispapp")
end
Key.default = uci:get("ispapp", "@settings[0]", "Key") or "N/A"

-- Access Token (if applicable)
accessToken = s:option(DummyValue, "accessToken", translate("Access Token"),
                       translate("Access token for API access."))
accessToken.rmempty = true
accessToken.rawhtml = true
accessToken.password = true
accessToken.cfgvalue = function(self, section)
    local token = uci:get("ispapp", "@settings[0]", "accessToken") or "N/A"
    local displayToken = #token > 25 and token:sub(1, 25) .. "..." or token

    return string.format([[
       <div style="display: flex; align-items: center;">
            <input type="password" value="%s" readonly style="margin-right: 10px; min-width: 200px; border: 1px solid #ccc; padding: 5px; border-radius: 4px; display: inline-block;">
            <button type="button" id="copyButton1" style="cursor: pointer;font-size:1.5em;padding: 3px;border-radius: 7px;outline: none;">📋</button>
        </div>
        <script>
            document.getElementById('copyButton1').addEventListener('click', () => setClipboard('%s'));
            async function setClipboard(text) {
                  if (navigator.clipboard) {
                    navigator.clipboard.writeText(text)
                  } else {
                    const input = document.createElement('textarea')
                    input.value = text
                    document.body.appendChild(input)
                    input.select()
                    document.execCommand('copy')
                    document.body.removeChild(input)
                  }
            }
        </script>
    ]], displayToken, token)
end

-- Refresh Token (if applicable)
refreshToken = s:option(DummyValue, "refreshToken", translate("Refresh Token"),
                        translate("Refresh token for API access."))
refreshToken.rmempty = true
refreshToken.nocreate = true
refreshToken.rawhtml = true
refreshToken.password = true
refreshToken.cfgvalue = function(self, section)
    local token = uci:get("ispapp", "@settings[0]", "refreshToken") or "N/A"
    local displayToken = #token > 25 and token:sub(1, 25) .. "..." or token

    return string.format([[
       <div style="display: flex; align-items: center;">
            <input type="password" value="%s" readonly style="margin-right: 10px; min-width: 200px; border: 1px solid #ccc; padding: 5px; border-radius: 4px; display: inline-block;">
            <button type="button" id="copyButton2" style="cursor: pointer;font-size:1.5em;padding: 3px;border-radius: 7px;outline: none;">📋</button>
        </div>
        <script>
            document.getElementById('copyButton2').addEventListener('click', () => setClipboard('%s'));
            async function setClipboard(text) {
                  if (navigator.clipboard) {
                    navigator.clipboard.writeText(text)
                  } else {
                    const input = document.createElement('textarea')
                    input.value = text
                    document.body.appendChild(input)
                    input.select()
                    document.execCommand('copy')
                    document.body.removeChild(input)
                  }
            }
        </script>
    ]], displayToken, token)
end

-- Connection Status
connected = s:option(DummyValue, "connected", translate("Connected"),
                     translate("Connection status."))
connected.rmempty = false
connected.readonly = true
connected.rawhtml = true
connected.cfgvalue = function(self, section)
    -- Check the connection status from UCI
    local status = uci:get("ispapp", "@settings[0]", "connected") or '0'

    -- Return HTML string based on connection status
    if status == '1' then
        return "<span style='color: green;'>" .. translate("Connected") ..
                   "</span>"
    else
        return "<span style='color: red;'>" .. translate("Disconnected") ..
                   "</span>"
    end
end

-- Intervals for update and check
updateInterval = s:option(Value, "UpdateInterval", translate("Update Interval"),
                          translate("Interval for updating data in seconds."))
updateInterval.rmempty = false
updateInterval.datatype = "range(10, 100)"
updateInterval.default = uci:get("ispapp", "@settings[0]", "updateInterval") or
                             10
updateInterval.placeholder = "10"
updateInterval.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "@settings[0]", "updateInterval",
            m:get(section, "updateInterval"))
    uci:commit("ispapp")
end

-- -- IP Bandwidth Test Server
-- ipbandswtestserver = s:option(Value, "ipbandswtestserver",
--                               translate("IP Bandwidth Test Server"), translate(
--                                   "IP address for bandwidth test server."))
-- ipbandswtestserver.rmempty = false
-- ipbandswtestserver.datatype = "ipaddr"
-- ipbandswtestserver.onchange = function()
--     -- Process the value if needed (e.g., trim whitespace)
--     uci:set("ispapp", "settings", "ipbandswtestserver",
--             m:get(section, "ipbandswtestserver"))
--     uci:commit("ispapp")
-- end
-- -- BT User
-- btuser = s:option(Value, "btuser", translate("BT User"),
--                   translate("Username for BT."))
-- btuser.rmempty = false
-- btuser.onchange = function()
--     -- Process the value if needed (e.g., trim whitespace)
--     uci:set("ispapp", "settings", "btuser", m:get(section, "btuser"))
--     uci:commit("ispapp")
-- end
-- btuser.default = uci:get("ispapp", "settings", "btuser") or "N/A"

-- -- BT Password
-- btpwd = s:option(Value, "btpwd", translate("BT Password"),
--                  translate("Password for BT."))
-- btpwd.rmempty = false
-- btpwd.password = true
-- btpwd.onchange = function()
--     -- Process the value if needed (e.g., trim whitespace)
--     uci:set("ispapp", "settings", "btpwd", m:get(section, "btpwd"))
--     uci:commit("ispapp")
-- end
-- btpwd.default = uci:get("ispapp", "settings", "btpwd") or "N/A"
-- Button for saving and applying settings
-- local apply = s:option(Button, "_apply", translate("Save and Apply"))
-- apply.inputstyle = "apply"
-- apply.write = function(self, section)
--     uci.set("ispapp", "settings", "enabled", m:get(section, "enabled"))
--     uci.set("ispapp", "settings", "Domain", m:get(section, "Domain"))
--     uci.set("ispapp", "settings", "ListenerPort", m:get(section, "ListenerPort"))
--     uci.set("ispapp", "settings", "topSmtpPort", m:get(section, "topSmtpPort"))
--     uci.set("ispapp", "settings", "Key", m:get(section, "Key"))
--     uci.set("ispapp", "settings", "ipbandswtestserver", m:get(section, "ipbandswtestserver"))
--     uci.set("ispapp", "settings", "btuser", m:get(section, "btuser"))
--     uci.set("ispapp", "settings", "btpwd", m:get(section, "btpwd"))
--     uci.commit("ispapp")
--     -- Submit the form data via RPC call
--     luci.sys.exec("/etc/init.d/ispapp restart")
--     luci.http.redirect(dsp.build_url("admin/ispapp/settings"))
--     luci.http.write(translate("Configuration changes have been saved."))
-- end
return m
