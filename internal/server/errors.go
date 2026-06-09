package server

import (
	"errors"
)

var (
	ErrBadRequest         = errors.New("Bad request")
	ErrNotFound           = errors.New("Not found summary")
	ErrNoHeader           = errors.New("While uploading wasn't firstly passed ExchangeData")
	ErrNotAuthorized      = errors.New("Cannot authorize")
	ErrNotPermitted       = errors.New("You don't have enough permissions to do this action")
	ErrNoContent          = errors.New("Part of data wasn't passed")
	ErrSomethingWentWrong = errors.New("Sory, something went wrong")
)
