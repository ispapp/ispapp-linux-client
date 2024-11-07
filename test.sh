curl  -X POST \
  'https://known-ample-gibbon.ngrok-free.app:443/initconfig?login=2d00ac77-bccc-430f-97a3-070264ccd132&key=etDvf0bISZ0ggkqd' \
  --header 'Accept: */*' \
  --header 'User-Agent: Thunder Client (https://www.thunderclient.com)' \
  --header 'Content-Type: application/json' \
  --data-raw '{
  "outsideIp": "196.112.59.233",
  "webshellSupport": true,
  "fw": "FIRMWARE-2167-202407161750",
  "lastConfigRequest": 1721296274,
  "osVersion": "19.07-SNAPSHOT",
  "hardwareCpuInfo": "ARMv7 Processor rev 4 (v7l), 4 cores",
  "hostname": "OpenWrt",
  "hardwareModelNumber": "Qualcomm Technologies, Inc. IPQ5332/AP-MI01.12",
  "hardwareMake": "qcom,ipq5332-ap-mi01.12",
  "bandwidthTestSupport": true,
  "clientInfo": "OpenWrt-19.07-SNAPSHOT",
  "os": "19.07-snapshot",
  "wirelessConfigured": [
    {
      ".id": "*1",
      "disabled": false,
      "mac-address": "c4:4b:d1:c0:06:88",
      "interface-type": "qcawificfg80211",
      "if": "wifinet0",
      "technology": "uci",
      "master-interface": "qcawificfg80211",
      "running": true,
      "name": "wifinet0",
      "ssid": "Longshot MLOv2 2.4Ghz",
      "key": "langshot",
      "band": "2.4ghz-ofdm",
      "security-profile": "*1",
      "hide-ssid": false
    },
    {
      ".id": "*3",
      "disabled": false,
      "mac-address": "c4:4b:d1:c0:06:89",
      "interface-type": "qcawificfg80211",
      "if": "wifinet1",
      "technology": "uci",
      "master-interface": "qcawificfg80211",
      "running": true,
      "name": "wifinet1",
      "ssid": "Longshot MLOv2",
      "key": "langshot",
      "band": "2.4/5/6ghz-eht-ofdma",
      "security-profile": "*3",
      "hide-ssid": false
    },
    {
      ".id": "*5",
      "disabled": false,
      "mac-address": "c4:4b:d1:c0:06:8a",
      "interface-type": "qcawificfg80211",
      "if": "wifinet2",
      "technology": "uci",
      "master-interface": "qcawificfg80211",
      "running": true,
      "name": "wifinet2",
      "ssid": "Longshot MLOv2",
      "key": "langshot",
      "band": "2.4/5/6ghz-eht-ofdma",
      "security-profile": "*5",
      "hide-ssid": false
    }
  ],
  "interfaces": [
    {
      "foundDescriptor": "ubus parsed",
      "defaultIf": "lo",
      "sentDrops": 0,
      "if": "lo",
      "sentErrors": 0,
      "recBytes": 3804,
      "sentPackets": 49,
      "recDrops": 0,
      "carrierChanges": 0,
      "recPackets": 49,
      "recErrors": 0,
      "mac": "00:00:00:00:00:00",
      "sentBytes": 3804
    },
    {
      "foundDescriptor": "ubus parsed",
      "defaultIf": "br-lan",
      "sentDrops": 0,
      "if": "eth1",
      "sentErrors": 0,
      "recBytes": 1119144,
      "sentPackets": 7712,
      "recDrops": 11,
      "carrierChanges": 0,
      "recPackets": 7702,
      "recErrors": 0,
      "mac": "c4:4b:d1:c0:06:87",
      "sentBytes": 2387368
    },
    {
      "foundDescriptor": "ubus parsed",
      "defaultIf": "br-lan",
      "sentDrops": 0,
      "if": "eth0",
      "sentErrors": 0,
      "recBytes": 0,
      "sentPackets": 0,
      "recDrops": 0,
      "carrierChanges": 0,
      "recPackets": 0,
      "recErrors": 0,
      "mac": "c4:4b:d1:c0:06:86",
      "sentBytes": 0
    },
    {
      "foundDescriptor": "ubus parsed",
      "defaultIf": "br-lan",
      "sentDrops": 0,
      "if": "br-lan",
      "sentErrors": 0,
      "recBytes": 982058,
      "sentPackets": 7336,
      "recDrops": 0,
      "carrierChanges": 0,
      "recPackets": 7684,
      "recErrors": 0,
      "mac": "c4:4b:d1:c0:06:86",
      "sentBytes": 2334218
    }
  ],
  "wirelessSupport": true,
  "hardwareSerialNumber": "0000000000000000",
  "osBuildDate": 1721295876,
  "sequenceNumber": 1,
  "Lng": -7.613300,
  "uptime": 420,
  "security-profiles": [
    {
      ".id": "*1",
      "wpa-pre-shared-key": "langshot",
      "technology": "wireless",
      "name": "wifinet0",
      "default": false,
      "authentication-types": [
        "sae"
      ],
      "mode": "sta",
      "wpa2-pre-shared-key": "langshot",
      "wpa3-pre-shared-key": "langshot"
    },
    {
      ".id": "*3",
      "wpa-pre-shared-key": "langshot",
      "technology": "wireless",
      "name": "wifinet1",
      "default": false,
      "authentication-types": [
        "sae"
      ],
      "mode": "sta",
      "wpa2-pre-shared-key": "langshot",
      "wpa3-pre-shared-key": "langshot"
    },
    {
      ".id": "*5",
      "wpa-pre-shared-key": "langshot",
      "technology": "wireless",
      "name": "wifinet2",
      "default": false,
      "authentication-types": [
        "sae"
      ],
      "mode": "sta",
      "wpa2-pre-shared-key": "langshot",
      "wpa3-pre-shared-key": "langshot"
    }
  ],
  "firmwareUpgradeSupport": true,
  "lat": 33.579200,
  "usingWebSocket": false,
  "hardwareModel": "Qualcomm Technologies, Inc. IPQ5332/AP-MI01.12"
}'