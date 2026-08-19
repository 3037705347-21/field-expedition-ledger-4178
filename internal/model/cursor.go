package model

import "time"

// RecordedAfterCursor excludes the record already consumed by a sync cursor.
func RecordedAfterCursor(recordedAt, cursor time.Time) bool {
	return recordedAt.After(cursor)
}
