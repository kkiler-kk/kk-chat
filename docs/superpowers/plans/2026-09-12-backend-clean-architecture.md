# kk-chat 后端 Clean Architecture 重构实现计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 将 server-go 从 gin+gorm+全局单例结构重写为 Echo + sqlc + Clean Architecture（domain/usecase/adapter）的全新后端，并重设计 HTTP/WS 协议与 MySQL/Redis/Mongo 存储键空间。

**Architecture:** 依赖规则 `main → adapter → usecase → domain`，依赖只能向内。domain 零外部依赖；usecase 只依赖 domain 与自己定义的端口接口；Echo/sqlc/redis/mongo/websocket 均为可替换适配器。所有依赖在 `cmd/server/main.go` 手动构造函数注入，无全局变量、无包级单例。

**Tech Stack:** Go 1.22+、Echo v4、sqlc(MySQL)、go-redis/v9、mongo-driver、slog、golang-jwt/v5、bcrypt、validator/v10、base64Captcha、BurntSushi/toml、gomail、gorilla/websocket。

**Spec:** `docs/superpowers/specs/2026-09-12-kk-chat-refactor-design.md`（执行前必读；协议、错误码、Redis 键名、WS 事件名以 spec 为准，本计划与 spec 冲突时以 spec 为准）

## Global Constraints

- **不写自动化测试**（用户明确决策，覆盖默认 TDD 流程）；每个任务的验证 = `go build ./... && go vet ./...` 通过。
- 旧代码（`main.go`、`internal/app/`、`internal/router/`、`internal/websocket/`、`internal/global/`、`internal/consts/`、`common/`、`main_test.go`）在 Task B16 之前保持原样不动，新旧代码共存于同一 module；**新代码禁止 import 任何旧包**。
- 旧数据不迁移：MySQL/Mongo/Redis 全部按新 schema/键空间重建。
- 模块名保持 `server-go`；新代码 import 路径形如 `server-go/internal/domain`。
- HTTP 响应统一信封 `{"code":0,"message":"ok","data":...,"request_id":"..."}`；错误码分段与 HTTP 状态码映射严格按 spec §3.2。
- WS 信封 `{"event":"...","data":{...},"ts":1234567890}`；事件名严格按 spec §3.4（`ping/pong/chat.message/chat.recent_updated/presence.changed/group.updated/system.kick`）。
- JSON 字段一律 snake_case；时间字段输出 RFC3339 字符串。
- 日志一律 slog 结构化输出；禁止 fmt.Println；禁止在日志中输出 token、密码、邮箱验证码明文。
- 密码哈希一律 bcrypt（cost 10）；JWT HS256。
- 每个任务完成即 git commit（conventional commits 风格）。
- 工作分支：`refactor/clean-architecture`（Task B1 创建，全部后端任务在此分支提交）。

---

### Task B1: 分支、依赖与基础设施（config + logger）

**Files:**
- Create: `server-go/internal/config/config.go`
- Create: `server-go/config.example.toml`
- Create: `server-go/internal/platform/logger/logger.go`
- Modify: `server-go/go.mod`

**Interfaces:**
- Produces:
  - `config.Load(path string) (*config.Config, error)` — TOML 加载 + 环境变量覆盖
  - `config.Config` 结构体（字段见下方代码）
  - `logger.New(level string, pretty bool) *slog.Logger`

- [ ] **Step 1: 创建分支**

```bash
cd /Users/kk/data/code/kk-chat && git checkout -b refactor/clean-architecture
```

- [ ] **Step 2: 添加新依赖到 go.mod**

```bash
cd /Users/kk/data/code/kk-chat/server-go
go get github.com/labstack/echo/v4@latest
go get github.com/redis/go-redis/v9@latest
go get go.mongodb.org/mongo-driver@v1.14.0
go get github.com/golang-jwt/jwt/v5@latest
go get golang.org/x/crypto@latest
go get github.com/go-playground/validator/v10@latest
go get github.com/mojocn/base64Captcha@latest
go get github.com/google/uuid@latest
```

注意：go.mod 的 `go` 指令改为 `go 1.22`（本机 Go 版本需 ≥1.22，先 `go version` 确认；若低于 1.22 需先升级本机 Go）。gin/gorm 等旧依赖保留到 Task B16 再清理。

- [ ] **Step 3: 写 internal/config/config.go**

```go
package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Server   ServerConfig   `toml:"server"`
	MySQL    MySQLConfig    `toml:"mysql"`
	Redis    RedisConfig    `toml:"redis"`
	Mongo    MongoConfig    `toml:"mongo"`
	JWT      JWTConfig      `toml:"jwt"`
	Email    EmailConfig    `toml:"email"`
	Message  MessageConfig  `toml:"message"`
	Presence PresenceConfig `toml:"presence"`
	Log      LogConfig      `toml:"log"`
}

type ServerConfig struct {
	Port       int    `toml:"port"`
	Mode       string `toml:"mode"`        // debug | release
	StaticPath string `toml:"static_path"` // 静态资源目录，如 "static"
}

type MySQLConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
	DBName   string `toml:"dbname"`
}

func (m MySQLConfig) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local&charset=utf8mb4",
		m.User, m.Password, m.Host, m.Port, m.DBName)
}

type RedisConfig struct {
	Addr     string `toml:"addr"`
	Password string `toml:"password"`
	DB       int    `toml:"db"`
}

type MongoConfig struct {
	URI      string `toml:"uri"`
	Database string `toml:"database"`
}

type JWTConfig struct {
	Secret     string `toml:"secret"`
	TTLMinutes int    `toml:"ttl_minutes"`
}

func (j JWTConfig) TTL() time.Duration { return time.Duration(j.TTLMinutes) * time.Minute }

type EmailConfig struct {
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	From     string `toml:"from"`
}

type MessageConfig struct {
	RetentionDays  int `toml:"retention_days"`   // mongo TTL，默认 90
	NonFriendLimit int `toml:"non_friend_limit"` // 非好友消息条数上限，默认 3
	LimitTTLHours  int `toml:"limit_ttl_hours"`  // 限流计数窗口，默认 24
}

type PresenceConfig struct {
	TTLSeconds       int `toml:"ttl_seconds"`       // presence 键 TTL，默认 120
	HeartbeatSeconds int `toml:"heartbeat_seconds"` // 客户端心跳间隔，默认 60
}

type LogConfig struct {
	Level  string `toml:"level"`  // debug | info | warn | error
	Pretty bool   `toml:"pretty"` // true=文本彩色输出，false=JSON
}

// Load 从 path 加载 TOML 配置，随后应用环境变量覆盖。
func Load(path string) (*Config, error) {
	var cfg Config
	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return nil, fmt.Errorf("加载配置文件 %s 失败: %w", path, err)
	}
	applyEnv(&cfg)
	setDefaults(&cfg)
	return &cfg, nil
}

// 环境变量覆盖：KK_SERVER_PORT / KK_MYSQL_HOST / KK_MYSQL_PASSWORD /
// KK_REDIS_ADDR / KK_MONGO_URI / KK_JWT_SECRET / KK_LOG_LEVEL
func applyEnv(cfg *Config) {
	if v := os.Getenv("KK_SERVER_PORT"); v != "" {
		if p, err := strconv.Atoi(v); err == nil {
			cfg.Server.Port = p
		}
	}
	if v := os.Getenv("KK_MYSQL_HOST"); v != "" {
		cfg.MySQL.Host = v
	}
	if v := os.Getenv("KK_MYSQL_PASSWORD"); v != "" {
		cfg.MySQL.Password = v
	}
	if v := os.Getenv("KK_REDIS_ADDR"); v != "" {
		cfg.Redis.Addr = v
	}
	if v := os.Getenv("KK_MONGO_URI"); v != "" {
		cfg.Mongo.URI = v
	}
	if v := os.Getenv("KK_JWT_SECRET"); v != "" {
		cfg.JWT.Secret = v
	}
	if v := os.Getenv("KK_LOG_LEVEL"); v != "" {
		cfg.Log.Level = v
	}
}

func setDefaults(cfg *Config) {
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}
	if cfg.Server.Mode == "" {
		cfg.Server.Mode = "debug"
	}
	if cfg.Server.StaticPath == "" {
		cfg.Server.StaticPath = "static"
	}
	if cfg.JWT.TTLMinutes == 0 {
		cfg.JWT.TTLMinutes = 720
	}
	if cfg.Message.RetentionDays == 0 {
		cfg.Message.RetentionDays = 90
	}
	if cfg.Message.NonFriendLimit == 0 {
		cfg.Message.NonFriendLimit = 3
	}
	if cfg.Message.LimitTTLHours == 0 {
		cfg.Message.LimitTTLHours = 24
	}
	if cfg.Presence.TTLSeconds == 0 {
		cfg.Presence.TTLSeconds = 120
	}
	if cfg.Presence.HeartbeatSeconds == 0 {
		cfg.Presence.HeartbeatSeconds = 60
	}
	if cfg.Log.Level == "" {
		cfg.Log.Level = "info"
	}
	if cfg.Redis.Addr == "" {
		cfg.Redis.Addr = "127.0.0.1:6379"
	}
	if cfg.Mongo.URI == "" {
		cfg.Mongo.URI = "mongodb://127.0.0.1:27017"
	}
	if cfg.Mongo.Database == "" {
		cfg.Mongo.Database = "kk_chat"
	}
}
```

- [ ] **Step 4: 写 config.example.toml**

```toml
[server]
port = 8080
mode = "debug"          # debug | release
static_path = "static"

[mysql]
host = "127.0.0.1"
port = 3306
user = "root"
password = ""
dbname = "kk_chat"

[redis]
addr = "127.0.0.1:6379"
password = ""
db = 0

[mongo]
uri = "mongodb://127.0.0.1:27017"
database = "kk_chat"

[jwt]
secret = "change-me-in-production"
ttl_minutes = 720

[email]
host = "smtp.example.com"
port = 465
username = ""
password = ""
from = ""

[message]
retention_days = 90
non_friend_limit = 3
limit_ttl_hours = 24

[presence]
ttl_seconds = 120
heartbeat_seconds = 60

[log]
level = "info"
pretty = true
```

- [ ] **Step 5: 写 internal/platform/logger/logger.go**

```go
package logger

import (
	"log/slog"
	"os"
	"strings"
)

// New 创建 slog.Logger。level: debug|info|warn|error；pretty=true 用文本格式，否则 JSON。
func New(level string, pretty bool) *slog.Logger {
	var lv slog.Level
	switch strings.ToLower(level) {
	case "debug":
		lv = slog.LevelDebug
	case "warn":
		lv = slog.LevelWarn
	case "error":
		lv = slog.LevelError
	default:
		lv = slog.LevelInfo
	}
	opts := &slog.HandlerOptions{Level: lv}
	var h slog.Handler
	if pretty {
		h = slog.NewTextHandler(os.Stdout, opts)
	} else {
		h = slog.NewJSONHandler(os.Stdout, opts)
	}
	return slog.New(h)
}
```

- [ ] **Step 6: 验证**

Run: `cd server-go && go build ./internal/config/... ./internal/platform/... && go vet ./internal/config/... ./internal/platform/...`
Expected: 无输出（成功）。旧代码仍须整体可编译：`go build ./...` 也应通过。

- [ ] **Step 7: Commit**

```bash
git add server-go/go.mod server-go/go.sum server-go/internal/config server-go/internal/platform/logger server-go/config.example.toml
git commit -m "feat(server): 新后端基础设施——config(TOML+env) 与 slog logger"
```

---

### Task B2: 数据库 schema 与 sqlc 代码生成

**Files:**
- Create: `server-go/migrations/schema.sql`
- Create: `server-go/sqlc.yaml`
- Create: `server-go/internal/adapter/persistence/mysql/query/users.sql`
- Create: `server-go/internal/adapter/persistence/mysql/query/friends.sql`
- Create: `server-go/internal/adapter/persistence/mysql/query/groups.sql`
- Create: `server-go/internal/adapter/persistence/mysql/sqlcgen/`（sqlc 生成，纳入版本控制）

**Interfaces:**
- Produces: 包 `sqlcgen`（`server-go/internal/adapter/persistence/mysql/sqlcgen`）：`New(db DBTX) *Queries`、`NewTx(tx *sql.Tx) *Queries`，及下方全部查询方法。Task B6 的 repo 实现直接消费这些方法。

- [ ] **Step 1: 安装 sqlc**

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.27.0
```

- [ ] **Step 2: 写 migrations/schema.sql**

```sql
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
```

- [ ] **Step 3: 写 sqlc.yaml**

```yaml
version: "2"
sql:
  - engine: "mysql"
    queries: "internal/adapter/persistence/mysql/query"
    schema: "migrations/schema.sql"
    gen:
      go:
        package: "sqlcgen"
        out: "internal/adapter/persistence/mysql/sqlcgen"
```

- [ ] **Step 4: 写 query/users.sql**

```sql
-- name: InsertUser :execresult
INSERT INTO users (identity, name, password_hash, email) VALUES (?, ?, ?, ?);

-- name: GetUserByID :one
SELECT * FROM users WHERE id = ?;

-- name: GetByIdentity :one
SELECT * FROM users WHERE identity = ?;

-- name: GetByEmail :one
SELECT * FROM users WHERE email = ?;

-- name: UpdateUser :exec
UPDATE users SET
  name       = COALESCE(sqlc.narg(name), name),
  phone      = COALESCE(sqlc.narg(phone), phone),
  email      = COALESCE(sqlc.narg(email), email),
  avatar     = COALESCE(sqlc.narg(avatar), avatar),
  signature  = COALESCE(sqlc.narg(signature), signature),
  birth_date = COALESCE(sqlc.narg(birth_date), birth_date)
WHERE id = sqlc.arg(id);

-- name: UpdatePassword :exec
UPDATE users SET password_hash = ? WHERE id = ?;

-- name: SearchUsers :many
SELECT * FROM users
WHERE identity LIKE CONCAT(?, '%') OR name LIKE CONCAT(?, '%') OR email LIKE CONCAT(?, '%')
LIMIT ?;

-- name: GetUsersByIDs :many
SELECT * FROM users WHERE id IN (sqlc.slice(ids));
```

- [ ] **Step 5: 写 query/friends.sql**

```sql
-- name: InsertFriend :exec
INSERT IGNORE INTO friends (user_id, friend_id) VALUES (?, ?);

-- name: CountFriendship :one
SELECT COUNT(*) FROM friends WHERE user_id = ? AND friend_id = ?;

-- name: ListFriendIDs :many
SELECT friend_id FROM friends WHERE user_id = ?;
```

- [ ] **Step 6: 写 query/groups.sql**

```sql
-- name: InsertGroup :execresult
INSERT INTO `groups` (name, avatar, owner_id) VALUES (?, ?, ?);

-- name: GetGroupByID :one
SELECT * FROM `groups` WHERE id = ?;

-- name: InsertGroupMember :exec
INSERT IGNORE INTO group_members (group_id, user_id, role) VALUES (?, ?, ?);

-- name: CountMembership :one
SELECT COUNT(*) FROM group_members WHERE group_id = ? AND user_id = ?;

-- name: ListGroupMemberIDs :many
SELECT user_id FROM group_members WHERE group_id = ?;

-- name: CountGroupMembers :one
SELECT COUNT(*) FROM group_members WHERE group_id = ?;

-- name: ListGroupsByUser :many
SELECT g.* FROM `groups` g
JOIN group_members m ON g.id = m.group_id
WHERE m.user_id = ?;

-- name: SearchGroupsByName :many
SELECT * FROM `groups` WHERE name LIKE CONCAT(?, '%') LIMIT ?;
```

- [ ] **Step 7: 生成代码**

```bash
cd server-go && sqlc generate
```

Expected: `internal/adapter/persistence/mysql/sqlcgen/` 下生成 `db.go/models.go/users.sql.go/friends.sql.go/groups.sql.go`。若 sqlc 对 `sqlc.narg` 的 DATE 类型生成 `interface{}`，将 UpdateUser 的 birth_date 参数在 B6 repo 层用 `sqlcgen.UpdateUserParams` 显式传值适配。

- [ ] **Step 8: 验证 + Commit**

Run: `go build ./internal/adapter/... && go vet ./internal/adapter/...` → 成功

```bash
git add server-go/migrations server-go/sqlc.yaml server-go/internal/adapter/persistence/mysql
git commit -m "feat(server): MySQL schema 与 sqlc 查询/生成代码"
```

---

### Task B3: domain 领域核心

**Files:**
- Create: `server-go/internal/domain/user.go`
- Create: `server-go/internal/domain/message.go`
- Create: `server-go/internal/domain/friend.go`
- Create: `server-go/internal/domain/group.go`
- Create: `server-go/internal/domain/errors.go`

**Interfaces:**
- Produces（Task B4/B6-B11 全部依赖）:
  - `domain.User{ID int64; Identity, Name, Password, Email, Phone, Avatar, Signature string; BirthDate *time.Time; IsAdmin bool; Status UserStatus; CreatedAt, UpdatedAt time.Time}`
  - `domain.UserStatus`（`UserStatusActive=1`、`UserStatusBanned=2`）
  - `domain.ConversationType`（`ConvPrivate=1`、`ConvGroup=2`）与方法 `String() string`（"private"/"group"）、`ParseConversationType(s string)`
  - `domain.ConversationID{Type ConversationType; Value string}` + `NewPrivateConversation(a, b int64) ConversationID`、`NewGroupConversation(groupID int64) ConversationID`、`(c ConversationID) String() string`、`ParseConversationID(s string) (ConversationID, error)`、`(c ConversationID) PrivateParticipants() (a, b int64, ok bool)`
  - `domain.Message{ID string; ConversationID ConversationID; SenderID int64; Content string; ContentType ContentType; CreatedAt time.Time}`、`domain.ContentType`（`ContentText="text"`、`ContentImage="image"`）
  - `domain.Group{ID int64; Name, Avatar string; OwnerID int64; CreatedAt time.Time}`、`domain.GroupRole`（`RoleOwner="owner"`、`RoleMember="member"`）、`domain.GroupMember{GroupID, UserID int64; Role GroupRole; JoinedAt time.Time}`
  - `domain.Friendship{UserID, FriendID int64; CreatedAt time.Time}`
  - 哨兵错误：`ErrUserNotFound ErrEmailTaken ErrIdentityTaken ErrBadCredentials ErrUserBanned ErrCaptchaInvalid ErrEmailCodeInvalid ErrNotFriend ErrMsgLimitExceeded ErrGroupNotFound ErrNotGroupMember ErrTokenInvalid ErrInvalidConversation`

- [ ] **Step 1: 写 domain/errors.go**

```go
package domain

import "errors"

var (
	ErrUserNotFound        = errors.New("用户不存在")
	ErrEmailTaken          = errors.New("邮箱已被注册")
	ErrIdentityTaken       = errors.New("账号已被占用")
	ErrBadCredentials      = errors.New("账号或密码错误")
	ErrUserBanned          = errors.New("账号已被封禁")
	ErrCaptchaInvalid      = errors.New("图形验证码错误或已过期")
	ErrEmailCodeInvalid    = errors.New("邮箱验证码错误或已过期")
	ErrNotFriend           = errors.New("不是好友关系")
	ErrMsgLimitExceeded    = errors.New("非好友关系发送消息超过上限，请先添加好友")
	ErrGroupNotFound       = errors.New("群组不存在")
	ErrNotGroupMember      = errors.New("不是群成员")
	ErrTokenInvalid        = errors.New("token 无效或已过期")
	ErrInvalidConversation = errors.New("无效的会话标识")
)
```

- [ ] **Step 2: 写 domain/user.go**

```go
package domain

import "time"

type UserStatus int

const (
	UserStatusActive UserStatus = 1
	UserStatusBanned UserStatus = 2
)

type User struct {
	ID        int64
	Identity  string
	Name      string
	Password  string // bcrypt hash，任何出参不得携带
	Email     string
	Phone     string
	Avatar    string
	Signature string
	BirthDate *time.Time
	IsAdmin   bool
	Status    UserStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u *User) Banned() bool { return u.Status == UserStatusBanned }
```

- [ ] **Step 3: 写 domain/message.go（ConversationID 是本设计的核心值对象）**

```go
package domain

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type ConversationType int

const (
	ConvPrivate ConversationType = 1
	ConvGroup   ConversationType = 2
)

func (t ConversationType) String() string {
	if t == ConvGroup {
		return "group"
	}
	return "private"
}

func ParseConversationType(s string) (ConversationType, error) {
	switch s {
	case "private":
		return ConvPrivate, nil
	case "group":
		return ConvGroup, nil
	}
	return 0, ErrInvalidConversation
}

type ContentType string

const (
	ContentText  ContentType = "text"
	ContentImage ContentType = "image"
)

// ConversationID 私聊双方共享同一会话键："u:{小ID}_{大ID}"；群聊："g:{群ID}"。
type ConversationID struct {
	Type  ConversationType
	Value string
}

func NewPrivateConversation(a, b int64) ConversationID {
	lo, hi := a, b
	if a > b {
		lo, hi = b, a
	}
	return ConversationID{Type: ConvPrivate, Value: fmt.Sprintf("u:%d_%d", lo, hi)}
}

func NewGroupConversation(groupID int64) ConversationID {
	return ConversationID{Type: ConvGroup, Value: fmt.Sprintf("g:%d", groupID)}
}

func (c ConversationID) String() string { return c.Value }

func ParseConversationID(s string) (ConversationID, error) {
	switch {
	case strings.HasPrefix(s, "u:"):
		parts := strings.Split(strings.TrimPrefix(s, "u:"), "_")
		if len(parts) != 2 {
			return ConversationID{}, ErrInvalidConversation
		}
		a, err1 := strconv.ParseInt(parts[0], 10, 64)
		b, err2 := strconv.ParseInt(parts[1], 10, 64)
		if err1 != nil || err2 != nil || a == 0 || b == 0 || a > b {
			return ConversationID{}, ErrInvalidConversation
		}
		return ConversationID{Type: ConvPrivate, Value: s}, nil
	case strings.HasPrefix(s, "g:"):
		id, err := strconv.ParseInt(strings.TrimPrefix(s, "g:"), 10, 64)
		if err != nil || id <= 0 {
			return ConversationID{}, ErrInvalidConversation
		}
		return ConversationID{Type: ConvGroup, Value: s}, nil
	}
	return ConversationID{}, ErrInvalidConversation
}

// PrivateParticipants 返回私聊会话的两个参与者ID；群聊会话 ok=false。
func (c ConversationID) PrivateParticipants() (a, b int64, ok bool) {
	if c.Type != ConvPrivate {
		return 0, 0, false
	}
	parts := strings.Split(strings.TrimPrefix(c.Value, "u:"), "_")
	a, _ = strconv.ParseInt(parts[0], 10, 64)
	b, _ = strconv.ParseInt(parts[1], 10, 64)
	return a, b, true
}

type Message struct {
	ID             string // mongo ObjectID hex
	ConversationID ConversationID
	SenderID       int64
	Content        string
	ContentType    ContentType
	CreatedAt      time.Time
}
```

- [ ] **Step 4: 写 domain/friend.go 与 domain/group.go**

```go
package domain

import "time"

type Friendship struct {
	UserID    int64
	FriendID  int64
	CreatedAt time.Time
}
```

```go
package domain

import "time"

type GroupRole string

const (
	RoleOwner  GroupRole = "owner"
	RoleMember GroupRole = "member"
)

type Group struct {
	ID        int64
	Name      string
	Avatar    string
	OwnerID   int64
	CreatedAt time.Time
}

type GroupMember struct {
	GroupID  int64
	UserID   int64
	Role     GroupRole
	JoinedAt time.Time
}
```

- [ ] **Step 5: 验证 + Commit**

Run: `go build ./internal/domain/... && go vet ./internal/domain/...` → 成功；并确认 domain 无外部依赖：`go list -deps ./internal/domain | grep -v '^internal/\|^unicode\|^unsafe\|^errors$\|^fmt$\|^math\|^reflect$\|^runtime\|^sort$\|^strconv$\|^strings$\|^sync\|^time$\|^io\|^os\|^path\|^slices$\|^cmp$\|^iter$\|^vendor' || true`（输出应为空或仅标准库）。

```bash
git add server-go/internal/domain
git commit -m "feat(server): domain 领域核心——实体、ConversationID 值对象、哨兵错误"
```

---

### Task B4: apperror 与端口定义（usecase/port）

**Files:**
- Create: `server-go/internal/usecase/apperror/apperror.go`
- Create: `server-go/internal/usecase/port/repo.go`
- Create: `server-go/internal/usecase/port/store.go`
- Create: `server-go/internal/usecase/port/notifier.go`
- Create: `server-go/internal/usecase/port/usecase.go`
- Create: `server-go/internal/usecase/port/dto.go`

**Interfaces:**
- Produces: 全部端口接口与 usecase 输入/输出 DTO。后续所有任务的契约都锚定在这里，**签名必须逐字一致**（见下方代码）。

- [ ] **Step 1: 写 apperror/apperror.go**

```go
package apperror

import (
	"errors"
	"fmt"
)

// 业务错误码（与 spec §3.2 一致）
const (
	CodeOK                 = 0
	CodeInvalidParam       = 10001
	CodeInternal           = 10002
	CodeCaptchaInvalid     = 20001
	CodeBadCredentials     = 20002
	CodeTokenInvalid       = 20003
	CodeUserBanned         = 20004
	CodeEmailTaken         = 30001
	CodeIdentityTaken      = 30002
	CodeUserNotFound       = 30003
	CodeEmailCodeInvalid   = 30004
	CodeAlreadyFriend      = 40001
	CodeNotFriendLimit     = 40002
	CodeGroupNotFound      = 50001
	CodeNotGroupMember     = 50002
	CodeConversationInvalid = 60001
)

// Error 是应用层统一错误类型：Code 业务码，Message 用户可读消息，Err 内部错误链。
type Error struct {
	Code    int
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%d] %s: %v", e.Code, e.Message, e.Err)
	}
	return fmt.Sprintf("[%d] %s", e.Code, e.Message)
}

func (e *Error) Unwrap() error { return e.Err }

func New(code int, message string) *Error {
	return &Error{Code: code, Message: message}
}

func Wrap(code int, message string, err error) *Error {
	return &Error{Code: code, Message: message, Err: err}
}

// From 提取错误链中的 *Error；不是则返回 nil,false。
func From(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}
```

- [ ] **Step 2: 写 port/repo.go（数据仓储端口，由 B6/B8 实现）**

```go
package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type UserRepo interface {
	Create(ctx context.Context, u *domain.User) error // 成功后回填 u.ID
	ByID(ctx context.Context, id int64) (*domain.User, error)
	ByIdentity(ctx context.Context, identity string) (*domain.User, error)
	ByEmail(ctx context.Context, email string) (*domain.User, error)
	ByIDs(ctx context.Context, ids []int64) ([]*domain.User, error)
	UpdateProfile(ctx context.Context, id int64, p UpdateProfile) error
	Search(ctx context.Context, keyword string, limit int) ([]*domain.User, error)
}

// UpdateProfile 局部更新：nil 字段不改。
type UpdateProfile struct {
	Name      *string
	Phone     *string
	Email     *string
	Avatar    *string
	Signature *string
	BirthDate *time.Time
}

type FriendRepo interface {
	Add(ctx context.Context, userID, friendID int64) error // 事务写双向两行
	Exists(ctx context.Context, a, b int64) (bool, error)
	ListFriends(ctx context.Context, userID int64) ([]*domain.User, error)
}

type GroupRepo interface {
	Create(ctx context.Context, g *domain.Group, memberIDs []int64) error // 事务：群+owner成员+初始成员
	ByID(ctx context.Context, id int64) (*domain.Group, error)
	Join(ctx context.Context, groupID, userID int64) error
	IsMember(ctx context.Context, groupID, userID int64) (bool, error)
	MemberIDs(ctx context.Context, groupID int64) ([]int64, error)
	MemberCount(ctx context.Context, groupID int64) (int, error)
	ListByUser(ctx context.Context, userID int64) ([]*domain.Group, error)
	SearchByName(ctx context.Context, keyword string, limit int) ([]*domain.Group, error)
}

type MessageRepo interface {
	Save(ctx context.Context, msg *domain.Message) error // 成功后回填 msg.ID
	// History 返回按 created_at 升序的一页消息：取 created_at < cursor 的最近 limit 条。
	// cursor 为零值表示从最新开始。
	History(ctx context.Context, convID domain.ConversationID, cursor time.Time, limit int) ([]domain.Message, error)
}
```

- [ ] **Step 3: 写 port/store.go（Redis/邮件等基础设施端口，由 B7/B5 实现）**

```go
package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type TokenStore interface {
	Save(ctx context.Context, jti string, userID int64, ttl time.Duration) error
	Exists(ctx context.Context, jti string) (bool, error)
	Delete(ctx context.Context, jti string) error
}

type CaptchaStore interface {
	Save(ctx context.Context, id, answer string, ttl time.Duration) error
	VerifyAndDelete(ctx context.Context, id, answer string) (bool, error)
}

type EmailCodeStore interface {
	Save(ctx context.Context, email, code string, ttl time.Duration) error
	VerifyAndDelete(ctx context.Context, email, code string) (bool, error)
}

type PresenceStore interface {
	Online(ctx context.Context, userID int64, ttl time.Duration) error
	Offline(ctx context.Context, userID int64) error
	IsOnline(ctx context.Context, userID int64) (bool, error)
}

type ConversationSummary struct {
	Type                domain.ConversationType `json:"type"`
	LastContent         string                  `json:"last_message_content"`
	LastContentType     domain.ContentType      `json:"last_content_type"`
	LastSenderID        int64                   `json:"last_sender_id"`
	LastSenderName      string                  `json:"last_sender_name"`
	LastTime            time.Time               `json:"last_time"`
}

type RecentChatStore interface {
	Touch(ctx context.Context, userID int64, convID domain.ConversationID, ts time.Time) error
	List(ctx context.Context, userID int64, limit int) ([]domain.ConversationID, error) // 按时间倒序
	SaveSummary(ctx context.Context, convID domain.ConversationID, s ConversationSummary) error
	Summary(ctx context.Context, convID domain.ConversationID) (*ConversationSummary, error)
}

type MsgLimitStore interface {
	Incr(ctx context.Context, convID domain.ConversationID, userID int64, ttl time.Duration) (int64, error)
}

type EmailSender interface {
	Send(ctx context.Context, to, subject, body string) error
}

type CaptchaGenerator interface {
	Generate(ctx context.Context) (id string, answer string, b64Image string, err error)
}

type Clock interface {
	Now() time.Time
}
```

- [ ] **Step 4: 写 port/notifier.go（推送端口，由 B12 实现）**

```go
package port

import "context"

// NotifierEvent 传输无关的推送事件；Event 取值见 spec §3.4。
type NotifierEvent struct {
	Event string
	Data  any
}

// Notifier 语义为 fire-and-forget：推送失败只记日志，绝不阻塞业务流程。
type Notifier interface {
	ToUser(ctx context.Context, userID int64, ev NotifierEvent)
	ToUsers(ctx context.Context, userIDs []int64, ev NotifierEvent)
}
```

- [ ] **Step 5: 写 port/usecase.go（输入端口）与 usecase/dto.go（输入输出 DTO）**

```go
package port

import (
	"context"
	"time"

	"server-go/internal/domain"
)

type AuthUseCase interface {
	Register(ctx context.Context, in RegisterInput) error
	Login(ctx context.Context, in LoginInput) (*LoginOutput, error)
	Logout(ctx context.Context, jti string, userID int64) error
	GenerateCaptcha(ctx context.Context) (id, b64Image string, err error)
	SendEmailCode(ctx context.Context, email string) error
}

type UserUseCase interface {
	Me(ctx context.Context, userID int64) (*UserInfo, error)
	UpdateMe(ctx context.Context, userID int64, in UpdateUserInput) error
	Detail(ctx context.Context, viewerID, targetID int64) (*UserDetailOutput, error)
	Search(ctx context.Context, viewerID int64, keyword string) ([]UserSearchItem, error)
}

type FriendUseCase interface {
	Add(ctx context.Context, userID, targetID int64) error
	List(ctx context.Context, userID int64) ([]FriendItem, error)
}

type GroupUseCase interface {
	Create(ctx context.Context, ownerID int64, in CreateGroupInput) (*GroupItem, error)
	Join(ctx context.Context, groupID, userID int64) error
	List(ctx context.Context, userID int64) ([]GroupItem, error)
	Search(ctx context.Context, keyword string) ([]GroupItem, error)
}

type ChatUseCase interface {
	SendMessage(ctx context.Context, senderID int64, in SendMessageInput) (*OutgoingMessage, error)
	History(ctx context.Context, userID int64, convID domain.ConversationID, cursor time.Time, limit int) ([]OutgoingMessage, error)
	RecentConversations(ctx context.Context, userID int64) ([]RecentConversation, error)
}

type PresenceUseCase interface {
	OnConnect(ctx context.Context, userID int64) error    // 写 presence + 广播上线
	Heartbeat(ctx context.Context, userID int64) error    // 续期 presence
	OnDisconnect(ctx context.Context, userID int64) error // 删 presence + 广播下线
}
```

DTO（`internal/usecase/dto.go`，package usecase 会与 port 循环吗？——不会：port 定义接口引用这些 DTO，因此 DTO 必须放在 port 包内或更内层。**决定：DTO 直接定义在 port 包中，文件 port/dto.go**）：

```go
package port

import "time"

type RegisterInput struct {
	Identity  string
	Name      string
	Password  string
	Email     string
	EmailCode string
}

type LoginInput struct {
	Account       string // identity 或 email
	Password      string
	CaptchaID     string
	CaptchaAnswer string
}

type UserInfo struct {
	ID        int64      `json:"id"`
	Identity  string     `json:"identity"`
	Name      string     `json:"name"`
	Avatar    string     `json:"avatar"`
	Email     string     `json:"email"`
	Phone     string     `json:"phone"`
	Signature string     `json:"signature"`
	BirthDate *time.Time `json:"birth_date,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type LoginOutput struct {
	Token string    `json:"token"`
	User  *UserInfo `json:"user"`
}

type UpdateUserInput struct {
	Name         *string
	Phone        *string
	Email        *string // 改邮箱必须同时给 EmailCode
	EmailCode    string
	Avatar       *string
	Signature    *string
	BirthDate    *time.Time
}

type UserDetailOutput struct {
	*UserInfo
	IsFriend bool `json:"is_friend"`
	IsSelf   bool `json:"is_self"`
}

type UserSearchItem struct {
	ID       int64  `json:"id"`
	Identity string `json:"identity"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	IsFriend bool   `json:"is_friend"`
}

type FriendItem struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Online bool   `json:"online"`
}

type GroupItem struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Avatar      string `json:"avatar"`
	OwnerID     int64  `json:"owner_id"`
	MemberCount int    `json:"member_count"`
}

type CreateGroupInput struct {
	Name      string
	MemberIDs []int64
}

type SendMessageInput struct {
	ConversationID string // 原始字符串，usecase 内 Parse
	Content        string
	ContentType    string // "text" | "image"
}

type OutgoingMessage struct {
	ID             string    `json:"id"`
	ConversationID string    `json:"conversation_id"`
	SenderID       int64     `json:"sender_id"`
	SenderName     string    `json:"sender_name"`
	SenderAvatar   string    `json:"sender_avatar"`
	Content        string    `json:"content"`
	ContentType    string    `json:"content_type"`
	CreatedAt      time.Time `json:"created_at"`
}

type RecentConversation struct {
	ConversationID string    `json:"conversation_id"`
	Type           string    `json:"type"` // private | group
	PeerID         int64     `json:"peer_id"`
	PeerName       string    `json:"peer_name"`
	PeerAvatar     string    `json:"peer_avatar"`
	Online         bool      `json:"online"`
	LastContent    string    `json:"last_message_content"`
	LastSenderName string    `json:"last_sender_name"`
	LastTime       time.Time `json:"last_time"`
}
```

注意：`port/usecase.go` 中的接口签名引用本包 DTO（`RegisterInput` 等），`dto.go` 与 `usecase.go` 同属 `package port`；`usecase.go` 需要 import `time` 与 `server-go/internal/domain`（History 用）。

- [ ] **Step 6: 验证 + Commit**

Run: `go build ./internal/usecase/... && go vet ./internal/usecase/...` → 成功

```bash
git add server-go/internal/usecase
git commit -m "feat(server): apperror 错误类型与 usecase 端口/DTO 契约"
```

---

### Task B5: platform 基础组件（token/password/random/captcha/mailer/clock）

**Files:**
- Create: `server-go/internal/platform/token/jwt.go`
- Create: `server-go/internal/platform/password/password.go`
- Create: `server-go/internal/platform/random/random.go`
- Create: `server-go/internal/platform/captchagen/captcha.go`
- Create: `server-go/internal/platform/mailer/mailer.go`
- Create: `server-go/internal/platform/clock/clock.go`
- Modify: `server-go/internal/usecase/port/store.go`（补 TokenIssuer 接口）

**Interfaces:**
- Produces:
  - `token.NewManager(secret string, ttl time.Duration) *token.Manager`；`(*Manager) Issue(ctx context.Context, userID int64) (tokenStr, jti string, err error)`（实现 `port.TokenIssuer`）；`(*Manager) Parse(tokenStr string) (*token.Claims, error)`；`token.Claims{UserID int64; JTI string}`
  - `password.Hash(pw string) (string, error)`、`password.Verify(hash, pw string) bool`
  - `random.Digits(n int) string`
  - `captchagen.New() *captchagen.Generator`，`(*Generator) Generate(ctx) (id, answer, b64Image string, err error)`（实现 `port.CaptchaGenerator`）
  - `mailer.New(host string, port int, user, pass, from string) *mailer.Mailer`，`(*Mailer) Send(ctx, to, subject, body string) error`（实现 `port.EmailSender`）
  - `clock.Real{}` 实现 `port.Clock`

- [ ] **Step 1: 在 port/store.go 末尾追加 TokenIssuer 接口**

```go
// TokenIssuer 由 platform/token.Manager 实现；usecase 登录时签发 token。
type TokenIssuer interface {
	Issue(ctx context.Context, userID int64) (tokenStr string, jti string, err error)
}
```

- [ ] **Step 2: 写 platform/token/jwt.go**

```go
package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"server-go/internal/domain"
)

type Claims struct {
	UserID int64
	JTI    string
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

func (m *Manager) Issue(_ context.Context, userID int64) (string, string, error) {
	jti := uuid.NewString()
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"exp": time.Now().Add(m.ttl).Unix(),
		"iat": time.Now().Unix(),
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", "", err
	}
	return t, jti, nil
}

func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !t.Valid {
		return nil, domain.ErrTokenInvalid
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}
	sub, err := claims.GetSubject()
	if err != nil {
		return nil, domain.ErrTokenInvalid
	}
	uidF, err := jwt.ParseWithClaims(sub, &jwt.MapClaims{}, func(*jwt.Token) (any, error) { return nil, nil })
	_ = uidF
	_ = err
	// sub 存的是数字，MapClaims 里为 float64
	uidNum, ok := claims["sub"].(float64)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}
	jti, _ := claims["jti"].(string)
	if jti == "" {
		return nil, errors.New("missing jti")
	}
	return &Claims{UserID: int64(uidNum), JTI: jti}, nil
}
```

注意：上面 Parse 中 `uidF/err` 两行是无效残留，**实现时删除**（保留 `claims["sub"].(float64)` 路径即可）。

- [ ] **Step 3: 写 password / random / clock**

```go
package password

import "golang.org/x/crypto/bcrypt"

func Hash(pw string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(pw), 10)
	return string(b), err
}

func Verify(hash, pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pw)) == nil
}
```

```go
package random

import (
	"crypto/rand"
	"math/big"
)

// Digits 返回 n 位随机数字字符串（crypto/rand）。
func Digits(n int) string {
	const digits = "0123456789"
	b := make([]byte, n)
	for i := range b {
		idx, _ := rand.Int(rand.Reader, big.NewInt(10))
		b[i] = digits[idx.Int64()]
	}
	return string(b)
}
```

```go
package clock

import "time"

type Real struct{}

func (Real) Now() time.Time { return time.Now() }
```

- [ ] **Step 4: 写 captchagen/captcha.go**

```go
package captchagen

import (
	"context"

	"github.com/mojocn/base64Captcha"
)

type Generator struct {
	driver base64Captcha.Driver
}

func New() *Generator {
	// 4 位数字验证码，240x80
	return &Generator{driver: base64Captcha.NewDriverDigit(80, 240, 4, 0.7, 80)}
}

func (g *Generator) Generate(_ context.Context) (id, answer, b64Image string, err error) {
	id, b64Image, answer, err = g.driver.GenerateIdQuestionAnswer()
	return
}
```

- [ ] **Step 5: 写 mailer/mailer.go**

```go
package mailer

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Mailer struct {
	host         string
	port         int
	user, pass   string
	from         string
}

func New(host string, port int, user, pass, from string) *Mailer {
	return &Mailer{host: host, port: port, user: user, pass: pass, from: from}
}

func (m *Mailer) Send(_ context.Context, to, subject, body string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)
	d := gomail.NewDialer(m.host, m.port, m.user, m.pass)
	if err := d.DialAndSend(msg); err != nil {
		return fmt.Errorf("发送邮件到 %s 失败: %w", to, err)
	}
	return nil
}
```

- [ ] **Step 6: 验证 + Commit**

Run: `go build ./internal/platform/... ./internal/usecase/... && go vet ./internal/platform/... ./internal/usecase/...` → 成功

```bash
git add server-go/internal/platform server-go/internal/usecase/port/store.go
git commit -m "feat(server): platform 基础组件——JWT/bcrypt/随机数/图形码/邮件/时钟"
```

---

### Task B6: MySQL 适配器（sqlc 之上的 repo 实现）

**Files:**
- Create: `server-go/internal/adapter/persistence/mysql/db.go`
- Create: `server-go/internal/adapter/persistence/mysql/user_repo.go`
- Create: `server-go/internal/adapter/persistence/mysql/friend_repo.go`
- Create: `server-go/internal/adapter/persistence/mysql/group_repo.go`
- Modify: `server-go/go.mod`（显式添加 `github.com/go-sql-driver/mysql`）

**Interfaces:**
- Consumes: `sqlcgen.New(db)/NewTx(tx)` 及 B2 生成的全部查询方法；`port.UserRepo/FriendRepo/GroupRepo` 接口；`domain.*`
- Produces:
  - `mysql.NewDB(ctx context.Context, dsn string) (*sql.DB, error)`
  - `mysql.NewUserRepo(db *sql.DB) port.UserRepo`
  - `mysql.NewFriendRepo(db *sql.DB) port.FriendRepo`
  - `mysql.NewGroupRepo(db *sql.DB) port.GroupRepo`

- [ ] **Step 1: `go get github.com/go-sql-driver/mysql@latest`，写 db.go**

```go
package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开 mysql 失败: %w", err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("连接 mysql 失败: %w", err)
	}
	return db, nil
}
```

- [ ] **Step 2: 写 user_repo.go（含完整映射函数，其余方法按行为说明实现）**

```go
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) port.UserRepo { return &UserRepo{db: db} }

func mapUser(row sqlcgen.User) *domain.User {
	u := &domain.User{
		ID:        row.ID,
		Identity:  row.Identity,
		Name:      row.Name,
		Password:  row.PasswordHash,
		Email:     row.Email,
		Phone:     row.Phone,
		Avatar:    row.Avatar,
		Signature: row.Signature,
		IsAdmin:   row.IsAdmin == 1,
		Status:    domain.UserStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.BirthDate.Valid {
		t := row.BirthDate.Time
		u.BirthDate = &t
	}
	return u
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	res, err := sqlcgen.New(r.db).InsertUser(ctx, sqlcgen.InsertUserParams{
		Identity: u.Identity, Name: u.Name, PasswordHash: u.Password, Email: u.Email,
	})
	if err != nil {
		return fmt.Errorf("插入用户: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取用户ID: %w", err)
	}
	u.ID = id
	return nil
}

func (r *UserRepo) ByID(ctx context.Context, id int64) (*domain.User, error) {
	row, err := sqlcgen.New(r.db).GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户: %w", err)
	}
	return mapUser(row), nil
}
```

其余方法逐一实现（同一文件）：
- `ByIdentity/ByEmail`：调用 `GetByIdentity/GetByEmail`，`sql.ErrNoRows → domain.ErrUserNotFound`，映射同 `ByID`。
- `ByIDs`：`GetUsersByIDs(ctx, ids)`；空 ids 直接返回空切片（避免生成 `IN ()` 非法 SQL）。
- `UpdateProfile(ctx, id, p)`：构造 `sqlcgen.UpdateUserParams`——每个 `*string` 字段用 `sql.NullString{String: *p.X, Valid: p.X != nil}`（若 sqlc 生成的是 `interface{}` 则 nil→nil、非nil→值），BirthDate 同理；`id` 必填。注意 email 唯一冲突时 MySQL 返回 1062 错误 → 转为 `domain.ErrEmailTaken`（用 `go-sql-driver/mysql.MySQLError.Number == 1062` 判断）。
- `Search(ctx, keyword, limit)`：`SearchUsers(ctx, keyword, keyword, keyword, int32(limit))`（LIKE 前缀匹配参数重复传 3 次），映射为切片。

- [ ] **Step 3: 写 friend_repo.go**

```go
package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type FriendRepo struct {
	db *sql.DB
}

func NewFriendRepo(db *sql.DB) port.FriendRepo { return &FriendRepo{db: db} }

// Add 在一个事务中写双向两行。
func (r *FriendRepo) Add(ctx context.Context, userID, friendID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Commit 后 Rollback 返回错误可忽略
	q := sqlcgen.NewTx(tx)
	if err := q.InsertFriend(ctx, sqlcgen.InsertFriendParams{UserID: userID, FriendID: friendID}); err != nil {
		return fmt.Errorf("写入好友关系: %w", err)
	}
	if err := q.InsertFriend(ctx, sqlcgen.InsertFriendParams{UserID: friendID, FriendID: userID}); err != nil {
		return fmt.Errorf("写入反向好友关系: %w", err)
	}
	return tx.Commit()
}

func (r *FriendRepo) Exists(ctx context.Context, a, b int64) (bool, error) {
	n, err := sqlcgen.New(r.db).CountFriendship(ctx, sqlcgen.CountFriendshipParams{UserID: a, FriendID: b})
	return n > 0, err
}

func (r *FriendRepo) ListFriends(ctx context.Context, userID int64) ([]*domain.User, error) {
	q := sqlcgen.New(r.db)
	ids, err := q.ListFriendIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询好友ID: %w", err)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := q.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("批量查询好友: %w", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapUser(row))
	}
	return out, nil
}
```

（sqlc 生成的 Params 结构体名/字段以实际生成为准，下同。）

- [ ] **Step 4: 写 group_repo.go**

```go
package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type GroupRepo struct {
	db *sql.DB
}

func NewGroupRepo(db *sql.DB) port.GroupRepo { return &GroupRepo{db: db} }

func mapGroup(row sqlcgen.Group) *domain.Group {
	return &domain.Group{
		ID: row.ID, Name: row.Name, Avatar: row.Avatar,
		OwnerID: row.OwnerID, CreatedAt: row.CreatedAt,
	}
}

// Create 事务：建群 + owner 成员行 + 初始成员行（INSERT IGNORE 天然去重）。
func (r *GroupRepo) Create(ctx context.Context, g *domain.Group, memberIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	q := sqlcgen.NewTx(tx)
	res, err := q.InsertGroup(ctx, sqlcgen.InsertGroupParams{Name: g.Name, Avatar: g.Avatar, OwnerID: g.OwnerID})
	if err != nil {
		return fmt.Errorf("插入群组: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取群ID: %w", err)
	}
	g.ID = id
	if err := q.InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
		GroupID: id, UserID: g.OwnerID, Role: string(domain.RoleOwner),
	}); err != nil {
		return fmt.Errorf("插入群主成员: %w", err)
	}
	for _, uid := range memberIDs {
		if err := q.InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
			GroupID: id, UserID: uid, Role: string(domain.RoleMember),
		}); err != nil {
			return fmt.Errorf("插入群成员: %w", err)
		}
	}
	return tx.Commit()
}

func (r *GroupRepo) ByID(ctx context.Context, id int64) (*domain.Group, error) {
	row, err := sqlcgen.New(r.db).GetGroupByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrGroupNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询群组: %w", err)
	}
	return mapGroup(row), nil
}

func (r *GroupRepo) Join(ctx context.Context, groupID, userID int64) error {
	return sqlcgen.New(r.db).InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
		GroupID: groupID, UserID: userID, Role: string(domain.RoleMember),
	})
}

func (r *GroupRepo) IsMember(ctx context.Context, groupID, userID int64) (bool, error) {
	n, err := sqlcgen.New(r.db).CountMembership(ctx, sqlcgen.CountMembershipParams{GroupID: groupID, UserID: userID})
	return n > 0, err
}

func (r *GroupRepo) MemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return sqlcgen.New(r.db).ListGroupMemberIDs(ctx, groupID)
}

func (r *GroupRepo) MemberCount(ctx context.Context, groupID int64) (int, error) {
	n, err := sqlcgen.New(r.db).CountGroupMembers(ctx, groupID)
	return int(n), err
}

func (r *GroupRepo) ListByUser(ctx context.Context, userID int64) ([]*domain.Group, error) {
	rows, err := sqlcgen.New(r.db).ListGroupsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户群列表: %w", err)
	}
	out := make([]*domain.Group, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapGroup(row))
	}
	return out, nil
}

func (r *GroupRepo) SearchByName(ctx context.Context, keyword string, limit int) ([]*domain.Group, error) {
	rows, err := sqlcgen.New(r.db).SearchGroupsByName(ctx, sqlcgen.SearchGroupsByNameParams{
		Name: keyword, Limit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("搜索群组: %w", err)
	}
	out := make([]*domain.Group, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapGroup(row))
	}
	return out, nil
}
```

- [ ] **Step 5: 验证 + Commit**

Run: `go build ./internal/adapter/... && go vet ./internal/adapter/...` → 成功

```bash
git add server-go/internal/adapter/persistence/mysql server-go/go.mod server-go/go.sum
git commit -m "feat(server): MySQL 适配器——sqlc 之上的 User/Friend/Group 仓储"
```

---

### Task B7: Redis 适配器（六个 store）

**Files:**
- Create: `server-go/internal/adapter/persistence/redis/client.go`
- Create: `server-go/internal/adapter/persistence/redis/token_store.go`
- Create: `server-go/internal/adapter/persistence/redis/captcha_store.go`
- Create: `server-go/internal/adapter/persistence/redis/email_code_store.go`
- Create: `server-go/internal/adapter/persistence/redis/presence_store.go`
- Create: `server-go/internal/adapter/persistence/redis/recent_chat_store.go`
- Create: `server-go/internal/adapter/persistence/redis/msg_limit_store.go`

**Interfaces:**
- Consumes: `port.TokenStore/CaptchaStore/EmailCodeStore/PresenceStore/RecentChatStore/MsgLimitStore` 接口、`domain.ConversationID`
- Produces:
  - `redisx.NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error)`（包名 `redisx`，避免与 go-redis 包名冲突）
  - `redisx.NewTokenStore(c *redis.Client, ttl time.Duration) port.TokenStore`
  - `redisx.NewCaptchaStore(c *redis.Client) port.CaptchaStore`
  - `redisx.NewEmailCodeStore(c *redis.Client) port.EmailCodeStore`
  - `redisx.NewPresenceStore(c *redis.Client, ttl time.Duration) port.PresenceStore`
  - `redisx.NewRecentChatStore(c *redis.Client) port.RecentChatStore`
  - `redisx.NewMsgLimitStore(c *redis.Client) port.MsgLimitStore`

- [ ] **Step 1: 写 client.go 与 token_store.go**

```go
package redisx

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func NewClient(ctx context.Context, addr, password string, db int) (*redis.Client, error) {
	c := redis.NewClient(&redis.Options{Addr: addr, Password: password, DB: db})
	if err := c.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接 redis 失败: %w", err)
	}
	return c, nil
}
```

```go
package redisx

// 键：token:{jti}，值：user_id 十进制字符串，TTL 由构造注入。
type TokenStore struct {
	c   *redis.Client
	ttl time.Duration
}

func NewTokenStore(c *redis.Client, ttl time.Duration) *TokenStore { return &TokenStore{c: c, ttl: ttl} }

func (s *TokenStore) Save(ctx context.Context, jti string, userID int64, ttl time.Duration) error {
	use := s.ttl
	if ttl > 0 {
		use = ttl
	}
	return s.c.Set(ctx, "token:"+jti, userID, use).Err()
}

func (s *TokenStore) Exists(ctx context.Context, jti string) (bool, error) {
	n, err := s.c.Exists(ctx, "token:"+jti).Result()
	return n > 0, err
}

func (s *TokenStore) Delete(ctx context.Context, jti string) error {
	return s.c.Del(ctx, "token:"+jti).Err()
}
```

- [ ] **Step 2: 写 captcha_store.go 与 email_code_store.go**

```go
package redisx

import (
	"context"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// 键：captcha:{id}，TTL 由调用方传入（5min）。
type CaptchaStore struct{ c *redis.Client }

func NewCaptchaStore(c *redis.Client) port.CaptchaStore { return &CaptchaStore{c: c} }

func (s *CaptchaStore) Save(ctx context.Context, id, answer string, ttl time.Duration) error {
	return s.c.Set(ctx, "captcha:"+id, answer, ttl).Err()
}

// VerifyAndDelete 用 GETDEL（Redis ≥ 6.2）；图形验证码不区分大小写。
func (s *CaptchaStore) VerifyAndDelete(ctx context.Context, id, answer string) (bool, error) {
	v, err := s.c.GetDel(ctx, "captcha:"+id).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return strings.EqualFold(v, answer), nil
}
```

```go
package redisx

// 键：verify:email:{email}，TTL 由调用方传入（10min）；纯数字验证码直接 == 比较。
type EmailCodeStore struct{ c *redis.Client }

func NewEmailCodeStore(c *redis.Client) port.EmailCodeStore { return &EmailCodeStore{c: c} }

func (s *EmailCodeStore) Save(ctx context.Context, email, code string, ttl time.Duration) error {
	return s.c.Set(ctx, "verify:email:"+email, code, ttl).Err()
}

func (s *EmailCodeStore) VerifyAndDelete(ctx context.Context, email, code string) (bool, error) {
	v, err := s.c.GetDel(ctx, "verify:email:"+email).Result()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return v == code, nil
}
```

- [ ] **Step 3: 写 presence_store.go**

```go
package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/usecase/port"
)

// 键：presence:{uid}，值 "1"，TTL 由构造注入（默认 120s）。
type PresenceStore struct {
	c   *redis.Client
	ttl time.Duration
}

func NewPresenceStore(c *redis.Client, ttl time.Duration) port.PresenceStore {
	return &PresenceStore{c: c, ttl: ttl}
}

// Online 写入/续期同一操作。
func (s *PresenceStore) Online(ctx context.Context, userID int64, ttl time.Duration) error {
	use := s.ttl
	if ttl > 0 {
		use = ttl
	}
	return s.c.Set(ctx, "presence:"+strconv.FormatInt(userID, 10), "1", use).Err()
}

func (s *PresenceStore) Offline(ctx context.Context, userID int64) error {
	return s.c.Del(ctx, "presence:"+strconv.FormatInt(userID, 10)).Err()
}

func (s *PresenceStore) IsOnline(ctx context.Context, userID int64) (bool, error) {
	n, err := s.c.Exists(ctx, "presence:"+strconv.FormatInt(userID, 10)).Result()
	return n > 0, err
}
```

- [ ] **Step 4: 写 recent_chat_store.go（键设计核心）**

```go
package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// recent:{uid}  → ZSET，member=会话键字符串，score=最后消息 Unix 秒
// conv:{会话键} → HASH：type / last_message_content / last_content_type /
//                last_sender_id / last_sender_name / last_time(RFC3339)
type RecentChatStore struct{ c *redis.Client }

func NewRecentChatStore(c *redis.Client) *RecentChatStore { return &RecentChatStore{c: c} }

func recentKey(uid int64) string  { return "recent:" + strconv.FormatInt(uid, 10) }
func convKey(c domain.ConversationID) string { return "conv:" + c.Value }

func (s *RecentChatStore) Touch(ctx context.Context, uid int64, convID domain.ConversationID, ts time.Time) error {
	return s.c.ZAdd(ctx, recentKey(uid), redis.Z{Score: float64(ts.Unix()), Member: convID.Value}).Err()
}

func (s *RecentChatStore) List(ctx context.Context, uid int64, limit int) ([]domain.ConversationID, error) {
	vals, err := s.c.ZRevRange(ctx, recentKey(uid), 0, int64(limit-1)).Result()
	if err != nil {
		return nil, err
	}
	out := make([]domain.ConversationID, 0, len(vals))
	for _, v := range vals {
		if cid, err := domain.ParseConversationID(v); err == nil {
			out = append(out, cid)
		}
	}
	return out, nil
}

func (s *RecentChatStore) SaveSummary(ctx context.Context, convID domain.ConversationID, sum port.ConversationSummary) error {
	return s.c.HSet(ctx, convKey(convID), map[string]any{
		"type":                 convID.Type.String(),
		"last_message_content": sum.LastContent,
		"last_content_type":    string(sum.LastContentType),
		"last_sender_id":       sum.LastSenderID,
		"last_sender_name":     sum.LastSenderName,
		"last_time":            sum.LastTime.Format(time.RFC3339),
	}).Err()
}

func (s *RecentChatStore) Summary(ctx context.Context, convID domain.ConversationID) (*port.ConversationSummary, error) {
	m, err := s.c.HGetAll(ctx, convKey(convID)).Result()
	if err != nil || len(m) == 0 {
		return nil, err
	}
	t, _ := time.Parse(time.RFC3339, m["last_time"])
	ct, _ := domain.ParseConversationType(m["type"])
	senderID, _ := strconv.ParseInt(m["last_sender_id"], 10, 64)
	return &port.ConversationSummary{
		Type:            ct,
		LastContent:     m["last_message_content"],
		LastContentType: domain.ContentType(m["last_content_type"]),
		LastSenderID:    senderID,
		LastSenderName:  m["last_sender_name"],
		LastTime:        t,
	}, nil
}
```

- [ ] **Step 5: 写 msg_limit_store.go**

```go
package redisx

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// 键：msg:limit:{会话键}:{uid}。INCR 后若结果为 1（首次），EXPIRE 设 TTL。
// Incr 返回累计值；usecase 判断 count > limit 则拒绝。
type MsgLimitStore struct{ c *redis.Client }

func NewMsgLimitStore(c *redis.Client) port.MsgLimitStore { return &MsgLimitStore{c: c} }

func (s *MsgLimitStore) Incr(ctx context.Context, convID domain.ConversationID, uid int64, ttl time.Duration) (int64, error) {
	key := "msg:limit:" + convID.Value + ":" + strconv.FormatInt(uid, 10)
	n, err := s.c.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if n == 1 {
		s.c.Expire(ctx, key, ttl)
	}
	return n, nil
}
```

- [ ] **Step 6: 验证 + Commit**

Run: `go build ./internal/adapter/persistence/... && go vet ./internal/adapter/persistence/...` → 成功

```bash
git add server-go/internal/adapter/persistence/redis
git commit -m "feat(server): Redis 适配器——token/验证码/presence/最近会话(ZSET)/限流六个 store"
```

---

### Task B8: Mongo 适配器（MessageRepo + TTL 索引）

**Files:**
- Create: `server-go/internal/adapter/persistence/mongo/client.go`
- Create: `server-go/internal/adapter/persistence/mongo/message_repo.go`

**Interfaces:**
- Consumes: `port.MessageRepo`、`domain.Message/ConversationID`
- Produces:
  - `mongox.NewClient(ctx context.Context, uri string) (*mongo.Client, error)`
  - `mongox.NewMessageRepo(db *mongo.Database, retentionDays int) *mongox.MessageRepo`（实现 `port.MessageRepo`）
  - `(*MessageRepo) EnsureIndexes(ctx context.Context) error`

- [ ] **Step 1: 写 client.go**

```go
package mongox

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("连接 mongo 失败: %w", err)
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo 失败: %w", err)
	}
	return c, nil
}
```

- [ ] **Step 2: 写 message_repo.go**

```go
package mongox

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// 集合 "messages" 文档结构：
// { _id, conversation_id: "u:1_2", type: "private", sender_id: NumberLong,
//   content: "...", content_type: "text", created_at: ISODate }
type messageDoc struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	ConversationID string             `bson:"conversation_id"`
	Type           string             `bson:"type"`
	SenderID       int64              `bson:"sender_id"`
	Content        string             `bson:"content"`
	ContentType    string             `bson:"content_type"`
	CreatedAt      time.Time          `bson:"created_at"`
}

type MessageRepo struct {
	col           *mongo.Collection
	retentionDays int
}

func NewMessageRepo(db *mongo.Database, retentionDays int) *MessageRepo {
	return &MessageRepo{col: db.Collection("messages"), retentionDays: retentionDays}
}

// EnsureIndexes 幂等创建：查询索引 (conversation_id, created_at DESC) + TTL 索引 (created_at, expireAfterSeconds=retentionDays*86400)。
func (r *MessageRepo) EnsureIndexes(ctx context.Context) error {
	_, err := r.col.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "conversation_id", Value: 1}, {Key: "created_at", Value: -1}},
		},
		{
			Keys:    bson.D{{Key: "created_at", Value: 1}},
			Options: options.Index().SetExpireAfterSeconds(int32(r.retentionDays * 86400)),
		},
	})
	return err
}

func (r *MessageRepo) Save(ctx context.Context, msg *domain.Message) error {
	doc := messageDoc{
		ConversationID: msg.ConversationID.Value,
		Type:           msg.ConversationID.Type.String(),
		SenderID:       msg.SenderID,
		Content:        msg.Content,
		ContentType:    string(msg.ContentType),
		CreatedAt:      msg.CreatedAt.UTC(),
	}
	res, err := r.col.InsertOne(ctx, doc)
	if err != nil {
		return err
	}
	if oid, ok := res.InsertedID.(primitive.ObjectID); ok {
		msg.ID = oid.Hex()
	}
	return nil
}

// History 取 created_at < cursor 的最近 limit 条（倒序查询后反转为升序返回）；cursor 零值表示不加时间过滤。
func (r *MessageRepo) History(ctx context.Context, convID domain.ConversationID, cursor time.Time, limit int) ([]domain.Message, error) {
	filter := bson.M{"conversation_id": convID.Value}
	if !cursor.IsZero() {
		filter["created_at"] = bson.M{"$lt": cursor.UTC()}
	}
	opts := options.Find().
		SetSort(bson.D{{Key: "created_at", Value: -1}}).
		SetLimit(int64(limit))
	cur, err := r.col.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cur.Close(ctx)
	var docs []messageDoc
	if err := cur.All(ctx, &docs); err != nil {
		return nil, err
	}
	out := make([]domain.Message, 0, len(docs))
	for i := len(docs) - 1; i >= 0; i-- { // 反转为升序
		d := docs[i]
		out = append(out, domain.Message{
			ID:             d.ID.Hex(),
			ConversationID: convID,
			SenderID:       d.SenderID,
			Content:        d.Content,
			ContentType:    domain.ContentType(d.ContentType),
			CreatedAt:      d.CreatedAt,
		})
	}
	return out, nil
}

var _ port.MessageRepo = (*MessageRepo)(nil)
```

- [ ] **Step 3: 验证 + Commit**

Run: `go build ./internal/adapter/persistence/... && go vet ./internal/adapter/persistence/...` → 成功

```bash
git add server-go/internal/adapter/persistence/mongo
git commit -m "feat(server): Mongo 适配器——MessageRepo 与会话索引/TTL 索引"
```

---

### Task B9: usecase——auth 与 user

**Files:**
- Create: `server-go/internal/usecase/auth.go`
- Create: `server-go/internal/usecase/user.go`
- Create: `server-go/internal/usecase/mapping.go`

**Interfaces:**
- Consumes: `port.*` 全部接口、`domain.*`、`apperror.*`、`platform/password`、`platform/random`
- Produces:
  - `usecase.NewAuth(deps AuthDeps) port.AuthUseCase`，`AuthDeps{Users port.UserRepo; Tokens port.TokenStore; Issuer port.TokenIssuer; Captchas port.CaptchaStore; CaptchaGen port.CaptchaGenerator; Codes port.EmailCodeStore; Mailer port.EmailSender; Clock port.Clock; TokenTTL time.Duration}`
  - `usecase.NewUser(deps UserDeps) port.UserUseCase`，`UserDeps{Users port.UserRepo; Friends port.FriendRepo; Codes port.EmailCodeStore}`
  - `usecase.ToUserInfo(u *domain.User) *port.UserInfo`（mapping.go，含 Email/Phone）
  - `usecase.ToPublicUserInfo(u *domain.User) *port.UserInfo`（脱敏版：Email/Phone 置空）

- [ ] **Step 1: 写 mapping.go**

```go
package usecase

import (
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

// ToUserInfo 完整映射（仅用于本人视角/登录响应）。
func ToUserInfo(u *domain.User) *port.UserInfo {
	return &port.UserInfo{
		ID: u.ID, Identity: u.Identity, Name: u.Name, Avatar: u.Avatar,
		Email: u.Email, Phone: u.Phone, Signature: u.Signature,
		BirthDate: u.BirthDate, CreatedAt: u.CreatedAt,
	}
}

// ToPublicUserInfo 脱敏映射：Email/Phone 不输出。
func ToPublicUserInfo(u *domain.User) *port.UserInfo {
	info := ToUserInfo(u)
	info.Email = ""
	info.Phone = ""
	return info
}
```

- [ ] **Step 2: 写 auth.go（结构体 + 构造函数 + 五个方法）**

```go
package usecase

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"server-go/internal/domain"
	"server-go/internal/platform/password"
	"server-go/internal/platform/random"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type AuthDeps struct {
	Users      port.UserRepo
	Tokens     port.TokenStore
	Issuer     port.TokenIssuer
	Captchas   port.CaptchaStore
	CaptchaGen port.CaptchaGenerator
	Codes      port.EmailCodeStore
	Mailer     port.EmailSender
	Notifier   port.Notifier
	Clock      port.Clock
	TokenTTL   time.Duration
	Logger     *slog.Logger
}

type authUseCase struct{ d AuthDeps }

func NewAuth(d AuthDeps) port.AuthUseCase { return &authUseCase{d: d} }
```

方法行为（错误一律包装为 `*apperror.Error`，内部 err 放 `Err` 字段）：

- `Register(ctx, in)`：
  1. `codes.VerifyAndDelete(ctx, in.Email, in.EmailCode)` → false ⇒ `apperror.New(CodeEmailCodeInvalid, domain.ErrEmailCodeInvalid.Error())`
  2. `users.ByEmail(ctx, in.Email)`：err 为 nil（找到了）⇒ `apperror.New(CodeEmailTaken, ...)`；err 非 `domain.ErrUserNotFound` ⇒ `apperror.Wrap(CodeInternal, "注册失败", err)`
  3. `users.ByIdentity(ctx, in.Identity)`：同上 ⇒ `CodeIdentityTaken`
  4. `password.Hash(in.Password)` 失败 ⇒ `Wrap(CodeInternal,...)`
  5. `users.Create(ctx, &domain.User{Identity, Name, Password: hash, Email, Status: UserStatusActive})`
- `Login(ctx, in)`：
  1. `captchas.VerifyAndDelete(ctx, in.CaptchaID, in.CaptchaAnswer)` → false ⇒ `apperror.New(CodeCaptchaInvalid, domain.ErrCaptchaInvalid.Error())`
  2. 账号定位：先 `users.ByIdentity(in.Account)`，若 `ErrUserNotFound` 再 `users.ByEmail(in.Account)`；都找不到 ⇒ `apperror.New(CodeBadCredentials, domain.ErrBadCredentials.Error())`（不暴露账号是否存在）
  3. `u.Banned()` ⇒ `apperror.New(CodeUserBanned, domain.ErrUserBanned.Error())`
  4. `password.Verify(u.Password, in.Password)` false ⇒ `CodeBadCredentials`
  5. `issuer.Issue(ctx, u.ID)` → `tokens.Save(ctx, jti, u.ID, d.TokenTTL)`
  6. 返回 `&port.LoginOutput{Token: tokenStr, User: ToUserInfo(u)}`
- `Logout(ctx, jti, userID)`：`tokens.Delete(ctx, jti)`（错误 ⇒ `Wrap(CodeInternal, "退出失败", err)`）；随后 `notifier.ToUser(ctx, userID, port.NotifierEvent{Event: "system.kick", Data: map[string]any{"reason": "logout"}})` —— 将该用户**所有在线端**踢下线（前端 WSClient 收到 system.kick 后自行断开并清理本地会话；这是 spec §3.3"踢人"能力的唯一触发点，登出即全局登出）
- `GenerateCaptcha(ctx)`：`captchaGen.Generate` → `captchas.Save(ctx, id, answer, 5*time.Minute)` → 返回 `(id, b64Image, nil)`
- `SendEmailCode(ctx, email)`：`code := random.Digits(6)` → `codes.Save(ctx, email, code, 10*time.Minute)` → `mailer.Send(ctx, email, "kk-chat 邮箱验证码", fmt.Sprintf("您的验证码是 <b>%s</b>，10 分钟内有效。", code))`；mailer 失败仅 `d.Logger.Error` 并返回 `Wrap(CodeInternal, "验证码发送失败", err)`。**任何日志不得输出 code 明文。**

- [ ] **Step 3: 写 user.go**

```go
package usecase

type UserDeps struct {
	Users   port.UserRepo
	Friends port.FriendRepo
	Codes   port.EmailCodeStore
}

type userUseCase struct{ d UserDeps }

func NewUser(d UserDeps) port.UserUseCase { return &userUseCase{d: d} }
```

方法行为：
- `Me(ctx, userID)`：`users.ByID` →（`ErrUserNotFound` ⇒ `New(CodeUserNotFound,...)`）→ `ToUserInfo`
- `UpdateMe(ctx, userID, in)`：若 `in.Email != nil`：`codes.VerifyAndDelete(ctx, *in.Email, in.EmailCode)` false ⇒ `CodeEmailCodeInvalid`；再 `users.ByEmail(*in.Email)` 若找到且 ID != userID ⇒ `CodeEmailTaken`。然后组装 `port.UpdateProfile` 调 `users.UpdateProfile`（`domain.ErrEmailTaken` ⇒ `CodeEmailTaken`）
- `Detail(ctx, viewerID, targetID)`：`users.ByID(targetID)` → 基础输出 `ToPublicUserInfo`；`IsSelf = viewerID == targetID`；若非本人：`friends.Exists(viewerID, targetID)` → `IsFriend`；`IsSelf || IsFriend` 时回填 Email/Phone（直接用 `ToUserInfo` 重映射）
- `Search(ctx, viewerID, keyword)`：keyword 去首尾空格，空 ⇒ 返回空切片；`users.Search(ctx, keyword, 20)`；`friends.ListFriends? 不需要`——逐个 `friends.Exists(viewerID, u.ID)` 标记 `IsFriend`（结果 ≤20 条，可接受）；映射为 `[]port.UserSearchItem`

- [ ] **Step 4: 验证 + Commit**

Run: `go build ./internal/usecase/... && go vet ./internal/usecase/...` → 成功

```bash
git add server-go/internal/usecase
git commit -m "feat(server): usecase——认证与用户模块（注册/登录/登出/验证码/资料/搜索）"
```

---

### Task B10: usecase——friend、group、presence

**Files:**
- Create: `server-go/internal/usecase/friend.go`
- Create: `server-go/internal/usecase/group.go`
- Create: `server-go/internal/usecase/presence.go`

**Interfaces:**
- Consumes: `port.*`、`apperror.*`、`domain.*`
- Produces:
  - `usecase.NewFriend(deps FriendDeps) port.FriendUseCase`，`FriendDeps{Users port.UserRepo; Friends port.FriendRepo; Presence port.PresenceStore}`
  - `usecase.NewGroup(deps GroupDeps) port.GroupUseCase`，`GroupDeps{Groups port.GroupRepo; Users port.UserRepo; Notifier port.Notifier; Clock port.Clock}`
  - `usecase.NewPresence(deps PresenceDeps) port.PresenceUseCase`，`PresenceDeps{Presence port.PresenceStore; Friends port.FriendRepo; Notifier port.Notifier; TTL time.Duration}`

- [ ] **Step 1: 写 friend.go**

```go
package usecase

import (
	"context"
	"errors"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type FriendDeps struct {
	Users    port.UserRepo
	Friends  port.FriendRepo
	Presence port.PresenceStore
}

type friendUseCase struct{ d FriendDeps }

func NewFriend(d FriendDeps) port.FriendUseCase { return &friendUseCase{d: d} }

func (u *friendUseCase) Add(ctx context.Context, userID, targetID int64) error {
	if userID == targetID {
		return apperror.New(apperror.CodeInvalidParam, "不能添加自己为好友")
	}
	if _, err := u.d.Users.ByID(ctx, targetID); err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return apperror.New(apperror.CodeUserNotFound, domain.ErrUserNotFound.Error())
		}
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	exists, err := u.d.Friends.Exists(ctx, userID, targetID)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	if exists {
		return apperror.New(apperror.CodeAlreadyFriend, "已经是好友了")
	}
	if err := u.d.Friends.Add(ctx, userID, targetID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "添加好友失败", err)
	}
	return nil
}

func (u *friendUseCase) List(ctx context.Context, userID int64) ([]port.FriendItem, error) {
	friends, err := u.d.Friends.ListFriends(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询好友失败", err)
	}
	out := make([]port.FriendItem, 0, len(friends))
	for _, f := range friends {
		online, err := u.d.Presence.IsOnline(ctx, f.ID)
		if err != nil {
			online = false // 在线状态查询失败降级为离线，不打断列表
		}
		out = append(out, port.FriendItem{ID: f.ID, Name: f.Name, Avatar: f.Avatar, Online: online})
	}
	return out, nil
}
```

- [ ] **Step 2: 写 group.go**

```go
package usecase

import (
	"context"
	"errors"
	"strings"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type GroupDeps struct {
	Groups   port.GroupRepo
	Users    port.UserRepo
	Notifier port.Notifier
	Clock    port.Clock
}

type groupUseCase struct{ d GroupDeps }

func NewGroup(d GroupDeps) port.GroupUseCase { return &groupUseCase{d: d} }

func (u *groupUseCase) Create(ctx context.Context, ownerID int64, in port.CreateGroupInput) (*port.GroupItem, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperror.New(apperror.CodeInvalidParam, "群名称不能为空")
	}
	// 成员去重并剔除 owner
	seen := map[int64]bool{ownerID: true}
	members := make([]int64, 0, len(in.MemberIDs))
	for _, id := range in.MemberIDs {
		if !seen[id] {
			seen[id] = true
			members = append(members, id)
		}
	}
	g := &domain.Group{Name: name, OwnerID: ownerID, CreatedAt: u.d.Clock.Now()}
	if err := u.d.Groups.Create(ctx, g, members); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "创建群组失败", err)
	}
	if len(members) > 0 {
		u.d.Notifier.ToUsers(ctx, members, port.NotifierEvent{
			Event: "group.updated", Data: map[string]any{"group_id": g.ID},
		})
	}
	return &port.GroupItem{
		ID: g.ID, Name: g.Name, Avatar: g.Avatar,
		OwnerID: g.OwnerID, MemberCount: len(members) + 1,
	}, nil
}

func (u *groupUseCase) Join(ctx context.Context, groupID, userID int64) error {
	if _, err := u.d.Groups.ByID(ctx, groupID); err != nil {
		if errors.Is(err, domain.ErrGroupNotFound) {
			return apperror.New(apperror.CodeGroupNotFound, domain.ErrGroupNotFound.Error())
		}
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	isMember, err := u.d.Groups.IsMember(ctx, groupID, userID)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	if isMember {
		return nil // 幂等
	}
	if err := u.d.Groups.Join(ctx, groupID, userID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	ids, err := u.d.Groups.MemberIDs(ctx, groupID)
	if err == nil {
		others := make([]int64, 0, len(ids))
		for _, id := range ids {
			if id != userID {
				others = append(others, id)
			}
		}
		u.d.Notifier.ToUsers(ctx, others, port.NotifierEvent{
			Event: "group.updated", Data: map[string]any{"group_id": groupID},
		})
	}
	return nil
}

func (u *groupUseCase) List(ctx context.Context, userID int64) ([]port.GroupItem, error) {
	groups, err := u.d.Groups.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询群列表失败", err)
	}
	return u.toItems(ctx, groups), nil
}

func (u *groupUseCase) Search(ctx context.Context, keyword string) ([]port.GroupItem, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []port.GroupItem{}, nil
	}
	groups, err := u.d.Groups.SearchByName(ctx, keyword, 20)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "搜索群组失败", err)
	}
	return u.toItems(ctx, groups), nil
}

func (u *groupUseCase) toItems(ctx context.Context, groups []*domain.Group) []port.GroupItem {
	out := make([]port.GroupItem, 0, len(groups))
	for _, g := range groups {
		count, err := u.d.Groups.MemberCount(ctx, g.ID)
		if err != nil {
			count = 0
		}
		out = append(out, port.GroupItem{
			ID: g.ID, Name: g.Name, Avatar: g.Avatar, OwnerID: g.OwnerID, MemberCount: count,
		})
	}
	return out
}
```

- [ ] **Step 3: 写 presence.go**

```go
package usecase

import (
	"context"
	"time"

	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type PresenceDeps struct {
	Presence port.PresenceStore
	Friends  port.FriendRepo
	Notifier port.Notifier
	TTL      time.Duration
}

type presenceUseCase struct{ d PresenceDeps }

func NewPresence(d PresenceDeps) port.PresenceUseCase { return &presenceUseCase{d: d} }

func (u *presenceUseCase) broadcast(ctx context.Context, userID int64, online bool) {
	friendIDs := func() []int64 {
		friendUsers, err := u.d.Friends.ListFriends(ctx, userID)
		if err != nil {
			return nil
		}
		ids := make([]int64, 0, len(friendUsers))
		for _, f := range friendUsers {
			ids = append(ids, f.ID)
		}
		return ids
	}()
	if len(friendIDs) == 0 {
		return
	}
	u.d.Notifier.ToUsers(ctx, friendIDs, port.NotifierEvent{
		Event: "presence.changed",
		Data:  map[string]any{"user_id": userID, "online": online},
	})
}

func (u *presenceUseCase) OnConnect(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Online(ctx, userID, u.d.TTL); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "上线失败", err)
	}
	u.broadcast(ctx, userID, true)
	return nil
}

func (u *presenceUseCase) Heartbeat(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Online(ctx, userID, u.d.TTL); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "心跳失败", err)
	}
	return nil
}

func (u *presenceUseCase) OnDisconnect(ctx context.Context, userID int64) error {
	if err := u.d.Presence.Offline(ctx, userID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "下线失败", err)
	}
	u.broadcast(ctx, userID, false)
	return nil
}
```

- [ ] **Step 4: 验证 + Commit**

Run: `go build ./internal/usecase/... && go vet ./internal/usecase/...` → 成功

```bash
git add server-go/internal/usecase
git commit -m "feat(server): usecase——好友/群组/在线状态模块"
```

---

### Task B11: usecase——chat（发消息/历史/最近会话）

**Files:**
- Create: `server-go/internal/usecase/chat.go`

**Interfaces:**
- Consumes: 全部相关端口 + `domain.ConversationID` + `apperror`
- Produces: `usecase.NewChat(deps ChatDeps) port.ChatUseCase`，

```go
type ChatDeps struct {
	Messages  port.MessageRepo
	Users     port.UserRepo
	Friends   port.FriendRepo
	Groups    port.GroupRepo
	Recent    port.RecentChatStore
	Limits    port.MsgLimitStore
	Presence  port.PresenceStore
	Notifier  port.Notifier
	Clock     port.Clock
	NonFriendLimit int           // 默认 3
	LimitTTL       time.Duration // 默认 24h
}
```

- [ ] **Step 1: 写 SendMessage**

行为（这是全后端最核心的方法，逐步实现）：

```go
func (u *chatUseCase) SendMessage(ctx context.Context, senderID int64, in port.SendMessageInput) (*port.OutgoingMessage, error) {
	// 1. 解析会话 ID
	convID, err := domain.ParseConversationID(in.ConversationID)
	//    err ⇒ apperror.New(CodeConversationInvalid, domain.ErrInvalidConversation.Error())
	// 2. 校验内容：Content 去空格非空；ContentType ∈ {text, image} ⇒ 否则 CodeInvalidParam
	// 3. sender := users.ByID(senderID)（ErrUserNotFound ⇒ CodeUserNotFound）
	// 4. 按类型确定接收者：
	//    private: a, b, _ := convID.PrivateParticipants()
	//             peer := 另一方；若 a、b 都不等于 senderID ⇒ CodeConversationInvalid
	//             ok, _ := friends.Exists(senderID, peer)
	//             若 !ok:
	//                 count := limits.Incr(ctx, convID, senderID, u.d.LimitTTL)
	//                 count > int64(u.d.NonFriendLimit) ⇒ apperror.New(CodeNotFriendLimit, domain.ErrMsgLimitExceeded.Error())
	//             recipients = []int64{peer}
	//    group:   gid := 从 convID.Value 去掉 "g:" 前缀 ParseInt
	//             groups.IsMember(gid, senderID) false ⇒ apperror.New(CodeNotGroupMember, ...)
	//             ids := groups.MemberIDs(gid)；recipients = ids 去掉 senderID
	// 5. msg := &domain.Message{ConversationID: convID, SenderID: senderID,
	//             Content: in.Content, ContentType: domain.ContentType(in.ContentType),
	//             CreatedAt: u.d.Clock.Now()}
	//    messages.Save(ctx, msg) 失败 ⇒ Wrap(CodeInternal, "消息发送失败", err)
	// 6. out := &port.OutgoingMessage{ID: msg.ID, ConversationID: convID.Value,
	//             SenderID: senderID, SenderName: sender.Name, SenderAvatar: sender.Avatar,
	//             Content, ContentType, CreatedAt}
	// 7. 推送 chat.message：notifier.ToUsers(ctx, recipients, {Event:"chat.message", Data: out})
	// 8. 更新最近会话（发送者自己 + 每个接收者都要 Touch）：
	//    summary := port.ConversationSummary{Type: convID.Type, LastContent: 截断 in.Content 到 50 rune,
	//              LastContentType: msg.ContentType, LastSenderID: senderID,
	//              LastSenderName: sender.Name, LastTime: msg.CreatedAt}
	//    recent.SaveSummary(ctx, convID, summary)
	//    all := append(recipients, senderID)
	//    for _, uid := range all { recent.Touch(ctx, uid, convID, msg.CreatedAt) }
	// 9. 推送 chat.recent_updated（Data 为按接收者视角构造的 port.RecentConversation，
	//    private: Peer* = 发送者信息，Online 无意义置 true（发送者刚发过消息）；
	//    group:   Peer* = 群信息（groups.ByID(gid)），Online=false）：
	//    private → notifier.ToUsers(recipients, ...) 加 ToUser(senderID, ...)
	//    group   → 同上
	// 10. return out, nil
}
```

（注释中的每一步都要落成真实代码；`截断到 50 rune` 用 `[]rune(s)` 切片防中文截半。）

- [ ] **Step 2: 写 History**

行为：
1. 权限校验：private ⇒ `a, b, _ := convID.PrivateParticipants()`，`senderID ∉ {a,b}` ⇒ `apperror.New(CodeConversationInvalid, "无权查看该会话")`；group ⇒ 解析 gid，`groups.IsMember` false ⇒ `CodeNotGroupMember`
2. `limit <= 0 || limit > 100` ⇒ limit = 50
3. `messages.History(ctx, convID, cursor, limit)`（cursor 零值 = 从最新）
4. 收集去重 senderIDs → `users.ByIDs` → `map[int64]*domain.User`
5. 映射 `[]port.OutgoingMessage`（SenderName/SenderAvatar 查 map，用户已注销则空串）

- [ ] **Step 3: 写 RecentConversations**

行为：
1. `recent.List(ctx, userID, 50)`
2. 遍历：`recent.Summary(ctx, convID)`（nil 则跳过该会话）
3. private ⇒ `a, b, _ := PrivateParticipants()`，peer = 另一方；`users.ByID(peer)`（NotFound 则 PeerName="已注销用户"）；`presence.IsOnline(peer)`
4. group ⇒ 解析 gid，`groups.ByID(gid)`（NotFound 跳过）；PeerID=gid, PeerName=g.Name, PeerAvatar=g.Avatar, Online=false
5. 组装 `port.RecentConversation{ConversationID: convID.Value, Type: convID.Type.String(), Peer*, Online, LastContent: summary.LastContent, LastSenderName: summary.LastSenderName, LastTime: summary.LastTime}`

- [ ] **Step 4: 验证 + Commit**

Run: `go build ./internal/usecase/... && go vet ./internal/usecase/...` → 成功

```bash
git add server-go/internal/usecase/chat.go
git commit -m "feat(server): usecase——聊天模块（会话键/非好友限流/批量推送/历史/最近会话）"
```

---

### Task B12: WebSocket 适配器（hub/client/升级入口）+ notify 桥接

**Files:**
- Create: `server-go/internal/adapter/ws/protocol.go`
- Create: `server-go/internal/adapter/ws/hub.go`
- Create: `server-go/internal/adapter/ws/client.go`
- Create: `server-go/internal/adapter/ws/handler.go`
- Create: `server-go/internal/adapter/notify/notifier.go`

**Interfaces:**
- Consumes: `port.Notifier/NotifierEvent/PresenceUseCase`、`token.Manager.Parse`、`port.TokenStore`
- Produces:
  - `ws.Envelope{Event string "json:\"event\""; Data any "json:\"data,omitempty\""; Ts int64 "json:\"ts\""}`
  - `ws.NewHub(logger *slog.Logger) *ws.Hub`；`(*Hub) Run(ctx context.Context)`（阻塞事件循环，main 里 `go hub.Run(ctx)`）；`(*Hub) SendToUser(userID int64, env Envelope)`；`(*Hub) SendToUsers(userIDs []int64, env Envelope)`；`(*Hub) Close()`
  - `ws.NewHandler(deps ws.HandlerDeps) echo.HandlerFunc`（升级入口，Task B14 挂到路由），`ws.HandlerDeps{Hub *ws.Hub; TokenMgr *token.Manager; Tokens port.TokenStore; Presence port.PresenceUseCase; Logger *slog.Logger}`
  - `notify.NewWSNotifier(hub *ws.Hub) port.Notifier`

- [ ] **Step 1: 写 protocol.go**

```go
package ws

// 事件名常量（与 spec §3.4 严格一致）
const (
	EventPing             = "ping"
	EventPong             = "pong"
	EventChatMessage      = "chat.message"
	EventRecentUpdated    = "chat.recent_updated"
	EventPresenceChanged  = "presence.changed"
	EventGroupUpdated     = "group.updated"
	EventSystemKick       = "system.kick"
)

// Envelope 是 WS 统一信封。
type Envelope struct {
	Event string `json:"event"`
	Data  any    `json:"data,omitempty"`
	Ts    int64  `json:"ts"`
}

// Incoming C→S 消息结构（当前仅 ping，预留 data）。
type Incoming struct {
	Event string         `json:"event"`
	Data  map[string]any `json:"data"`
}
```

- [ ] **Step 2: 写 hub.go**

```go
package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"
)

// Hub 维护 userID → clients 映射（支持多端），串行处理注册/注销，广播按快照进行。
type Hub struct {
	mu      sync.RWMutex
	clients map[int64]map[*Client]struct{}
	logger  *slog.Logger
	closed  chan struct{}
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients: make(map[int64]map[*Client]struct{}),
		logger:  logger,
		closed:  make(chan struct{}),
	}
}

func (h *Hub) register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.clients[c.UserID] == nil {
		h.clients[c.UserID] = make(map[*Client]struct{})
	}
	h.clients[c.UserID][c] = struct{}{}
}

func (h *Hub) unregister(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if set, ok := h.clients[c.UserID]; ok {
		delete(set, c)
		if len(set) == 0 {
			delete(h.clients, c.UserID)
		}
	}
}

func (h *Hub) SendToUser(userID int64, env Envelope) {
	if env.Ts == 0 {
		env.Ts = time.Now().Unix()
	}
	b, err := json.Marshal(env)
	if err != nil {
		h.logger.Error("ws 信封序列化失败", "event", env.Event, "err", err)
		return
	}
	h.mu.RLock()
	defer h.mu.RUnlock()
	for c := range h.clients[userID] {
		select {
		case c.send <- b:
		default: // 发送缓冲满：丢弃并记录，避免阻塞业务
			h.logger.Warn("ws 发送缓冲满，丢弃消息", "user_id", userID, "event", env.Event)
		}
	}
}

func (h *Hub) SendToUsers(userIDs []int64, env Envelope) {
	for _, id := range userIDs {
		h.SendToUser(id, env)
	}
}

// Kick 向用户所有连接推送 system.kick 并关闭。
func (h *Hub) Kick(userID int64, reason string) {
	h.SendToUser(userID, Envelope{Event: EventSystemKick, Data: map[string]any{"reason": reason}})
	h.mu.RLock()
	targets := make([]*Client, 0)
	for c := range h.clients[userID] {
		targets = append(targets, c)
	}
	h.mu.RUnlock()
	for _, c := range targets {
		c.close()
	}
}

// Close 通知所有连接关闭（优雅停机时调用）。
func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, set := range h.clients {
		for c := range set {
			c.close()
		}
	}
	h.clients = make(map[int64]map[*Client]struct{})
}

// Run 保留为兼容接口：当前 Hub 采用锁+快照模型无需事件循环，Run 阻塞至 ctx 取消后 Close。
func (h *Hub) Run(ctx context.Context) {
	<-ctx.Done()
	h.Close()
}
```

- [ ] **Step 3: 写 client.go**

```go
package ws

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 150 * time.Second             // > 2 个心跳周期（60s*2）
	pingPeriod     = 60 * time.Second
	maxMessageSize = 8192
)

type Client struct {
	UserID int64
	conn   *websocket.Conn
	send   chan []byte
	hub    *Hub
	logger *slog.Logger
	// onPing 每次收到客户端 ping 时回调（用于 presence 心跳续期）；onClose 连接终止时回调一次。
	onPing  func(ctx context.Context, userID int64)
	onClose func(ctx context.Context, userID int64)
	once    sync.Once // close 幂等（import "sync"）
}

func (h *Hub) newClient(userID int64, conn *websocket.Conn,
	onPing, onClose func(context.Context, int64)) *Client {
	return &Client{
		UserID: userID, conn: conn, send: make(chan []byte, 64),
		hub: h, logger: h.logger, onPing: onPing, onClose: onClose,
	}
}

func (c *Client) close() {
	c.once.Do(func() { close(c.send) })
}

// readPump 读循环：收到 ping 事件 → onPing + 回 pong；协议级 ping/pong 由 gorilla 处理。
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister(c)
		c.once.Do(func() { close(c.send) })
		c.conn.Close()
		if c.onClose != nil {
			c.onClose(ctx, c.UserID)
		}
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			return
		}
		var in Incoming
		if err := json.Unmarshal(data, &in); err != nil {
			c.logger.Warn("ws 消息解析失败", "user_id", c.UserID, "err", err)
			continue
		}
		if in.Event == EventPing {
			if c.onPing != nil {
				c.onPing(ctx, c.UserID)
			}
			c.SendToClient(Envelope{Event: EventPong, Ts: time.Now().Unix()})
		}
	}
}

// SendToClient 序列化后入发送通道（缓冲满丢弃并告警）。
func (c *Client) SendToClient(env Envelope) {
	b, err := json.Marshal(env)
	if err != nil {
		return
	}
	select {
	case c.send <- b:
	default:
		c.logger.Warn("ws 客户端缓冲满", "user_id", c.UserID)
	}
}

// writePump 写循环：从 send 通道取消息写出；pingPeriod 定时发协议级 ping。
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() { ticker.Stop(); c.conn.Close() }()
	for {
		select {
		case msg, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage,
					websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return
			}
			if err := c.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
```

注意：`once sync.Once` 字段需要在 import 中加 `"sync"`；writePump 里 `<-c.send` 被 close 后取零值 `ok=false` 分支发 CloseMessage 退出——`hub.unregister` 之后才允许 close(send)，readPump 的 defer 顺序已保证。

- [ ] **Step 4: 写 handler.go（升级入口 + 鉴权）**

```go
package ws

import (
	"context"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"log/slog"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/platform/token"
	"server-go/internal/usecase/port"
)

type HandlerDeps struct {
	Hub      *Hub
	TokenMgr *token.Manager
	Tokens   port.TokenStore
	Presence port.PresenceUseCase
	Logger   *slog.Logger
}

// NewHandler 返回 echo.HandlerFunc：?token= 鉴权 → 升级 → 注册 → 读写循环。
func NewHandler(d HandlerDeps) echo.HandlerFunc {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(r *http.Request) bool { return true }, // 开发环境放开；生产按配置收紧
	}
	return func(c echo.Context) error {
		raw := c.QueryParam("token")
		if raw == "" {
			return echo.NewHTTPError(http.StatusUnauthorized, "missing token")
		}
		claims, err := d.TokenMgr.Parse(raw)
		if err != nil {
			return echo.NewHTTPError(http.StatusUnauthorized, "invalid token")
		}
		ok, err := d.Tokens.Exists(c.Request().Context(), claims.JTI)
		if err != nil || !ok {
			return echo.NewHTTPError(http.StatusUnauthorized, "token expired")
		}
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			d.Logger.Error("ws 升级失败", "err", err)
			return nil
		}
		ctx := c.Request().Context()
		onPing := func(ctx context.Context, uid int64) { _ = d.Presence.Heartbeat(ctx, uid) }
		onClose := func(ctx context.Context, uid int64) { _ = d.Presence.OnDisconnect(ctx, uid) }
		client := d.Hub.newClient(claims.UserID, conn, onPing, onClose)
		d.Hub.register(client)
		if err := d.Presence.OnConnect(ctx, claims.UserID); err != nil {
			d.Logger.Warn("presence 上线失败", "user_id", claims.UserID, "err", err)
		}
		go client.readPump(context.WithoutCancel(ctx))
		go client.writePump()
		return nil
	}
}
```

注意：`middleware` import 若未使用请删除（升级入口自鉴权，不走 HTTP 中间件）。`context.WithoutCancel` 用于读循环生命周期长于请求 ctx 的场景（Go 1.21+）。

- [ ] **Step 5: 写 notify/notifier.go（port.Notifier 实现）**

```go
package notify

import (
	"context"
	"time"

	"server-go/internal/adapter/ws"
	"server-go/internal/usecase/port"
)

type WSNotifier struct{ hub *ws.Hub }

func NewWSNotifier(hub *ws.Hub) *WSNotifier { return &WSNotifier{hub: hub} }

func (n *WSNotifier) ToUser(_ context.Context, userID int64, ev port.NotifierEvent) {
	n.hub.SendToUser(userID, ws.Envelope{Event: ev.Event, Data: ev.Data, Ts: time.Now().Unix()})
}

func (n *WSNotifier) ToUsers(_ context.Context, userIDs []int64, ev port.NotifierEvent) {
	n.hub.SendToUsers(userIDs, ws.Envelope{Event: ev.Event, Data: ev.Data, Ts: time.Now().Unix()})
}

var _ port.Notifier = (*WSNotifier)(nil)
```

- [ ] **Step 6: 验证 + Commit**

Run: `go build ./internal/adapter/... && go vet ./internal/adapter/...` → 成功

```bash
git add server-go/internal/adapter/ws server-go/internal/adapter/notify
git commit -m "feat(server): WebSocket 适配器（hub/client/鉴权升级）与 Notifier 桥接"
```

---

### Task B13: HTTP 适配器——DTO、统一响应与中间件

**Files:**
- Create: `server-go/internal/adapter/http/dto/request.go`
- Create: `server-go/internal/adapter/http/dto/response.go`
- Create: `server-go/internal/adapter/http/middleware/error.go`
- Create: `server-go/internal/adapter/http/middleware/auth.go`
- Create: `server-go/internal/adapter/http/middleware/logging.go`
- Modify: `server-go/internal/adapter/ws/handler.go`（如有 middleware 循环 import：ws 不 import middleware，见 Step 4 注意事项）

**Interfaces:**
- Produces:
  - `dto.Response{Code int "json:\"code\""; Message string "json:\"message\""; Data any "json:\"data,omitempty\""; RequestID string "json:\"request_id,omitempty\""}`
  - `dto.OK(c echo.Context, data any) error`（code=0, message="ok", request_id 取 `c.Response().Header().Get(echo.HeaderXRequestID)`）
  - 请求 DTO（validator tag 见 Step 1）
  - `middleware.ErrorHandler(logger *slog.Logger) echo.MiddlewareFunc`
  - `middleware.JWTAuth(tokenMgr *token.Manager, tokens port.TokenStore) echo.MiddlewareFunc`
  - `middleware.UserIDFromContext(ctx context.Context) (int64, bool)`、`middleware.JTIFromContext(ctx context.Context) (string, bool)`
  - `middleware.RequestLogger(logger *slog.Logger) echo.MiddlewareFunc`
  - `middleware.Validate(err error) error`（validator 错误 → apperror 10001）

- [ ] **Step 1: 写 dto/request.go（validator tag 全部就位）**

```go
package dto

type RegisterReq struct {
	Identity  string `json:"identity" validate:"required,min=3,max=32"`
	Name      string `json:"name" validate:"required,min=1,max=64"`
	Password  string `json:"password" validate:"required,min=6,max=64"`
	Email     string `json:"email" validate:"required,email,max=128"`
	EmailCode string `json:"email_code" validate:"required,len=6"`
}

type LoginReq struct {
	Account       string `json:"account" validate:"required"`
	Password      string `json:"password" validate:"required"`
	CaptchaID     string `json:"captcha_id" validate:"required"`
	CaptchaAnswer string `json:"captcha_answer" validate:"required"`
}

type EmailCodeReq struct {
	Email string `json:"email" validate:"required,email"`
}

type UpdateUserReq struct {
	Name      *string `json:"name" validate:"omitempty,min=1,max=64"`
	Phone     *string `json:"phone" validate:"omitempty,max=32"`
	Email     *string `json:"email" validate:"omitempty,email"`
	EmailCode string  `json:"email_code" validate:"omitempty,len=6"`
	Avatar    *string `json:"avatar" validate:"omitempty,max=255"`
	Signature *string `json:"signature" validate:"omitempty,max=255"`
	BirthDate *string `json:"birth_date"` // "2006-01-02"，handler 解析为 time.Time
}

type AddFriendReq struct {
	UserID int64 `json:"user_id" validate:"required"`
}

type CreateGroupReq struct {
	Name      string  `json:"name" validate:"required,min=1,max=64"`
	MemberIDs []int64 `json:"member_ids"`
}

type SendMessageReq struct {
	ConversationID string `json:"conversation_id" validate:"required"`
	Content        string `json:"content" validate:"required,min=1,max=4096"`
	ContentType    string `json:"content_type" validate:"required,oneof=text image"`
}

type SearchQuery struct {
	Keyword string `query:"search" validate:"omitempty,max=64"`
}

type HistoryQuery struct {
	Cursor string `query:"cursor"` // RFC3339，可空
	Limit  int    `query:"limit" validate:"omitempty,min=1,max=100"`
}
```

- [ ] **Step 2: 写 dto/response.go + 校验辅助**

```go
package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"server-go/internal/usecase/apperror"
)

type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func OK(c echo.Context, data any) error {
	return c.JSON(200, Response{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

var validate = validator.New(validator.WithRequiredStructEnabled())

// Validate 结构体校验；失败 → apperror 10001，message 形如 "字段xx: 校验规则"。
func Validate(v any) error {
	if err := validate.Struct(v); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			f := ve[0]
			return apperror.New(apperror.CodeInvalidParam,
				"参数 "+f.Field()+" 校验失败("+f.Tag()+")")
		}
		return apperror.New(apperror.CodeInvalidParam, "参数格式错误")
	}
	return nil
}
```

- [ ] **Step 3: 写 middleware/error.go（统一错误处理，核心映射表）**

```go
package middleware

import (
	"errors"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
)

// codeToStatus 业务码 → HTTP 状态码（spec §3.2）。
var codeToStatus = map[int]int{
	apperror.CodeInvalidParam:        http.StatusBadRequest,
	apperror.CodeInternal:            http.StatusInternalServerError,
	apperror.CodeCaptchaInvalid:      http.StatusBadRequest,
	apperror.CodeBadCredentials:      http.StatusUnauthorized,
	apperror.CodeTokenInvalid:        http.StatusUnauthorized,
	apperror.CodeUserBanned:          http.StatusForbidden,
	apperror.CodeEmailTaken:          http.StatusConflict,
	apperror.CodeIdentityTaken:       http.StatusConflict,
	apperror.CodeUserNotFound:        http.StatusNotFound,
	apperror.CodeEmailCodeInvalid:    http.StatusBadRequest,
	apperror.CodeAlreadyFriend:       http.StatusConflict,
	apperror.CodeNotFriendLimit:      http.StatusTooManyRequests,
	apperror.CodeGroupNotFound:       http.StatusNotFound,
	apperror.CodeNotGroupMember:      http.StatusForbidden,
	apperror.CodeConversationInvalid: http.StatusBadRequest,
}

// domainSentinelToCode 兜底：未包装的 domain 哨兵错误映射。
var domainSentinelToCode = map[error]int{
	domain.ErrUserNotFound:        apperror.CodeUserNotFound,
	domain.ErrEmailTaken:          apperror.CodeEmailTaken,
	domain.ErrIdentityTaken:       apperror.CodeIdentityTaken,
	domain.ErrBadCredentials:      apperror.CodeBadCredentials,
	domain.ErrUserBanned:          apperror.CodeUserBanned,
	domain.ErrCaptchaInvalid:      apperror.CodeCaptchaInvalid,
	domain.ErrEmailCodeInvalid:    apperror.CodeEmailCodeInvalid,
	domain.ErrGroupNotFound:       apperror.CodeGroupNotFound,
	domain.ErrNotGroupMember:      apperror.CodeNotGroupMember,
	domain.ErrMsgLimitExceeded:    apperror.CodeNotFriendLimit,
	domain.ErrTokenInvalid:        apperror.CodeTokenInvalid,
	domain.ErrInvalidConversation: apperror.CodeConversationInvalid,
}

func ErrorHandler(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			err := next(c)
			if err == nil {
				return nil
			}
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			code, message, status := apperror.CodeInternal, "服务器内部错误", http.StatusInternalServerError

			if ae, ok := apperror.From(err); ok {
				code, message = ae.Code, ae.Message
				if s, ok2 := codeToStatus[ae.Code]; ok2 {
					status = s
				}
				logger.Warn("业务错误", "code", code, "msg", message,
					"request_id", reqID, "path", c.Path(), "err", ae.Err)
			} else if he, ok := err.(*echo.HTTPError); ok { // 框架错误（404/405/401等）
				status = he.Code
				code = status * 100 // 例 404 → 40400，前端仅按非 0 处理
				message, _ = he.Message.(string)
			} else {
				for sentinel, c2 := range domainSentinelToCode {
					if errors.Is(err, sentinel) {
						code, message = c2, sentinel.Error()
						if s, ok2 := codeToStatus[c2]; ok2 {
							status = s
						}
						break
					}
				}
				if code == apperror.CodeInternal {
					logger.Error("未处理错误", "request_id", reqID, "path", c.Path(),
						"err", err, "stack", string(debug.Stack()))
				}
			}
			if c.Response().Committed {
				return nil
			}
			return c.JSON(status, dto.Response{Code: code, Message: message, RequestID: reqID})
		}
	}
}
```

- [ ] **Step 4: 写 middleware/auth.go + middleware/logging.go**

```go
package middleware

import (
	"context"
	"strings"

	"github.com/labstack/echo/v4"
	"server-go/internal/platform/token"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type ctxKey struct{ name string }

var (
	userIDKey = &ctxKey{"uid"}
	jtiKey    = &ctxKey{"jti"}
)

// JWTAuth 解析 Authorization: Bearer <token>，校验签名 + Redis 白名单，
// 将 uid/jti 写入请求 context（类型化私有 key）。
func JWTAuth(tokenMgr *token.Manager, tokens port.TokenStore) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			raw := strings.TrimPrefix(header, "Bearer ")
			if raw == "" || raw == header {
				return apperror.New(apperror.CodeTokenInvalid, "缺少 Authorization 头")
			}
			claims, err := tokenMgr.Parse(raw)
			if err != nil {
				return apperror.New(apperror.CodeTokenInvalid, "token 无效或已过期")
			}
			ok, err := tokens.Exists(c.Request().Context(), claims.JTI)
			if err != nil || !ok {
				return apperror.New(apperror.CodeTokenInvalid, "登录状态已失效，请重新登录")
			}
			ctx := context.WithValue(c.Request().Context(), userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, jtiKey, claims.JTI)
			c.SetRequest(ctx)
			return next(c)
		}
	}
}

func UserIDFromContext(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(userIDKey).(int64)
	return v, ok
}

func JTIFromContext(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(jtiKey).(string)
	return v, ok
}
```

```go
package middleware

import (
	"log/slog"
	"runtime/debug"
	"time"

	"github.com/labstack/echo/v4"
	"server-go/internal/usecase/apperror"
)

// RequestLogger 访问日志：method/path/status/latency/request_id/ip。
func RequestLogger(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			reqID := c.Response().Header().Get(echo.HeaderXRequestID)
			logger.Info("http",
				"method", c.Request().Method,
				"path", c.Request().URL.Path,
				"status", c.Response().Status,
				"latency_ms", time.Since(start).Milliseconds(),
				"request_id", reqID,
				"ip", c.RealIP(),
			)
			return err
		}
	}
}

// Recover 捕获 panic → 记录堆栈 → 交给 ErrorHandler 统一渲染 500。
func Recover(logger *slog.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if r := recover(); r != nil {
					logger.Error("panic recovered",
						"panic", r, "stack", string(debug.Stack()),
						"path", c.Request().URL.Path)
					c.Error(apperror.New(apperror.CodeInternal, "服务器内部错误"))
				}
			}()
			return next(c)
		}
	}
}
```

（request_id 由 echo 内置 `echoMiddleware.RequestID()` 生成，路由装配时置于链条最前，见 B14 server.go。）

- [ ] **Step 5: 验证 + Commit**

Run: `go build ./internal/adapter/http/... && go vet ./internal/adapter/http/...` → 成功

```bash
git add server-go/internal/adapter/http
git commit -m "feat(server): HTTP 适配器——DTO/统一响应/错误处理/JWT鉴权/日志中间件"
```

---

### Task B14: HTTP handlers 与路由装配

**Files:**
- Create: `server-go/internal/adapter/http/handler/auth.go`
- Create: `server-go/internal/adapter/http/handler/user.go`
- Create: `server-go/internal/adapter/http/handler/friend.go`
- Create: `server-go/internal/adapter/http/handler/group.go`
- Create: `server-go/internal/adapter/http/handler/chat.go`
- Create: `server-go/internal/adapter/http/handler/file.go`
- Create: `server-go/internal/adapter/http/server.go`

**Interfaces:**
- Consumes: `port.*UseCase` 接口、`dto.*`、`middleware.*`、`ws.NewHandler`
- Produces:
  - `handler.NewAuth(uc port.AuthUseCase) *handler.Auth`（方法 `Register/Login/Logout/Captcha/EmailCode`，均为 `func(echo.Context) error`）
  - `handler.NewUser(uc port.UserUseCase) *handler.User`（`Me/UpdateMe/Detail/Search`）
  - `handler.NewFriend(uc port.FriendUseCase) *handler.Friend`（`Add/List`）
  - `handler.NewGroup(uc port.GroupUseCase) *handler.Group`（`Create/Join/List/Search`）
  - `handler.NewChat(uc port.ChatUseCase) *handler.Chat`（`SendMessage/History/Conversations`）
  - `handler.NewFile(staticPath string, logger *slog.Logger) *handler.File`（`UploadImage`）
  - `http.NewServer(deps http.ServerDeps) *echo.Echo`，

```go
type ServerDeps struct {
	Cfg      *config.Config
	Logger   *slog.Logger
	TokenMgr *token.Manager
	Tokens   port.TokenStore
	Auth     port.AuthUseCase
	User     port.UserUseCase
	Friend   port.FriendUseCase
	Group    port.GroupUseCase
	Chat     port.ChatUseCase
	Presence port.PresenceUseCase
	Hub      *ws.Hub
}
```

- [ ] **Step 1: 写 handler/auth.go（完整示范模板，其余 handler 同构）**

```go
package handler

import (
	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type Auth struct {
	uc port.AuthUseCase
}

func NewAuth(uc port.AuthUseCase) *Auth { return &Auth{uc: uc} }

func (h *Auth) Register(c echo.Context) error {
	var req dto.RegisterReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	err := h.uc.Register(c.Request().Context(), port.RegisterInput{
		Identity: req.Identity, Name: req.Name, Password: req.Password,
		Email: req.Email, EmailCode: req.EmailCode,
	})
	if err != nil {
		return err // ErrorHandler 统一渲染
	}
	return dto.OK(c, nil)
}

func (h *Auth) Login(c echo.Context) error {
	var req dto.LoginReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	out, err := h.uc.Login(c.Request().Context(), port.LoginInput{
		Account: req.Account, Password: req.Password,
		CaptchaID: req.CaptchaID, CaptchaAnswer: req.CaptchaAnswer,
	})
	if err != nil {
		return err
	}
	return dto.OK(c, out)
}

func (h *Auth) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	jti, _ := middleware.JTIFromContext(ctx)
	uid, _ := middleware.UserIDFromContext(ctx)
	if err := h.uc.Logout(ctx, jti, uid); err != nil {
		return err
	}
	return dto.OK(c, nil)
}

func (h *Auth) Captcha(c echo.Context) error {
	id, b64, err := h.uc.GenerateCaptcha(c.Request().Context())
	if err != nil {
		return err
	}
	return dto.OK(c, map[string]string{"captcha_id": id, "image": b64})
}

func (h *Auth) EmailCode(c echo.Context) error {
	var req dto.EmailCodeReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	if err := h.uc.SendEmailCode(c.Request().Context(), req.Email); err != nil {
		return err
	}
	return dto.OK(c, nil)
}
```

- [ ] **Step 2: 写其余 handler（模板同上：Bind → Validate → usecase → OK/return err）**

- `user.go`：
  - `Me`：`uid, _ := middleware.UserIDFromContext(c.Request().Context())`；`uc.Me(ctx, uid)`
  - `UpdateMe`：绑定 `dto.UpdateUserReq`；`BirthDate *string` 按 `"2006-01-02"` 解析为 `*time.Time`（解析失败 → 10001）；组装 `port.UpdateUserInput`（Email 非 nil 时要求 EmailCode 非空，否则 10001"修改邮箱需要验证码"）
  - `Detail`：`id, err := strconv.ParseInt(c.Param("id"), 10, 64)`，err → 10001；viewerID 从 ctx 取（**该路由挂 JWTAuth**，spec 中 GET /users/{id} 未标 🔒——为支持"未登录看公开资料"不挂鉴权：viewerID 取不到时用 0，`uc.Detail` 已按 0 处理为非好友视角）。**决定：Detail 路由不挂 JWTAuth，viewerID 可选**；但 JWTAuth 是路由级中间件无法"可选"，因此 Detail 挂一个**宽松鉴权**中间件 `middleware.OptionalJWT(tokenMgr, tokens)`：token 有效则写入 ctx，无效/缺失不报错直接放行。在 middleware/auth.go 中补充实现（逻辑与 JWTAuth 相同，但所有失败路径 `return next(c)`）。
  - `Search`：绑定 `dto.SearchQuery`；`uc.Search(ctx, uid, req.Keyword)`
- `friend.go`：`Add`（绑定 AddFriendReq → `uc.Add(ctx, uid, req.UserID)`）、`List`（`uc.List(ctx, uid)`）
- `group.go`：`Create`（绑定 CreateGroupReq → `uc.Create(ctx, uid, port.CreateGroupInput{...})`）、`Join`（`gid := c.Param("id")` ParseInt → `uc.Join(ctx, gid, uid)`）、`List`、`Search`（SearchQuery）
- `chat.go`：
  - `SendMessage`：绑定 SendMessageReq → `uc.SendMessage(ctx, uid, port.SendMessageInput{...})`
  - `History`：`convRaw := c.Param("id")` → `domain.ParseConversationID`（err → 60001）；绑定 HistoryQuery；cursor 非空按 RFC3339 解析（err → 10001）→ `uc.History(ctx, uid, convID, cursor, limit)`
  - `Conversations`：`uc.RecentConversations(ctx, uid)`
- `file.go` `UploadImage`：
  1. `file, err := c.FormFile("file")`（err → 10001"缺少文件字段 file"）
  2. 校验：`file.Size <= 5<<20`（→ 10001"图片不能超过 5MB"）；扩展名 ∈ {.jpg,.jpeg,.png,.gif,.webp}（小写化比较，→ 10001"仅支持 jpg/png/gif/webp"）
  3. 目标路径 `{staticPath}/uploads/{yyyyMM}/{uuid}{ext}`；`os.MkdirAll(dir, 0o755)`；`c.Save(file, dst)`
  4. `dto.OK(c, map[string]string{"url": "/{staticPath}/uploads/{yyyyMM}/{uuid}{ext}"})`

- [ ] **Step 3: 写 server.go（路由装配，与 spec §3.1 表逐行对应）**

```go
package http

import (
	"github.com/labstack/echo/v4"
	echoMiddleware "github.com/labstack/echo/v4/middleware"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/handler"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/adapter/ws"
	"server-go/internal/config"
	// ... port/token/ws 等 import 按 ServerDeps 需要补全
)

func NewServer(d ServerDeps) *echo.Echo {
	e := echo.New()
	e.HideBanner = true
	if d.Cfg.Server.Mode == "release" {
		e.Debug = false
	}

	// 中间件链（顺序敏感）：RequestID → Recover → 访问日志 → CORS → 统一错误处理
	e.Use(echoMiddleware.RequestID())
	e.Use(middleware.Recover(d.Logger))
	e.Use(middleware.RequestLogger(d.Logger))
	e.Use(echoMiddleware.CORSWithConfig(echoMiddleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "Content-Type"},
	}))
	e.Use(middleware.ErrorHandler(d.Logger))

	// 静态资源（上传的图片）
	e.Static("/"+d.Cfg.Server.StaticPath, "./"+d.Cfg.Server.StaticPath)

	authH := handler.NewAuth(d.Auth)
	userH := handler.NewUser(d.User)
	friendH := handler.NewFriend(d.Friend)
	groupH := handler.NewGroup(d.Group)
	chatH := handler.NewChat(d.Chat)
	fileH := handler.NewFile(d.Cfg.Server.StaticPath, d.Logger)
	jwtMW := middleware.JWTAuth(d.TokenMgr, d.Tokens)

	v1 := e.Group("/api/v1")
	// 认证（公开）
	v1.POST("/auth/register", authH.Register)
	v1.POST("/auth/login", authH.Login)
	v1.POST("/auth/captcha", authH.Captcha)
	v1.POST("/auth/email-code", authH.EmailCode)
	v1.POST("/auth/logout", authH.Logout, jwtMW)

	// 用户
	v1.GET("/users/me", userH.Me, jwtMW)
	v1.PATCH("/users/me", userH.UpdateMe, jwtMW)
	v1.GET("/users/:id", userH.Detail, middleware.OptionalJWT(d.TokenMgr, d.Tokens))
	v1.GET("/users", userH.Search, jwtMW)

	// 好友
	v1.POST("/friends", friendH.Add, jwtMW)
	v1.GET("/friends", friendH.List, jwtMW)

	// 群组
	v1.POST("/groups", groupH.Create, jwtMW)
	v1.POST("/groups/:id/members", groupH.Join, jwtMW)
	v1.GET("/groups", groupH.List, jwtMW) // 带 ?search=xx 时内部走搜索（spec §3.1）

	// 聊天
	v1.POST("/messages", chatH.SendMessage, jwtMW)
	v1.GET("/conversations", chatH.Conversations, jwtMW)
	v1.GET("/conversations/:id/messages", chatH.History, jwtMW)

	// 文件
	v1.POST("/files/images", fileH.UploadImage, jwtMW)

	// WebSocket（自鉴权，见 ws.HandlerDeps；hub 引用由 main 装配时通过 d 传入，
	// ServerDeps 需增加 Hub *ws.Hub 字段）
	v1.GET("/ws", ws.NewHandler(ws.HandlerDeps{
		Hub: d.Hub, TokenMgr: d.TokenMgr, Tokens: d.Tokens,
		Presence: d.Presence, Logger: d.Logger,
	}))

	// 健康检查
	v1.GET("/healthz", func(c echo.Context) error { return dto.OK(c, map[string]string{"status": "up"}) })
	return e
}
```

注意：Echo 的路由级中间件支持"后置参数"写法（`v1.POST(path, h, mw...)` 中 mw 声明在 handler 之后是 Echo 的正确语法）。`groupH.List` 内部按 `c.QueryParam("search")` 分流：非空 → `uc.Search(ctx, search)`，空 → `uc.List(ctx, uid)`（spec §3.1 的 `GET /groups?search=xx` 与 `GET /groups` 共用一个 handler）。`groupH.Search` 方法保留（供 List 内部复用逻辑），不单独挂路由。

- [ ] **Step 4: 验证 + Commit**

Run: `go build ./internal/adapter/... && go vet ./internal/adapter/...` → 成功

```bash
git add server-go/internal/adapter/http
git commit -m "feat(server): HTTP handlers 与路由装配（RESTful /api/v1）"
```

---

### Task B15: cmd/server 组装根与优雅关闭

**Files:**
- Create: `server-go/cmd/server/main.go`
- Create: `server-go/config.toml`（从 config.example.toml 复制，填本地真实值；加入 .gitignore）
- Modify: `server-go/.gitignore`（追加 `config.toml`；若仓库根已有 .gitignore 则在 server-go 下新建）

**Interfaces:**
- Consumes: 之前全部任务的构造函数
- Produces: 可执行入口 `go run ./cmd/server --config config.toml`

- [ ] **Step 1: 写 cmd/server/main.go（完整装配顺序）**

```go
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"server-go/internal/adapter/notify"
	httppkg "server-go/internal/adapter/http"
	"server-go/internal/adapter/persistence/mongo"
	"server-go/internal/adapter/persistence/mysql"
	"server-go/internal/adapter/persistence/redis"
	"server-go/internal/adapter/ws"
	"server-go/internal/config"
	"server-go/internal/platform/clock"
	"server-go/internal/platform/captchagen"
	"server-go/internal/platform/logger"
	"server-go/internal/platform/mailer"
	"server-go/internal/platform/token"
	"server-go/internal/usecase"
	"server-go/internal/usecase/apperror"
	_ "server-go/internal/usecase/port"
)

func main() {
	configPath := flag.String("config", "config.toml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log := logger.New(cfg.Log.Level, cfg.Log.Pretty)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ---- 基础设施 ----
	db, err := mysql.NewDB(ctx, cfg.MySQL.DSN())
	must(err, "mysql", log)
	defer db.Close()

	rdb, err := redis.NewClient(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	must(err, "redis", log)
	defer rdb.Close()

	mongoCli, err := mongo.NewClient(ctx, cfg.Mongo.URI)
	must(err, "mongo", log)
	defer func() {
		dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoCli.Disconnect(dctx)
	}()
	mongoDB := mongoCli.Database(cfg.Mongo.Database)

	// ---- 出站适配器 ----
	userRepo := mysql.NewUserRepo(db)
	friendRepo := mysql.NewFriendRepo(db)
	groupRepo := mysql.NewGroupRepo(db)
	msgRepo := mongo.NewMessageRepo(mongoDB, cfg.Message.RetentionDays)
	must(msgRepo.EnsureIndexes(ctx), "mongo 索引", log)

	jwtMgr := token.NewManager(cfg.JWT.Secret, cfg.JWT.TTL())
	tokenStore := redis.NewTokenStore(rdb, cfg.JWT.TTL())
	captchaStore := redis.NewCaptchaStore(rdb)
	codeStore := redis.NewEmailCodeStore(rdb)
	presenceStore := redis.NewPresenceStore(rdb, time.Duration(cfg.Presence.TTLSeconds)*time.Second)
	recentStore := redis.NewRecentChatStore(rdb)
	limitStore := redis.NewMsgLimitStore(rdb)
	mail := mailer.New(cfg.Email.Host, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, cfg.Email.From)
	captchaGen := captchagen.New()
	clk := clock.Real{}

	hub := ws.NewHub(log)
	go hub.Run(ctx)
	notifier := notify.NewWSNotifier(hub)

	// ---- usecase ----
	authUC := usecase.NewAuth(usecase.AuthDeps{
		Users: userRepo, Tokens: tokenStore, Issuer: jwtMgr,
		Captchas: captchaStore, CaptchaGen: captchaGen, Codes: codeStore,
		Mailer: mail, Notifier: notifier, Clock: clk, TokenTTL: cfg.JWT.TTL(), Logger: log,
	})
	userUC := usecase.NewUser(usecase.UserDeps{Users: userRepo, Friends: friendRepo, Codes: codeStore})
	friendUC := usecase.NewFriend(usecase.FriendDeps{Users: userRepo, Friends: friendRepo, Presence: presenceStore})
	groupUC := usecase.NewGroup(usecase.GroupDeps{Groups: groupRepo, Users: userRepo, Notifier: notifier, Clock: clk})
	presenceUC := usecase.NewPresence(usecase.PresenceDeps{
		Presence: presenceStore, Friends: friendRepo, Notifier: notifier,
		TTL: time.Duration(cfg.Presence.TTLSeconds) * time.Second,
	})
	chatUC := usecase.NewChat(usecase.ChatDeps{
		Messages: msgRepo, Users: userRepo, Friends: friendRepo, Groups: groupRepo,
		Recent: recentStore, Limits: limitStore, Presence: presenceStore,
		Notifier: notifier, Clock: clk,
		NonFriendLimit: cfg.Message.NonFriendLimit,
		LimitTTL:       time.Duration(cfg.Message.LimitTTLHours) * time.Hour,
	})

	// ---- HTTP 服务器 ----
	e := httppkg.NewServer(httppkg.ServerDeps{
		Cfg: cfg, Logger: log, TokenMgr: jwtMgr, Tokens: tokenStore,
		Auth: authUC, User: userUC, Friend: friendUC, Group: groupUC,
		Chat: chatUC, Presence: presenceUC, Hub: hub,
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: e}
	go func() {
		log.Info("HTTP 服务启动", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP 服务异常退出", "err", err)
			stop()
		}
	}()

	// ---- 优雅关闭 ----
	<-ctx.Done()
	log.Info("收到退出信号，开始优雅关闭…")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error("HTTP 关闭超时", "err", err)
	}
	hub.Close()
	log.Info("服务已退出")
}

func must(err error, what string, log *slog.Logger) {
	if err != nil {
		log.Error("初始化失败", "component", what, "err", err)
		os.Exit(1)
	}
}
```

注意：persistence 三个适配器目录的包名分别是 `mysql`、`redisx`、`mongox`（import 目录路径后以各自包名调用）。`internal/adapter/http` 的包名保持 `http` 即可——main 中用 import 别名 `httppkg "server-go/internal/adapter/http"` 规避与标准库 `net/http` 的冲突（上方代码已如此）。`adapter/ws` → `adapter/http/middleware` 无循环依赖（middleware 不依赖 ws）。main 中若有未使用的 import（如 `apperror`）删除。

- [ ] **Step 2: 准备 config.toml 与 .gitignore**

```bash
cd server-go && cp config.example.toml config.toml
printf 'config.toml\nstatic/uploads/\n' > .gitignore
```

编辑 config.toml 填本地 MySQL 密码等真实值（该文件不入库）。

- [ ] **Step 3: 建库并验证启动**

```bash
mysql -uroot -p < migrations/schema.sql   # 需要本地 MySQL；密码按本机情况
go build ./cmd/server && go vet ./...
./server --config config.toml              # 本地需已启动 MySQL/Redis/Mongo
curl -s localhost:8080/api/v1/healthz
```

Expected: `{"code":0,"message":"ok","data":{"status":"up"},"request_id":"..."}`。验证后 Ctrl+C 应打印"服务已退出"。

- [ ] **Step 4: Commit**

```bash
git add server-go/cmd server-go/.gitignore server-go/internal
git commit -m "feat(server): cmd/server 组装根——手动 DI 与优雅关闭"
```

---

### Task B16: 删除旧代码、依赖清理与后端冒烟验证

**Files:**
- Delete: `server-go/main.go`、`server-go/main_test.go`、`server-go/internal/app/`、`server-go/internal/router/`、`server-go/internal/global/`、`server-go/internal/consts/`、`server-go/internal/websocket/`、`server-go/common/`
- Delete: `server-go/mainifest/`、`server-go/resource/`（先检查其中是否有仍被引用的运行期资源，如邮件模板/静态文件；有则移入 `server-go/static/` 再删）
- Modify: `server-go/go.mod`（`go mod tidy` 后 gin/gorm/mapstructure/zerolog/satori-uuid/cron 等旧依赖自动移除）
- Modify: `README.md`（仓库根：更新后端启动方式为 `go run ./cmd/server --config config.toml`，注明需先执行 migrations/schema.sql）

- [ ] **Step 1: 删除旧代码**

```bash
cd server-go
git rm -r main.go main_test.go internal/app internal/router internal/global internal/consts internal/websocket common
ls mainifest resource  # 检查内容，按需搬运后 git rm -r
```

- [ ] **Step 2: 依赖清理与全量验证**

```bash
go mod tidy
go build ./... && go vet ./...
gofmt -l . | grep -v sqlcgen   # 输出应为空
```

- [ ] **Step 3: 冒烟验证（需本地 MySQL/Redis/Mongo 已启动且 schema 已建）**

```bash
go run ./cmd/server --config config.toml &
# 1. 图形验证码
curl -s -X POST localhost:8080/api/v1/auth/captcha | head -c 200
# 2. 注册（验证码流程依赖邮件配置；无邮件环境可先在 redis-cli SET verify:email:test@test.com 123456 手工注入）
curl -s -X POST localhost:8080/api/v1/auth/register -H 'Content-Type: application/json' \
  -d '{"identity":"kk001","name":"KK","password":"123456","email":"test@test.com","email_code":"123456"}'
# 3. 登录（captcha_id/answer 用第 1 步返回值）
curl -s -X POST localhost:8080/api/v1/auth/login -H 'Content-Type: application/json' \
  -d '{"account":"kk001","password":"123456","captcha_id":"<id>","captcha_answer":"<answer>"}'
# 4. 带 token 访问
TOKEN=<上一步 data.token>
curl -s localhost:8080/api/v1/users/me -H "Authorization: Bearer $TOKEN"
curl -s localhost:8080/api/v1/conversations -H "Authorization: Bearer $TOKEN"
```

Expected: 各步返回 `{"code":0,...}`；错误路径（错密码/无 token）返回对应业务码与 401/400 状态。

- [ ] **Step 4: Commit**

```bash
git add -A server-go README.md
git commit -m "refactor(server)!: 移除 gin/gorm 旧实现，后端切换至 Clean Architecture 新架构"
```

---

## 后端完成标准

- `go build ./... && go vet ./... && gofmt -l`（除 sqlcgen）全部干净
- B16 冒烟清单全通过
- 旧目录（internal/app、internal/websocket、internal/global、internal/consts、common、main.go、main_test.go、mainifest、resource）不复存在
- go.mod 无 gin/gorm/zerolog 等旧依赖
