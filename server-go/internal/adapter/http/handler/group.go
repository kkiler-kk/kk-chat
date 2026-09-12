package handler

import (
	"strconv"

	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type Group struct {
	uc port.GroupUseCase
}

func NewGroup(uc port.GroupUseCase) *Group { return &Group{uc: uc} }

func (h *Group) Create(c echo.Context) error {
	var req dto.CreateGroupReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	item, err := h.uc.Create(c.Request().Context(), uid, port.CreateGroupInput{
		Name: req.Name, MemberIDs: req.MemberIDs,
	})
	if err != nil {
		return err
	}
	return dto.OK(c, item)
}

func (h *Group) Join(c echo.Context) error {
	gid, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || gid <= 0 {
		return apperror.New(apperror.CodeInvalidParam, "传入错误群id")
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	if err := h.uc.Join(c.Request().Context(), gid, uid); err != nil {
		return err
	}
	return dto.OK(c, nil)
}

// List 带 ?search=xx 时走搜索（spec §3.1 GET /groups?search=xx）。
func (h *Group) List(c echo.Context) error {
	ctx := c.Request().Context()
	if search := c.QueryParam("search"); search != "" {
		items, err := h.uc.Search(ctx, search)
		if err != nil {
			return err
		}
		return dto.OK(c, items)
	}
	uid, _ := middleware.UserIDFromContext(ctx)
	items, err := h.uc.List(ctx, uid)
	if err != nil {
		return err
	}
	return dto.OK(c, items)
}
