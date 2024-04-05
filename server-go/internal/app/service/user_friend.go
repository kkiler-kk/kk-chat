package service

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"server-go/internal/app/core/db"
	"server-go/internal/app/dao"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
)

var UserFriendService = &userFriendService{}

type userFriendService struct {
}

// AddFriend @Title 添加好友
func (u *userFriendService) AddFriend(input request.UserFriendAddReqInput) (err error) {
	var userFriend = &models.UserFriends{
		UserId:   input.UserId,
		FriendId: input.FriendId,
	}
	if input.UserId == input.FriendId { // 自己加自己好友
		_, err = dao.UserFriends.Insert(userFriend)
	} else {
		err = db.Instance().Transaction(func(tx *gorm.DB) error { // 开启事务 保证数据一致性
			_, err = dao.UserFriends.Insert(userFriend) // 加对方好友
			if err != nil {
				return err
			}
			userFriend = &models.UserFriends{
				UserId:   input.FriendId,
				FriendId: input.UserId,
			}
			_, err = dao.UserFriends.Insert(userFriend) // 让对方加一好友 双向绑定! 第一次是添加进数据库 1->2 第二次是 2->1 这样就不会有的单向好友
			if err != nil {
				return err
			}
			return nil
		})
	}
	return
}

// List @Title返回好友列表
func (u *userFriendService) List(c *gin.Context) (userList []*response.UserFriendListModelRes, err error) {
	id, _ := c.Get("id")
	var userFriend = &models.UserFriends{UserId: id.(int64)}
	friendList, err := dao.UserFriends.FindUserFriends(userFriend)
	for _, friend := range friendList {
		userList = append(userList, &response.UserFriendListModelRes{
			ID:           friend.FriendUser.ID,
			Name:         friend.FriendUser.Name,
			Avatar:       friend.FriendUser.Avatar,
			IsLogout:     friend.FriendUser.IsLogout == 2, // 2 说明在线
			LoginOutTime: friend.FriendUser.LoginOutTime,
		})
	}
	return
}
