package mongo

import (
	"context"
	"time"

	"github.com/heroticket/internal/service/notice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoCommand struct {
	client *mongo.Client
	dbname string
}

func NewCommand(client *mongo.Client, dbname string) notice.Command {
	return &mongoCommand{
		client: client,
		dbname: dbname,
	}
}

func (c *mongoCommand) CreateNotice(ctx context.Context, n *notice.Notice) (*notice.Notice, error) {
	coll := c.collection()

	// find last inserted id
	var last notice.Notice
	opts := options.FindOne().SetSort(primitive.D{{Key: "_id", Value: -1}})
	err := coll.FindOne(ctx, primitive.M{}, opts).Decode(&last)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	// set id
	if err == mongo.ErrNoDocuments {
		n.ID = 1
	} else {
		n.ID = last.ID + 1
	}

	n.CreatedAt = time.Now().Unix()
	n.UpdatedAt = time.Now().Unix()

	_, err = coll.InsertOne(ctx, n)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func (c *mongoCommand) DeleteNotice(ctx context.Context, id int64) error {
	coll := c.collection()

	filter := primitive.M{"_id": id}

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

	filter := primitive.M{"_id": params.ID}

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

	value = append(value, primitive.E{Key: "updated_at", Value: time.Now().Unix()})

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
	return c.client.Database(c.dbname).Collection("notices")
}
