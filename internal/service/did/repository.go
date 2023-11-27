package did

import "context"

type Query interface {
	FindClaim(ctx context.Context, userID, contractAddress string) (*Claim, error)
}

type Command interface {
	SaveClaim(ctx context.Context, params SaveClaimParams) (*Claim, error)
}

type Repository interface {
	Query
	Command
}
