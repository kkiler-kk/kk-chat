package token

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"server-go/internal/domain"
)

type Claims struct {
	UserID int64
	JTI    string
}

type Manager struct {
	secret []byte
	ttl    time.Duration
}

func NewManager(secret string, ttl time.Duration) *Manager {
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Issue 签发 token，返回 (tokenStr, jti)。实现 port.TokenIssuer。
func (m *Manager) Issue(_ context.Context, userID int64) (string, string, error) {
	jti := uuid.NewString()
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": userID,
		"jti": jti,
		"exp": now.Add(m.ttl).Unix(),
		"iat": now.Unix(),
	}
	t, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(m.secret)
	if err != nil {
		return "", "", err
	}
	return t, jti, nil
}

// Parse 校验签名与有效期，返回 Claims；任何失败一律 domain.ErrTokenInvalid（不泄露原因）。
func (m *Manager) Parse(tokenStr string) (*Claims, error) {
	t, err := jwt.Parse(tokenStr, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, domain.ErrTokenInvalid
		}
		return m.secret, nil
	}, jwt.WithValidMethods([]string{"HS256"}))
	if err != nil || !t.Valid {
		return nil, domain.ErrTokenInvalid
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}
	uidNum, ok := claims["sub"].(float64)
	if !ok {
		return nil, domain.ErrTokenInvalid
	}
	jti, _ := claims["jti"].(string)
	if jti == "" {
		return nil, errors.New("missing jti")
	}
	return &Claims{UserID: int64(uidNum), JTI: jti}, nil
}
