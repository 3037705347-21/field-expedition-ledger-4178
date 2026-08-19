package model

import (
	"strings"

	"example.com/field-expedition-ledger/internal/policy"
)

type ExpeditionFilter struct {
	Status string
	Region string
	Search string
}

func (f ExpeditionFilter) Matches(expedition Expedition) bool {
	if f.Status != "" && policy.CanonicalStatus(string(expedition.Status)) != policy.CanonicalStatus(f.Status) {
		return false
	}
	if f.Region != "" && !strings.EqualFold(strings.TrimSpace(expedition.Region), strings.TrimSpace(f.Region)) {
		return false
	}
	if f.Search != "" {
		query := strings.ToLower(strings.TrimSpace(f.Search))
		haystack := strings.ToLower(expedition.Name + " " + expedition.Region + " " + expedition.Lead)
		return strings.Contains(haystack, query)
	}
	return true
}
