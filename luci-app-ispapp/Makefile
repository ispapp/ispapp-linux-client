#
# Copyright (C) 2008-2014 The LuCI Team <luci@lists.subsignal.org>
#
# This is free software, licensed under the Apache License, Version 2.0 .
#

include $(TOPDIR)/rules.mk

LUCI_TITLE:=Ispapp Agent
LUCI_DEPENDS:=+luasocket +luasec +rpcd-mod-iwinfo +rpcd
LUCI_PKGARCH:=all

PKG_MAINTAINER:=kmoz000 <abdelkarim@ispapp.co>
PKG_NAME:=luci-app-ispapp
PKG_VERSION:=1.0.0
PKG_RELEASE:=2024092
PKG_LICENSE:=Apache-2.0

define Package/luci-app-ispapp
	SECTION:=net
	CATEGORY:=Network
	TITLE:=LuCI Support of Ispapp.co Agent
	URL:=https://ispapp.co
endef
define Package/$(PKG_NAME)/conffiles
/etc/config/ispapp
endef
define Package/$(PKG_NAME)/description
LuCI Support of Ispapp.co Agent
endef

define Package/$(PKG_NAME)/postinst
#!/bin/sh
[ -n "${IPKG_INSTROOT}" ] || {
	if [ -f /etc/uci-defaults/50_luci-ispapp ]; then
		( . /etc/uci-defaults/50_luci-ispapp ) && rm -f /etc/uci-defaults/50_luci-ispapp
	fi
	exit 0
}
endef

include $(TOPDIR)/feeds/luci/luci.mk

# call BuildPackage - OpenWrt buildroot signature
