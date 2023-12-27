# Copyright (C) 2023 ispapp.co

include $(TOPDIR)/rules.mk
INSTALL = install
PREFIX  = /usr/bin
PKG_NAME:=ispappd
PKG_VERSION:=20231230
PKG_RELEASE:=1.0.0
include $(INCLUDE_DIR)/package.mk
define Package/$(PKG_NAME)/description
LuCI Support of Ispapp.co Agent daemon
endef
define Package/$(PKG_NAME)
  SECTION:=net
  CATEGORY:=Network
  TITLE:=Ispapp daemon
  # DEPENDS:=+gcc
endef
define Build/Prepare
  # Your build preparation commands, if any
endef

define Build/Compile
  $(MAKE) -C $(PKG_BUILD_DIR) all
endef
define Package/ispappd/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/usr/sbin
endef

$(eval $(call BuildPackage,ispappd))