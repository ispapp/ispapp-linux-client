'use strict';
'require form';
'require rpc';
'require view';

// RPC call to save the configuration
var save_config = rpc.declare({
    object: 'luci.ispapp',
    method: 'write_ispapp_config',
    expect: { 'response': 'json' }
});

// RPC call to get the current configuration
var load_config = rpc.declare({
    object: 'luci.ispapp',
    method: 'read_ispapp_config',
    expect: { 'response': 'json' }
});

// Add some custom CSS for styling
function addCustomStyles() {
    var style = document.createElement('style');
    style.innerHTML = `
        .form-section {
            background-color: #f3e5f5; /* Light purple background */
            border-radius: 8px;
            padding: 15px;
            margin-bottom: 15px;
            border: 1px solid #e1bee7; /* Light magenta border */
        }
        .form-section h2 {
            color: #ad1457; /* Magenta color for section titles */
        }
        .form-section .control-label {
            color: #6a1b9a; /* Darker purple for labels */
        }
        .form-section input[type="text"],
        .form-section input[type="password"],
        .form-section select {
            border: 1px solid #d1c4e9; /* Light purple input borders */
            border-radius: 4px;
            padding: 8px;
        }
        .form-section input[type="checkbox"] {
            accent-color: #ad1457; /* Magenta for checkboxes */
        }
        .form-section .btn-save-apply {
            background-color: #ad1457; /* Magenta button background */
            color: white;
            border: none;
            padding: 10px 20px;
            border-radius: 5px;
            cursor: pointer;
        }
        .form-section .btn-save-apply:hover {
            background-color: #880e4f; /* Darker magenta on hover */
        }
    `;
    document.head.appendChild(style);
}

// Call the function to add custom styles
addCustomStyles();

return view.extend({
    load: function() {
        return Promise.all([
            load_config()  // Load the current configuration
        ]);
    },

    render: function(data) {
        var config = data[0] || {};
        var m, s, o;

        m = new form.Map('ispapp', _('ISPApp Configuration'),
            _('Configure ISPApp settings.'));

        // Section for configuration settings
        s = m.section(form.TypedSection, 'settings', _('Settings'));
        s.anonymous = true;
        s.className = 'form-section'; // Apply custom CSS class

        // Options for each configuration setting
        o = s.option(form.Flag, 'enabled', _('Enabled'), _('Enable ISPApp.'));
        o.default = config.enabled || '0';
        o.rmempty = false;

        o = s.option(form.Value, 'login', _('Login'), _('MAC address for login.'));
        o.default = config.login || '';
        o.rmempty = false;

        o = s.option(form.Value, 'topDomain', _('Top Domain'), _('Domain name for top domain.'));
        o.default = config.topDomain || '';
        o.rmempty = false;

        o = s.option(form.Value, 'topListenerPort', _('Top Listener Port'), _('Port for listening.'));
        o.default = config.topListenerPort || '';
        o.rmempty = false;

        o = s.option(form.Value, 'topSmtpPort', _('Top SMTP Port'), _('SMTP port.'));
        o.default = config.topSmtpPort || '';
        o.rmempty = false;

        o = s.option(form.Value, 'topKey', _('Top Key'), _('Key for authentication.'));
        o.default = config.topKey || '';
        o.rmempty = true;

        o = s.option(form.Value, 'ipbandswtestserver', _('IP Bandwidth Test Server'), _('IP address for bandwidth test server.'));
        o.default = config.ipbandswtestserver || '';
        o.rmempty = false;

        o = s.option(form.Value, 'btuser', _('BT User'), _('Username for BT.'));
        o.default = config.btuser || '';
        o.rmempty = false;

        o = s.option(form.Value, 'btpwd', _('BT Password'), _('Password for BT.'));
        o.default = config.btpwd || '';
        o.rmempty = false;

        // Handle save and apply actions
        var buttons = E('div', { 'class': 'form-section' }, [
            E('button', { 'class': 'btn-save-apply', 'click': this.handleSaveApply.bind(this) }, _('Save and Apply'))
        ]);

        return m.render().add(buttons);
    },

    handleSaveApply: function(ev) {
        // Collect form data
        var formData = form.Map.getValues('ispapp');

        // Submit the data via RPC
        return save_config({ config: formData })
            .then(function(response) {
                if (response.error) {
                    ui.addNotification(null, E('p', _('Unable to save changes: %s').format(response.error)));
                } else {
                    ui.addNotification(null, E('p', _('Configuration changes have been saved.')));
                }
            }).catch(function(e) {
                ui.addNotification(null, E('p', _('Unable to save changes: %s').format(e.message)));
            });
    },

    handleSave: null,
    handleReset: null
});
