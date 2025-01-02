-- /usr/lib/lua/luci/controller
-- -@diagnostic disable-next-line: deprecated
module("luci.controller.ispapp", package.seeall)


function index()
    -- Top-level menu entry for ISPApp in LuCI
    entry({"admin", "ispapp"}, firstchild(), _("ISPApp"), 87).acl_depends = {"luci-app-ispapp"}

    -- Sub-menu entries under ISPApp
    entry({"admin", "ispapp", "overview"}, view("ispapp/overview"), _("Overview"), 10)
    entry({"admin", "ispapp", "settings"}, view("ispapp/settings"), _("Settings"), 20)
    entry({"admin", "ispapp", "logread"}, view("ispapp/logread"), _("Log View"), 30)
end
