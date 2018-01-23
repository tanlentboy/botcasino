/*
Navicat MySQL Data Transfer

Source Server         : mysql
Source Server Version : 50720
Source Host           : 209.250.228.79:3306
Source Database       : casino_logs

Target Server Type    : MYSQL
Target Server Version : 50720
File Encoding         : 65001

Date: 2018-01-08 17:42:37
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for history
-- ----------------------------
DROP TABLE IF EXISTS `history`;
CREATE TABLE `history` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL COMMENT '用户ID',
  `describe` varchar(255) NOT NULL COMMENT '描述信息',
  `inserted_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '插入日期',
  PRIMARY KEY (`id`),
  KEY `history_idx_user_id` (`user_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for transfer
-- ----------------------------
DROP TABLE IF EXISTS `transfer`;
CREATE TABLE `transfer` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
  `trx_id` varchar(32) NOT NULL COMMENT '操作id',
  `block_num` bigint(20) unsigned NOT NULL COMMENT '区块高度',
  `asset` varchar(64) NOT NULL COMMENT '资产名称',
  `asset_id` varchar(64) NOT NULL COMMENT '资产id',
  `amount` int(10) unsigned NOT NULL COMMENT '资产数量',
  `from_id` varchar(32) NOT NULL COMMENT '来源用户id',
  `to_id` varchar(32) NOT NULL COMMENT '目标用户id',
  `from_name` varchar(64) NOT NULL COMMENT '来源用户名',
  `to_name` varchar(64) NOT NULL COMMENT '目标用户名',
  `memo` varchar(255) DEFAULT NULL COMMENT '备注信息',
  `nonce` varchar(255) DEFAULT NULL COMMENT '随机数',
  `timestamp` varchar(32) NOT NULL COMMENT '时间戳',
  `inserted_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `transfer_idx_block_num` (`block_num`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Table structure for withdraw
-- ----------------------------
DROP TABLE IF EXISTS `withdraw`;
CREATE TABLE `withdraw` (
  `id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '订单id',
  `user_id` int(10) unsigned NOT NULL COMMENT '用户id',
  `to` varchar(64) NOT NULL COMMENT '帐户名',
  `asset_id` varchar(64) NOT NULL COMMENT '资产id',
  `amount` int(10) unsigned NOT NULL COMMENT '资产金额',
  `fee` int(10) unsigned NOT NULL COMMENT '手续费',
  `real` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '真实手续费',
  `status` tinyint(4) unsigned NOT NULL COMMENT '操作状态(0等待提现/1正在提现/3提现成功/4提现失败)',
  `reason` varchar(255) DEFAULT NULL COMMENT '失败原因',
  `inserted_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '插入时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `withdraw_idx_id` (`id`) USING BTREE,
  KEY `withdraw_idx_user_id` (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=100000 DEFAULT CHARSET=utf8mb4;
