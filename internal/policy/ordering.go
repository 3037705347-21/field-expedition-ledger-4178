package policy

import "strings"

func ExpeditionTieBreak(leftID, rightID string) bool {
	left := strings.ToLower(strings.TrimSpace(leftID))
	right := strings.ToLower(strings.TrimSpace(rightID))
	return left < right
}
