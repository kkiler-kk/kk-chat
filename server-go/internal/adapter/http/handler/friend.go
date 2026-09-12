package handler

import (
	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type Friend struct {
	uc port.FriendUseCase
}

func NewFriend(uc port.FriendUseCase) *Friend { return &Friend{uc: uc} }

func (h *Friend) Add(c echo.Context) error {
	var req dto.AddFriendReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	if err := h.uc.Add(c.Request().Context(), uid, req.UserID); err != nil {
		return err
	}
	return dto.OK(c, nil)
}

func (h *Friend) List(c echo.Context) error {
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	items, err := h.uc.List(c.Request().Context(), uid)
	if err != nil {
		return err
	}
	return dto.OK(c, items)
}
