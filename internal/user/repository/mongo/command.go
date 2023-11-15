package mongo

import (
	"context"
	"time"

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

	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

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

func (c *mongoCommand) UpdateUser(ctx context.Context, params *user.UserUpdateParams) error {
	coll := c.collection()

	if params.DID == "" {
		return user.ErrInvalidID
	}

	filter := bson.M{"_id": params.DID}

	value := bson.D{}

	if params.Name == "" && params.WalletAddress == "" && params.TBAAddress == "" {
		return user.ErrNothingToUpdate
	}

	if params.Name != "" {
		value = append(value, bson.E{Key: "name", Value: params.Name})
	}

	if params.WalletAddress != "" {
		value = append(value, bson.E{Key: "wallet_address", Value: params.WalletAddress})
	}

	if params.TBAAddress != "" {
		value = append(value, bson.E{Key: "tba_address", Value: params.TBAAddress})
	}

	value = append(value, bson.E{Key: "updated_at", Value: time.Now()})

	update := bson.D{
		{
			Key:   "$set",
			Value: value,
		},
	}

	res, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if res.ModifiedCount == 0 {
		return user.ErrUserNotFound
	}

	return nil
}

func (c *mongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection(c.collname)
}
