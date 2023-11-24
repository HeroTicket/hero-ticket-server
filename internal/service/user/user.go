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
	Identifier     string    `json:"identifier" bson:"_id"`
	AccountAddress string    `json:"accountAddress" bson:"accountAddress"`
	TbaAddress     string    `json:"tbaAddress" bson:"tbaAddress"`
	Name           string    `json:"name" bson:"name"`
	Bio            string    `json:"bio" bson:"bio"`
	Avatar         string    `json:"avatar" bson:"avatar"`
	Banner         string    `json:"banner" bson:"banner"`
	IsAdmin        bool      `json:"isAdmin" bson:"isAdmin"`
	Verified       bool      `json:"verified" bson:"verified"`
	CreatedAt      time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt" bson:"updatedAt"`
}

type CreateUserParams struct {
	Identifier     string
	AccountAddress string
	TBAAddress     string
	Name           string
	Avatar         string
	IsAdmin        bool
}

type UpdateUserParams struct {
	Identifier     string
	AccountAddress string
	TBAAddress     string
	Name           string
	Bio            string
	Avatar         string
	Banner         string
	IsAdmin        bool
	Verified       bool
}
