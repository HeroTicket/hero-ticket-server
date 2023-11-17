package mongo

import (
	"context"
	"errors"
	"time"

	"github.com/heroticket/internal/did"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoCommand struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewCommand(client *mongo.Client, dbname, collname string) did.Command {
	return &mongoCommand{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

func (c *mongoCommand) CreateVerifier(ctx context.Context, v *did.Verifier) (*did.Verifier, error) {
	coll := c.collection()

	v.CreatedAt = time.Now()
	v.UpdatedAt = time.Now()

	result, err := coll.InsertOne(ctx, v)
	if err != nil {
		return nil, err
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("failed to convert InsertedID to primitive.ObjectID")
	}

	v.ID = objectID.Hex()

	return v, nil
}

func (c *mongoCommand) DeleteVerifier(ctx context.Context, id string) error {
	coll := c.collection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := primitive.M{"_id": objectID}

	_, err = coll.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	return nil
}

func (c *mongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection(c.collname)
}
