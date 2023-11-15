package mongo

import (
	"context"

	"github.com/heroticket/internal/user"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCommand struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewMongoCommand(client *mongo.Client, dbname, collname string) user.Command {
	return &mongoCommand{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

func (c *mongoCommand) CreateUser(ctx context.Context, u *user.User) (*user.User, error) {
	coll := c.collection()

	_, err := coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	return u, nil
}

func (c *mongoCommand) DeleteUser(ctx context.Context, did string) error {
	coll := c.collection()

	filter := bson.M{"_id": did}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *mongoCommand) UpdateUser(ctx context.Context, u *user.User) (*user.User, error) {
	coll := c.collection()

	filter := bson.M{"_id": u.DID}

	update := bson.D{
		{
			Key: "$set",
			Value: bson.D{
				{Key: "", Value: ""},
			},
		},
	}

	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	if res.ModifiedCount == 0 {
		return nil, user.ErrUserNotFound
	}

	return u, nil
}

func (c *mongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection(c.collname)
}
