package location

import "github.com/gin-gonic/gin"

// GetClientIp 获取客户端IP
func GetClientIp(c *gin.Context) string {
	return c.ClientIP()
}
