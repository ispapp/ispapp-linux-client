#!/bin/bash

export RSYNC_PASSWORD=''

while true; do
    rsync -avz -oHostKeyAlgorithms=+ssh-rsa --progress -e ./luci-app-ispapp/root/usr/libexec/rpcd/ispapp root@192.168.100.1:/usr/libexec/rpcd/ispapp
    sleep 5  # Wait for 5 seconds before syncing again
done
