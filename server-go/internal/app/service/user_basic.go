package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"server-go/common/utility/authUtil"
	"server-go/common/utility/redisUtil"
	"server-go/common/utility/utils"
	"server-go/internal/app/core/config"
	"server-go/internal/app/dao"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/app/models/response"
	"server-go/internal/consts"
	"time"
)

var UserBasicService = &userBasicService{}

type userBasicService struct {
}

// CreateUserBasic @Title 创建用户
func (u *userBasicService) CreateUserBasic(c *gin.Context, input request.UserBasicCreateReqInput) (err error) {
	salt := utils.RandStr(6) // 加密盐
	var userBasicModel = &models.UserBasic{
		Name:          input.UserName,
		Identity:      input.Identity,
		Password:      utils.EncryptPassword(input.OldPassword, salt),
		Email:         input.Email,
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LoginOutTime:  time.Now(),
		Salt:          salt,
		IsLogout:      1,
		IsForbidden:   2,
		IsAdmin:       2,
	}
	code, err := redisUtil.Get(c, consts.CodeEmail+userBasicModel.Email)
	if err != nil || input.VerificationCode != code {
		return errors.New("验证码过期 or 错误")
	}
	temp, _ := dao.UserBasic.FindUniqueUserBasic(0, []string{"", userBasicModel.Email, ""})
	if temp.ID != 0 {
		return errors.New("邮箱已经存在了！")
	}
	temp, _ = dao.UserBasic.FindUniqueUserBasic(0, []string{userBasicModel.Identity, "", ""})
	if temp.ID != 0 {
		return errors.New("id已经存在了！")
	}
	_, err = dao.UserBasic.Insert(userBasicModel)
	_ = redisUtil.Del(c, consts.CodeEmail+userBasicModel.Email) // 添加成功 删除邮箱验证码
	return
}

// Login @Title 登录
func (u *userBasicService) Login(c *gin.Context, input request.UserBasicLoginReqInput) (res *response.LoginModelRes, err error) {
	user, err := dao.UserBasic.FindUniqueUserBasic(0, []string{input.UserName, input.UserName, input.UserName})
	if err != nil {
		return res, errors.New("账号不存在")
	}
	fmt.Println("user1", user)
	user, err = dao.UserBasic.Login(user.ID, utils.EncryptPassword(input.Password, user.Salt)) // 登录验证
	if err != nil {
		return res, errors.New("账号或密码错误！")
	}

	// 更新用户登录时间
	go u.UpdateLoginTime(user.ID)

	token, err := authUtil.GenToken(user.ID)
	jsonStr, _ := json.Marshal(user)
	_ = redisUtil.Set(c, config.Instance().Token.CacheKey+token, string(jsonStr), config.Instance().Token.TimeOut*time.Minute)
	res = &response.LoginModelRes{
		Token: token,
		User: &response.UserDetailModelRes{
			ID:        user.ID,
			Name:      user.Name,
			Avatar:    user.Avatar,
			Phone:     user.Phone,
			Email:     user.Email,
			Identity:  user.Identity,
			CreatedAt: user.CreatedAt,
		},
	}
	return res, err
}

// UserDetail @Title 查看用户详细
func (u *userBasicService) UserDetail(c *gin.Context, id int64) (res *response.UserDetailModelRes, err error) {
	user, err := dao.UserBasic.FindUniqueUserBasic(id, []string{"", "", ""})
	if err != nil {
		log.Error().Msg(err.Error())
		return
	}
	res = &response.UserDetailModelRes{
		ID:        user.ID,
		Avatar:    user.Avatar,
		Name:      user.Name,
		Phone:     user.Phone,
		Email:     user.Email,
		Identity:  user.Identity,
		CreatedAt: user.CreatedAt,
	}
	return
}

func (u *userBasicService) SearchUser(c *gin.Context, search string) (res []*response.UserOrGroupModelRes, err error) {
	id, _ := c.Get("id")
	userList, err := dao.UserBasic.FindUniqueUserBasicLikeList(0, []string{search, search, search})
	if err != nil {
		log.Err(err).Msg("尚未查询到用户")
		return
	}
	var userFriend *models.UserFriends
	for _, user := range userList {
		userFriend = &models.UserFriends{
			UserId:   id.(int64),
			FriendId: user.ID,
		}
		isFriend := false
		result, _ := dao.UserFriends.FindUserFriends(userFriend) // 查询是否为好友关系
		if len(result) > 0 {                                     // > 0 说明是好友关系
			isFriend = true
		}
		res = append(res, &response.UserOrGroupModelRes{
			Id:       user.ID,
			Avatar:   user.Avatar,
			Identity: user.Identity,
			Name:     user.Name,
			IsFriend: isFriend,
		})
	}
	return
}

// Logout @Title 用户退出
func (u *userBasicService) Logout(c *gin.Context) (err error) {
	token, _ := c.Get("token")
	id, _ := c.Get("id")
	err = redisUtil.Del(c, fmt.Sprintf("%s%s", config.Instance().Token.CacheKey, token)) // 退出登录 删除 redis token缓存
	var user = &models.UserBasic{
		ID:           id.(int64),
		IsLogout:     1,
		LoginOutTime: time.Now(),
	}
	err = dao.UserBasic.Update(user)
	return
}

// UpdateUserBasic @Title 更新用户信息
func (u *userBasicService) UpdateUserBasic(c *gin.Context, input request.UserBasicUpdateReqInput) (err error) {
	var userBasicModel = &models.UserBasic{
		ID:       input.ID,
		Name:     input.UserName,
		Phone:    input.Phone,
		Email:    input.Email,
		Avatar:   input.Avatar,
		Identity: input.Identity,
	}
	code, err := redisUtil.Get(c, consts.CodeEmail+userBasicModel.Email)
	if code == "" && input.VerificationCode == "" { // 说明人家没想改邮箱
		userBasicModel.Email = ""
	} else if err != nil || input.VerificationCode != code {
		return errors.New("验证码过期 or 错误")
	}
	err = dao.UserBasic.Update(userBasicModel)
	return
}

// UpdateHeartTime @Title 更新心跳
func (u *userBasicService) UpdateHeartTime(id int64) {
	var user = &models.UserBasic{ID: id, HeartbeatTime: time.Now()}
	_ = dao.UserBasic.Update(user)
}

// UpdateLogout @Title 用户断开websocket 更新退出时间
func (u *userBasicService) UpdateLogout(id int64) {
	var user = &models.UserBasic{ID: id, LoginOutTime: time.Now(), IsLogout: 1}
	_ = dao.UserBasic.Update(user)
}

// UpdateLoginTime @Title 用户断开websocket 更新用户在线
func (u *userBasicService) UpdateLoginTime(id int64) {
	var user = &models.UserBasic{ID: id, LoginTime: time.Now(), IsLogout: 2}
	_ = dao.UserBasic.Update(user)
}
