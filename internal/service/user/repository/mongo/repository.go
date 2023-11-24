package mongo

import (
	"context"

	"github.com/heroticket/internal/service/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	user.Query
	user.Command
	client *mongo.Client
	dbname string
}

func New(ctx context.Context, client *mongo.Client, dbname, collname string) (user.Repository, error) {
	cmd := NewMongoCommand(client, dbname, collname)
	qry := NewMongoQuery(client, dbname, collname)
	repo := &mongoRepository{
		Query:   qry,
		Command: cmd,
		client:  client,
		dbname:  dbname,
	}

	_, err := cmd.collection().Indexes().CreateOne(
		ctx,
		mongo.IndexModel{
			Keys:    bson.M{"accountAddress": 1},
			Options: options.Index().SetUnique(true),
		},
	)

	return repo, err
}
