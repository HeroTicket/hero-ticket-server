package mongo

import (
	"context"

	"github.com/heroticket/internal/did"
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

// CreateVerifier implements did.Command.
func (*mongoCommand) CreateVerifier(ctx context.Context, verifier *did.Verifier) (*did.Verifier, error) {
	panic("unimplemented")
}

// DeleteVerifier implements did.Command.
func (*mongoCommand) DeleteVerifier(ctx context.Context, id string) error {
	panic("unimplemented")
}
