package mongo

import (
	"context"

	"github.com/heroticket/internal/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type mongoTx struct {
	client *mongo.Client
}

func NewTx(client *mongo.Client) db.Tx {
	return &mongoTx{
		client: client,
	}
}

func (tx *mongoTx) Exec(ctx context.Context, fn db.TxFn) (interface{}, error) {
	var s mongo.Session

	if ctx == nil { // If context is nil, create a new one
		ctx = context.Background()
	} else { // If context is not nil, get the session from context
		s = mongo.SessionFromContext(ctx)
	}
	if s == nil { // If session is nil, create a new one
		ns, err := tx.client.StartSession()
		if err != nil {
			return nil, err
		}

		s = ns
	}
	defer s.EndSession(ctx) // End session after function is done

	// Start transaction
	return s.WithTransaction(ctx, func(ctx mongo.SessionContext) (interface{}, error) {
		return fn(ctx)
	})
}
