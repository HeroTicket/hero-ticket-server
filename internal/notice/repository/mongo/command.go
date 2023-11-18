package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/heroticket/internal/notice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCommand struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewCommand(client *mongo.Client, dbname, collname string) notice.Command {
	return &mongoCommand{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

func (c *mongoCommand) CreateNotice(ctx context.Context, n *notice.Notice) (*notice.Notice, error) {
	coll := c.collection()

	n.CreatedAt = time.Now()
	n.UpdatedAt = time.Now()

	result, err := coll.InsertOne(ctx, n)
	if err != nil {
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert InsertedID to primitive.ObjectID")
	}

	n.ID = objectID.Hex()

	return n, nil
}

func (c *mongoCommand) DeleteNotice(ctx context.Context, id string) error {
	coll := c.collection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := primitive.M{"_id": objectID}

	result, err := coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return notice.ErrNotFound
	}

	return nil
}

func (c *mongoCommand) UpdateNotice(ctx context.Context, params *notice.NoticeUpdateParams) error {
	coll := c.collection()

	objectID, err := primitive.ObjectIDFromHex(params.ID)
	if err != nil {
		return err
	}

	filter := primitive.M{"_id": objectID}

	value := primitive.D{}

	if params.Title == "" && params.Content == "" {
		return notice.ErrNothingToUpdate
	}

	if params.Title != "" {
		value = append(value, primitive.E{Key: "title", Value: params.Title})
	}

	if params.Content != "" {
		value = append(value, primitive.E{Key: "content", Value: params.Content})
	}

	value = append(value, primitive.E{Key: "updated_at", Value: time.Now()})

	update := primitive.D{
		{
			Key:   "$set",
			Value: value,
		},
	}

	result, err := coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return notice.ErrNotFound
	}

	return nil
}

func (c *mongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection(c.collname)
}
