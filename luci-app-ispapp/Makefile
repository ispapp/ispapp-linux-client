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
PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_VERSION:=1.0.0
PKG_RELEASE:=2024092
PKG_LICENSE:=Apache-2.0
include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk
define Package/luci-app-ispapp
	SECTION:=net
	CATEGORY:=Network
	TITLE:=LuCI Support of Ispapp.co Agent
	URL:=https://ispapp.co
	DEPENDS:=$(GO_ARCH_DEPENDS)
endef
define Build/Prepare
	$(CP) ./files/* $(PKG_BUILD_DIR)
endef
define Package/$(PKG_NAME)/description
LuCI Support of Ispapp.co Agent
endef
define Package/luci-app-ispapp/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/luci-app-ispapp $(1)/usr/sbin/ispappd
endef
include $(TOPDIR)/feeds/luci/luci.mk

# call BuildPackage - OpenWrt buildroot signature
