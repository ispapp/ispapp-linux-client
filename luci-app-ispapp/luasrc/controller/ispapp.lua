-- stub lua controller for 19.07 backward compatibility

module("luci.controller.ispapp", package.seeall)

function index()
	entry({"admin", "ispapp"}, firstchild(), _("ISPApp"), 60)
	entry({"admin", "ispapp", "overview"}, view("ispapp/overview"), _("Overview"), 10)
	entry({"admin", "ispapp", "logread"}, view("ispapp/settings"), _("Settings"), 20)
	entry({"admin", "ispapp", "logread"}, view("ispapp/logread"), _("Log View"), 30)
end