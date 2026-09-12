package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	gomysql "github.com/go-sql-driver/mysql"
	"server-go/internal/adapter/persistence/mysql/sqlcgen"
	"server-go/internal/domain"
	"server-go/internal/usecase/port"
)

type UserRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) port.UserRepo { return &UserRepo{db: db} }

func mapUser(row sqlcgen.User) *domain.User {
	u := &domain.User{
		ID:        row.ID,
		Identity:  row.Identity,
		Name:      row.Name,
		Password:  row.PasswordHash,
		Email:     row.Email,
		Phone:     row.Phone,
		Avatar:    row.Avatar,
		Signature: row.Signature,
		IsAdmin:   row.IsAdmin == 1,
		Status:    domain.UserStatus(row.Status),
		CreatedAt: row.CreatedAt,
		UpdatedAt: row.UpdatedAt,
	}
	if row.BirthDate.Valid {
		t := row.BirthDate.Time
		u.BirthDate = &t
	}
	return u
}

func (r *UserRepo) Create(ctx context.Context, u *domain.User) error {
	res, err := sqlcgen.New(r.db).InsertUser(ctx, sqlcgen.InsertUserParams{
		Identity: u.Identity, Name: u.Name, PasswordHash: u.Password, Email: u.Email,
	})
	if err != nil {
		if isDuplicate(err) {
			return domain.ErrEmailTaken // identity/email 唯一键冲突
		}
		return fmt.Errorf("插入用户: %w", err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		return fmt.Errorf("获取用户ID: %w", err)
	}
	u.ID = id
	return nil
}

func (r *UserRepo) ByID(ctx context.Context, id int64) (*domain.User, error) {
	row, err := sqlcgen.New(r.db).GetUserByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("查询用户: %w", err)
	}
	return mapUser(row), nil
}

func (r *UserRepo) ByIdentity(ctx context.Context, identity string) (*domain.User, error) {
	row, err := sqlcgen.New(r.db).GetByIdentity(ctx, identity)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("按账号查询用户: %w", err)
	}
	return mapUser(row), nil
}

func (r *UserRepo) ByEmail(ctx context.Context, email string) (*domain.User, error) {
	row, err := sqlcgen.New(r.db).GetByEmail(ctx, email)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("按邮箱查询用户: %w", err)
	}
	return mapUser(row), nil
}

func (r *UserRepo) ByIDs(ctx context.Context, ids []int64) ([]*domain.User, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	rows, err := sqlcgen.New(r.db).GetUsersByIDs(ctx, ids)
	if err != nil {
		return nil, fmt.Errorf("批量查询用户: %w", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapUser(row))
	}
	return out, nil
}

func (r *UserRepo) UpdateProfile(ctx context.Context, id int64, p port.UpdateProfile) error {
	params := sqlcgen.UpdateUserParams{ID: id}
	if p.Name != nil {
		params.Name = sql.NullString{String: *p.Name, Valid: true}
	}
	if p.Phone != nil {
		params.Phone = sql.NullString{String: *p.Phone, Valid: true}
	}
	if p.Email != nil {
		params.Email = sql.NullString{String: *p.Email, Valid: true}
	}
	if p.Avatar != nil {
		params.Avatar = sql.NullString{String: *p.Avatar, Valid: true}
	}
	if p.Signature != nil {
		params.Signature = sql.NullString{String: *p.Signature, Valid: true}
	}
	if p.BirthDate != nil {
		params.BirthDate = sql.NullTime{Time: *p.BirthDate, Valid: true}
	}
	if err := sqlcgen.New(r.db).UpdateUser(ctx, params); err != nil {
		if isDuplicate(err) {
			return domain.ErrEmailTaken
		}
		return fmt.Errorf("更新用户资料: %w", err)
	}
	return nil
}

func (r *UserRepo) Search(ctx context.Context, keyword string, limit int) ([]*domain.User, error) {
	rows, err := sqlcgen.New(r.db).SearchUsers(ctx, sqlcgen.SearchUsersParams{
		CONCAT: keyword, CONCAT_2: keyword, CONCAT_3: keyword, Limit: int32(limit),
	})
	if err != nil {
		return nil, fmt.Errorf("搜索用户: %w", err)
	}
	out := make([]*domain.User, 0, len(rows))
	for _, row := range rows {
		out = append(out, mapUser(row))
	}
	return out, nil
}

// isDuplicate 判断是否 MySQL 1062 唯一键冲突。
func isDuplicate(err error) bool {
	var me *gomysql.MySQLError
	return errors.As(err, &me) && me.Number == 1062
}

var _ port.UserRepo = (*UserRepo)(nil)
