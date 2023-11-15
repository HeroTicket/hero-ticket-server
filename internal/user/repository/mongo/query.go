package mongo

import (
	"context"

	"github.com/heroticket/internal/user"
)

type mongoQuery struct {
}

// FindUserByDID implements user.Query.
func (*mongoQuery) FindUserByDID(ctx context.Context, did string) (*user.User, error) {
	panic("unimplemented")
}

// FindUserByWalletAddress implements user.Query.
func (*mongoQuery) FindUserByWalletAddress(ctx context.Context, walletAddress string) (*user.User, error) {
	panic("unimplemented")
}

// FindUsers implements user.Query.
func (*mongoQuery) FindUsers(ctx context.Context) ([]*user.User, error) {
	panic("unimplemented")
}

func NewMongoQuery() user.Query {
	return &mongoQuery{}
}
