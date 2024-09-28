#
# Copyright (C) 2008-2014 The LuCI Team <luci@lists.subsignal.org>
#
# This is free software, licensed under the Apache License, Version 2.0 .
#

include $(TOPDIR)/rules.mk

LUCI_TITLE:=Ispapp Agent
LUCI_DEPENDS:=+ispappd +luci-base
LUCI_PKGARCH:=all

PKG_MAINTAINER:=Karim O. <abdelkarim@ispapp.co>
PKG_NAME:=luci-app-ispapp
PKG_VERSION:=1.0.0
PKG_RELEASE:=20231230
PKG_LICENSE:=Apache-2.0
PKG_MAINTAINER:=Ouazmir Abdelkarim <abdelkarim@ispapp.co>

define Package/luci-app-ispapp
	SECTION:=net
	CATEGORY:=Network
	TITLE:=LuCI Support of Ispapp.co Agent
	URL:=https://ispapp.co
	DEPENDS:=+ispappd +luci-lua-runtime
endef
define Package/luci-app-ispapp/conffiles
/etc/config/ispapp
endef
define Package/$(PKG_NAME)/description
LuCI Support of Ispapp.co Agent
endef

define Build/Prepare
endef

define Build/Compile
endef

define Package/$(PKG_NAME)/install
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) ./root/etc/config/* $(1)/etc/config
	
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) ./root/etc/init.d/* $(1)/etc/init.d

	$(INSTALL_DIR) $(1)/etc/uci-defaults
	$(INSTALL_CONF) ./root/etc/uci-defaults/* $(1)/etc/uci-defaults/
	
	$(INSTALL_DIR) $(1)/usr/lib/lua/luci
	cp -pR ./luasrc/* $(1)/usr/lib/lua/luci/
	
	$(INSTALL_DIR) $(1)/usr/share/luci/menu.d/
	cp -pR ./root/usr/share/luci/menu.d/* $(1)/usr/share/luci/menu.d/
endef


define Package/$(PKG_NAME)/postinst
#!/bin/sh
[ -n "${IPKG_INSTROOT}" ] || {
	( . /etc/uci-defaults/50_luci-ispapp ) && rm -f /etc/uci-defaults/50_luci-ispapp
	exit 0
}
endef

define Package/$(PKG_NAME)/prerm
#!/bin/sh
if [ -z "$${IPKG_INSTROOT}" ]; then
     /etc/init.d/ispapp disable
     /etc/init.d/ispapp stop
fi
exit 0
endef

include $(TOPDIR)/feeds/luci/luci.mk

# call BuildPackage - OpenWrt buildroot signature
