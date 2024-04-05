package router

import (
	"github.com/gin-gonic/gin"
	"server-go/internal/app/control"
)

func FileRouter(r *gin.RouterGroup) {
	fileRouter := r.Group("file")
	fileRouter.POST("/upload/image", control.FileControl.UploadFileImage) // 上传单张image 文件接口
}
