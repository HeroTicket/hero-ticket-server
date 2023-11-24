package mongo

import (
	"github.com/heroticket/internal/service/notice"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	notice.Command
	notice.Query
	client *mongo.Client
	dbname string
}

func New(client *mongo.Client, dbname string) notice.Repository {
	return &mongoRepository{
		Command: NewCommand(client, dbname),
		Query:   NewQuery(client, dbname),
		client:  client,
		dbname:  dbname,
	}
}
