package main

import (
	"context"
	"syscall"

	"os"

	"github.com/heroticket/internal/cache/redis"
	"github.com/heroticket/internal/db"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/http"
	"github.com/heroticket/internal/http/rest"
	"github.com/heroticket/internal/service/auth"
	"github.com/heroticket/internal/service/did"
	"github.com/heroticket/internal/service/jwt"
	"github.com/heroticket/internal/service/notice"
	noticerepo "github.com/heroticket/internal/service/notice/repository/mongo"
	"github.com/heroticket/internal/service/user"
	userrepo "github.com/heroticket/internal/service/user/repository/mongo"
	"github.com/heroticket/internal/shutdown"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

//	@title			Hero Ticket API
//	@version		1.0
//	@description	API for Hero Ticket DApp
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @host		api.heroticket.xyz
// @BasePath	/v1
func main() {
	logger, _ := zap.NewProduction(zap.Fields(zap.String("service", "hero-ticket")))
	defer logger.Sync()

	zap.ReplaceGlobals(logger)

	zap.L().Info("starting server")

	defer func() {
		if r := recover(); r != nil {
			zap.L().Info("recovered from panic", zap.Any("r", r))
		}
	}()

	client, err := mongo.New(context.Background(), os.Getenv("MONGO_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to mongo")

	authRedis, err := redis.New(context.Background(), os.Getenv("AUTH_REDIS_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to auth redis")

	didRedis, err := redis.New(context.Background(), os.Getenv("DID_REDIS_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to did redis")

	noticeRepo := noticerepo.New(client, "hero-ticket", "notices")
	userRepo := userrepo.NewMongoRepository(client, "hero-ticket", "users")

	authSvc := auth.New(auth.AuthServiceConfig{
		ReqCache: redis.NewCache(authRedis),
	})

	didSvc := did.New(did.DidServiceConfig{
		RPCUrl:    os.Getenv("RPC_URL"),
		IssuerUrl: "https://issuer.heroticket.xyz",
		Username:  "user-issuer",
		Password:  "password-issuer",

		QrCache: redis.NewCache(didRedis),
		Client:  did.DefaultClient,
	})

	jwtSvc := jwt.New("secret")

	noticeSvc := notice.New(noticeRepo)
	userSvc := user.New(userRepo)
	tx := mongo.NewTx(client)

	server := newServer(
		initUserController(authSvc, didSvc, jwtSvc, userSvc, tx),
		initNoticeController(noticeSvc, userSvc),
		initTicketController(),
	)

	go func() {
		if err := server.Run(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()

	stop := shutdown.GracefulShutdown(func() {
		if err := server.Shutdown(context.Background()); err != nil {
			panic(err)
		}

		zap.L().Info("server stopped")

		if err := client.Disconnect(context.Background()); err != nil {
			panic(err)
		}

		zap.L().Info("disconnected from mongo")
	}, syscall.SIGINT, syscall.SIGTERM)

	<-stop
}

func initUserController(authSvc auth.Service, didSvc did.Service, jwtSvc jwt.Service, userSvc user.Service, tx db.Tx) *rest.UserCtrl {
	return rest.NewUserCtrl(authSvc, didSvc, jwtSvc, userSvc, tx, os.Getenv("BASE_URL"))
}

func initNoticeController(noticeSvc notice.Service, userSvc user.Service) *rest.NoticeCtrl {
	return rest.NewNoticeCtrl(noticeSvc, userSvc)
}

func initTicketController() *rest.TicketCtrl {
	return rest.NewTicketCtrl()
}

func newServer(ctrls ...http.Controller) *http.Server {
	return http.NewServer(http.DefaultConfig(), ctrls...)
}
