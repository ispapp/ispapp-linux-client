'use strict';
'require rpc';
'require view';
'require ui';

// RPC calls to get various information
var get_status = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_service_status',
    expect: { 'response': 'json' }
});

var get_last_edit_time = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_last_edit_time',
    expect: { 'response': 'json' }
});

var get_active_time = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_active_time',
    expect: { 'response': 'json' }
});

var get_cpu_usage = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_cpu_usage',
    expect: { 'response': 'json' }
});

var get_network_speed = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_network_speed',
    expect: { 'response': 'json' }
});

var get_device_mode = rpc.declare({
    object: 'luci.ispapp',
    method: 'get_device_mode',
    expect: { 'response': 'json' }
});

// Function to handle editing settings
function editSettings() {
    window.location = UCI.createURL('admin/ispapp/settings');
}

// Function to handle sign in to ISPApp cloud
function signInToISPApp() {
    window.open('https://ispapp.co/signin', '_blank');
}

return view.extend({
    load: function () {
        return Promise.all([
            get_status(),
            get_last_edit_time(),
            get_active_time(),
            get_cpu_usage(),
            get_network_speed(),
            get_device_mode()
        ]);
    },

    render: function (data) {
        var status = data[0] || {};
        var lastEditTime = data[1] || {};
        var activeTime = data[2] || {};
        var cpuUsage = data[3] || {};
        var networkSpeed = data[4] || {};
        var deviceMode = data[5] || {};

        var statusMessage = (status.service_status === "Running" ? 'Running' : 'Not Running');

        // Create status display
        var statusDisplay = E('div', { 'class': 'status-overview' }, [
            E('h2', {}, _('Service Status')),
            E('p', {}, _('Service is currently: ') + statusMessage),
            E('p', {}, _('Last Edit Time: ') + (lastEditTime.last_edit_time || 'Unknown')),
            E('p', {}, _('Active Time: ') + (activeTime.active_time || 'Unknown')),
            E('p', {}, _('CPU Usage: ') + (cpuUsage.cpu_usage || 'Unknown')),
            E('p', {}, _('Network Speed: ') + (networkSpeed.network_speed || 'Unknown')),
            E('p', {}, _('Device Mode: ') + (deviceMode.device_mode || 'Unknown'))
        ]);

        // Create buttons for actions
        var buttons = E('div', { 'class': 'overview-buttons' }, [
            E('button', { 'class': 'btn-edit-settings', 'click': editSettings }, _('Edit Settings')),
            E('button', { 'class': 'btn-sign-in', 'click': signInToISPApp }, _('Sign In to ISPApp Cloud'))
        ]);

        // Combine status display and buttons
        return E('div', { 'class': 'overview-container' }, [
            statusDisplay,
            buttons
        ]);
    }
});
