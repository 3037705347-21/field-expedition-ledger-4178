package model

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNotFound         = errors.New("not found")
	ErrInvalidState     = errors.New("invalid state")
	ErrClosedExpedition = errors.New("expedition is closed")
	ErrInvalidFilter    = errors.New("invalid filter")
)
