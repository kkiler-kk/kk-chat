package usecase

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"server-go/internal/domain"
	"server-go/internal/platform/password"
	"server-go/internal/platform/random"
	"server-go/internal/usecase/apperror"
	"server-go/internal/usecase/port"
)

type AuthDeps struct {
	Users      port.UserRepo
	Tokens     port.TokenStore
	Issuer     port.TokenIssuer
	Captchas   port.CaptchaStore
	CaptchaGen port.CaptchaGenerator
	Codes      port.EmailCodeStore
	Mailer     port.EmailSender
	Notifier   port.Notifier
	Clock      port.Clock
	TokenTTL   time.Duration
	Logger     *slog.Logger
}

type authUseCase struct{ d AuthDeps }

func NewAuth(d AuthDeps) port.AuthUseCase { return &authUseCase{d: d} }

func (u *authUseCase) Register(ctx context.Context, in port.RegisterInput) error {
	ok, err := u.d.Codes.VerifyAndDelete(ctx, in.Email, in.EmailCode)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "注册失败", err)
	}
	if !ok {
		return apperror.New(apperror.CodeEmailCodeInvalid, domain.ErrEmailCodeInvalid.Error())
	}
	if _, err := u.d.Users.ByEmail(ctx, in.Email); err == nil {
		return apperror.New(apperror.CodeEmailTaken, domain.ErrEmailTaken.Error())
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return apperror.Wrap(apperror.CodeInternal, "注册失败", err)
	}
	if _, err := u.d.Users.ByIdentity(ctx, in.Identity); err == nil {
		return apperror.New(apperror.CodeIdentityTaken, domain.ErrIdentityTaken.Error())
	} else if !errors.Is(err, domain.ErrUserNotFound) {
		return apperror.Wrap(apperror.CodeInternal, "注册失败", err)
	}
	hash, err := password.Hash(in.Password)
	if err != nil {
		return apperror.Wrap(apperror.CodeInternal, "注册失败", err)
	}
	user := &domain.User{
		Identity: in.Identity,
		Name:     in.Name,
		Password: hash,
		Email:    in.Email,
		Status:   domain.UserStatusActive,
	}
	if err := u.d.Users.Create(ctx, user); err != nil {
		if errors.Is(err, domain.ErrEmailTaken) {
			return apperror.New(apperror.CodeEmailTaken, domain.ErrEmailTaken.Error())
		}
		return apperror.Wrap(apperror.CodeInternal, "注册失败", err)
	}
	return nil
}

func (u *authUseCase) Login(ctx context.Context, in port.LoginInput) (*port.LoginOutput, error) {
	ok, err := u.d.Captchas.VerifyAndDelete(ctx, in.CaptchaID, in.CaptchaAnswer)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "登录失败", err)
	}
	if !ok {
		return nil, apperror.New(apperror.CodeCaptchaInvalid, domain.ErrCaptchaInvalid.Error())
	}
	// 账号定位：identity 优先，其次 email；找不到统一返回"账号或密码错误"（不暴露账号是否存在）
	user, err := u.d.Users.ByIdentity(ctx, in.Account)
	if errors.Is(err, domain.ErrUserNotFound) {
		user, err = u.d.Users.ByEmail(ctx, in.Account)
	}
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, apperror.New(apperror.CodeBadCredentials, domain.ErrBadCredentials.Error())
		}
		return nil, apperror.Wrap(apperror.CodeInternal, "登录失败", err)
	}
	if user.Banned() {
		return nil, apperror.New(apperror.CodeUserBanned, domain.ErrUserBanned.Error())
	}
	if !password.Verify(user.Password, in.Password) {
		return nil, apperror.New(apperror.CodeBadCredentials, domain.ErrBadCredentials.Error())
	}
	tokenStr, jti, err := u.d.Issuer.Issue(ctx, user.ID)
	if err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "登录失败", err)
	}
	if err := u.d.Tokens.Save(ctx, jti, user.ID, u.d.TokenTTL); err != nil {
		return nil, apperror.Wrap(apperror.CodeInternal, "登录失败", err)
	}
	return &port.LoginOutput{Token: tokenStr, User: ToUserInfo(user)}, nil
}

func (u *authUseCase) Logout(ctx context.Context, jti string, userID int64) error {
	if err := u.d.Tokens.Delete(ctx, jti); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "退出失败", err)
	}
	// 将该用户所有在线端踢下线（前端收到 system.kick 后断开并清理本地会话）
	u.d.Notifier.ToUser(ctx, userID, port.NotifierEvent{
		Event: "system.kick",
		Data:  map[string]any{"reason": "logout"},
	})
	return nil
}

func (u *authUseCase) GenerateCaptcha(ctx context.Context) (id, b64Image string, err error) {
	id, answer, b64Image, err := u.d.CaptchaGen.Generate(ctx)
	if err != nil {
		return "", "", apperror.Wrap(apperror.CodeInternal, "生成验证码失败", err)
	}
	if err := u.d.Captchas.Save(ctx, id, answer, 5*time.Minute); err != nil {
		return "", "", apperror.Wrap(apperror.CodeInternal, "生成验证码失败", err)
	}
	return id, b64Image, nil
}

func (u *authUseCase) SendEmailCode(ctx context.Context, email string) error {
	code := random.Digits(6)
	if err := u.d.Codes.Save(ctx, email, code, 10*time.Minute); err != nil {
		return apperror.Wrap(apperror.CodeInternal, "验证码发送失败", err)
	}
	body := fmt.Sprintf("您的 kk-chat 验证码是 <b>%s</b>，10 分钟内有效。如非本人操作请忽略。", code)
	if err := u.d.Mailer.Send(ctx, email, "kk-chat 邮箱验证码", body); err != nil {
		u.d.Logger.Error("发送验证码邮件失败", "err", err) // 不输出验证码明文
		return apperror.Wrap(apperror.CodeInternal, "验证码发送失败", err)
	}
	return nil
}

var _ port.AuthUseCase = (*authUseCase)(nil)
