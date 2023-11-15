package mongo

import (
	"github.com/heroticket/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	user.Query
	user.Command
	client *mongo.Client
	dbname string
}

func NewMongoRepository(client *mongo.Client, dbname, collname string) user.Repository {
	return &mongoRepository{
		Query:   NewMongoQuery(client, dbname, collname),
		Command: NewMongoCommand(client, dbname, collname),
		client:  client,
		dbname:  dbname,
	}
}
