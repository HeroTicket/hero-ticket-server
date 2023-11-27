package main

import (
	"context"
	"os"
	"syscall"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/heroticket/internal/app/shutdown"
	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/logger"
	"github.com/heroticket/internal/web3"
	"github.com/heroticket/pkg/contracts/heroticket"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			logger.Info("recovered from panic", "panic", r)
		}
	}()

	err := logger.New(false, "service", "heroticket-subscriber")
	if err != nil {
		logger.Panic("failed to initialize logger", "error", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var configFile string

	if os.Getenv("GO_ENV") != "production" {
		configFile = "config.dev.json"
	}

	cfg, err := config.NewSubscriberConfig(configFile)
	if err != nil {
		logger.Panic("failed to load config", "error", err)
	}

	client, err := web3.NewClient(ctx, cfg.RpcUrl)
	if err != nil {
		logger.Panic("failed to connect to rpc", "error", err)
	}

	filterer, err := heroticket.NewHeroticketFilterer(common.HexToAddress(cfg.ContractAddress), client)
	if err != nil {
		logger.Panic("failed to initialize filterer", "error", err)
	}

	mintedChan := make(chan *heroticket.HeroticketMinted)

	sub, err := filterer.WatchMinted(&bind.WatchOpts{}, mintedChan)
	if err != nil {
		logger.Panic("failed to watch minted", "error", err)
	}

	go func() {
		for {
			select {
			case err := <-sub.Err():
				if err != nil {
					logger.Error("failed to watch minted", "error", err)
					return
				}
			case event := <-mintedChan:
				println(event.TokenId.String())
			}
		}
	}()

	<-shutdown.GracefulShutdown(func() {
		sub.Unsubscribe()
		close(mintedChan)
	}, syscall.SIGINT, syscall.SIGTERM)
}
