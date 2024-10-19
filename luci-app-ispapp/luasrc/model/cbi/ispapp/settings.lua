local dsp = require "luci.dispatcher"
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
enabled = s:option(Flag, "enabled", translate("Enabled"), translate("Enable ISPApp."))
enabled.rmempty = false
enabled.datatype = "bool"
enabled.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "enabled", m:get(section, "enabled"))
    uci:commit("ispapp")
end
enabled.default = uci:get("ispapp", "settings", "enabled") or 0

-- MAC Address for Login
login = s:option(Value, "login", translate("Login"), translate("MAC address for login."))
login.rmempty = false
login.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "login", "login", m:get(section, "login"))
    uci:commit("ispapp")
end
login.default = uci:get("ispapp", "settings", "login") or "N/A"

-- Top Domain
topDomain = s:option(Value, "topDomain", translate("Top Domain"), translate("Domain name for top domain."))
topDomain.rmempty = false
topDomain.datatype = "host"
topDomain.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "topDomain", m:get(section, "topDomain"))
    uci:commit("ispapp")
end
topDomain.default = uci:get("ispapp", "settings", "topDomain") or "N/A"

-- Top Listener Port
topListenerPort = s:option(Value, "topListenerPort", translate("Top Listener Port"), translate("Port for listening."))
topListenerPort.rmempty = false
topListenerPort.datatype = "port"
topListenerPort.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "topListenerPort", m:get(section, "topListenerPort"))
    uci:commit("ispapp")
end
topListenerPort.default = uci:get("ispapp", "settings", "topSmtpPort") or "N/A"
-- Top SMTP Port
topSmtpPort = s:option(Value, "topSmtpPort", translate("Top SMTP Port"), translate("SMTP port."))
topSmtpPort.rmempty = false
topSmtpPort.datatype = "port"
topSmtpPort.default = uci:get("ispapp", "settings", "topSmtpPort") or "N/A"
topSmtpPort.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "topSmtpPort", m:get(section, "topSmtpPort"))
    uci:commit("ispapp")
end

-- Top Key
topKey = s:option(Value, "topKey", translate("Top Key"), translate("Key for authentication."))
topKey.rmempty = true
topKey.datatype = "string"
topKey.password = true
topKey.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "topKey", m:get(section, "topKey"))
    uci:commit("ispapp")
end
topKey.default = uci:get("ispapp", "settings", "topKey") or "N/A"

-- Access Token (if applicable)
accessToken = s:option(DummyValue, "accessToken", translate("Access Token"), translate("Access token for API access."))
accessToken.rmempty = true
accessToken.rawhtml = true
accessToken.password = true
accessToken.cfgvalue = function(self, section)
    local token = uci:get("ispapp", "settings", "accessToken") or "N/A"
    
    -- Return HTML string with the token and a copy button
    return string.format([[
       <div style="display: flex; align-items: center;">
            <span style="margin-right: 10px; min-width: 200px; border: 1px solid #ccc; padding: 5px; border-radius: 4px; display: inline-block;">%s</span>
            <button type="button" id="copyButton1" style="cursor: pointer;">Copy</button>
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
                  alert('accessToken copied')
            }
        </script>
    ]], token, token)
end

-- Refresh Token (if applicable)
refreshToken = s:option(DummyValue, "refreshToken", translate("Refresh Token"), translate("Refresh token for API access."))
refreshToken.rmempty = true
refreshToken.nocreate = true
refreshToken.rawhtml = true
refreshToken.password = true
refreshToken.cfgvalue = function(self, section)
    local token = uci:get("ispapp", "settings", "refreshToken") or "N/A"
    
    -- Return HTML string with the token and a copy button
    return string.format([[
       <div style="display: flex; align-items: center;">
            <span style="margin-right: 10px; min-width: 200px; border: 1px solid #ccc; padding: 5px; border-radius: 4px; display: inline-block;">%s</span>
            <button type="button" id="copyButton2" style="cursor: pointer;">Copy</button>
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
                  alert('refreshToken copied')
            }
        </script>
    ]], token, token)
end

-- Connection Status
connected = s:option(DummyValue, "connected", translate("Connected"), translate("Connection status."))
connected.rmempty = false
connected.readonly = true
connected.rawhtml = true
connected.cfgvalue = function(self, section)
    -- Check the connection status from UCI
    local status = uci:get("ispapp", "settings", "connected") or '0'
    
    -- Return HTML string based on connection status
    if status == '1' then
        return "<span style='color: green;'>" .. translate("Connected") .. "</span>"
    else
        return "<span style='color: red;'>" .. translate("Disconnected") .. "</span>"
    end
end

-- IP Bandwidth Test Server
ipbandswtestserver = s:option(Value, "ipbandswtestserver", translate("IP Bandwidth Test Server"), translate("IP address for bandwidth test server."))
ipbandswtestserver.rmempty = false
ipbandswtestserver.datatype = "ipaddr"
ipbandswtestserver.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "ipbandswtestserver", m:get(section, "ipbandswtestserver"))
    uci:commit("ispapp")
end
topSmtpPort.default = uci:get("ispapp", "settings", "ipbandswtestserver") or "N/A"

-- BT User
btuser = s:option(Value, "btuser", translate("BT User"), translate("Username for BT."))
btuser.rmempty = false
btuser.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "btuser", m:get(section, "btuser"))
    uci:commit("ispapp")
end
btuser.default = uci:get("ispapp", "settings", "btuser") or "N/A"

-- BT Password
btpwd = s:option(Value, "btpwd", translate("BT Password"), translate("Password for BT."))
btpwd.rmempty = false
btpwd.password = true
btpwd.onchange = function()
    -- Process the value if needed (e.g., trim whitespace)
    uci:set("ispapp", "settings", "btpwd", m:get(section, "btpwd"))
    uci:commit("ispapp")
end
btpwd.default = uci:get("ispapp", "settings", "btpwd") or "N/A"
-- Button for saving and applying settings
-- local apply = s:option(Button, "_apply", translate("Save and Apply"))
-- apply.inputstyle = "apply"
-- apply.write = function(self, section)
--     uci.set("ispapp", "settings", "enabled", m:get(section, "enabled"))
--     uci.set("ispapp", "settings", "topDomain", m:get(section, "topDomain"))
--     uci.set("ispapp", "settings", "topListenerPort", m:get(section, "topListenerPort"))
--     uci.set("ispapp", "settings", "topSmtpPort", m:get(section, "topSmtpPort"))
--     uci.set("ispapp", "settings", "topKey", m:get(section, "topKey"))
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
