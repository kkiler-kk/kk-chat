package control

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"server-go/internal/app/core/config"
	"server-go/internal/app/models/response"
	"server-go/internal/common/utility"
)

var FileControl *fileControl

type fileControl struct {
}

func (f *fileControl) UploadFileImage(c *gin.Context) {
	fmt.Println("uploadImages Hello World")
	file, _ := c.FormFile("file")
	dst := config.Instance().Server.ServerRoot + "/file/"
	nowTime := utility.CreateDateDir(dst)
	dst = dst + nowTime + "/" + file.Filename
	err := c.SaveUploadedFile(file, dst)
	if err != nil {
		response.ErrorResp(c).SetMsg("上传头像失败!").WriteJsonExit()
		return
	}
	url := "http://" + c.Request.Host + "/" + dst
	response.SuccessResp(c).SetData(url).WriteJsonExit()
}
