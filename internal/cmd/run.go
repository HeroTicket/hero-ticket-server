package cmd

import (
	"context"
	"net/http"
	"os"
	"strings"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/heroticket/internal/app"
	"github.com/heroticket/internal/app/rest"
	"github.com/heroticket/internal/app/shutdown"
	"github.com/heroticket/internal/cache/redis"
	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/did"
	drepo "github.com/heroticket/internal/service/did/repository/mongo"
	"github.com/heroticket/internal/service/ipfs"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/notice"
	nrepo "github.com/heroticket/internal/service/notice/repository/mongo"
	"github.com/heroticket/internal/service/ticket"
	trepo "github.com/heroticket/internal/service/ticket/repository/mongo"
	"github.com/heroticket/internal/service/user"
	urepo "github.com/heroticket/internal/service/user/repository/mongo"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/pkg/contracts/heroticket"
)

func Run() {
	defer func() {
		if r := recover(); r != nil {
			logger.Info("recovered from panic", "panic", r)
		}
	}()

	err := logger.New(false, "service", "heroticket")
	handleErr(err)

	var configFile string

	if os.Getenv("GO_ENV") != "production" {
		configFile = "config.dev.json"
	}

	cfg, err := config.NewServerConfig(configFile)
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

	auths := auth.New(auth.AuthServiceConfig{
		IPFSUrl:         cfg.Auth.IPFSUrl,
		RPCUrl:          cfg.RpcUrl,
		ContractAddress: cfg.Auth.ContractAddress,
		ResolverPrefix:  cfg.Auth.ResolverPrefix,
		KeyDir:          cfg.Auth.KeyDir,
		ReqCache:        authCache,
	})

	didRepo, err := drepo.New(ctx, mongoClient, cfg.Did.DbName)
	handleErr(err)

	dids := did.New(did.DidServiceConfig{
		RPCUrl:    cfg.RpcUrl,
		IssuerUrl: cfg.Did.IssuerUrl,
		Username:  cfg.Did.Username,
		Password:  cfg.Did.Password,
		QrCache:   didCache,
		Repo:      didRepo,
	})

	ipfss := ipfs.New(ipfs.IpfsServiceConfig{
		ApiKey: cfg.Ipfs.ApiKey,
		Secret: cfg.Ipfs.Secret,
	})

	jwts := jwt.New(cfg.Jwt.AccessTokenKey, cfg.Jwt.RefreshTokenKey,
		jwt.WithAudience(cfg.Jwt.Audience), jwt.WithIssuer(cfg.Jwt.Issuer))

	notices := notice.New(nrepo.New(mongoClient, cfg.Notice.DbName))

	ethclient, err := web3.NewClient(ctx, cfg.RpcUrl)
	handleErr(err)

	heroticketContract, err := heroticket.NewHeroticket(web3.HexToAddress(cfg.Ticket.ContractAddress), ethclient)
	handleErr(err)

	pvk, err := web3.ParsePrivateKey(cfg.Ticket.PrivateKey)
	handleErr(err)

	ticketRepo, err := trepo.New(ctx, mongoClient, cfg.Ticket.DbName)
	handleErr(err)

	tickets := ticket.New(ethclient, heroticketContract, pvk, ticketRepo, cfg.Ticket.MoralisApiKey)

	userRepo, err := urepo.New(ctx, mongoClient, cfg.User.DbName)
	handleErr(err)

	users := user.New(userRepo)

	_ = mongo.NewTx(mongoClient)

	// find admin user
	_, err = users.FindAdmin(ctx)
	if err != nil {
		if err == user.ErrUserNotFound {
			resp, err := dids.CreateIdentity(ctx, did.CreateIdentityRequest{
				DidMetadata: did.DidMetadata{
					Blockchain: "polygon",
					Method:     "polygonid",
					Network:    "mumbai",
					Type:       did.BJJ,
				},
			})
			handleErr(err)

			adminAddress := crypto.PubkeyToAddress(pvk.PublicKey)

			_, err = users.CreateUser(ctx, user.CreateUserParams{
				ID:             resp.Identifier,
				AccountAddress: strings.ToLower(adminAddress.Hex()),
				Name:           "admin",
				Avatar:         "https://ipfs.io/ipfs/QmfFbvLH37DebBqmVBm7V8ecfzgjFPnPeHRYiYk1PNoW84/6level.png",
				IsAdmin:        true,
			})
			handleErr(err)
		}
	}

	claimCtrl := rest.NewClaimCtrl(dids, jwts, tickets, users)
	noticeCtrl := rest.NewNoticeCtrl(notices, users)
	profileCtrl := rest.NewProfileCtrl(tickets, users)
	ticketCtrl := rest.NewTicketCtrl(auths, ipfss, jwts, tickets, users, cfg.ServerUrl)
	userCtrl := rest.NewUserCtrl(auths, jwts, users, tickets, cfg.ServerUrl)

	srv := app.New(app.DefaultConfig(), claimCtrl, noticeCtrl, profileCtrl, ticketCtrl, userCtrl)

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

		logger.Info("Successfully shutdown server")

		err = mongoClient.Disconnect(ctx)
		handleErr(err)

		logger.Info("Successfully disconnected from MongoDB")

		logger.Sync()
	}, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func handleErr(err error) {
	if err != nil {
		logger.Panic(err.Error())
	}
}
