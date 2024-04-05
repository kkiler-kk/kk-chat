package config

import (
	"github.com/BurntSushi/toml"
	"github.com/rs/zerolog/log"
	"time"
)

type Conf struct {
	Server   Server   `toml:"server"`
	Logger   Logger   `toml:"logger"`
	Database Database `toml:"database"`
	Token    Token    `toml:"token"`
	Redis    Redis    `toml:"redis"`
	Mongo    Mongo    `toml:"mongo"`
	Mail     Mail     `toml:"mail"`
	System   System   `toml:"system"`
}

type Server struct {
	Port              int    `toml:"port"`
	ServerRoot        string `toml:"serverRoot"`
	NameToUriType     uint   `toml:"nameToUriType"`
	MaxHeaderBytes    string `toml:"maxHeaderBytes"`
	ClientMaxBodySize string `toml:"clientMaxBodySize"`
	LogPath           string `toml:"logPath"`
	LogStdout         bool   `toml:"logStdout"`
	ErrorStack        bool   `toml:"errorStack"`
	ErrorLogEnabled   bool   `toml:"errorLogEnabled"`
	ErrorLogPattern   string `toml:"errorLogPattern"`
	AccessLogEnabled  bool   `toml:"accessLogEnabled"`
	AccessLogPattern  string `toml:"accessLogPattern"`
}
type Mail struct {
	Host     string `toml:"host"`
	UserName string `toml:"username"`
	Password string `toml:"password"`
	Port     int    `toml:"port"`
}

type Logger struct {
	Path   string `toml:"path"`
	File   string `toml:"file"`
	Level  string `toml:"level"`
	Stdout bool   `toml:"stdout"`
}
type Database struct {
	DatabaseLogger  DatabaseLogger  `toml:"logger"`
	DatabaseDefault DatabaseDefault `toml:"default"`
}
type DatabaseLogger struct {
	Level  string `toml:"level"`
	Stdout bool   `toml:"stdout"`
	Path   string `toml:"path"`
}
type DatabaseDefault struct {
	Host        string `toml:"host"`
	Port        uint   `toml:"port"`
	User        string `toml:"user"`
	Pass        string `toml:"pass"`
	Name        string `toml:"chat"`
	Link        string `toml:"link"`
	MaxIdle     uint   `toml:"maxIdle"`
	MaxOpen     uint   `toml:"maxOpen"`
	MaxLifetime string `toml:"maxLifetime"`
}

type Token struct {
	CacheKey    string        `toml:"cacheKey"`
	TimeOut     time.Duration `toml:"timeOut"`
	MaxRefresh  uint          `toml:"maxRefresh"`
	MultilLogin bool          `toml:"multilLogin"`
	EncryptKey  string        `toml:"encryptKey"`
	CacheModel  string        `toml:"cacheModel"`
}

type Redis struct {
	Host            string `toml:"host"`
	Port            uint   `toml:"port"`
	DB              int    `toml:"db"`
	Pass            string `toml:"pass"`
	IdleTimeout     string `toml:"idleTimeout"`
	MaxConnLifetime string `toml:"maxConnLifetime"`
	WaitTimeout     string `toml:"waitTimeout"`
	DialTimeout     string `toml:"dialTimeout"`
	ReadTimeout     string `toml:"readTimeout"`
	WriteTimeout    string `toml:"writeTimeout"`
	MaxActive       uint   `toml:"maxActive"`
}
type Mongo struct {
	Host        string  `toml:"host"`
	Port        uint    `toml:"port"`
	Pass        string  `toml:"pass"`
	MaxPoolSize *uint64 `toml:"maxPoolSize"`
	MinPoolSize *uint64 `toml:"minPoolSize"`
}
type System struct {
	NotCheckAuthAdminIds []uint      `toml:"notCheckAuthAdminIds"`
	DataDir              string      `toml:"dataDir"`
	SystemCache          SystemCache `toml:"cache"`
}
type SystemCache struct {
	Model  string `toml:"model"`
	Prefix string `toml:"prefix"`
}

var config *Conf

func Instance() *Conf {
	if config == nil {
		Init()
	}
	return config
}

func Init() {

	var path = "./mainifest/config/app.toml"
	//读取配置文件
	if _, err := toml.DecodeFile(path, &config); err != nil {
		log.Err(err).Str("path", path).Msg("配置文件读取异常文件存在切格式路径正确吗？")
		return
	}
}
