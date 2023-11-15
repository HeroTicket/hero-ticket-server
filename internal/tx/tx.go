package tx

import "context"

type Tx interface {
	Exec(ctx context.Context, fn func(ctx context.Context) (interface{}, error)) (interface{}, error)
}
