package user

import "context"

type Query interface {
	FindUsers(ctx context.Context) ([]*User, error)
	FindUserByDID(ctx context.Context, did string) (*User, error)
	FindUserByWalletAddress(ctx context.Context, walletAddress string) (*User, error)
}

type Command interface {
	CreateUser(ctx context.Context, user *User) (*User, error)
	UpdateUser(ctx context.Context, user *User) (*User, error)
	DeleteUser(ctx context.Context, did string) error
}

type Repository interface {
	Query
	Command
}
