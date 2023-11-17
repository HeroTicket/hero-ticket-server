package did

import "context"

type Query interface {
	FindVerifiers(ctx context.Context) ([]*Verifier, error)
	FindVerifierByID(ctx context.Context, id string) (*Verifier, error)
	FindMatchedVerifier(ctx context.Context, did, walletAddress, contractAddress string) (*Verifier, error)
}

type Command interface {
	CreateVerifier(ctx context.Context, verifier *Verifier) (*Verifier, error)
	DeleteVerifier(ctx context.Context, id string) error
}

type Repository interface {
	Query
	Command
}
