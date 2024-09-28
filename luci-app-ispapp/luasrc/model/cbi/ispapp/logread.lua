m = SimpleForm("logread", translate("ISP App Log"), translate("Displays logs for ispapp"))

log_section = m:section(SimpleSection)
log_section.template = "admin_system/logread"
log_section.title = translate("ISP App Log Viewer")
log_section.cfgvalue = function(self, section)
    return luci.sys.exec("logread | grep 'ispapp'")
end

return m
