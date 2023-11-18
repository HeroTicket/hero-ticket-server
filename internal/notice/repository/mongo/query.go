package mongo

import (
	"context"

	"github.com/heroticket/internal/notice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (q *mongoQuery) GetNotice(ctx context.Context, id string) (*notice.Notice, error) {
	coll := q.collection()

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	filter := primitive.M{"_id": objectID}

	var n notice.Notice

	err = coll.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (q *mongoQuery) GetNotices(ctx context.Context, page, limit int64) ([]*notice.Notice, error) {
	coll := q.collection()

	skip := (page - 1) * limit

	findOptions := &options.FindOptions{
		Skip:  &skip,
		Limit: &limit,
	}

	cursor, err := coll.Find(ctx, primitive.M{}, findOptions)
	if err != nil {
		return nil, err
	}

	var notices []*notice.Notice

	for cursor.Next(ctx) {
		var n notice.Notice

		err := cursor.Decode(&n)
		if err != nil {
			return nil, err
		}

		notices = append(notices, &n)
	}

	return notices, nil
}

func (q *mongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection(q.collname)
}
