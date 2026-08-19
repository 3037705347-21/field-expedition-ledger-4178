package model

import "time"

// RecordedOnOrAfter is used by incremental observation queries.
func RecordedOnOrAfter(recordedAt, cursor time.Time) bool {
	return !recordedAt.Before(cursor)
}
