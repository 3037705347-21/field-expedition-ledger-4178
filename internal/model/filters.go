package model

import (
	"strings"
	"time"
)

type ExpeditionFilter struct {
	Status string
	Region string
	Search string
}

func ParseBeforeFilter(value string) (time.Time, error) {
	moment, err := time.Parse(time.RFC3339, strings.TrimSpace(value))
	if err != nil {
		return time.Time{}, ErrInvalidInput
	}
	return moment, nil
}

func (f ExpeditionFilter) Matches(expedition Expedition) bool {
	if f.Status != "" && string(expedition.Status) != f.Status {
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
