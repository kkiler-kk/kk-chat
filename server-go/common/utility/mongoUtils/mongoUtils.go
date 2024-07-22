package mongoUtils

import (
	"context"
	uuid "github.com/satori/go.uuid"
	"go.mongodb.org/mongo-driver/bson"
	"server-go/internal/app/core/mongo"
	"server-go/internal/app/models"
	"time"
)

// Insert @Title 插入信息到mongodb
func Insert(c context.Context, database, id string, content string, expire int64) (err error) {
	collection := mongo.Instance().Database(database).Collection(id)
	message := models.Trainer{
		ID:        uuid.NewV4().String(),
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire, // 加上过期时间
	}
	_, err = collection.InsertOne(c, message)
	return
}

func ListMessAge(c context.Context, database, id string) (result []models.Trainer, err error) {
	collection := mongo.Instance().Database(database).Collection(id)
	find, err := collection.Find(c, bson.D{})
	if err != nil {
		return nil, err
	}
	for find.Next(c) {
		var demo models.Trainer
		err = find.Decode(&demo)
		if err != nil {
			return
		}
		result = append(result, demo)
	}
	return result, nil
}
