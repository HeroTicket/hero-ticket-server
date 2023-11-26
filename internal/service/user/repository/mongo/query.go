package mongo

import (
	"context"

	"github.com/heroticket/internal/service/user"
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

func (q *MongoQuery) FindAdmin(ctx context.Context) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"isAdmin": true}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (q *MongoQuery) FindUserByID(ctx context.Context, id string) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"_id": id}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (q *MongoQuery) FindUserByAccountAddress(ctx context.Context, accountAddress string) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"accountAddress": accountAddress}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (q *MongoQuery) FindUsers(ctx context.Context) ([]*user.User, error) {
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

func (q *MongoQuery) FindUserByName(ctx context.Context, name string) (*user.User, error) {
	coll := q.collection()

	filter := bson.M{"name": name}

	var u user.User

	if err := coll.FindOne(ctx, filter).Decode(&u); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}

	return &u, nil
}

func (q *MongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection("users")
}
