package usecase

import (
	"context"
	"errors"
	"strings"

	"server-go/internal/domain"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type GroupDeps struct {
	Groups   port.GroupRepo
	Notifier port.Notifier
	Clock    port.Clock
}

type groupUseCase struct{ d GroupDeps }

func NewGroup(d GroupDeps) port.GroupUseCase { return &groupUseCase{d: d} }

func (u *groupUseCase) Create(ctx context.Context, ownerID int64, in port.CreateGroupInput) (*port.GroupItem, error) {
	name := strings.TrimSpace(in.Name)
	if name == "" {
		return nil, apperror.New(apperror.CodeInvalidParam, "群名称不能为空")
	}
	// 成员去重并剔除 owner
	seen := map[int64]bool{ownerID: true}
	members := make([]int64, 0, len(in.MemberIDs))
	for _, id := range in.MemberIDs {
		if !seen[id] {
			seen[id] = true
			members = append(members, id)
		}
	}
	g := &domain.Group{Name: name, OwnerID: ownerID, CreatedAt: u.d.Clock.Now()}
	if err := u.d.Groups.Create(ctx, g, members); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "创建群组失败", err)
	}
	if len(members) > 0 {
		u.d.Notifier.ToUsers(ctx, members, port.NotifierEvent{
			Event: "group.updated", Data: map[string]any{"group_id": g.ID},
		})
	}
	return &port.GroupItem{
		ID: g.ID, Name: g.Name, Avatar: g.Avatar,
		OwnerID: g.OwnerID, MemberCount: len(members) + 1,
	}, nil
}

func (u *groupUseCase) Join(ctx context.Context, groupID, userID int64) error {
	if _, err := u.d.Groups.ByID(ctx, groupID); err != nil {
		if errors.Is(err, domain.ErrGroupNotFound) {
			return apperror.New(apperror.CodeGroupNotFound, domain.ErrGroupNotFound.Error())
		}
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	isMember, err := u.d.Groups.IsMember(ctx, groupID, userID)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	if isMember {
		return nil // 幂等
	}
	if err := u.d.Groups.Join(ctx, groupID, userID); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "加入群组失败", err)
	}
	ids, err := u.d.Groups.MemberIDs(ctx, groupID)
	if err == nil {
		others := make([]int64, 0, len(ids))
		for _, id := range ids {
			if id != userID {
				others = append(others, id)
			}
		}
		u.d.Notifier.ToUsers(ctx, others, port.NotifierEvent{
			Event: "group.updated", Data: map[string]any{"group_id": groupID},
		})
	}
	return nil
}

func (u *groupUseCase) List(ctx context.Context, userID int64) ([]port.GroupItem, error) {
	groups, err := u.d.Groups.ListByUser(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "查询群列表失败", err)
	}
	return u.toItems(ctx, groups), nil
}

func (u *groupUseCase) Search(ctx context.Context, keyword string) ([]port.GroupItem, error) {
	keyword = strings.TrimSpace(keyword)
	if keyword == "" {
		return []port.GroupItem{}, nil
	}
	groups, err := u.d.Groups.SearchByName(ctx, keyword, 20)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "搜索群组失败", err)
	}
	return u.toItems(ctx, groups), nil
}

func (u *groupUseCase) toItems(ctx context.Context, groups []*domain.Group) []port.GroupItem {
	out := make([]port.GroupItem, 0, len(groups))
	for _, g := range groups {
		count, err := u.d.Groups.MemberCount(ctx, g.ID)
		if err != nil {
			count = 0
		}
		out = append(out, port.GroupItem{
			ID: g.ID, Name: g.Name, Avatar: g.Avatar, OwnerID: g.OwnerID, MemberCount: count,
		})
	}
	return out
}

var _ port.GroupUseCase = (*groupUseCase)(nil)
