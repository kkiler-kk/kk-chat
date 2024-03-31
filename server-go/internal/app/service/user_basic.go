package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"server-go/internal/app/core/config"
	"server-go/internal/app/dao"
	"server-go/internal/app/models"
	"server-go/internal/app/models/request"
	"server-go/internal/consts"
	"server-go/internal/utility/authUtil"
	"server-go/internal/utility/redisUtil"
	"server-go/internal/utility/utils"
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
	redisUtil.Del(c, consts.CodeEmail+userBasicModel.Email) // 添加成功 删除邮箱验证码
	return
}

// Login @Title 登录
func (u *userBasicService) Login(c *gin.Context, input request.UserBasicLoginReqInput) (token string, err error) {
	user, err := dao.UserBasic.FindUniqueUserBasic(0, []string{input.UserName, input.UserName, input.UserName})
	if err != nil {
		return "", errors.New("账号不存在")
	}
	fmt.Println("user1", user)
	user, err = dao.UserBasic.Login(user.ID, utils.EncryptPassword(input.Password, user.Salt)) // 登录验证
	if err != nil {
		return "", errors.New("账号或密码错误！")
	}
	token, err = authUtil.GenToken(user.ID)
	jsonStr, _ := json.Marshal(user)
	redisUtil.Set(c, config.Instance().Token.CacheKey+token, string(jsonStr), config.Instance().Token.TimeOut*time.Minute)
	return token, err
}
