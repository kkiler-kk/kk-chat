package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type GroupRepo struct {
	db *sql.DB
}

func NewGroupRepo(db *sql.DB) port.GroupRepo { return &GroupRepo{db: db} }

func mapGroup(row sqlcgen.Group) *domain.Group {
	return &domain.Group{
		ID: row.ID, Name: row.Name, Avatar: row.Avatar,
		OwnerID: row.OwnerID, CreatedAt: row.CreatedAt,
	}
}

// Create 事务：建群 + owner 成员行 + 初始成员行（INSERT IGNORE 天然去重）。
func (r *GroupRepo) Create(ctx context.Context, g *domain.Group, memberIDs []int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck
	q := sqlcgen.New(tx)
	res, err := q.InsertGroup(ctx, sqlcgen.InsertGroupParams{Name: g.Name, Avatar: g.Avatar, OwnerID: g.OwnerID})
	if err != nil {
		return fmt.Errorf("插入群组: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取群ID: %w", err)
	}
	g.ID = id
	if err := q.InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
		GroupID: id, UserID: g.OwnerID, Role: string(domain.RoleOwner),
	}); err != nil {
		return fmt.Errorf("插入群主成员: %w", err)
	}
	for _, uid := range memberIDs {
		if err := q.InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
			GroupID: id, UserID: uid, Role: string(domain.RoleMember),
		}); err != nil {
			return fmt.Errorf("插入群成员: %w", err)
		}
	}
	return tx.Commit()
}

func (r *GroupRepo) ByID(ctx context.Context, id int64) (*domain.Group, error) {
	row, err := sqlcgen.New(r.db).GetGroupByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrGroupNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询群组: %w", err)
	}
	return mapGroup(row), nil
}

func (r *GroupRepo) Join(ctx context.Context, groupID, userID int64) error {
	return sqlcgen.New(r.db).InsertGroupMember(ctx, sqlcgen.InsertGroupMemberParams{
		GroupID: groupID, UserID: userID, Role: string(domain.RoleMember),
	})
}

func (r *GroupRepo) IsMember(ctx context.Context, groupID, userID int64) (bool, error) {
	n, err := sqlcgen.New(r.db).CountMembership(ctx, sqlcgen.CountMembershipParams{GroupID: groupID, UserID: userID})
	if err != nil {
		return false, fmt.Errorf("查询群成员: %w", err)
	}
	return n > 0, nil
}

func (r *GroupRepo) MemberIDs(ctx context.Context, groupID int64) ([]int64, error) {
	return sqlcgen.New(r.db).ListGroupMemberIDs(ctx, groupID)
}

func (r *GroupRepo) MemberCount(ctx context.Context, groupID int64) (int, error) {
	n, err := sqlcgen.New(r.db).CountGroupMembers(ctx, groupID)
	return int(n), err
}

func (r *GroupRepo) ListByUser(ctx context.Context, userID int64) ([]*domain.Group, error) {
	rows, err := sqlcgen.New(r.db).ListGroupsByUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询用户群列表: %w", err)
	}
	out := make([]*domain.Group, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapGroup(row))
	}
	return out, nil
}

func (r *GroupRepo) SearchByName(ctx context.Context, keyword string, limit int) ([]*domain.Group, error) {
	rows, err := sqlcgen.New(r.db).SearchGroupsByName(ctx, sqlcgen.SearchGroupsByNameParams{
		CONCAT: keyword, Limit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("搜索群组: %w", err)
	}
	out := make([]*domain.Group, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapGroup(row))
	}
	return out, nil
}

var _ port.GroupRepo = (*GroupRepo)(nil)
