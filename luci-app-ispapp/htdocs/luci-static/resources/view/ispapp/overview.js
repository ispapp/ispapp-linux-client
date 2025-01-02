'use strict';
'require view';
'require form';
'require rpc';
'require fs';
'require ui';

var callGetLastEditTime = rpc.declare({
    object: 'ispapp',
    method: 'get_last_edit_time',
    expect: { '': {} }
});

var callGetActiveTime = rpc.declare({
    object: 'ispapp',
    method: 'get_active_time',
    expect: { '': {} }
});

var callGetCpuUsage = rpc.declare({
    object: 'ispapp',
    method: 'get_cpu_usage',
    expect: { '': {} }
});

var callProcessStats = rpc.declare({
    object: 'ispapp',
    method: 'process_stats',
    expect: { '': {} }
});

var callGetDeviceMode = rpc.declare({
    object: 'ispapp',
    method: 'get_device_mode',
    expect: { '': {} }
});

var callCheckDomain = rpc.declare({
    object: 'ispapp',
    method: 'check_domain',
    expect: { '': {} }
});

return view.extend({
    load: function() {
        return Promise.all([
            callGetLastEditTime(),
            callGetActiveTime(),
            callGetCpuUsage(),
            callProcessStats(),
            callGetDeviceMode(),
            callCheckDomain()
        ]);
    },

    render: function(data) {
        var lastEditTime = data[0];
        var activeTime = data[1];
        var cpuUsage = data[2];
        var processStats = data[3];
        var deviceMode = data[4];
        var serverLiveStatus = data[5];

        var m, s, o;

        m = new form.Map('ispapp', _('ISPApp Overview'), _('Overview of ISPApp service status and statistics.'));

        s = m.section(form.TypedSection, 'overview');
        s.anonymous = true;
        s.rmempty = false;

        o = s.option(form.DummyValue, 'service_status', _('Service Status'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return L.resolveDefault(fs.exec_direct('/bin/sh', ['-c', '[ -s /var/run/ispappd.pid ] && echo "running" || echo "stopped"']), 'stopped').then(function(res) {
                if (res.trim() === 'running') {
                    return '<span style="color:green; font-weight:bold;">Running</span>';
                } else {
                    return '<span style="color:red; font-weight:bold;">Stopped</span>';
                }
            });
        };

        o = s.option(form.Button, '_start', _('Start Service'));
        o.inputstyle = 'apply';
        o.onclick = function() {
            return fs.exec_direct('/etc/init.d/ispapp', ['start']).then(function() {
                ui.changes.apply();
            });
        };

        o = s.option(form.Button, '_stop', _('Stop Service'));
        o.inputstyle = 'remove';
        o.onclick = function() {
            return fs.exec_direct('/etc/init.d/ispapp', ['stop']).then(function() {
                ui.changes.apply();
            });
        };

        o = s.option(form.Button, '_restart', _('Restart Service'));
        o.inputstyle = 'reset';
        o.onclick = function() {
            return fs.exec_direct('/etc/init.d/ispapp', ['restart']).then(function() {
                ui.changes.apply();
            });
        };

        o = s.option(form.DummyValue, 'last_edit_time', _('Last Edit Time'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (lastEditTime && lastEditTime.last_edit_time || _('Unknown')) + '</span>';
        };

        o = s.option(form.DummyValue, 'active_time', _('Active Time'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (activeTime && activeTime.active_time || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'cpu_usage', _('CPU Usage'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (cpuUsage && cpuUsage.cpu_usage || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'memory_usage', _('Memory Usage (VmRSS)'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (processStats && processStats.VmRSS || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'threads', _('Threads'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (processStats && processStats.Threads || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'cpu_affinity', _('CPU Affinity'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (processStats && processStats.Cpus_allowed || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'process_state', _('Process State'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (processStats && processStats.State || 'N/A') + '</span>';
        };

        o = s.option(form.DummyValue, 'device_mode', _('Device Mode'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            return '<span style="font-style:italic;">' + (deviceMode && deviceMode.device_mode || _('Unknown')) + '</span>';
        };

        o = s.option(form.DummyValue, 'server_live_status', _('Server Live Status'));
        o.rawhtml = true;
        o.cfgvalue = function() {
            if (serverLiveStatus && serverLiveStatus.status === 200) {
                return '<span style="color:green; font-weight:bold;">Live</span>';
            } else {
                return '<span style="color:red; font-weight:bold;">Down</span>';
            }
        };

        return m.render();
    }
});
