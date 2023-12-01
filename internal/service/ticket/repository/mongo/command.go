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

func (c *MongoCommand) CreateTicketCollection(ctx context.Context, params ticket.CreateTicketCollectionParams) (*ticket.TicketCollection, error) {
	coll := c.collection()

	var tc ticket.TicketCollection

	tc.IssuerAddress = params.IssuerAddress
	tc.ContractAddress = params.ContractAddress
	tc.Name = params.Name
	tc.Symbol = params.Symbol
	tc.Description = params.Description
	tc.Organizer = params.Organizer
	tc.Location = params.Location
	tc.Date = params.Date
	tc.BannerUrl = params.BannerUrl
	tc.TicketUrl = params.TicketUrl
	tc.EthPrice = params.EthPrice
	tc.TokenPrice = params.TokenPrice
	tc.TotalSupply = params.TotalSupply
	tc.Remaining = params.Remaining
	tc.SaleStartAt = params.SaleStartAt
	tc.SaleEndAt = params.SaleEndAt
	tc.CreatedAt = time.Now().Unix()
	tc.UpdatedAt = time.Now().Unix()

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

// func (c *MongoCommand) CreateTBA(ctx context.Context, params ticket.TbaAddresses) (*ticket.TbaAddresses, error) {
// 	coll := c.collection()

// 	var tba ticket.TbaAddresses

// 	tba.ID = params.ID
// 	tba.OwnerAddress = params.OwnerAddress
// 	tba.TbaAddress = params.TbaAddress
// 	tba.TokenID = params.TokenID
// 	tba.Image = params.Image

// 	_, err := coll.InsertOne(ctx, tba)
// 	if err != nil {
// 		return nil, err
// 	}

// 	return &tba, nil
// }

func (c *MongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection("tickets")
}
