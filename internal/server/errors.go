package server

import (
	"errors"
)

var (
	ErrBadRequest    = errors.New("Bad request")
	ErrNotFound      = errors.New("Not found summary")
	ErrNoHeader      = errors.New("While uploading wasn't firstly passed ExchangeData")
	ErrNotAuthorized = errors.New("Cannot authorize")
)
