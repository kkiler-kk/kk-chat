package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type ChatDeps struct {
	Messages       port.MessageRepo
	Users          port.UserRepo
	Friends        port.FriendRepo
	Groups         port.GroupRepo
	Recent         port.RecentChatStore
	Limits         port.MsgLimitStore
	Presence       port.PresenceStore
	Notifier       port.Notifier
	Clock          port.Clock
	NonFriendLimit int           // 默认 3
	LimitTTL       time.Duration // 默认 24h
}

type chatUseCase struct{ d ChatDeps }

func NewChat(d ChatDeps) port.ChatUseCase { return &chatUseCase{d: d} }

// truncateRunes 截断到 n 个 rune，防中文截半。
func truncateRunes(s string, n int) string {
	r := []rune(s)
	if len(r) <= n {
		return s
	}
	return string(r[:n]) + "…"
}

func (u *chatUseCase) SendMessage(ctx context.Context, senderID int64, in port.SendMessageInput) (*port.OutgoingMessage, error) {
	// 1. 解析会话 ID
	convID, err := domain.ParseConversationID(in.ConversationID)
	if err != nil {
		return nil, apperror.New(apperror.CodeConversationInvalid, domain.ErrInvalidConversation.Error())
	}
	// 2. 校验内容
	content := strings.TrimSpace(in.Content)
	if content == "" && in.ContentType == string(domain.ContentText) {
		return nil, apperror.New(apperror.CodeInvalidParam, "消息内容不能为空")
	}
	if in.ContentType != string(domain.ContentText) && in.ContentType != string(domain.ContentImage) {
		return nil, apperror.New(apperror.CodeInvalidParam, "不支持的消息类型")
	}
	// 3. 发送者
	sender, err := u.d.Users.ByID(ctx, senderID)
	if err != nil {
		return nil, u.mapUserErr(err)
	}
	// 4. 按会话类型确定接收者
	var recipients []int64
	switch convID.Type {
	case domain.ConvPrivate:
		a, b, _ := convID.PrivateParticipants()
		if senderID != a && senderID != b {
			return nil, apperror.New(apperror.CodeConversationInvalid, "无权在该会话发言")
		}
		peer := b
		if senderID == b {
			peer = a
		}
		isFriend, err := u.d.Friends.Exists(ctx, senderID, peer)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
		}
		if !isFriend {
			count, err := u.d.Limits.Incr(ctx, convID, senderID, u.d.LimitTTL)
			if err != nil {
				return nil, apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
			}
			if count > int64(u.d.NonFriendLimit) {
				return nil, apperror.New(apperror.CodeNotFriendLimit, domain.ErrMsgLimitExceeded.Error())
			}
		}
		recipients = []int64{peer}
	case domain.ConvGroup:
		gid, ok := convID.GroupID()
		if !ok {
			return nil, apperror.New(apperror.CodeConversationInvalid, domain.ErrInvalidConversation.Error())
		}
		isMember, err := u.d.Groups.IsMember(ctx, gid, senderID)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
		}
		if !isMember {
			return nil, apperror.New(apperror.CodeNotGroupMember, domain.ErrNotGroupMember.Error())
		}
		ids, err := u.d.Groups.MemberIDs(ctx, gid)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
		}
		recipients = make([]int64, 0, len(ids))
		for _, id := range ids {
			if id != senderID {
				recipients = append(recipients, id)
			}
		}
	}
	// 5. 落库
	msg := &domain.Message{
		ConversationID: convID,
		SenderID:       senderID,
		Content:        in.Content,
		ContentType:    domain.ContentType(in.ContentType),
		CreatedAt:      u.d.Clock.Now(),
	}
	if err := u.d.Messages.Save(ctx, msg); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
	}
	// 6. 出参
	out := &port.OutgoingMessage{
		ID:             msg.ID,
		ConversationID: convID.Value,
		SenderID:       senderID,
		SenderName:     sender.Name,
		SenderAvatar:   sender.Avatar,
		Content:        msg.Content,
		ContentType:    string(msg.ContentType),
		CreatedAt:      msg.CreatedAt,
	}
	// 7. 推送 chat.message 给接收者
	if len(recipients) > 0 {
		u.d.Notifier.ToUsers(ctx, recipients, port.NotifierEvent{Event: "chat.message", Data: out})
	}
	// 8. 更新最近会话（发送者 + 每个接收者都 Touch）
	summary := port.ConversationSummary{
		Type:            convID.Type,
		LastContent:     truncateRunes(in.Content, 50),
		LastContentType: msg.ContentType,
		LastSenderID:    senderID,
		LastSenderName:  sender.Name,
		LastTime:        msg.CreatedAt,
	}
	_ = u.d.Recent.SaveSummary(ctx, convID, summary)
	all := append(append([]int64{}, recipients...), senderID)
	for _, uid := range all {
		_ = u.d.Recent.Touch(ctx, uid, convID, msg.CreatedAt)
	}
	// 9. 推送 chat.recent_updated（payload 按各自视角构造：私聊双方的 Peer 互为对方，群聊对所有人 Peer 都是群）
	switch convID.Type {
	case domain.ConvPrivate:
		for _, uid := range recipients {
			payload := u.buildRecentConversation(ctx, convID, uid, summary)
			u.d.Notifier.ToUser(ctx, uid, port.NotifierEvent{Event: "chat.recent_updated", Data: payload})
		}
		selfPayload := u.buildRecentConversation(ctx, convID, senderID, summary)
		u.d.Notifier.ToUser(ctx, senderID, port.NotifierEvent{Event: "chat.recent_updated", Data: selfPayload})
	case domain.ConvGroup:
		shared := u.buildRecentConversation(ctx, convID, senderID, summary)
		ev := port.NotifierEvent{Event: "chat.recent_updated", Data: shared}
		if len(recipients) > 0 {
			u.d.Notifier.ToUsers(ctx, recipients, ev)
		}
		u.d.Notifier.ToUser(ctx, senderID, ev)
	}
	return out, nil
}

// buildRecentConversation 按 viewer 视角组装会话摘要出参：
// 私聊 Peer 为 viewer 的对方；群聊 Peer 为群信息。
func (u *chatUseCase) buildRecentConversation(ctx context.Context, convID domain.ConversationID, viewerID int64, sum port.ConversationSummary) port.RecentConversation {
	rc := port.RecentConversation{
		ConversationID: convID.Value,
		Type:           convID.Type.String(),
		Online:         false,
		LastContent:    sum.LastContent,
		LastSenderName: sum.LastSenderName,
		LastTime:       sum.LastTime,
	}
	switch convID.Type {
	case domain.ConvPrivate:
		a, b, _ := convID.PrivateParticipants()
		peerID := b
		if viewerID == b {
			peerID = a
		}
		rc.PeerID = peerID
		if peer, err := u.d.Users.ByID(ctx, peerID); err == nil {
			rc.PeerName = peer.Name
			rc.PeerAvatar = peer.Avatar
		}
		if online, err := u.d.Presence.IsOnline(ctx, peerID); err == nil {
			rc.Online = online
		}
	case domain.ConvGroup:
		if gid, ok := convID.GroupID(); ok {
			if g, err := u.d.Groups.ByID(ctx, gid); err == nil {
				rc.PeerID = g.ID
				rc.PeerName = g.Name
				rc.PeerAvatar = g.Avatar
			}
		}
	}
	return rc
}

func (u *chatUseCase) History(ctx context.Context, userID int64, convID domain.ConversationID, cursor time.Time, limit int) ([]port.OutgoingMessage, error) {
	// 权限校验
	switch convID.Type {
	case domain.ConvPrivate:
		a, b, _ := convID.PrivateParticipants()
		if userID != a && userID != b {
			return nil, apperror.New(apperror.CodeConversationInvalid, "无权查看该会话")
		}
	case domain.ConvGroup:
		gid, ok := convID.GroupID()
		if !ok {
			return nil, apperror.New(apperror.CodeConversationInvalid, domain.ErrInvalidConversation.Error())
		}
		isMember, err := u.d.Groups.IsMember(ctx, gid, userID)
		if err != nil {
			return nil, apperror.Wrap(apperror.CodeInternal, "查询消息失败", err)
		}
		if !isMember {
			return nil, apperror.New(apperror.CodeNotGroupMember, domain.ErrNotGroupMember.Error())
		}
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	msgs, err := u.d.Messages.History(ctx, convID, cursor, limit)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询消息失败", err)
	}
	// 批量取发送者信息，避免 N+1
	idSet := make(map[int64]struct{}, len(msgs))
	ids := make([]int64, 0, len(msgs))
	for _, m := range msgs {
		if _, ok := idSet[m.SenderID]; !ok {
			idSet[m.SenderID] = struct{}{}
			ids = append(ids, m.SenderID)
		}
	}
	senders := make(map[int64]*domain.User, len(ids))
	if users, err := u.d.Users.ByIDs(ctx, ids); err == nil {
		for _, usr := range users {
			senders[usr.ID] = usr
		}
	}
	out := make([]port.OutgoingMessage, 0, len(msgs))
	for _, m := range msgs {
		om := port.OutgoingMessage{
			ID:             m.ID,
			ConversationID: convID.Value,
			SenderID:       m.SenderID,
			Content:        m.Content,
			ContentType:    string(m.ContentType),
			CreatedAt:      m.CreatedAt,
		}
		if s, ok := senders[m.SenderID]; ok {
			om.SenderName = s.Name
			om.SenderAvatar = s.Avatar
		}
		out = append(out, om)
	}
	return out, nil
}

func (u *chatUseCase) RecentConversations(ctx context.Context, userID int64) ([]port.RecentConversation, error) {
	convIDs, err := u.d.Recent.List(ctx, userID, 50)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询最近会话失败", err)
	}
	out := make([]port.RecentConversation, 0, len(convIDs))
	for _, convID := range convIDs {
		sum, err := u.d.Recent.Summary(ctx, convID)
		if err != nil || sum == nil {
			continue // 摘要缺失的会话跳过
		}
		rc := port.RecentConversation{
			ConversationID: convID.Value,
			Type:           convID.Type.String(),
			LastContent:    sum.LastContent,
			LastSenderName: sum.LastSenderName,
			LastTime:       sum.LastTime,
		}
		switch convID.Type {
		case domain.ConvPrivate:
			a, b, _ := convID.PrivateParticipants()
			peerID := b
			if userID == b {
				peerID = a
			}
			rc.PeerID = peerID
			peer, err := u.d.Users.ByID(ctx, peerID)
			if err != nil {
				if errors.Is(err, domain.ErrUserNotFound) {
					rc.PeerName = "已注销用户"
				} else {
					continue
				}
			} else {
				rc.PeerName = peer.Name
				rc.PeerAvatar = peer.Avatar
			}
			online, err := u.d.Presence.IsOnline(ctx, peerID)
			rc.Online = err == nil && online
		case domain.ConvGroup:
			gid, ok := convID.GroupID()
			if !ok {
				continue
			}
			g, err := u.d.Groups.ByID(ctx, gid)
			if err != nil {
				continue // 群已不存在则跳过
			}
			rc.PeerID = g.ID
			rc.PeerName = g.Name
			rc.PeerAvatar = g.Avatar
		}
		out = append(out, rc)
	}
	return out, nil
}

func (u *chatUseCase) mapUserErr(err error) error {
	if errors.Is(err, domain.ErrUserNotFound) {
		return apperror.New(apperror.CodeUserNotFound, domain.ErrUserNotFound.Error())
	}
	return apperror.Wrap(apperror.CodeInternal, "消息发送失败", err)
}

var _ port.ChatUseCase = (*chatUseCase)(nil)
