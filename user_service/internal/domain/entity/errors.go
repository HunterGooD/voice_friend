package entity

import (
	"github.com/pkg/errors"
)

var (
	ErrUserAlreadyExists = errors.New("User already exists")
	ErrNotFound          = errors.New("Not found")
	ErrUnauthorized      = errors.New("Unauthorization user")
	ErrDataNotValid      = errors.New("Not valid user data")
	ErrInvalidPassword   = errors.New("Invalid password")
	ErrInternal          = errors.New("Unexpected error")
)
