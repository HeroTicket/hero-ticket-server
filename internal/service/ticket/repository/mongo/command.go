package mongo

import (
	"context"
	"time"

	"github.com/heroticket/internal/service/ticket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MongoCommand struct {
	client *mongo.Client
	dbname string
}

func NewMongoCommand(client *mongo.Client, dbname string) *MongoCommand {
	return &MongoCommand{
		client: client,
		dbname: dbname,
	}
}

func (c *MongoCommand) CreateTicket(ctx context.Context, params ticket.CreateTicketParams) (*ticket.Ticket, error) {
	coll := c.collection()

	var t ticket.Ticket

	t.ID = params.ID
	t.AccountAddress = params.AccountAddress
	t.TbaAddress = params.TbaAddress
	t.Name = params.Name
	t.Avatar = params.Avatar
	t.IsAdmin = params.IsAdmin
	t.CreatedAt = time.Now().Unix()
	t.UpdatedAt = time.Now().Unix()

	_, err := coll.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *MongoCommand) DeleteTicket(ctx context.Context, id string) error {
	coll := c.collection()

	filter := bson.M{"_id": id}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *MongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection("users")
}
