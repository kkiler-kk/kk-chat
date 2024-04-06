package mongoUtils

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/core/mongo"
	"server-go/internal/app/models"
	"time"
)

// Insert @Title 插入信息到mongodb
func Insert(c *gin.Context, database, id string, content string, expire int64) (err error) {
	collection := mongo.Instance().Database(database).Collection(id)
	message := models.Trainer{
		Content:   content,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire, // 加上过期时间
	}
	_, err = collection.InsertOne(c, message)
	return
}
