#
# Copyright (C) 2023 OpenWrt.org
#
# This is free software, licensed under the GNU General Public License v2.
# See /LICENSE for more information.
#

include $(TOPDIR)/rules.mk

PKG_NAME:=ath12k-firmware
PKG_SOURCE_DATE:=2024-03-14
PKG_SOURCE_VERSION:=6ed267e22501be2cdd63ea564d407fb64eb879dd
PKG_MIRROR_HASH:=d93edc651b641a97ec4f5ed329e15ac110928f863b4afbaae7fcb25953d04544
PKG_RELEASE:=2
DL_SUBDIR:=$(PKG_NAME)-$(PKG_VERSION)
PKG_SOURCE_PROTO:=git
PKG_SOURCE_URL:=https://github.com/quic/upstream-wifi-fw.git

PKG_LICENSE_FILES:=LICENSE.qca_firmware

PKG_MAINTAINER:=Robert Marko <robimarko@gmail.com>

include $(INCLUDE_DIR)/package.mk

RSTRIP:=:
STRIP:=:


define Package/ath12k-firmware-default
  SECTION:=firmware
  CATEGORY:=Firmware
  URL:=$(PKG_SOURCE_URL)
  DEPENDS:=
endef

define Package/ath12k-firmware-ipq5322
$(Package/ath12k-firmware-default)
  TITLE:=IPQ5322 ath12k firmware
endef

define Package/ath12k-firmware-ipq5322_qcn6432
$(Package/ath12k-firmware-default)
  TITLE:=IPQ5322_QCN6432 ath12k firmware
endef

define Package/ath12k-firmware-ipq5322_qcn6432_qcn6432
$(Package/ath12k-firmware-default)
  TITLE:=IPQ5322_QCN6432_QCN6432 ath12k firmware
endef

define Package/ath12k-firmware-qcn9274_hw1
$(Package/ath12k-firmware-default)
  TITLE:=QCN9274 hw1.0 ath12k firmware
endef

define Package/ath12k-firmware-qcn9274_hw2
$(Package/ath12k-firmware-default)
  TITLE:=QCN9274 hw2.0 ath12k firmware
endef

define Package/ath12k-firmware-qcn9274_hw2_split_phy
$(Package/ath12k-firmware-default)
  TITLE:=QCN9274 hw2.0_split_phy ath12k firmware
endef

define Build/Compile

endef

define Package/ath12k-firmware-ipq5322/install
	$(INSTALL_DIR) $(1)/lib/firmware/IPQ5322/hw1.0
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/IPQ5322/hw1.0/* \
		$(1)/lib/firmware/IPQ5322/hw1.0/
endef

define Package/ath12k-firmware-ipq5322_qcn6432/install
	$(INSTALL_DIR) $(1)/lib/firmware/IPQ5322_QCN6432/hw1.0
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/IPQ5322_QCN6432/hw1.0/* \
		$(1)/lib/firmware/IPQ5322_QCN6432/hw1.0/
endef

define Package/ath12k-firmware-ipq5322_qcn6432_qcn6432/install
	$(INSTALL_DIR) $(1)/lib/firmware/IPQ5322_QCN6432_QCN6432/hw1.0
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/IPQ5322_QCN6432_QCN6432/hw1.0/* \
		$(1)/lib/firmware/IPQ5322_QCN6432_QCN6432/hw1.0/
endef

define Package/ath12k-firmware-qcn9274_hw1/install
	$(INSTALL_DIR) $(1)/lib/firmware/ath12k/QCN9274/hw1.0
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/QCN9274/hw1.0/* \
		$(1)/lib/firmware/ath12k/QCN9274/hw1.0/
endef

define Package/ath12k-firmware-qcn9274_hw2/install
	$(INSTALL_DIR) $(1)/lib/firmware/ath12k/QCN9274/hw2.0
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/QCN9274/hw2.0/* \
		$(1)/lib/firmware/ath12k/QCN9274/hw2.0/
endef

define Package/ath12k-firmware-qcn9274_hw2_split_phy/install
	$(INSTALL_DIR) $(1)/lib/firmware/ath12k/QCN9274/hw2.0_split_phy
	$(INSTALL_DATA) \
		$(PKG_BUILD_DIR)/ath12k-firmware/QCN9274/hw2.0_split_phy/* \
		$(1)/lib/firmware/ath12k/QCN9274/hw2.0_split_phy/
endef

$(eval $(call BuildPackage,ath12k-firmware-ipq5322))
$(eval $(call BuildPackage,ath12k-firmware-ipq5322_qcn6432))
$(eval $(call BuildPackage,ath12k-firmware-ipq5322_qcn6432_qcn6432))
$(eval $(call BuildPackage,ath12k-firmware-qcn9274_hw1))
$(eval $(call BuildPackage,ath12k-firmware-qcn9274_hw2))
$(eval $(call BuildPackage,ath12k-firmware-qcn9274_hw2_split_phy))