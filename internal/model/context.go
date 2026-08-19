package model

import "context"

func ContextReady(ctx context.Context) error {
	return ctx.Err()
}
