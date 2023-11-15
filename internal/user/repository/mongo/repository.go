package mongo

import (
	"github.com/heroticket/internal/user"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	user.Query
	user.Command
	client *mongo.Client
	name   string
}

func NewMongoRepository(client *mongo.Client, name string) user.Repository {
	return &mongoRepository{
		client: client,
		name:   name,
	}
}
