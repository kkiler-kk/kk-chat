package handler

import (
	"time"

	"github.com/labstack/echo/v4"
	"server-go/internal/adapter/http/dto"
	"server-go/internal/adapter/http/middleware"
	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type Chat struct {
	uc port.ChatUseCase
}

func NewChat(uc port.ChatUseCase) *Chat { return &Chat{uc: uc} }

func (h *Chat) SendMessage(c echo.Context) error {
	var req dto.SendMessageReq
	if err := c.Bind(&req); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "请求体格式错误")
	}
	if err := dto.Validate(&req); err != nil {
		return err
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	out, err := h.uc.SendMessage(c.Request().Context(), uid, port.SendMessageInput{
		ConversationID: req.ConversationID,
		Content:        req.Content,
		ContentType:    req.ContentType,
	})
	if err != nil {
		return err
	}
	return dto.OK(c, out)
}

func (h *Chat) History(c echo.Context) error {
	convID, err := domain.ParseConversationID(c.Param("id"))
	if err != nil {
		return apperror.New(apperror.CodeConversationInvalid, domain.ErrInvalidConversation.Error())
	}
	var q dto.HistoryQuery
	if err := c.Bind(&q); err != nil {
		return apperror.New(apperror.CodeInvalidParam, "参数格式错误")
	}
	if err := dto.Validate(&q); err != nil {
		return err
	}
	var cursor time.Time
	if q.Cursor != "" {
		cursor, err = time.Parse(time.RFC3339, q.Cursor)
		if err != nil {
			return apperror.New(apperror.CodeInvalidParam, "cursor 应为 RFC3339 时间")
		}
	}
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	msgs, err := h.uc.History(c.Request().Context(), uid, convID, cursor, q.Limit)
	if err != nil {
		return err
	}
	return dto.OK(c, msgs)
}

func (h *Chat) Conversations(c echo.Context) error {
	uid, _ := middleware.UserIDFromContext(c.Request().Context())
	list, err := h.uc.RecentConversations(c.Request().Context(), uid)
	if err != nil {
		return err
	}
	return dto.OK(c, list)
}
