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

func New(client *mongo.Client, dbname, collname string) notice.Repository {
	return &mongoRepository{
		Command: NewCommand(client, dbname, collname),
		Query:   NewQuery(client, dbname, collname),
		client:  client,
		dbname:  dbname,
	}
}
