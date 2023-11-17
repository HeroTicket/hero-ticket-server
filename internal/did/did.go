package did

import (
	"errors"
	"time"
)

var (
	ErrRequestNotFound = errors.New("request not found")
)

var DefaultCacheExpiry = 10 * time.Minute

type Verifier struct {
	ID              string    `json:"id" bson:"_id,omitempty"`
	DID             string    `json:"did" bson:"did"`
	WalletAddress   string    `json:"wallet_address" bson:"wallet_address"`
	ContractAddress string    `json:"contract_address" bson:"contract_address"`
	ExpiresAt       time.Time `json:"expires_at" bson:"expires_at"`
	CreatedAt       time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" bson:"updated_at"`
}
