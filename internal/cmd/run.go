package cmd

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/heroticket/internal/cache/redis"
	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/http"
	"github.com/heroticket/internal/http/rest"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/ipfs"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/notice"
	nrepo "github.com/heroticket/internal/service/notice/repository/mongo"
	"github.com/heroticket/internal/service/user"
	urepo "github.com/heroticket/internal/service/user/repository/mongo"
	"github.com/heroticket/internal/shutdown"
	"go.uber.org/zap"
)

func Run() {
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "heroticket")))
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	defer func() {
		if r := recover(); r != nil {
			logger.Info("recovered from panic", zap.Any("r", r))
		}
	}()

	var configFile string

	goEnv := os.Getenv("GO_ENV")
	if goEnv == "" {
		configFile = "config.dev.json"
	}

	cfg, err := config.NewConfig(configFile)
	if err != nil {
		panic(err)
	}

	logger.Info("Successfully loaded config")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongo.New(ctx, cfg.MongoUrl)
	handleErr(err)

	logger.Info("Successfully connected to MongoDB")

	authRedis, err := redis.New(ctx, cfg.Auth.RedisUrl)
	handleErr(err)

	logger.Info("Successfully connected to Redis for Auth")

	didRedis, err := redis.New(ctx, cfg.Did.RedisUrl)
	handleErr(err)

	logger.Info("Successfully connected to Redis for DID")

	authCache := redis.NewCache(authRedis)
	didCache := redis.NewCache(didRedis)

	authSvc := auth.New(auth.AuthServiceConfig{
		IPFSUrl:         cfg.Auth.IPFSUrl,
		RPCUrl:          cfg.RpcUrl,
		ContractAddress: cfg.Auth.ContractAddress,
		ResolverPrefix:  cfg.Auth.ResolverPrefix,
		KeyDir:          cfg.Auth.KeyDir,
		ReqCache:        authCache,
	})

	_ = did.New(did.DidServiceConfig{
		RPCUrl:    cfg.RpcUrl,
		IssuerUrl: cfg.Did.IssuerUrl,
		Username:  cfg.Did.Username,
		Password:  cfg.Did.Password,
		QrCache:   didCache,
	})

	_ = ipfs.New(ipfs.IpfsServiceConfig{
		ApiKey: cfg.Ipfs.ApiKey,
		Secret: cfg.Ipfs.Secret,
	})

	jwtSvc := jwt.New(cfg.Jwt.AccessTokenKey, cfg.Jwt.RefreshTokenKey,
		jwt.WithAudience(cfg.Jwt.Audience), jwt.WithIssuer(cfg.Jwt.Issuer))

	noticeSvc := notice.New(nrepo.New(mongoClient, cfg.Notice.DbName))

	// TODO: add ticket service
	userRepo, err := urepo.New(ctx, mongoClient, cfg.User.DbName)
	handleErr(err)

	userSvc := user.New(userRepo)

	tx := mongo.NewTx(mongoClient)

	userCtrl := rest.NewUserCtrl(authSvc, jwtSvc, userSvc, tx, cfg.ServerUrl)
	noticeCtrl := rest.NewNoticeCtrl(noticeSvc, userSvc)
	ticketCtrl := rest.NewTicketCtrl()

	srv := http.NewServer(http.DefaultConfig(), userCtrl, noticeCtrl, ticketCtrl)

	logger.Info("Starting server")

	go func() {
		if err := srv.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := shutdown.GracefulShutdown(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		err := srv.Shutdown(ctx)
		handleErr(err)

		err = mongoClient.Disconnect(ctx)
		handleErr(err)

		logger.Info("Successfully shutdown server")
	}, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
