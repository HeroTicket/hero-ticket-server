package mongo

import (
	"context"

	"github.com/heroticket/internal/service/ticket"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
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

func (q *MongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection("tickets")
}
