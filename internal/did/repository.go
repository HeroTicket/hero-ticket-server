package did

import "context"

type Query interface {
}

type Command interface{}

type Repository interface {
	Query
	Command
	Exec(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
}
