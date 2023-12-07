package main

import (
	"context"
	"fmt"

	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/service/notice"
	nmongo "github.com/heroticket/internal/service/notice/repository/mongo"
)

func main() {
	cfg, err := config.NewServerConfig("./configs/server/config.json")
	if err != nil {
		panic(err)
	}

	client, err := mongo.New(context.Background(), cfg.MongoUrl)
	if err != nil {
		panic(err)
	}

	repo := nmongo.New(client, cfg.Notice.DbName)

	res, err := repo.CreateNotice(context.Background(), &notice.Notice{
		Title:   "Welcome to Hero Ticket",
		Content: "Welcome to Hero Ticket, we are happy to have you here.",
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(res)
}
