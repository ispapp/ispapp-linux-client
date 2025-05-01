'use strict';
'require view';
'require form';
'require request';
'require rpc';
'require fs';
'require uci';
'require ui';

var callCheckConnection = rpc.declare({
    object: 'ispapp',
    method: 'checkconnection',
    params: {}
});
var signup = rpc.declare({
    object: 'ispapp',
    method: 'signup',
    params: {}
});

return view.extend({
    render: async function() {
        const command_exist = async (cmd)=>{
            return await fs.exec("which", [cmd]).then((res)=>{
                if (res.stdout) {
                    return true;
                }
                return false;
            });
        }
        const version = {
            v2: false,
            v3: false,
        }
        version.v2 = await command_exist('iperf')
        version.v3 = await command_exist('iperf3')
        var m, s, o;
        uci.load('ispapp');
        m = new form.Map('ispapp', _('ISPApp Configuration'), _('Configure ISPApp settings.'));
        s = m.section(form.TypedSection, 'settings', _('Settings'));
        s.anonymous = true;
        s.rmempty = false;

        o = s.option(form.Value, 'login', _('Login'), _('MAC address for login.'));
        o.rmempty = false;
        o.default = uci.get('ispapp', '@settings[0]', 'login') || 'N/A';
        o.onchange = function() {
            uci.set('ispapp', '@settings[0]', 'login', this.formvalue(this.section));
        };

        o = s.option(form.Value, 'Domain', _('Domain'), _('Domain name for top domain.'));
        o.rmempty = false;
        o.datatype = 'host';
        o.default = uci.get('ispapp', '@settings[0]', 'Domain') || 'N/A';
        o.onchange = function() {
            uci.set('ispapp', '@settings[0]', 'Domain', this.formvalue(this.section));
        };

        o = s.option(form.Value, 'ListenerPort', _('Listener Port'), _('Port for listening.'));
        o.rmempty = false;
        o.datatype = 'port';
        o.default = uci.get('ispapp', '@settings[0]', 'ListenerPort') || 'N/A';
        o.onchange = function() {
            uci.set('ispapp', '@settings[0]', 'ListenerPort', this.formvalue(this.section));
        };

        o = s.option(form.Value, 'Key', _('Key'), _('Key for authentication.'));
        o.rmempty = true;
        o.datatype = 'string';
        o.password = true;
        o.default = uci.get('ispapp', '@settings[0]', 'Key') || 'N/A';
        o.onchange = function() {
            var key = this.formvalue(this.section);
            uci.set('ispapp', '@settings[0]', 'Key', key);
            fs.exec('fw_setenv Key ' + key);
        };

        o = s.option(form.Value, 'accessToken', _('Access Token'), _('Access token for API access.'));
        o.rmempty = true;
        o.password = true;
        o.readonly = true;
        o.datatype = 'string';
        o.default = uci.get('ispapp', '@settings[0]', 'accessToken') || 'N/A';

        o = s.option(form.Value, 'refreshToken', _('Refresh Token'), _('Refresh token for API access.'));
        o.rmempty = true;
        o.password = true;
        o.readonly = true;
        o.datatype = 'string';
        o.default = uci.get('ispapp', '@settings[0]', 'refreshToken') || 'N/A';

        o = s.option(form.DummyValue, 'connected', _('Connected'), _('Connection status.'));
        o.rmempty = false;
        o.readonly = true;
        o.rawhtml = true;
        o.cfgvalue = async function() {
            return await callCheckConnection().then(function(connected){
                uci.set('ispapp', '@settings[0]', 'connected', connected && connected.code == 200);
                return connected.code == 200 ? "<span style='color: green;'>Connected</span>" : "<span style='color: red;'>Disconnected</span>";
            })
        };

        o = s.option(form.Button, '_refreshToken', _('Refresh connection'));
        o.inputstyle = 'reload';
        o.onclick = function() {
            callCheckConnection().then(function(response) {
                if (response && response.code === 200) {
                    var tokens = JSON.parse(response.body);
                    uci.set('ispapp', '@settings[0]', 'refreshToken', tokens.refreshToken);
                    uci.set('ispapp', '@settings[0]', 'accessToken', tokens.accessToken);
                    ui.addNotification(null, E('p', _('Refresh token updated successfully.')), 'info');
                } else {
                    var tokens = response && JSON.parse(response.body) || { error: null };
                    ui.addNotification(null, E('p', tokens.error || _('Failed to update refresh token.')), 'error');
                }
            });
        };

        o = s.option(form.Value, 'UpdateInterval', _('Update Interval'), _('Interval for updating data in seconds.'));
        o.rmempty = false;
        o.datatype = 'range(30, 100)';
        o.default = uci.get('ispapp', '@settings[0]', 'updateInterval') || 10;
        o.placeholder = '10';
        o.onchange = function() {
            uci.set('ispapp', '@settings[0]', 'updateInterval', this.formvalue(this.section));
            fs.exec('/etc/init.d/ispapp restart');
        };

        var ops = s.option(form.ListValue, 'IperfServer', 'Iperf Server', 'Select the Iperf server to use for testing.');
        ops.default = uci.get('ispapp', '@settings[0]', 'IperfServer') || 'N/A';
        ops.optional = true;
		ops.rmempty = true;
        
        request.get('/luci-static/resources/iperf').then(data=>data.json()).then(async function(data) {
            var json_values = data;
            if (json_values && typeof json_values === 'object' && json_values.length > 0){
                json_values.forEach(element => {
                    // need add value to option in other way than
                    if((version.v3 && version.v3 == element.V3) || (version.v2 && version.v2 == element.V2)){
                        ops.value(element?.IP_HOST, `${element.COUNTRY}-${element.SITE}(v${element.V3 ? 3 : 2})`);
                    }
                });
            } else {
                ops.value('N/A', 'N/A');
            }
        })
        o = s.option(form.Button, '_apply', _('Device Signup'));
        o.inputstyle = 'apply';
        o.onclick = function() {
            var Domain = uci.get('ispapp', '@settings[0]', 'Domain');
            var ListenerPort = uci.get('ispapp', '@settings[0]', 'ListenerPort');
            var Key = uci.get('ispapp', '@settings[0]', 'Key');
            var updateInterval = uci.get('ispapp', '@settings[0]', 'updateInterval');
            uci.set('ispapp', '@settings[0]', 'Domain', Domain);
            uci.set('ispapp', '@settings[0]', 'ListenerPort', ListenerPort);
            uci.set('ispapp', '@settings[0]', 'Key', Key);
            uci.set('ispapp', '@settings[0]', 'updateInterval', updateInterval);
            uci.apply()
            signup().then(function(response) {
                console.log(response);
                if (response && response.code === 200) {
                    ui.addNotification(null, E('p', _('ISPApp is connected.')), 'info');
                } else {
                    ui.addNotification(null, E('p', _('ISPApp is not connected.')), 'error');
                }
            }).catch((err)=>{
                console.error(err);
                ui.addNotification(null, E('p', _('ISPApp is not connected.')), 'error');
            });
        };

        return m.render();
    }
});
