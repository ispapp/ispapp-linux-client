#!/usr/bin/bash

set -ef

GROUP=

group() {
	endgroup
	echo "::group::  $1"
	GROUP=1
}

endgroup() {
	if [ -n "$GROUP" ]; then
		echo "::endgroup::"
	fi
	GROUP=
}

trap 'endgroup' ERR

# Variables
SDK_IMAGE="openwrt/sdk:x86-64-23.05.5"
ROOTFS_IMAGE="openwrt/rootfs:x86-64-23.05.5"
PACKAGE_DIR="ispapp"
BUILD_DIR="/builder"
HTTP_PORT="80"
FORWARD_PORT="8909"
FEEDNAME="ispapp"
BUILD_LOG="${BUILD_LOG:-1}"

# Docker container names
SDK_CONTAINER="sdk_build"
ROOTFS_CONTAINER="openwrt_instance"

# Prepare environment for feeds and packages
echo "Building luci-app-ispapp and ispappd packages in SDK..."

docker run --rm -v "$(pwd)"/bin/:"$BUILD_DIR/bin" -v "$(pwd)/$PACKAGE_DIR/":"$BUILD_DIR/$PACKAGE_DIR" $SDK_IMAGE bash -c "
    set -e
    echo 'src-link $FEEDNAME "$BUILD_DIR/$PACKAGE_DIR/"' >> feeds.conf
    echo 'src-git packages https://git.openwrt.org/feed/packages.git' >> feeds.conf
    echo 'src-git luci https://git.openwrt.org/project/luci.git' >> feeds.conf
    echo 'src-git routing https://git.openwrt.org/feed/routing.git' >> feeds.conf
    echo 'src-git telephony https://git.openwrt.org/feed/telephony.git' >> feeds.conf
    # Custom Feeds
    LC_ALL=C
    ./scripts/feeds update
    ./scripts/feeds install -p $FEEDNAME -f
    make defconfig
    make package/applications/luci-app-ispapp/compile V=s -j \$(nproc)
    make package/utils/ispappd/compile V=s -j \$(nproc)
    mkdir -p ./packages
    cp bin/packages/x86_64/*ispapp*.ipk ./packages/
"

# Reinstall the packages in the OpenWrt root filesystem
echo "Installing packages in rootfs..."
docker run --rm -v "$PACKAGE_DIR":/packages $ROOTFS_IMAGE opkg install /packages/*.ipk --force-reinstall

# Run the OpenWrt rootfs with HTTP forwarding
echo "Starting OpenWrt with HTTP port forwarding..."
docker run -d -p $FORWARD_PORT:$HTTP_PORT --name $ROOTFS_CONTAINER $ROOTFS_IMAGE /sbin/init

# Output success message
if [ $? -eq 0 ]; then
    echo "Packages installed and OpenWrt is running. HTTP forwarded from port $HTTP_PORT to $FORWARD_PORT."
    echo "To stop OpenWrt, use: docker stop $ROOTFS_CONTAINER"
else
    echo "Failed to start OpenWrt."
    exit 1
fi
