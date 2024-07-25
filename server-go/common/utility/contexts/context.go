package contexts

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"server-go/common/utility/redisUtil"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models"
)

func GetUser(c *gin.Context) *models.UserBasic {
	var (
		result = &models.UserBasic{}
	)
	token, exists := c.Get("token")
	if exists {
		return result
	}
	jsonStr, err := redisUtil.Get(c, config.Instance().Token.CacheKey+token.(string))
	if err != nil {
		return result
	}
	err = json.Unmarshal([]byte(jsonStr), &result)
	return result
}
func SetUser(c *gin.Context, user *models.GroupBasic) {

}
