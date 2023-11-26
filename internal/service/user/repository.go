package user

import "context"

type Query interface {
	FindAdmin(ctx context.Context) (*User, error)
	FindUsers(ctx context.Context) ([]*User, error)
	FindUserByID(ctx context.Context, id string) (*User, error)
	FindUserByAccountAddress(ctx context.Context, accountAddress string) (*User, error)
	FindUserByName(ctx context.Context, name string) (*User, error)
}

type Command interface {
	CreateUser(ctx context.Context, params CreateUserParams) (*User, error)
	UpdateUser(ctx context.Context, params UpdateUserParams) error
	DeleteUser(ctx context.Context, id string) error
}

type Repository interface {
	Query
	Command
}
