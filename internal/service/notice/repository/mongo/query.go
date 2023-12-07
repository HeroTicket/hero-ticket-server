package mongo

import (
	"context"
	"math"
	"strconv"

	"github.com/heroticket/internal/service/notice"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoQuery struct {
	client *mongo.Client
	dbname string
}

func NewQuery(client *mongo.Client, dbname string) notice.Query {
	return &mongoQuery{
		client: client,
		dbname: dbname,
	}
}

func (q *mongoQuery) GetNotice(ctx context.Context, id string) (*notice.Notice, error) {
	coll := q.collection()

	idInt, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		return nil, err
	}

	filter := primitive.M{"_id": idInt}

	var n notice.Notice

	err = coll.FindOne(ctx, filter).Decode(&n)
	if err != nil {
		return nil, err
	}

	return &n, nil
}

func (q *mongoQuery) GetNotices(ctx context.Context, page, limit int64) (*notice.Notices, error) {
	coll := q.collection()

	total, err := coll.CountDocuments(ctx, primitive.M{})
	if err != nil {
		return nil, err
	}

	pagination := q.pagination(total, page, limit)

	if total == 0 {
		return &notice.Notices{
			Items:      []*notice.Notice{},
			Pagination: pagination,
		}, nil
	}

	skip := (pagination.CurrentPage - 1) * pagination.Limit

	findOptions := &options.FindOptions{
		Skip:  &skip,
		Limit: &pagination.Limit,
		Sort:  primitive.D{{Key: "_id", Value: -1}},
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

	return &notice.Notices{
		Items:      notices,
		Pagination: pagination,
	}, nil
}

func (q *mongoQuery) collection() *mongo.Collection {
	return q.client.Database(q.dbname).Collection("notices")
}

func (q *mongoQuery) pagination(total, page, limit int64) *notice.Pagination {
	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 1
	}

	if limit > 20 {
		limit = 20
	}

	pages := int64(math.Ceil(float64(total) / float64(limit)))

	if page > pages {
		page = pages
	}

	// 2 pages before and after the current page
	start := max(1, page-2)
	end := min(pages, page+2)

	if end < start {
		end = start
	}

	hasNext := page < pages
	hasPrev := page > 1

	return &notice.Pagination{
		Total:       total,
		Pages:       pages,
		CurrentPage: page,
		Limit:       limit,
		Start:       start,
		End:         end,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
	}
}
