package authUtil

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"server-go/common/utility/utils"
	"server-go/internal/app/core/config"
	"time"
)

type MyClaims struct {
	ID int64 `json:"id"`
	jwt.RegisteredClaims
}

// 定义密钥
const sign_key = "kk-chat"

// 定义过期时间
var expiresAt = config.Instance().Token.TimeOut * time.Minute

// 定义secret
var MySecret = []byte(config.Instance().Token.EncryptKey)

// 生成jwt
func GenToken(ID int64) (string, error) {
	claim := MyClaims{
		ID: ID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Auth_Server",                                   // 签发者
			Subject:   "Tom",                                           // 签发对象
			Audience:  jwt.ClaimStrings{"Android_APP", "IOS_APP"},      //签发受众
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiresAt)),   //过期时间
			NotBefore: jwt.NewNumericDate(time.Now().Add(time.Second)), //最早使用时间
			IssuedAt:  jwt.NewNumericDate(time.Now()),                  //签发时间
			ID:        utils.RandStr(10),                               // wt ID, 类似于盐值
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claim).SignedString([]byte(sign_key))
	return token, err
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(sign_key), nil //返回签名密钥
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("claim invalid")
	}

	claims, ok := token.Claims.(*MyClaims)
	if !ok {
		return nil, errors.New("invalid claim type")
	}
	return claims, nil
}
