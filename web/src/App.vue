<template>
    <div id="main">
        <!--  -->
        <h2>Ispapp openwrt Agent</h2>
        <div class="cbi-map-descr">
            IspappAgent connects your device to Ispapp.co platform for easy monitoring and control.<br>
            For more details, please visit <a href="https://ispapp.co" target="_blank">ispapp.co</a>.
        </div>
        <!--  -->
        <div class="cbi-section">
            <h3>Service Status</h3>
            <div class="cbi-section-node">
                <div class="cbi-value cbi-value-last">
                    <label class="cbi-value-title">Enabled</label>
                    <div class="cbi-value-field">
                        <div class="cbi-checkbox">
                            <input name="enabled" type="checkbox" :value="false" v-model="config.enabled">
                        </div>
                    </div>
                </div>
                <div class="cbi-value cbi-value-last">
                    <label class="cbi-value-title">Service Status</label>
                    <div class="cbi-value-field">
                        <a v-if="loading">Loading...</a>
                        <template v-else>
                            <a v-if="status?.BackendState == 'running'" style="color:green"> {{ status?.BackendState }}</a>
                            <a v-else style="color:red">Not running</a>
                        </template>
                    </div>
                </div>
                <div class="cbi-value cbi-value-last" v-if="status.BackendState == 'running'">
                    <label class="cbi-value-title">Health stats:</label>
                    <div v-for="hstt in status.Health" class="cbi-value-field">
                            <template>
                                ({{ hstt }})
                            </template>
                    </div>
                </div>
            </div>
        </div>
        <!--  -->
        <div class="cbi-section">
            <h3>Global Settings</h3>
            <div class="cbi-section-node">
                <!-- Device Name -->
                <div class="cbi-value">
                    <label class="cbi-value-title">Device Name</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="hostname" v-model.trim="config.login"
                                placeholder="e.g., 00:00:00:00:00:00">
                        </div>
                        <div class="cbi-value-description">
                            Leave empty to use the connection interface MAC.
                        </div>
                    </div>
                </div>

                <!-- Top Listener Port -->
                <div class="cbi-value">
                    <label class="cbi-value-title">Top Listener Port</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="topListenerPort"
                                v-model.trim="config.topListenerPort" placeholder="e.g., 8550">
                        </div>
                        <div class="cbi-value-description">
                            Specify the top listener port.
                        </div>
                    </div>
                </div>

                <!-- Top SMTP Port -->
                <div class="cbi-value">
                    <label class="cbi-value-title">Top SMTP Port</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="topSmtpPort" v-model.trim="config.topSmtpPort"
                                placeholder="e.g., 8465">
                        </div>
                        <div class="cbi-value-description">
                            Specify the top SMTP port.
                        </div>
                    </div>
                </div>

                <!-- Top Key -->
                <div class="cbi-value">
                    <label class="cbi-value-title">Top Key</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="topKey" v-model.trim="config.topKey"
                                placeholder="e.g., (leave empty)">
                        </div>
                        <div class="cbi-value-description">
                            Description for the Top Key setting.
                        </div>
                    </div>
                </div>

                <!-- IP Band Switch Test Server -->
                <div class="cbi-value">
                    <label class="cbi-value-title">IP Band Switch Test Server</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="ipbandswtestserver"
                                v-model.trim="config.ipbandswtestserver" placeholder="e.g., 3.239.254.95">
                        </div>
                        <div class="cbi-value-description">
                            Specify the IP Band Switch Test Server.
                        </div>
                    </div>
                </div>

                <!-- BT User -->
                <div class="cbi-value">
                    <label class="cbi-value-title">BT User</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="btuser" v-model.trim="config.btuser"
                                placeholder="e.g., btest">
                        </div>
                        <div class="cbi-value-description">
                            Specify the BT User.
                        </div>
                    </div>
                </div>

                <!-- BT Password -->
                <div class="cbi-value">
                    <label class="cbi-value-title">BT Password</label>
                    <div class="cbi-value-field">
                        <div>
                            <input type="text" class="cbi-input-text" name="btpwd" v-model.trim="config.btpwd"
                                placeholder="e.g., 0XSYIGkRlP6MUQJMZMdrogi2">
                        </div>
                        <div class="cbi-value-description">
                            Specify the BT Password.
                        </div>
                    </div>
                </div>

            </div>
        </div>
        <span class="cbi-page-actions control-group">
            <button class="btn cbi-button cbi-button-apply" @click="onSubmit" :disabled="disabled">Save and Apply</button>
        </span>
    </div>
</template>
<script setup lang="ts">
import { onMounted, ref } from 'vue';
const BASEURL = "/cgi-bin/luci/admin/services/ispapp"
const loading = ref(true);
const disabled = ref(false);
const config = ref<ResponseConfig>({});
const status = ref<ResponseStatus>({});

const request = (input: string, init?: RequestInit | undefined) => {
    const uri = `${BASEURL}${input}`;
    return fetch(uri, init);
};

const getStatus = async () => {
    try {
        const resp = await request("/status", {
            method: "GET",
        });
        const res = (await resp.json()) as ResponseStatus;
        if (res) {
            status.value = {...status.value, ...res};
        }
    } catch (error) {
        status.value = {...status.value, BackendState:"Error getting status responce from agent!"}
    }
};

const getConfig = async () => {
    try {
        const resp = await request("/config", {
            method: "GET",
        });
        const res = (await resp.json()) as ResponseConfig;
        if (res) {
            config.value = res;
        }
    } catch (error) {
        console.log(error);
    }
};

const getData = async () => {
    try {
        await Promise.all([getConfig(), getStatus()]);
    } catch (error) {
    } finally {
        loading.value = false;
    }
};

const getInterval = () => {
    setInterval(() => {
        getStatus();
    }, 5000);
};

onMounted(() => {
    getInterval();
    getData();
});

const onSubmit = async () => {
    try {
        const resp = await request("/config", {
            method: "POST",
            headers: {
                'Content-Type': 'application/json;charset=utf-8',
            },
            body: JSON.stringify(config.value),
        });
        if (resp) {
        }
    } catch (error) {
        console.log(error);
    } finally {
        location.reload();
    }
};
</script>
<style lang="scss" scoped></style>