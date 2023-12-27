local http = require "luci.http"
module("luci.controller.ispapp", package.seeall)

function index()
	if not nixio.fs.access("/etc/config/ispapp") then
        return
    end
	entry({"admin", "services", "ispapp"},				call("ispappd_template"), _("Ispapp"), 21).dependent = true
	entry({"admin", "services", "ispapp", "config"}, 	call("ispappd_config"))
	entry({"admin", "services", "ispapp", "status"}, 	call("ispappd_status"))
	entry({"admin", "services", "ispapp", "logout"}, 	call("ispappd_logout"))
end

function ispappd_template()
    luci.template.render("ispapp/main")
end


function getIspappConfig()
    local uci = require "luci.model.uci".cursor()

    local enabled = uci:get("ispapp", "settings", "enabled") or "0"
    local login = uci:get("ispapp", "settings", "login") or "00:00:00:00:00:00"
    local topDomain = uci:get("ispapp", "settings", "topDomain") or "qwer.ispapp.co"
    local topListenerPort = uci:get("ispapp", "settings", "topListenerPort") or "8550"
    local topSmtpPort = uci:get("ispapp", "settings", "topSmtpPort") or "8465"
    local topKey = uci:get("ispapp", "settings", "topKey") or ""
    local ipbandswtestserver = uci:get("ispapp", "settings", "ipbandswtestserver") or "3.239.254.95"
    local btuser = uci:get("ispapp", "settings", "btuser") or "btest"
    local btpwd = uci:get("ispapp", "settings", "btpwd") or "0XSYIGkRlP6MUQJMZMdrogi2"

    local result = {
        enabled = (enabled == "1"),
        login = login,
        topDomain = topDomain,
        topListenerPort = topListenerPort,
        topSmtpPort = topSmtpPort,
        topKey = topKey,
        ipbandswtestserver = ipbandswtestserver,
        btuser = btuser,
        btpwd = btpwd
    }
    return result
end


function submitIspappConfig(req)
    local uci = require "luci.model.uci".cursor()
    
    if req.enabled ~= nil then uci:set("ispapp", "@settings[0]", "enabled", req.enabled) end
    if req.login ~= nil then uci:set("ispapp", "@settings[0]", "login", req.login) end
    if req.topDomain ~= nil then uci:set("ispapp", "@settings[0]", "topDomain", req.topDomain) end
    if req.topListenerPort ~= nil then uci:set("ispapp", "@settings[0]", "topListenerPort", req.topListenerPort) end
    if req.topKey ~= nil then uci:set("ispapp", "@settings[0]", "topKey", req.topKey) end
    if req.ipbandswtestserver ~= nil then uci:set("ispapp", "@settings[0]", "ipbandswtestserver", req.ipbandswtestserver) end
    if req.btuser ~= nil then uci:set("ispapp", "@settings[0]", "btuser", req.btuser) end
    if req.btpwd ~= nil then uci:set("ispapp", "@settings[0]", "btpwd", req.btpwd) end
    
    uci:commit("ispapp")
end


function ispappd_config()
	local http = require "luci.http"
	http.prepare_content("application/json")
	local method = http.getenv("REQUEST_METHOD")
	if method == "post" or method == "POST" then
		local content = http.content()
		local jsonc = require "luci.jsonc"
		local json_parse = jsonc.parse
		local req = json_parse(content)
		if req == nil or next(req) == nil then
			luci.http.write_json({
				error =  "invalid request"
			})
			return 
		end
		submitIspappConfig(req)
		if req.enabled == true then
			luci.util.exec("/etc/init.d/ispapp start")
		else
			luci.util.exec("/etc/init.d/ispapp stop")
		end
	end
	local response = getIspappConfig()
    luci.http.write_json(response)
end

function ispappd_status()
	local sys  = require "luci.sys"
	local http = require "luci.http" 
    -- http.prepare_content("text/plain;charset=utf-8")
	http.prepare_content("application/json")
	local text = sys.exec("ispappd status")
    http.write(text)
end