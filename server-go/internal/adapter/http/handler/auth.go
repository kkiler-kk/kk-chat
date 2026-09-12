package handler

import (
	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type Auth struct {
	uc port.AuthUseCase
}

func NewAuth(uc port.AuthUseCase) *Auth { return &Auth{uc: uc} }

func (h *Auth) Register(c echo.Context) error {
	var req dto.RegisterReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	err := h.uc.Register(c.Request().Context(), port.RegisterInput{
		Identity: req.Identity, Name: req.Name, Password: req.Password,
		Email: req.Email, EmailCode: req.EmailCode,
	})
	if err != nil {
		return err // ErrorHandler 统一渲染
	}
	return dto.OK(c, nil)
}

func (h *Auth) Login(c echo.Context) error {
	var req dto.LoginReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	out, err := h.uc.Login(c.Request().Context(), port.LoginInput{
		Account: req.Account, Password: req.Password,
		CaptchaID: req.CaptchaID, CaptchaAnswer: req.CaptchaAnswer,
	})
	if err != nil {
		return err
	}
	return dto.OK(c, out)
}

func (h *Auth) Logout(c echo.Context) error {
	ctx := c.Request().Context()
	jti, _ := middleware.JTIFromContext(ctx)
	uid, _ := middleware.UserIDFromContext(ctx)
	if err := h.uc.Logout(ctx, jti, uid); err != nil {
		return err
	}
	return dto.OK(c, nil)
}

func (h *Auth) Captcha(c echo.Context) error {
	id, b64, err := h.uc.GenerateCaptcha(c.Request().Context())
	if err != nil {
		return err
	}
	return dto.OK(c, map[string]string{"captcha_id": id, "image": b64})
}

func (h *Auth) EmailCode(c echo.Context) error {
	var req dto.EmailCodeReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	if err := h.uc.SendEmailCode(c.Request().Context(), req.Email); err != nil {
		return err
	}
	return dto.OK(c, nil)
}
