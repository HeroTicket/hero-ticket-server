package mongo

import (
	"context"

	"github.com/heroticket/internal/service/ticket"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	ticket.Query
	ticket.Command
	client *mongo.Client
	dbname string
}

func New(ctx context.Context, client *mongo.Client, dbname string) (ticket.Repository, error) {
	cmd := NewMongoCommand(client, dbname)
	qry := NewMongoQuery(client, dbname)
	repo := &mongoRepository{
		Query:   qry,
		Command: cmd,
		client:  client,
		dbname:  dbname,
	}

	return repo, nil
}
