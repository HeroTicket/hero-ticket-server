package did

import (
	"errors"
	"time"
)

var (
	ErrRequestNotFound = errors.New("request not found")
)

var DefaultCacheExpiry = 10 * time.Minute
