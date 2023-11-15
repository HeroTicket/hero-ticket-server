package mongo

import (
	"context"

	"github.com/heroticket/internal/user"
)

type mongoCommand struct {
}

// CreateUser implements user.Command.
func (*mongoCommand) CreateUser(ctx context.Context, user *user.User) (*user.User, error) {
	panic("unimplemented")
}

// DeleteUser implements user.Command.
func (*mongoCommand) DeleteUser(ctx context.Context, did string) error {
	panic("unimplemented")
}

// UpdateUser implements user.Command.
func (*mongoCommand) UpdateUser(ctx context.Context, user *user.User) (*user.User, error) {
	panic("unimplemented")
}

func NewMongoCommand() user.Command {
	return &mongoCommand{}
}
