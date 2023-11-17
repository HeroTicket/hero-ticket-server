package mongo

import (
	"context"

	"github.com/heroticket/internal/did"
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

// FindMatchedVerifier implements did.Query.
func (*mongoQuery) FindMatchedVerifier(ctx context.Context, did string, walletAddress string, contractAddress string) (*did.Verifier, error) {
	panic("unimplemented")
}

// FindVerifierByID implements did.Query.
func (*mongoQuery) FindVerifierByID(ctx context.Context, id string) (*did.Verifier, error) {
	panic("unimplemented")
}

// FindVerifiers implements did.Query.
func (*mongoQuery) FindVerifiers(ctx context.Context) ([]*did.Verifier, error) {
	panic("unimplemented")
}
