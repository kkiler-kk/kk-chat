package db

import (
	"github.com/rs/zerolog/log"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	logOld "log"
	"os"
	"server-go/internal/app/core/config"
	"time"
)

var conn *gorm.DB

func Instance() *gorm.DB {
	if conn == nil {
		InitMySQL()
	}
	return conn
}

func InitMySQL() {
	newLogger := logger.New(
		logOld.New(os.Stdout, "\r\n", logOld.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open(config.Instance().Database.DatabaseDefault.Link), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		log.Err(err).Str("host", config.Instance().Database.DatabaseDefault.Host).Str("user", config.Instance().Database.DatabaseDefault.User).Msg("连接数据库失败")
		return
	}
	conn = db
}
