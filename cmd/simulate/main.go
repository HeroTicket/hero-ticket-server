package main

import (
	"context"
	"fmt"

	"github.com/heroticket/internal/config"
	"github.com/heroticket/internal/db/mongo"
	"github.com/heroticket/internal/service/ticket"
	nmongo "github.com/heroticket/internal/service/ticket/repository/mongo"
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

	repo, err := nmongo.New(context.Background(), client, cfg.Ticket.DbName)
	if err != nil {
		panic(err)
	}

	res, err := repo.FindTicketCollections(context.Background(), ticket.TicketCollectionFilter{})
	if err != nil {
		panic(err)
	}

	for _, r := range res {
		fmt.Println(r)
	}
}
