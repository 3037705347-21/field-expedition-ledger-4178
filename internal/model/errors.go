package model

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrInvalidFilter    = errors.New("invalid filter")
	ErrNotFound         = errors.New("not found")
	ErrInvalidState     = errors.New("invalid state")
	ErrClosedExpedition = errors.New("expedition is closed")
)
