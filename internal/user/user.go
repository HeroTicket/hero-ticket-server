package user

import (
	"errors"
	"time"
)

var (
	ErrNothingToUpdate = errors.New("nothing to update")
	ErrUserNotFound    = errors.New("user not found")
)

type User struct {
	DID           string    `json:"did" bson:"_id"`
	WalletAddress string    `json:"wallet_address" bson:"wallet_address"`
	Name          string    `json:"name" bson:"name"`
	IsAdmin       bool      `json:"is_admin" bson:"is_admin"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}
