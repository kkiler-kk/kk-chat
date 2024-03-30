package router

import "github.com/gin-gonic/gin"

func Router() *gin.Engine {
	r := gin.Default()
	return r
}
