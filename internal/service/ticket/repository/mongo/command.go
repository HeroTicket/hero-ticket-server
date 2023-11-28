package mongo

import (
	"context"

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

func (c *MongoCommand) CreateTicket(ctx context.Context, params ticket.Ticket) (*ticket.Ticket, error) {
	coll := c.collection()

	var t ticket.Ticket

	t.ID = params.ID
	t.Address = params.Address
	t.OwnerAddress = params.OwnerAddress
	t.TokenID = params.TokenID
	t.Name = params.Name
	t.Symbol = params.Symbol
	t.Image = params.Image
	t.PurchasedAt = params.PurchasedAt

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

func (c *MongoCommand) CreateTicketCollection(ctx context.Context, params ticket.TicketCollection) (*ticket.TicketCollection, error) {
	coll := c.collection()

	var tc ticket.TicketCollection

	tc.ID = params.ID
	tc.CreatorID = params.CreatorID
	tc.Address = params.Address
	tc.Name = params.Name
	tc.Symbol = params.Symbol
	tc.Description = params.Description
	tc.Organizer = params.Organizer
	tc.Location = params.Location
	tc.Date = params.Date
	tc.BannerImage = params.BannerImage
	tc.TicketImage = params.TicketImage
	tc.Price = params.Price
	tc.TotalSupply = params.TotalSupply
	tc.Remaining = params.Remaining
	tc.CreatedAt = params.CreatedAt
	tc.UpdatedAt = params.UpdatedAt

	_, err := coll.InsertOne(ctx, tc)
	if err != nil {
		return nil, err
	}

	return &tc, nil
}

func (c *MongoCommand) DeleteTicketCollection(ctx context.Context, id string) error {
	coll := c.collection()

	filter := bson.M{"_id": id}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *MongoCommand) SaveTicket(ctx context.Context, params ticket.SaveTicketParams) (*ticket.Ticket, error) {
	coll := c.collection()

	var t ticket.Ticket

	t.Address = params.Address
	t.OwnerAddress = params.OwnerAddress
	t.TokenID = params.TokenID
	t.PurchasedAt = params.PurchasedAt

	_, err := coll.InsertOne(ctx, t)
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (c *MongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection("users")
}
