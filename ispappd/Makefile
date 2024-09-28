##############################################
# OpenWrt Makefile for ispappd program
##############################################
include $(TOPDIR)/rules.mk

# Name and release number of this package
PKG_NAME:=ispappd
PKG_VERSION:=1.0.1
PKG_RELEASE:=1
PKG_LICENSE:=Apache-2.0

include $(INCLUDE_DIR)/package.mk

# Specify package information for ispappd
define Package/ispappd
	SECTION:=utils
	CATEGORY:=Utilities
	TITLE:=ispappd -- ISPApp Daemon
	DEPENDS:=+libpthread +libcurl +jansson +libc
	MAINTAINER:=Ouazmir karim <abdelkarim@ispapp.co>
endef


# Specify the description for ispappd
define Package/$(PKG_NAME)/description
	The ISPApp Log Collection Daemon is a robust, lightweight service 
	designed specifically for ISPApp-enabled routers. It facilitates 
	the seamless aggregation, management, and transmission of critical 
	system and network logs from individual routers to the ISPApp Cloud 
	infrastructure.
endef

# No configuration step needed for this package
define Build/$(PKG_NAME)/Configure
endef

# Compilation step
define Build/Compile
	$(TARGET_CC) $(TARGET_CFLAGS) -I$(STAGING_DIR)/usr/include \
		-o $(PKG_BUILD_DIR)/ispappd $(PKG_BUILD_DIR)/ispappd.c \
		-lcurl -ljansson -lpthread $(TARGET_LDFLAGS)
endef

# Installation step
define Package/$(PKG_NAME)/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/usr/sbin/ispappd
endef

# This line executes the necessary commands to compile ispappd.
$(eval $(call BuildPackage,ispappd))
