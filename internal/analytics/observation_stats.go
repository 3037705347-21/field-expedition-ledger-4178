package analytics

import (
	"sort"
	"strings"

	"example.com/field-expedition-ledger/internal/model"
)

type ObservationStats struct {
	Count             int            `json:"count"`
	MinimumElevation  float64        `json:"minimum_elevation"`
	MaximumElevation  float64        `json:"maximum_elevation"`
	AverageElevation  float64        `json:"average_elevation"`
	AverageConfidence float64        `json:"average_confidence"`
	Sites             []string       `json:"sites"`
	RockTypes         map[string]int `json:"rock_types"`
	TagCounts         map[string]int `json:"tag_counts"`
	ConfidenceBand    string         `json:"confidence_band"`
}

func SummarizeObservations(items []model.Observation) ObservationStats {
	stats := ObservationStats{
		Sites:     []string{},
		RockTypes: map[string]int{},
		TagCounts: map[string]int{},
	}
	if len(items) == 0 {
		stats.ConfidenceBand = "unrated"
		return stats
	}
	stats.Count = len(items)
	stats.MinimumElevation = items[0].ElevationM
	stats.MaximumElevation = items[0].ElevationM
	seenSites := make(map[string]struct{})
	totalElevation := 0.0
	totalConfidence := 0.0
	for _, item := range items {
		if item.ElevationM < stats.MinimumElevation {
			stats.MinimumElevation = item.ElevationM
		}
		if item.ElevationM > stats.MaximumElevation {
			stats.MaximumElevation = item.ElevationM
		}
		totalElevation += item.ElevationM
		totalConfidence += item.Confidence
		rock := strings.ToLower(strings.TrimSpace(item.RockType))
		if rock != "" {
			stats.RockTypes[rock]++
		}
		site := strings.ToUpper(strings.TrimSpace(item.SiteCode))
		if site != "" {
			seenSites[site] = struct{}{}
		}
		for _, tag := range item.Tags {
			cleanTag := strings.ToLower(strings.TrimSpace(tag))
			if cleanTag != "" {
				stats.TagCounts[cleanTag]++
			}
		}
	}
	for site := range seenSites {
		stats.Sites = append(stats.Sites, site)
	}
	sort.Strings(stats.Sites)
	stats.AverageElevation = totalElevation / float64(stats.Count)
	stats.AverageConfidence = totalConfidence / float64(stats.Count)
	stats.ConfidenceBand = ConfidenceBand(stats.AverageConfidence)
	return stats
}

func ConfidenceBand(value float64) string {
	switch {
	case value >= .9:
		return "high"
	case value >= .7:
		return "medium"
	case value > 0:
		return "low"
	default:
		return "unrated"
	}
}

func ElevationRange(items []model.Observation) (float64, float64, bool) {
	if len(items) == 0 {
		return 0, 0, false
	}
	minimum := items[0].ElevationM
	maximum := items[0].ElevationM
	for _, item := range items[1:] {
		if item.ElevationM < minimum {
			minimum = item.ElevationM
		}
		if item.ElevationM > maximum {
			maximum = item.ElevationM
		}
	}
	return minimum, maximum, true
}

func AverageConfidence(items []model.Observation) float64 {
	if len(items) == 0 {
		return 0
	}
	total := 0.0
	for _, item := range items {
		total += item.Confidence
	}
	return total / float64(len(items))
}

func UniqueRockTypes(items []model.Observation) []string {
	set := make(map[string]struct{})
	for _, item := range items {
		value := strings.ToLower(strings.TrimSpace(item.RockType))
		if value != "" {
			set[value] = struct{}{}
		}
	}
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func GroupBySite(items []model.Observation) map[string][]model.Observation {
	grouped := make(map[string][]model.Observation)
	for _, item := range items {
		site := strings.ToUpper(strings.TrimSpace(item.SiteCode))
		grouped[site] = append(grouped[site], item)
	}
	for site := range grouped {
		sort.SliceStable(grouped[site], func(i, j int) bool {
			return grouped[site][i].RecordedAt.Before(grouped[site][j].RecordedAt)
		})
	}
	return grouped
}
