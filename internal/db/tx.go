package db

import "context"

type TxFn func(ctx context.Context) (interface{}, error)

type Tx interface {
	Exec(ctx context.Context, fn TxFn) (interface{}, error)
}
