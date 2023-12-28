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
# target: ispappd
# ispappd:
# 	gcc -g -lpthread  -o ispappd ./src/*.c
# This specifies the directory where we're going to build the program.  
# The root build directory, $(BUILD_DIR), is by default the build_mipsel 
# directory in your OpenWrt SDK directory

# Include necessary build system package.mk file
include $(INCLUDE_DIR)/package.mk

# Specify package information for ispappd
define Package/ispappd
	SECTION:=utils
	CATEGORY:=Utilities
	TITLE:= ISPApp Daemon
  DEPENDS:=+libpthread +librt
  MAINTAINER:=Karim O. <abdelkarim@ispapp.co>
endef

# Specify the description for ispappd
define Package/ispappd/description
	ISP App Daemon is a sample package for OpenWRT.
endef
PKG_BUILD_DIR := $(BUILD_DIR)/$(PKG_NAME)
# Specify what needs to be done to prepare for building the package.
define Build/Prepare
	mkdir -p $(PKG_BUILD_DIR)
	$(CP) ./src/* $(PKG_BUILD_DIR)/
endef

define Build/Compile
endef

# Specify where and how to install ispappd
define Package/ispappd/install
	$(INSTALL_DIR) $(1)/bin
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/ispappd $(1)/bin/
endef

# This line executes the necessary commands to compile ispappd.
$(eval $(call BuildPackage,ispappd))
