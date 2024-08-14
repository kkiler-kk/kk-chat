package response

import "server-go/internal/consts"

var (
	SuccessCode       = NewCommonResp(consts.StatusOK, "OK", OperOther)
	StatusNotAccess   = NewCommonResp(consts.StatusNotAccess, "对不起，您暂无权限访问", OperOther)
	StatusServerError = NewCommonResp(consts.StatusServerError, "服务器错误！请联系管理员", OperOther)
)
