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
	Pretty bool   `toml:"pretty"` // true=文本输出，false=JSON
}

// Load 从 path 加载 TOML 配置，随后应用环境变量覆盖与默认值。
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
