package authUtil

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"server-go/internal/app/core/config"
	"time"
)

type MyClaims struct {
	ID int64 `json:"id"`
	jwt.StandardClaims
}

// 定义过期时间 2小时

// 定义secret
var MySecret = []byte(config.Instance().Token.EncryptKey)

// 生成jwt
func GenToken(ID int64) (string, error) {
	c := MyClaims{
		ID,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(config.Instance().Token.TimeOut).Unix(), // config 获取是过期时间
			Issuer:    "kk-chat",
		},
	}
	//使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, c)
	//使用指定的secret签名并获得完成的编码后的字符串token
	return token.SignedString(MySecret)
}

// 解析JWT
func ParseToken(tokenString string) (*MyClaims, error) {
	//解析token
	token, err := jwt.ParseWithClaims(tokenString, &MyClaims{}, func(token *jwt.Token) (interface{}, error) {
		return MySecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*MyClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
