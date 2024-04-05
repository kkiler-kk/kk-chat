package main

import (
	"context"
	"fmt"
	"server-go/internal/app/models/response"
	"server-go/internal/common/utility/redisUtil"
	"server-go/internal/consts"
	"testing"
	"time"
)

func TestName(t *testing.T) {
	c := context.Background()
	model1 := response.RecentChatListRes{
		ID:        1,
		Name:      "KK",
		Content:   "123",
		Avatar:    "12231",
		Type:      1,
		CreatedAt: time.Now(),
	}
	model2 := response.RecentChatListRes{
		ID:        2,
		Name:      "KK",
		Content:   "Content",
		Avatar:    "123",
		Type:      2,
		CreatedAt: time.Now(),
	}
	err := redisUtil.HSet(c, "ws_1", consts.RecentChat, model1, model2)
	fmt.Println("err: ", err)
	result, err := redisUtil.HGet(c, "3123", consts.RecentChat).Result()
	fmt.Println(err)
	fmt.Println("result", result)
}
