package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"server-go/internal/consts"
)

type ApiResp struct {
	c *gin.Context
	r *CommonResp
}

// 返回一个成功的消息体
func SuccessResp(c *gin.Context) *ApiResp {
	var a = ApiResp{
		r: SuccessCode,
		c: c,
	}
	return &a
}

// 返回一个错误的消息体
func ErrorResp(c *gin.Context) *ApiResp {
	msg := CommonResp{
		Code: consts.StatusServerError,
		Type: OperOther,
		Msg:  "操作失败",
	}
	var a = ApiResp{
		r: &msg,
		c: c,
	}
	return &a
}

// 返回一个拒绝访问的消息体
func ForbiddenResp(c *gin.Context) *ApiResp {
	var a = ApiResp{
		r: StatusNotAccess,
		c: c,
	}
	return &a
}

// JWT认证失败
func UnauthorizedResp(c *gin.Context) *ApiResp {
	msg := CommonResp{
		Code: 401,
		Type: OperOther,
		Msg:  "鉴权失败",
	}
	var a = ApiResp{
		r: &msg,
		c: c,
	}
	return &a
}

// 设置消息体的内容
func (resp *ApiResp) SetMsg(msg string) *ApiResp {
	resp.r.Msg = msg
	return resp
}

// 设置消息体的编码
func (resp *ApiResp) SetCode(code int) *ApiResp {
	resp.r.Code = code
	return resp
}

// 设置消息体的数据
func (resp *ApiResp) SetData(data interface{}) *ApiResp {
	resp.r.Data = data
	return resp
}

// 设置消息体的业务类型
func (resp *ApiResp) SetType(btype OperationType) *ApiResp {
	resp.r.Type = btype
	return resp
}

// 设置消息体的业务类型
func (resp *ApiResp) SetCount(count int) *ApiResp {
	resp.r.Count = count
	return resp
}

// 记录操作日志到数据库
//func (resp *ApiResp) Log(title string, inParam interface{}) *ApiResp {
//	var inContentStr string
//	switch inParam.(type) {
//	case string, []byte:
//		inContentStr = gconv.String(inParam)
//	}
//	// Else use json.Marshal function to encode the parameter.
//	if b, err := json.Marshal(inParam); err != nil {
//		inContentStr = ""
//	} else {
//		inContentStr = string(b)
//	}
//	var errMsg string
//	if resp.r.Code == 500 {
//		errMsg = resp.r.Msg
//	}
//	err := service.CreateOperLog(resp.c, models.OperForm{Title: title, InContent: inContentStr, ErrorMsg: errMsg, OutContent: resp.r})
//	if err != nil {
//		log.Instance().Error(err.Error())
//	}
//	return resp
//}

// 输出json到客户端
func (resp *ApiResp) WriteJsonExit() {
	resp.c.JSON(http.StatusOK, resp.r)
	resp.c.Abort()
}

func (resp *ApiResp) WriteErrJsonExit(errCode int) {
	resp.c.JSON(errCode, resp.r)
	resp.c.Abort()
}
