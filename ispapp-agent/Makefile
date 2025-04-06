include $(TOPDIR)/rules.mk

PKG_NAME:=ispapp-agent
PKG_VERSION:=1.0.0
PKG_RELEASE:=1

PKG_MAINTAINER:=ISP App Team <support@ispapp.com>
PKG_LICENSE:=MIT
PKG_LICENSE_FILES:=LICENSE

PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/ispapp/ispapp-linux-client.git
PKG_SOURCE_VERSION:=HEAD
PKG_SOURCE_BRANCH:=websocket

PKG_BUILD_DEPENDS:=golang/host
PKG_BUILD_PARALLEL:=1
PKG_BUILD_FLAGS:=no-mips16

GO_PKG:=ispapp-agent
GO_PKG_BUILD_PKG:=$(GO_PKG)/cmd/agent

include $(INCLUDE_DIR)/package.mk
include $(TOPDIR)/feeds/packages/lang/golang/golang-package.mk

define Package/ispapp-agent
	SECTION:=net
	CATEGORY:=Network
	TITLE:=ISP App OpenWrt Agent
	URL:=https://github.com/ispapp/ispapp-linux-client
	DEPENDS:=$(GO_ARCH_DEPENDS)
	AUTOREMOVE:=1
endef

define Package/ispapp-agent/description
	ISP App Agent for OpenWrt devices. Monitors and reports
	system and network information to the ISP App platform.
endef

define Package/ispapp-agent/conffiles
/etc/config/ispapp
endef

define Package/ispapp-agent/install
	$(INSTALL_DIR) $(1)/usr/bin
	$(INSTALL_BIN) $(GO_PKG_BUILD_BIN_DIR)/agent $(1)/usr/bin/ispapp-agent
	
	$(INSTALL_DIR) $(1)/etc/config
	$(INSTALL_CONF) $(PKG_BUILD_DIR)/files/ispapp $(1)/etc/config/ispapp
	
	$(INSTALL_DIR) $(1)/etc/init.d
	$(INSTALL_BIN) $(PKG_BUILD_DIR)/files/ispapp-agent.init $(1)/etc/init.d/ispapp-agent
endef

$(eval $(call BuildPackage,ispapp-agent))
