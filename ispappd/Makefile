include $(TOPDIR)/rules.mk

PKG_NAME:=ispappd
PKG_VERSION:=1.0
PKG_RELEASE:=1

PKG_BUILD_DIR:=$(BUILD_DIR)/$(PKG_NAME)-$(PKG_VERSION)

include $(INCLUDE_DIR)/package.mk

define Package/ispappd
  SECTION:=utils
  CATEGORY:=Utilities
  TITLE:=ISP App Daemon
  DEPENDS:=+libpthread
endef

define Package/ispappd/description
  ISP App Daemon is a sample package for OpenWRT.
endef

define Build/Prepare
  mkdir -p $(PKG_BUILD_DIR)
  $(CP) ./src/* $(PKG_BUILD_DIR)/
endef

define Build/Compile
  $(MAKE) -C $(PKG_BUILD_DIR) $(TARGET_CONFIGURE_OPTS) CFLAGS="$(TARGET_CFLAGS)"
endef

define Package/ispappd/install
  $(INSTALL_DIR) $(1)/usr/bin
  $(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/usr/bin/
endef

$(eval $(call BuildPackage,ispappd))
