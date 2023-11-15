package mongo

import (
	"context"

	"github.com/heroticket/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoQuery struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewMongoQuery(client *mongo.Client, dbname, collname string) user.Query {
	return &mongoQuery{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

func (q *mongoQuery) FindUserByDID(ctx context.Context, did string) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"_id": did}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (q *mongoQuery) FindUserByWalletAddress(ctx context.Context, walletAddress string) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"wallet_address": walletAddress}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		return nil, err
	}

	return &u, nil
}

func (q *mongoQuery) FindUsers(ctx context.Context) ([]*user.User, error) {
	coll := q.collection()

	filter := bson.M{}

	cur, err := coll.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var users []*user.User

	for cur.Next(ctx) {
		var u user.User

		if err := cur.Decode(&u); err != nil {
			return nil, err
		}

		users = append(users, &u)
	}

	return users, nil
}

func (q *mongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection(q.collname)
}
