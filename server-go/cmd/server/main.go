package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	httpad "server-go/internal/adapter/http"
	"server-go/internal/adapter/notify"
	mongox "server-go/internal/adapter/persistence/mongo"
	"server-go/internal/adapter/persistence/mysql"
	redisx "server-go/internal/adapter/persistence/redis"
	"server-go/internal/adapter/ws"
	"server-go/internal/config"
	"server-go/internal/platform/captchagen"
	"server-go/internal/platform/clock"
	"server-go/internal/platform/logger"
	"server-go/internal/platform/mailer"
	"server-go/internal/platform/token"
	"server-go/internal/usecase"
)

func main() {
	configPath := flag.String("config", "config.toml", "配置文件路径")
	flag.Parse()

	cfg, err := config.Load(*configPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	log := logger.New(cfg.Log.Level, cfg.Log.Pretty)
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// ---- 基础设施 ----
	db, err := mysql.NewDB(ctx, cfg.MySQL.DSN())
	must(err, "mysql", log)
	defer db.Close()

	rdb, err := redisx.NewClient(ctx, cfg.Redis.Addr, cfg.Redis.Password, cfg.Redis.DB)
	must(err, "redis", log)
	defer rdb.Close()

	mongoCli, err := mongox.NewClient(ctx, cfg.Mongo.URI)
	must(err, "mongo", log)
	defer func() {
		dctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = mongoCli.Disconnect(dctx)
	}()
	mongoDB := mongoCli.Database(cfg.Mongo.Database)

	// ---- 出站适配器 ----
	userRepo := mysql.NewUserRepo(db)
	friendRepo := mysql.NewFriendRepo(db)
	groupRepo := mysql.NewGroupRepo(db)
	msgRepo := mongox.NewMessageRepo(mongoDB, cfg.Message.RetentionDays)
	must(msgRepo.EnsureIndexes(ctx), "mongo 索引", log)

	jwtMgr := token.NewManager(cfg.JWT.Secret, cfg.JWT.TTL())
	tokenStore := redisx.NewTokenStore(rdb, cfg.JWT.TTL())
	captchaStore := redisx.NewCaptchaStore(rdb)
	codeStore := redisx.NewEmailCodeStore(rdb)
	presenceStore := redisx.NewPresenceStore(rdb, time.Duration(cfg.Presence.TTLSeconds)*time.Second)
	recentStore := redisx.NewRecentChatStore(rdb)
	limitStore := redisx.NewMsgLimitStore(rdb)
	mail := mailer.New(cfg.Email.Host, cfg.Email.Port, cfg.Email.Username, cfg.Email.Password, cfg.Email.From)
	captchaGen := captchagen.New()
	clk := clock.Real{}

	hub := ws.NewHub(log)
	go hub.Run(ctx)
	notifier := notify.NewWSNotifier(hub)

	// ---- usecase ----
	authUC := usecase.NewAuth(usecase.AuthDeps{
		Users: userRepo, Tokens: tokenStore, Issuer: jwtMgr,
		Captchas: captchaStore, CaptchaGen: captchaGen, Codes: codeStore,
		Mailer: mail, Notifier: notifier, Clock: clk,
		TokenTTL: cfg.JWT.TTL(), Logger: log,
	})
	userUC := usecase.NewUser(usecase.UserDeps{Users: userRepo, Friends: friendRepo, Codes: codeStore})
	friendUC := usecase.NewFriend(usecase.FriendDeps{Users: userRepo, Friends: friendRepo, Presence: presenceStore})
	groupUC := usecase.NewGroup(usecase.GroupDeps{Groups: groupRepo, Users: userRepo, Notifier: notifier, Clock: clk})
	presenceUC := usecase.NewPresence(usecase.PresenceDeps{
		Presence: presenceStore, Friends: friendRepo, Notifier: notifier,
		TTL: time.Duration(cfg.Presence.TTLSeconds) * time.Second,
	})
	chatUC := usecase.NewChat(usecase.ChatDeps{
		Messages: msgRepo, Users: userRepo, Friends: friendRepo, Groups: groupRepo,
		Recent: recentStore, Limits: limitStore, Presence: presenceStore,
		Notifier: notifier, Clock: clk,
		NonFriendLimit: cfg.Message.NonFriendLimit,
		LimitTTL:       time.Duration(cfg.Message.LimitTTLHours) * time.Hour,
	})

	// ---- HTTP 服务器 ----
	e := httpad.NewServer(httpad.ServerDeps{
		Cfg: cfg, Logger: log, TokenMgr: jwtMgr, Tokens: tokenStore,
		Auth: authUC, User: userUC, Friend: friendUC, Group: groupUC,
		Chat: chatUC, Presence: presenceUC, Hub: hub,
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Server.Port), Handler: e}
	go func() {
		log.Info("HTTP 服务启动", "addr", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("HTTP 服务异常退出", "err", err)
			stop()
		}
	}()

	// ---- 优雅关闭 ----
	<-ctx.Done()
	log.Info("收到退出信号，开始优雅关闭…")
	shutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error("HTTP 关闭超时", "err", err)
	}
	hub.Close()
	log.Info("服务已退出")
}

func must(err error, what string, log *slog.Logger) {
	if err != nil {
		log.Error("初始化失败", "component", what, "err", err)
		os.Exit(1)
	}
}
