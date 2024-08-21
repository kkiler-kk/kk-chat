/*
 Navicat Premium Data Transfer

 Source Server         : localhost
 Source Server Type    : MySQL
 Source Server Version : 90000 (9.0.0)
 Source Host           : localhost:3306
 Source Schema         : chat

 Target Server Type    : MySQL
 Target Server Version : 90000 (9.0.0)
 File Encoding         : 65001

 Date: 15/08/2024 21:08:52
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for group_basic
-- ----------------------------
DROP TABLE IF EXISTS `group_basic`;
CREATE TABLE `group_basic` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `owner_id` bigint DEFAULT NULL COMMENT '群主',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '头像',
  `type` int DEFAULT NULL COMMENT '类型',
  `identity` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT 'id 唯一标识(可能由系统分配或用户自定义)',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '群昵称',
  `desc` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '描述',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `identity_unique` (`identity`) USING BTREE COMMENT 'id唯一索引'
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='群聊分组';

-- ----------------------------
-- Records of group_basic
-- ----------------------------
BEGIN;
INSERT INTO `group_basic` (`id`, `owner_id`, `avatar`, `type`, `identity`, `name`, `desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, 'http://127.0.0.1:9345/resource/public/file/default/member.jpg', 0, '123', 'tst', '', '2024-04-08 15:53:22', '2024-04-08 15:53:22', NULL);
INSERT INTO `group_basic` (`id`, `owner_id`, `avatar`, `type`, `identity`, `name`, `desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 1, 'http://127.0.0.1:9345/resource/public/file/default/member.jpg', 0, '231', 'test', '', '2024-04-08 15:55:10', '2024-04-08 15:55:10', NULL);
INSERT INTO `group_basic` (`id`, `owner_id`, `avatar`, `type`, `identity`, `name`, `desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 1, 'http://127.0.0.1:9345/resource/public/file/default/member.jpg', 0, '23', 'test', '', '2024-04-08 15:56:29', '2024-04-08 15:56:29', NULL);
INSERT INTO `group_basic` (`id`, `owner_id`, `avatar`, `type`, `identity`, `name`, `desc`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 1, 'http://127.0.0.1:9345/resource/public/file/default/member.jpg', 0, '12', 'test', '', '2024-04-08 15:57:27', '2024-04-08 15:57:27', NULL);
COMMIT;

-- ----------------------------
-- Table structure for group_member
-- ----------------------------
DROP TABLE IF EXISTS `group_member`;
CREATE TABLE `group_member` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `group_id` bigint DEFAULT NULL COMMENT '群id',
  `user_id` bigint DEFAULT NULL COMMENT '用户id',
  `type` int DEFAULT NULL COMMENT '类型',
  `is_admin` int DEFAULT NULL COMMENT '是否是管理员 1 yes 2 no ',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `member_unique` (`group_id`,`user_id`) USING BTREE COMMENT '成员唯一索引'
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='群成员';

-- ----------------------------
-- Records of group_member
-- ----------------------------
BEGIN;
INSERT INTO `group_member` (`id`, `group_id`, `user_id`, `type`, `is_admin`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 4, 1, 0, 0, '2024-04-08 15:57:29', '2024-04-08 15:57:29', NULL);
INSERT INTO `group_member` (`id`, `group_id`, `user_id`, `type`, `is_admin`, `created_at`, `updated_at`, `deleted_at`) VALUES (2, 4, 9, 0, 0, '2024-04-08 15:57:29', '2024-04-08 15:57:29', NULL);
INSERT INTO `group_member` (`id`, `group_id`, `user_id`, `type`, `is_admin`, `created_at`, `updated_at`, `deleted_at`) VALUES (3, 4, 10, 0, 0, '2024-04-08 15:57:29', '2024-04-08 15:57:29', NULL);
INSERT INTO `group_member` (`id`, `group_id`, `user_id`, `type`, `is_admin`, `created_at`, `updated_at`, `deleted_at`) VALUES (4, 4, 11, 0, 0, '2024-04-08 15:57:29', '2024-04-08 15:57:29', NULL);
INSERT INTO `group_member` (`id`, `group_id`, `user_id`, `type`, `is_admin`, `created_at`, `updated_at`, `deleted_at`) VALUES (5, 4, 12, 0, 0, '2024-04-08 15:57:29', '2024-04-08 15:57:29', NULL);
COMMIT;

-- ----------------------------
-- Table structure for message
-- ----------------------------
DROP TABLE IF EXISTS `message`;
CREATE TABLE `message` (
  `from_id` bigint NOT NULL COMMENT '发送者',
  `target_id` bigint NOT NULL COMMENT '接收者',
  `msg_type` int DEFAULT NULL COMMENT '消息类型 群聊 私聊 广播',
  `media_type` int DEFAULT NULL COMMENT '消息类型 文字 图片 音频',
  `content` text CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci COMMENT '内容',
  `pic` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '图片',
  `url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `desc` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL,
  `amount` int DEFAULT NULL,
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC;

-- ----------------------------
-- Records of message
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user_basic
-- ----------------------------
DROP TABLE IF EXISTS `user_basic`;
CREATE TABLE `user_basic` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT 'id',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '姓名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '密码',
  `salt` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '加密盐',
  `phone` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '手机号',
  `email` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '邮箱',
  `identity` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT 'id 唯一标识(可能由系统分配或用户自定义)',
  `client_ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '客户端ip',
  `client_port` int DEFAULT NULL COMMENT '客户端端口',
  `login_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '登录时间',
  `heartbeat_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '心跳时间',
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '头像',
  `login_out_time` datetime DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '最后登录时间',
  `is_logout` int DEFAULT '1' COMMENT '是否退出 1 yes 2 no ',
  `is_forbidden` int DEFAULT '2' COMMENT '是否禁止 1 yes 2 no',
  `is_admin` int DEFAULT '2' COMMENT '是否是管理员 1 yes 2 no ',
  `device_info` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci DEFAULT NULL COMMENT '设备信息（mac window 苹果 安卓）',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `unique_email` (`email`) USING BTREE COMMENT '唯一邮箱',
  UNIQUE KEY `unique_id` (`identity`) USING BTREE COMMENT '唯一id',
  UNIQUE KEY `unique_phone` (`phone`) USING BTREE COMMENT '唯一手机号'
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户基础信息表';

-- ----------------------------
-- Records of user_basic
-- ----------------------------
BEGIN;
INSERT INTO `user_basic` (`id`, `name`, `password`, `salt`, `phone`, `email`, `identity`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `avatar`, `login_out_time`, `is_logout`, `is_forbidden`, `is_admin`, `device_info`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 'KK', 'bfb82ef6ecad2068c39b4ff71f52010f', 'h7JKoA', '18857392027', 'kk996icu@qq.com', 'kk996icu', '', 0, '2024-08-15 13:08:47', '2024-08-15 21:08:48', 'http://127.0.0.1:9345/resource/public/file/20240721/p961412293.jpg', '2024-08-15 13:08:47', 2, 2, 2, '', '2024-03-28 21:50:54', '2024-08-15 21:08:48', NULL);
INSERT INTO `user_basic` (`id`, `name`, `password`, `salt`, `phone`, `email`, `identity`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `avatar`, `login_out_time`, `is_logout`, `is_forbidden`, `is_admin`, `device_info`, `created_at`, `updated_at`, `deleted_at`) VALUES (9, '张凡', '7ae54e850e11269ba46dd0cd23c492f0', '4IGQFP', NULL, 'fan996icu@qq.com', 'chat-fan996icu', '', 0, '2024-04-07 09:30:36', '2024-04-07 09:30:36', 'http://127.0.0.1:9345/resource/public/file/20240403/kk1 (1).jpg', '2024-04-07 17:30:36', 1, 2, 1, '', '2024-04-01 14:09:37', '2024-04-07 17:30:36', NULL);
INSERT INTO `user_basic` (`id`, `name`, `password`, `salt`, `phone`, `email`, `identity`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `avatar`, `login_out_time`, `is_logout`, `is_forbidden`, `is_admin`, `device_info`, `created_at`, `updated_at`, `deleted_at`) VALUES (10, '鑫华', '4740b7b89de2291bbc6380b5247d4b89', 'NwBVzI', NULL, '2960732922@qq.com', 'chat-xin996icu', '', 0, '2024-08-14 13:18:17', '2024-08-14 21:18:18', 'http://127.0.0.1:9345/resource/public/file/20240727/10-240H4163326.jpg', '2024-08-14 13:18:17', 1, 2, 1, '', '2024-04-01 14:12:47', '2024-08-14 21:18:18', NULL);
INSERT INTO `user_basic` (`id`, `name`, `password`, `salt`, `phone`, `email`, `identity`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `avatar`, `login_out_time`, `is_logout`, `is_forbidden`, `is_admin`, `device_info`, `created_at`, `updated_at`, `deleted_at`) VALUES (11, 'THY', '4417f0b8ac3fa15cb3c84096acafef05', 'HwDMYq', NULL, '3471435758@qq.com', 'chat-THY996icu', '', 0, '2024-08-14 15:17:11', '2024-08-14 23:17:12', 'http://127.0.0.1:9345/resource/public/file/20240403/kk1 (3).jpg', '2024-08-14 15:17:11', 2, 2, 1, '', '2024-04-01 14:13:34', '2024-08-14 23:17:12', NULL);
INSERT INTO `user_basic` (`id`, `name`, `password`, `salt`, `phone`, `email`, `identity`, `client_ip`, `client_port`, `login_time`, `heartbeat_time`, `avatar`, `login_out_time`, `is_logout`, `is_forbidden`, `is_admin`, `device_info`, `created_at`, `updated_at`, `deleted_at`) VALUES (12, 'AA', '68c8d7dd62fc9cda29cc71425dd4c619', 'NNJoqu', NULL, '995420691@qq.com', 'AA996icu', '', 0, '2024-07-18 15:12:17', '2024-07-18 23:12:18', 'http://127.0.0.1:9345/resource/public/file/20240403/kk1 (4).jpg', '2024-07-18 15:12:17', 2, 2, 1, '', '2024-04-03 21:11:27', '2024-07-18 23:12:18', NULL);
COMMIT;

-- ----------------------------
-- Table structure for user_black_list
-- ----------------------------
DROP TABLE IF EXISTS `user_black_list`;
CREATE TABLE `user_black_list` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT 'id',
  `user_id` bigint DEFAULT NULL COMMENT '用户id',
  `bad_id` bigint DEFAULT NULL COMMENT '被拉黑用户',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`) USING BTREE,
  KEY `user_bad` (`user_id`,`bad_id`) USING BTREE COMMENT '用户id 和被拉黑人id 不能重合'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户黑名单表';

-- ----------------------------
-- Records of user_black_list
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for user_friends
-- ----------------------------
DROP TABLE IF EXISTS `user_friends`;
CREATE TABLE `user_friends` (
  `id` bigint NOT NULL AUTO_INCREMENT COMMENT '主键',
  `user_id` bigint NOT NULL COMMENT '用户id',
  `friend_id` bigint NOT NULL COMMENT '朋友id',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '修改时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `type` int DEFAULT '1' COMMENT '类型 1 好友 2 群聊',
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE KEY `user_id` (`user_id`,`friend_id`) USING BTREE COMMENT '唯一'
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci ROW_FORMAT=DYNAMIC COMMENT='用户好友表';

-- ----------------------------
-- Records of user_friends
-- ----------------------------
BEGIN;
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (1, 1, 1, '2024-04-03 16:13:30', '2024-04-03 16:13:30', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (2, 1, 9, '2024-04-03 16:13:32', '2024-04-03 16:13:32', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (3, 9, 1, '2024-04-03 16:13:32', '2024-04-03 16:13:32', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (5, 1, 10, '2024-04-03 20:53:05', '2024-04-03 20:53:05', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (6, 10, 1, '2024-04-03 20:53:05', '2024-04-03 20:53:05', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (7, 1, 11, '2024-04-03 20:53:06', '2024-04-03 20:53:06', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (8, 11, 1, '2024-04-03 20:53:06', '2024-04-03 20:53:06', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (9, 1, 12, '2024-04-03 21:13:53', '2024-04-03 21:13:53', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (10, 12, 1, '2024-04-03 21:13:54', '2024-04-03 21:13:54', NULL, 1);
INSERT INTO `user_friends` (`id`, `user_id`, `friend_id`, `created_at`, `updated_at`, `deleted_at`, `type`) VALUES (11, 9, 9, '2024-04-07 12:38:21', '2024-04-07 12:38:21', NULL, 1);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
