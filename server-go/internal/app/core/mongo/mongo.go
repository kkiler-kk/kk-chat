package mongo

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"server-go/internal/app/core/config"
)

var mongoCli *mongo.Client

func Instance() *mongo.Client {
	if mongoCli == nil {
		InitMongoDB()
	}
	return mongoCli
}

func InitMongoDB() {
	// 这里uri使用副本集模式，如果你的MongoDB是其他模式，改为上面其他模式的uri即可
	uri := fmt.Sprintf("mongodb://%s:%d", config.Instance().Mongo.Host, config.Instance().Mongo.Port)
	opts := options.Client().ApplyURI(uri).SetAuth(options.Credential{
		Username: config.Instance().Mongo.Name,
		Password: config.Instance().Mongo.Pass,
	})
	opts.MaxPoolSize = config.Instance().Mongo.MaxPoolSize
	opts.MinPoolSize = config.Instance().Mongo.MinPoolSize
	client, err := mongo.Connect(context.TODO())
	if err != nil {
		panic(err)
	}
	if err = client.Ping(context.TODO(), readpref.Primary()); err != nil { // ping
		panic(err)
	}
	mongoCli = client
	log.Info().Msg("mongodb connect success~")
}
