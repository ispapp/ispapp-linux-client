#!/bin/bash

# Check if the correct number of arguments are provided
if [ "$#" -ne 3 ]; then
    echo "Usage: $0 <router_ip> <username> <password>"
    exit 1
fi

ROUTER_IP=$1
USERNAME=$2
PASSWORD=$3

# Define the packages to install
PACKAGES="luasec luasocket"

# Function to check SSH connectivity
check_ssh() {
    sshpass -p $PASSWORD ssh -o StrictHostKeyChecking=no $USERNAME@$ROUTER_IP "echo SSH connection successful" > /dev/null 2>&1
}

# Function to check SSH connectivity with legacy SHA
check_ssh_legacy() {
    sshpass -p $PASSWORD ssh -o StrictHostKeyChecking=no -o KexAlgorithms=+diffie-hellman-group1-sha1 -o HostKeyAlgorithms=+ssh-rsa $USERNAME@$ROUTER_IP "echo SSH connection successful" > /dev/null 2>&1
}

# Function to get the architecture of the device
get_arch() {
    sshpass -p $PASSWORD ssh $USERNAME@$ROUTER_IP "uname -m"
}

# Function to check internet connectivity
check_internet() {
    sshpass -p $PASSWORD ssh $USERNAME@$ROUTER_IP "ping -c 1 google.com > /dev/null 2>&1"
}

# Function to install packages using opkg
install_packages_opkg() {
    for PACKAGE in $PACKAGES; do
        sshpass -p $PASSWORD ssh $USERNAME@$ROUTER_IP "opkg update && opkg install $PACKAGE"
    done
}

# Function to download and install packages manually
install_packages_manual() {
    for PACKAGE in $PACKAGES; do
        PACKAGE_URL="https://downloads.openwrt.org/releases/packages-19.07/$ARCH/packages/$PACKAGE.ipk"
        wget $PACKAGE_URL -O /tmp/$PACKAGE.ipk
        sshpass -p $PASSWORD scp /tmp/$PACKAGE.ipk $USERNAME@$ROUTER_IP:/tmp/
        sshpass -p $PASSWORD ssh $USERNAME@$ROUTER_IP "opkg install /tmp/$PACKAGE.ipk && rm /tmp/$PACKAGE.ipk"
        rm /tmp/$PACKAGE.ipk
    done
}

# Check SSH connectivity
if ! check_ssh; then
    if ! check_ssh_legacy; then
        echo "SSH connection failed"
        exit 1
    fi
fi

# Get the architecture of the device
ARCH=$(get_arch)

# Check internet connectivity and install packages
if check_internet; then
    install_packages_opkg
else
    install_packages_manual
fi

# Download and install luci-app-ispapp
ISPAPP_URL="https://github.com/ispapp/ispapp-linux-client/releases/download/latest/luci-app-ispapp-$ARCH.ipk"
wget $ISPAPP_URL -O /tmp/luci-app-ispapp-$ARCH.ipk
sshpass -p $PASSWORD scp /tmp/luci-app-ispapp-$ARCH.ipk $USERNAME@$ROUTER_IP:/tmp/
sshpass -p $PASSWORD ssh $USERNAME@$ROUTER_IP "opkg install /tmp/luci-app-ispapp-$ARCH.ipk && rm /tmp/luci-app-ispapp-$ARCH.ipk"
rm /tmp/luci-app-ispapp-$ARCH.ipk

echo "Installation completed."