package mongox

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewClient 连接并探活 MongoDB。
func NewClient(ctx context.Context, uri string) (*mongo.Client, error) {
	c, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("连接 mongo 失败: %w", err)
	}
	if err := c.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("ping mongo 失败: %w", err)
	}
	return c, nil
}
