package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog/log"
	"server-go/internal/app/core/config"
)

var redisCli *redis.Client

func Instance() *redis.Client {
	if redisCli == nil {
		InitRedis()
	}
	return redisCli
}

func InitRedis() {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Instance().Redis.Host, config.Instance().Redis.Port), // redis地址
		Password: config.Instance().Redis.Pass,                                                     // 密码
		DB:       config.Instance().Redis.DB,                                                       // 使用默认数据库
	})
	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		log.Err(err).Str("host", config.Instance().Database.DatabaseDefault.Host).Str("user", config.Instance().Database.DatabaseDefault.User).Msg("连接redis失败")
		return
	}
	redisCli = client
	log.Info().Msg("链接Redis成功~")
}
