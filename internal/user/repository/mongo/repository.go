package mongo

import (
	"github.com/heroticket/internal/user"
)

type mongoRepository struct {
	user.Query
	user.Command
}

func NewMongoRepository() user.Repository {
	return &mongoRepository{}
}
