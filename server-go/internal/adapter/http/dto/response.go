package dto

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"server-go/internal/usecase/apperror"
)

// Response 统一响应信封（成功与失败同构）。
type Response struct {
	Code      int    `json:"code"`
	Message   string `json:"message"`
	Data      any    `json:"data,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

func OK(c echo.Context, data any) error {
	return c.JSON(200, Response{
		Code:      apperror.CodeOK,
		Message:   "ok",
		Data:      data,
		RequestID: c.Response().Header().Get(echo.HeaderXRequestID),
	})
}

var validate = validator.New(validator.WithRequiredStructEnabled())

// Validate 结构体校验；失败 → apperror 10001。
func Validate(v any) error {
	if err := validate.Struct(v); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			f := ve[0]
			return apperror.New(apperror.CodeInvalidParam,
				"参数 "+f.Field()+" 校验失败("+f.Tag()+")")
		}
		return apperror.New(apperror.CodeInvalidParam, "参数格式错误")
	}
	return nil
}
