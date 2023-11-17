package mongo

import (
	"context"

	"github.com/heroticket/internal/did"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoQuery struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewQuery(client *mongo.Client, dbname, collname string) did.Query {
	return &mongoQuery{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

func (q *mongoQuery) FindMatchedVerifier(ctx context.Context, didStr string, walletAddress string, contractAddress string) (*did.Verifier, error) {
	coll := q.collection()

	filter := bson.M{
		"did":              didStr,
		"wallet_address":   walletAddress,
		"contract_address": contractAddress,
	}

	var v did.Verifier

	err := coll.FindOne(ctx, filter).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (q *mongoQuery) FindVerifierByID(ctx context.Context, id string) (*did.Verifier, error) {
	coll := q.collection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectID}

	var v did.Verifier

	err = coll.FindOne(ctx, filter).Decode(&v)
	if err != nil {
		return nil, err
	}

	return &v, nil
}

func (q *mongoQuery) FindVerifiers(ctx context.Context) ([]*did.Verifier, error) {
	coll := q.collection()

	cursor, err := coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var verifiers []*did.Verifier

	for cursor.Next(ctx) {
		var v did.Verifier
		if err := cursor.Decode(&v); err != nil {
			return nil, err
		}

		verifiers = append(verifiers, &v)
	}

	return verifiers, nil
}

func (q *mongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection(q.collname)
}
