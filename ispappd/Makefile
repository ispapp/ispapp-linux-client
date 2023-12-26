# Copyright (C) 2023 ispapp.co

include $(TOPDIR)/rules.mk
INSTALL = install
PREFIX  = /usr/bin
PKG_NAME:=ispappd
CC := gcc
CFLAGS := -Wall -Wextra
include $(INCLUDE_DIR)/package.mk
define Package/ispappd
  SECTION:=net
  CATEGORY:=Network
  TITLE:=Ispapp daemon
  DEPENDS:=+gcc
endef
define Build/Prepare
  # Your build preparation commands, if any
endef

define Build/Compile
  $(MAKE) -C $(PKG_BUILD_DIR) $(TARGET_CONFIGURE_OPTS) all
endef
define Package/ispappd/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/usr/sbin
endef

$(eval $(call BuildPackage,ispappd))