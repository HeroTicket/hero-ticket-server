package mongo

import (
	"context"
	"time"

	"github.com/heroticket/internal/service/user"
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

func (c *MongoCommand) CreateUser(ctx context.Context, params user.CreateUserParams) (*user.User, error) {
	coll := c.collection()

	var u user.User

	u.ID = params.ID
	u.AccountAddress = params.AccountAddress
	u.TbaAddress = params.TbaAddress
	u.Name = params.Name
	u.Avatar = params.Avatar
	if params.TbaTokenBalance != "" {
		u.TbaTokenBalance = params.TbaTokenBalance
	} else {
		u.TbaTokenBalance = "0"
	}
	u.IsAdmin = params.IsAdmin
	u.CreatedAt = time.Now().Unix()
	u.UpdatedAt = time.Now().Unix()

	_, err := coll.InsertOne(ctx, u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func (c *MongoCommand) DeleteUser(ctx context.Context, id string) error {
	coll := c.collection()

	filter := bson.M{"_id": id}

	_, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *MongoCommand) UpdateUser(ctx context.Context, params user.UpdateUserParams) error {
	coll := c.collection()

	if params.ID == "" {
		return user.ErrInvalidID
	}

	filter := bson.M{"_id": params.ID}

	value := bson.D{}

	if params.Name == "" && params.AccountAddress == "" && params.TBAAddress == "" {
		return user.ErrNothingToUpdate
	}

	if params.Name != "" {
		value = append(value, bson.E{Key: "name", Value: params.Name})
	}

	if params.AccountAddress != "" {
		value = append(value, bson.E{Key: "accountAddress", Value: params.AccountAddress})
	}

	if params.TBAAddress != "" {
		value = append(value, bson.E{Key: "tbaAddress", Value: params.TBAAddress})
	}

	if params.Bio != "" {
		value = append(value, bson.E{Key: "bio", Value: params.Bio})
	}

	if params.Avatar != "" {
		value = append(value, bson.E{Key: "avatar", Value: params.Avatar})
	}

	if params.Banner != "" {
		value = append(value, bson.E{Key: "banner", Value: params.Banner})
	}

	if params.TBATokenBalance != "" {
		value = append(value, bson.E{Key: "tbaTokenBalance", Value: params.TBATokenBalance})
	}

	value = append(value, bson.E{Key: "updated_at", Value: time.Now().Unix()})

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

func (c *MongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection("users")
}
