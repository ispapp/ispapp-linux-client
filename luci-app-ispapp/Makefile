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

GO_PKG:=github.com/ispapp/ispappd
GO_PKG_BUILD_DIR:=$(PKG_BUILD_DIR)
GO_PKG_INSTALL_ALL:=1
GO_PKG_LDFLAGS:=-s -w
GO_PKG_LDFLAGS+= -linkmode external -extldflags '-static'

define Package/luci-app-ispapp
	SECTION:=net
	CATEGORY:=Network
	TITLE:=LuCI Support of Ispapp.co Agent
	URL:=https://ispapp.co
endef

define Build/Prepare
	$(INSTALL_DIR) $(PKG_BUILD_DIR)
	$(CP) ./files/* $(PKG_BUILD_DIR)/
endef

define Build/Compile
	cd $(PKG_BUILD_DIR) && \
	$(GO_VARS) \
	CGO_ENABLED=0 \
	go build -v -o ispappd \
		-ldflags "$(GO_PKG_LDFLAGS)" \
		.
endef

define Package/$(PKG_NAME)/description
LuCI Support of Ispapp.co Agent
endef

define Package/$(PKG_NAME)/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/usr/sbin/
endef

include $(TOPDIR)/feeds/luci/luci.mk

# call BuildPackage - OpenWrt buildroot signature
