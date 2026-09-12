package handler

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type User struct {
	uc port.UserUseCase
}

func NewUser(uc port.UserUseCase) *User { return &User{uc: uc} }

func (h *User) Me(c echo.Context) error {
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	info, err := h.uc.Me(c.Request().Context(), uid)
	if err != nil {
		return err
	}
	return dto.OK(c, info)
}

func (h *User) UpdateMe(c echo.Context) error {
	var req dto.UpdateUserReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	in := port.UpdateUserInput{
		Name: req.Name, Phone: req.Phone, Email: req.Email,
		EmailCode: req.EmailCode, Avatar: req.Avatar, Signature: req.Signature,
	}
	if req.BirthDate != nil && *req.BirthDate != "" {
		t, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return apperror.New(apperror.CodeInvalidParam, "生日格式应为 YYYY-MM-DD")
		}
		in.BirthDate = &t
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	if err := h.uc.UpdateMe(c.Request().Context(), uid, in); err != nil {
		return err
	}
	return dto.OK(c, nil)
}

// Detail 走 OptionalJWT：viewerID 取不到时为 0（游客视角）。
func (h *User) Detail(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		return apperror.New(apperror.CodeInvalidParam, "传入错误id")
	}
	viewerID, _ := middleware.UserIDFromContext(c.Request().Context())
	out, err := h.uc.Detail(c.Request().Context(), viewerID, id)
	if err != nil {
		return err
	}
	return dto.OK(c, out)
}

func (h *User) Search(c echo.Context) error {
	var q dto.SearchQuery
	if err := c.Bind(&q); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "参数格式错误")
	}
	if err := dto.Validate(&q); err != nil {
		return err
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	items, err := h.uc.Search(c.Request().Context(), uid, q.Keyword)
	if err != nil {
		return err
	}
	return dto.OK(c, items)
}
