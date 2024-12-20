#!/bin/bash

export RSYNC_PASSWORD='admin'

while true; do
    rsync -avz --progress ./luci-app-ispapp/root/usr/libexec/rpcd/ispapp root@192.168.137.22:/usr/libexec/rpcd/ispapp
    rsync -avz --progress ./luci-app-ispapp/luasrc/model/cbi/ispapp/overview.lua root@192.168.137.22:/usr/lib/lua/luci/model/cbi/ispapp/overview.lua
    rsync -avz --progress ./luci-app-ispapp/luasrc/model/cbi/ispapp/settings.lua root@192.168.137.22:/usr/lib/lua/luci/model/cbi/ispapp/settings.lua
    rsync -avz --progress ./luci-app-ispapp/root/bin/ispappd root@192.168.137.22:/bin/ispappd
    rsync -avz --progress ./luci-app-ispapp/root/usr/lib/lua/assoclist.awk root@192.168.137.22:/usr/lib/lua/assoclist.awk
    sleep 10  # Wait for 5 seconds before syncing again
done