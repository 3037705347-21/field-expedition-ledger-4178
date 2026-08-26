package model

import (
	"fmt"
	"sync/atomic"
	"time"
)

var sequence atomic.Uint64

func NewID(prefix string) string {
	n := sequence.Add(1)
	return fmt.Sprintf("%s-%d-%04d", prefix, time.Now().UTC().UnixNano(), n)
}
