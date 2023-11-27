package mongo

import (
	"context"
	"time"

	"github.com/heroticket/internal/service/did"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoRepository struct {
	did.Query
	did.Command
	client *mongo.Client
	dbname string
}

func New(ctx context.Context, client *mongo.Client, dbname string) (did.Repository, error) {
	cmd := NewMongoCommand(client, dbname)
	qry := NewMongoQuery(client, dbname)
	repo := &mongoRepository{
		Query:   qry,
		Command: cmd,
		client:  client,
		dbname:  dbname,
	}

	return repo, nil
}

type mongoQuery struct {
	client *mongo.Client
	dbname string
}

func NewMongoQuery(client *mongo.Client, dbname string) did.Query {
	return &mongoQuery{
		client: client,
		dbname: dbname,
	}
}

func (q *mongoQuery) FindClaim(ctx context.Context, userID string, contractAddress string) (*did.Claim, error) {
	coll := q.collection()

	filter := bson.M{"userId": userID, "contractAddress": contractAddress}

	var claim did.Claim

	err := coll.FindOne(ctx, filter).Decode(&claim)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, did.ErrClaimNotFound
		}
	}

	return &claim, nil
}

func (q *mongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection("claims")
}

type mongoCommand struct {
	client *mongo.Client
	dbname string
}

func NewMongoCommand(client *mongo.Client, dbname string) did.Command {
	return &mongoCommand{
		client: client,
		dbname: dbname,
	}
}

func (c *mongoCommand) SaveClaim(ctx context.Context, params did.SaveClaimParams) (*did.Claim, error) {
	coll := c.collection()

	claim := &did.Claim{
		ID:              params.ID,
		UserID:          params.UserID,
		ContractAddress: params.ContractAddress,
		CreatedAt:       time.Now().Unix(),
		UpdateAt:        time.Now().Unix(),
	}

	_, err := coll.InsertOne(ctx, claim)
	if err != nil {
		return nil, err
	}

	return claim, nil
}

func (c *mongoCommand) collection() *mongo.Collection {
	return c.client.Database(c.dbname).Collection("claims")
}
