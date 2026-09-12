package mysql

import (
	"context"
	"database/sql"
	"fmt"

	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type FriendRepo struct {
	db *sql.DB
}

func NewFriendRepo(db *sql.DB) port.FriendRepo { return &FriendRepo{db: db} }

// Add 在一个事务中写双向两行。
func (r *FriendRepo) Add(ctx context.Context, userID, friendID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("开启事务: %w", err)
	}
	defer tx.Rollback() //nolint:errcheck // Commit 后 Rollback 返回错误可忽略
	q := sqlcgen.New(tx)
	if err := q.InsertFriend(ctx, sqlcgen.InsertFriendParams{UserID: userID, FriendID: friendID}); err != nil {
		return fmt.Errorf("写入好友关系: %w", err)
	}
	if err := q.InsertFriend(ctx, sqlcgen.InsertFriendParams{UserID: friendID, FriendID: userID}); err != nil {
		return fmt.Errorf("写入反向好友关系: %w", err)
	}
	return tx.Commit()
}

func (r *FriendRepo) Exists(ctx context.Context, a, b int64) (bool, error) {
	n, err := sqlcgen.New(r.db).CountFriendship(ctx, sqlcgen.CountFriendshipParams{UserID: a, FriendID: b})
	if err != nil {
		return false, fmt.Errorf("查询好友关系: %w", err)
	}
	return n > 0, nil
}

func (r *FriendRepo) ListFriends(ctx context.Context, userID int64) ([]*domain.User, error) {
	q := sqlcgen.New(r.db)
	ids, err := q.ListFriendIDs(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("查询好友ID: %w", err)
	}
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := q.GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("批量查询好友: %w", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapUser(row))
	}
	return out, nil
}

var _ port.FriendRepo = (*FriendRepo)(nil)
