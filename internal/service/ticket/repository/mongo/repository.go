package mongo

import (
	"context"

	"github.com/heroticket/internal/service/ticket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	_, err := cmd.collection().Indexes().CreateMany(
		ctx,
		[]mongo.IndexModel{
			{
				Keys:    bson.M{"Address": 1},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.M{"OwnerAddress": 1},
				Options: options.Index().SetUnique(true),
			},
		},
	)

	return repo, err
}
