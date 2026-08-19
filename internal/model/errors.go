package model

import "errors"

var (
	ErrInvalidInput     = errors.New("invalid input")
	ErrNotFound         = errors.New("not found")
	ErrInvalidState     = errors.New("invalid state")
	ErrClosedExpedition = errors.New("expedition is closed")
)

func IsInputError(err error) bool {
	return err == ErrInvalidInput
}
