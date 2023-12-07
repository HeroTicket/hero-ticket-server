package mongo

import (
	"context"

	"github.com/heroticket/internal/service/ticket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoQuery struct {
	client *mongo.Client
	dbname string
}

func NewMongoQuery(client *mongo.Client, dbname string) *MongoQuery {
	return &MongoQuery{
		client: client,
		dbname: dbname,
	}
}

func (q *MongoQuery) FindTicketByAddress(ctx context.Context, address string) (*ticket.Ticket, error) {
	coll := q.collection()

	filter := bson.M{"address": address}

	var t ticket.Ticket

	if err := coll.FindOne(ctx, filter).Decode(&t); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ticket.ErrTicketNotFound
		}
		return nil, err
	}

	return &t, nil
}

func (q *MongoQuery) FindTicketByOwnerAddress(ctx context.Context, ownerAddress string) ([]*ticket.Ticket, error) {
	coll := q.collection()

	filter := bson.M{"ownerAddress": ownerAddress}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var tickets []*ticket.Ticket

	for cur.Next(ctx) {
		var t ticket.Ticket

		if err := cur.Decode(&t); err != nil {
			return nil, err
		}

		tickets = append(tickets, &t)
	}

	return tickets, nil
}

func (q *MongoQuery) FindTicketCollectionByContractAddress(ctx context.Context, contractAddress string) (*ticket.TicketCollection, error) {
	coll := q.collection()

	filter := bson.M{"contractAddress": contractAddress}

	var tc ticket.TicketCollection

	if err := coll.FindOne(ctx, filter).Decode(&tc); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ticket.ErrTicketCollectionNotFound
		}
		return nil, err
	}

	return &tc, nil
}

func (q *MongoQuery) FindTicketCollections(ctx context.Context, filter ticket.TicketCollectionFilter) ([]*ticket.TicketCollection, error) {
	coll := q.collection()

	var query bson.M

	// TODO: add filter
	if filter.IssuerAddress != "" {
		query = bson.M{"issuerAddress": filter.IssuerAddress}
	}

	// decreasing order
	opts := &options.FindOptions{
		Sort: bson.M{"createdAt": -1},
	}

	cur, err := coll.Find(ctx, query, opts)
	if err != nil {
		return nil, err
	}

	var tcs []*ticket.TicketCollection

	for cur.Next(ctx) {
		var tc ticket.TicketCollection

		if err := cur.Decode(&tc); err != nil {
			return nil, err
		}

		tcs = append(tcs, &tc)
	}

	return tcs, nil
}

func (q *MongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection("tickets")
}
