package mongo

import (
	"github.com/heroticket/internal/did"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	did.Query
	did.Command
	client *mongo.Client
	dbname string
}

func NewRepository(client *mongo.Client, dbname, collname string) did.Repository {
	return &mongoRepository{
		Query:   NewQuery(client, dbname, collname),
		Command: NewCommand(client, dbname, collname),
		client:  client,
		dbname:  dbname,
	}
}
