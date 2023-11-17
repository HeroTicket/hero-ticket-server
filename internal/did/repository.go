package did

import "context"

type Query interface {
	FindVerifiers(ctx context.Context) ([]*Verifier, error)
	FindVerifierByID(ctx context.Context, id string) (*Verifier, error)
	FindMatchedVerifier(ctx context.Context, didStr, walletAddress, contractAddress string) (*Verifier, error)
}

type Command interface {
	CreateVerifier(ctx context.Context, v *Verifier) (*Verifier, error)
	DeleteVerifier(ctx context.Context, id string) error
}

type Repository interface {
	Query
	Command
}
