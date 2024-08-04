##############################################
# OpenWrt Makefile for ispappd program
#
# Most of the variables used here are defined in
# the include directives below. We just need to 
# specify a basic description of the package, 
# where to build our program, where to find 
# the source files, and where to install the 
# compiled program on the router. 
# 
# Be very careful of spacing in this file.
# Indents should be tabs, not spaces, and 
# there should be no trailing whitespace in
# lines that are not commented.
# 
##############################################
# Include OpenWrt build system configuration
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
	DEPENDS:=+libpthread
  MAINTAINER:=Ouazmir karim <abdelkarim@ispapp.co>
endef
# Specify the description for ispappd
define Package/$(PKG_NAME)/description
	ISP App Daemon is a sample package for OpenWRT.
endef
define Build/Configure
endef

define Build/Compile
	$(TARGET_CC) $(TARGET_CFLAGS) \
		-o $(PKG_BUILD_DIR)/iwcap $(PKG_BUILD_DIR)/iwcap.c
endef


define Package/iwcap/install
	$(INSTALL_DIR) $(1)/usr/sbin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/iwcap $(1)/usr/sbin/iwcap
endef

# This line executes the necessary commands to compile ispappd.
$(eval $(call BuildPackage,ispappd))
