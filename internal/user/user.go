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
	WalletAddress string    `json:"wallet_address" bson:"wallet_address"`
	TBAAddress    string    `json:"tba_address" bson:"tba_address"`
	Name          string    `json:"name" bson:"name"`
	IsAdmin       bool      `json:"is_admin" bson:"is_admin"`
	CreatedAt     time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" bson:"updated_at"`
}

type UserUpdateParams struct {
	DID           string
	WalletAddress string
	TBAAddress    string
	Name          string
}
