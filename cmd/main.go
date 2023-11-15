package main

import (
	"context"

	"os"
	"syscall"

	"github.com/heroticket/internal/did"
	dredis "github.com/heroticket/internal/did/cache/redis"
	"github.com/heroticket/internal/http"
	"github.com/heroticket/internal/http/rest"
	"github.com/heroticket/internal/infra/mongo"
	"github.com/heroticket/internal/infra/redis"
	"github.com/heroticket/internal/jwt"
	"github.com/heroticket/pkg/shutdown"
	"go.uber.org/zap"

	_ "github.com/joho/godotenv/autoload"
)

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

	cache, err := redis.NewCache(context.Background(), os.Getenv("REDIS_URL"))
	if err != nil {
		panic(err)
	}

	zap.L().Info("connected to redis")

	didSvc := did.New(dredis.New(cache), os.Getenv("RPC_URL_MUMBAI"))
	jwtSvc, _ := jwt.New("secret1", "secret2")

	server := newServer(
		initUserController(didSvc, jwtSvc),
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

func initUserController(didSvc did.Service, jwtSvc jwt.Service) *rest.UserCtrl {
	return rest.NewUserCtrl(didSvc, jwtSvc)
}

func initTicketController() *rest.TicketCtrl {
	return rest.NewTicketCtrl()
}

func newServer(ctrls ...http.Controller) *http.Server {
	return http.NewServer(http.DefaultConfig(), ctrls...)
}
