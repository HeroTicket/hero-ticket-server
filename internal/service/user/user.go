package user

import (
	"errors"
	"time"
)

var (
	ErrInvalidID        = errors.New("invalid id")
	ErrNothingToUpdate  = errors.New("nothing to update")
	ErrUserNotFound     = errors.New("user not found")
	ErrTBAAlreadyExists = errors.New("tba address already exists")
)

type User struct {
	DID           string    `json:"did" bson:"_id"`
	WalletAddress string    `json:"walletAddress" bson:"walletAddress"`
	TbaAddress    string    `json:"tbaAddress" bson:"tbaAddress"`
	Nonce         int64     `json:"nonce" bson:"nonce"`
	IsAdmin       bool      `json:"isAdmin" bson:"isAdmin"`
	CreatedAt     time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt" bson:"updatedAt"`
}

type UserUpdateParams struct {
	DID           string
	WalletAddress string
	TBAAddress    string
	Name          string
}

type Membership struct {
	Identifier string    `json:"identifier" bson:"_id"`
	ClaimID    string    `json:"claimId,omitempty" bson:"claimId,omitempty"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}

type TicketOwnership struct {
	ID         string    `json:"id" bson:"_id"`
	Identifier string    `json:"identifier" bson:"identifier"`
	ClaimID    string    `json:"claimId" bson:"claimId"`
	CreatedAt  time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt" bson:"updatedAt"`
}
