package response

type CommonResp struct {
	Code  int           `json:"code"` //响应编码 0 成功 500 错误 403 无权限  -1  失败
	Msg   string        `json:"msg"`  //消息
	Data  interface{}   `json:"data"` //数据内容
	Type  OperationType `json:"type"` //业务类型
	Count int           `json:"count,omitempty"`
}

func NewCommonResp(code int, msg string, Type OperationType) *CommonResp {
	return &CommonResp{
		Code: code,
		Msg:  msg,
		Type: Type,
	}
}

type OperationType int

// 业务类型
const (
	OperOther OperationType = 0 //0其它
	OperAdd   OperationType = 1 //1新增
	OperEdit  OperationType = 2 //2修改
	OperDel   OperationType = 3 //3删除
)

// 响应结果
const (
	SUCCESS      = 0   // 成功
	ERROR        = 500 //错误
	UNAUTHORIZED = 403 //无权限
	FAIL         = -1  //失败
)
