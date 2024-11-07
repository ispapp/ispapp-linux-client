-- /usr/lib/lua/luci/controller

---@diagnostic disable-next-line: deprecated
module("luci.controller.ispapp", package.seeall)

local fs = require("nixio.fs")
local util  = require("luci.util")
local templ = require("luci.template")
local i18n  = require("luci.i18n")

function index()
    -- Ensure the configuration file exists before proceeding
    if not fs.access("/etc/config/ispapp") then
        return
    end

    -- Top-level menu entry for ISPApp in LuCI
    entry({"admin", "ispapp"}, firstchild(), _("ISPApp"), 60).dependent = false

    -- Sub-menu entries under ISPApp
    entry({"admin", "ispapp", "overview"}, cbi("ispapp/overview"), _("Overview"), 10).leaf = true
    entry({"admin", "ispapp", "settings"}, cbi("ispapp/settings"), _("Settings"), 20).leaf = true
    entry({"admin", "ispapp", "logread"}, call("logread"), _("Log View"), 30).leaf = true
end

function logread()
	local logfile = util.exec("logread -e 'ispapp'")
	templ.render("ispapp/logread", {title = i18n.translate("ISPApp agent Logfile"), content = logfile})
end