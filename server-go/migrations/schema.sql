CREATE DATABASE IF NOT EXISTS kk_chat DEFAULT CHARSET utf8mb4;
USE kk_chat;

CREATE TABLE IF NOT EXISTS users (
  id            BIGINT AUTO_INCREMENT PRIMARY KEY,
  identity      VARCHAR(32)  NOT NULL,
  name          VARCHAR(64)  NOT NULL,
  password_hash VARCHAR(100) NOT NULL,
  email         VARCHAR(128) NOT NULL,
  phone         VARCHAR(32)  NOT NULL DEFAULT '',
  avatar        VARCHAR(255) NOT NULL DEFAULT '',
  signature     VARCHAR(255) NOT NULL DEFAULT '',
  birth_date    DATE NULL,
  is_admin      TINYINT NOT NULL DEFAULT 0,
  status        TINYINT NOT NULL DEFAULT 1 COMMENT '1正常 2封禁',
  created_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY uk_identity (identity),
  UNIQUE KEY uk_email (email)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS friends (
  user_id    BIGINT NOT NULL,
  friend_id  BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (user_id, friend_id)
) ENGINE=InnoDB COMMENT='好友关系，双向各存一行';

CREATE TABLE IF NOT EXISTS `groups` (
  id         BIGINT AUTO_INCREMENT PRIMARY KEY,
  name       VARCHAR(64)  NOT NULL,
  avatar     VARCHAR(255) NOT NULL DEFAULT '',
  owner_id   BIGINT NOT NULL,
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY idx_name (name)
) ENGINE=InnoDB;

CREATE TABLE IF NOT EXISTS group_members (
  group_id  BIGINT NOT NULL,
  user_id   BIGINT NOT NULL,
  role      VARCHAR(16) NOT NULL DEFAULT 'member' COMMENT 'owner | member',
  joined_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (group_id, user_id),
  KEY idx_user (user_id)
) ENGINE=InnoDB;
