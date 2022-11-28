package digdaggo

import (
	"errors"
)

// Generic Http Status Error
var (
	ErrUnauthorized = errors.New("unauthorized")

	ErrNotFound = errors.New("not found")

	ErrForbidden = errors.New("you are not allowed to this operation")

	ErrServer = errors.New("internal Server Error")

	ErrClient = errors.New("client Error")
)
