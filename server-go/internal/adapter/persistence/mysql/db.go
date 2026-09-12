package mysql

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

// NewDB 打开并探活 MySQL 连接池。
func NewDB(ctx context.Context, dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("打开 mysql 失败: %w", err)
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(10)
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("连接 mysql 失败: %w", err)
	}
	return db, nil
}
