package model

import "time"

import "example.com/field-expedition-ledger/internal/policy"

func NormalizeTime(value time.Time) time.Time {
	if value.IsZero() {
		return time.Time{}
	}
	return policy.UTCDate(value)
}
