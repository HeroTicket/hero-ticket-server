package mongo

import (
	"context"

	"github.com/heroticket/internal/notice"
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

// CreateNotice implements notice.Command.
func (*mongoCommand) CreateNotice(ctx context.Context, n *notice.Notice) (*notice.Notice, error) {
	panic("unimplemented")
}

// DeleteNotice implements notice.Command.
func (*mongoCommand) DeleteNotice(ctx context.Context, id string) error {
	panic("unimplemented")
}

// UpdateNotice implements notice.Command.
func (*mongoCommand) UpdateNotice(ctx context.Context, id string, n *notice.Notice) (*notice.Notice, error) {
	panic("unimplemented")
}
