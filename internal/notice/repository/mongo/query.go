package mongo

import (
	"context"

	"github.com/heroticket/internal/notice"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoQuery struct {
	client   *mongo.Client
	dbname   string
	collname string
}

func NewQuery(client *mongo.Client, dbname, collname string) notice.Query {
	return &mongoQuery{
		client:   client,
		dbname:   dbname,
		collname: collname,
	}
}

// GetNotice implements notice.Query.
func (*mongoQuery) GetNotice(ctx context.Context, id string) (*notice.Notice, error) {
	panic("unimplemented")
}

// GetNotices implements notice.Query.
func (*mongoQuery) GetNotices(ctx context.Context, page int, limit int) ([]*notice.Notice, error) {
	panic("unimplemented")
}
