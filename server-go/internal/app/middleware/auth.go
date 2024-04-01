package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server-go/internal/common/utility/authUtil"
	"strings"
)

// AuthMiddlewareUpdate @Title 访问接口需的登录验证
func AuthMiddlewareUpdate(c *gin.Context) {
	authHeader := c.Request.Header.Get("Authorization")
	if "" == authHeader {
		authHeader = c.Query("token")
		authHeader = "Bearer " + authHeader
		if "" == authHeader || authHeader == "undefined" {

			c.JSON(http.StatusOK, gin.H{
				"code": 2003,
				"msg":  "请重新登录",
			})
			c.Abort()
			return
		}
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		c.JSON(http.StatusOK, gin.H{
			"code": 2004,
			"msg":  "请求头中的auth格式错误",
		})
		//阻止调用后续的函数
		c.Abort()
		return
	}
	mc, err := authUtil.ParseToken(parts[1])
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": 2005,
			"msg":  "无效的token",
		})
		c.Abort()
		return
	}
	c.Set("id", mc.ID)
	c.Set("token", parts[1])
	c.Next()
}
