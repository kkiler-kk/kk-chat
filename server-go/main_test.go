package main

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	//c := context.Background()
	//var recentModel = &response.RecentChatListRes{ // 发送者先改变一次
	//	ID:        2,
	//	Name:      "namqweqe2323",
	//	Content:   "weweqe",
	//	Avatar:    "123",
	//	Type:      1,
	//	CreatedAt: time.Now(),
	//}
	//redisUtil.HSet(c, consts.RecentChat, "123", recentModel)
	//str, _ := redisUtil.HGet(c, consts.RecentChat, "123").Result()
	//var userRecentList []*response.RecentChatListRes
	//json.Unmarshal([]byte(str), &userRecentList)
	//fmt.Println("str", str)
	//fmt.Println(userRecentList)
	args := []string{"23", "23"}
	fmt.Println(args)
	MuTest(args...)
}
func MuTest(args ...string) {
	for _, e := range args {
		fmt.Println(e)
	}
}
