package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/analysis"
	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

type ExpeditionInsight struct {
	ExpeditionID        string         `json:"expedition_id"`
	ObservationDateList []string       `json:"observation_days"`
	SiteCoverage        map[string]int `json:"site_coverage"`
	ConfidenceBands     map[string]int `json:"confidence_bands"`
	MaterialNames       []string       `json:"material_names"`
	ReviewWarnings      []string       `json:"review_warnings"`
	ReadyForClose       bool           `json:"ready_for_close"`
	ObservationWindow   string         `json:"observation_window"`
	ObservationDays     int            `json:"observation_day_count"`
	RepeatedSite        bool           `json:"repeated_site"`
	SiteNames           []string       `json:"site_names"`
	CoverageText        []string       `json:"coverage_text"`
	ReviewStatus        string         `json:"review_status"`
	SortedSiteCounts    []string       `json:"sorted_site_counts"`
	WindowHasData       bool           `json:"window_has_data"`
	MatchedMaterials    int            `json:"matched_materials"`
	MatchedTags         int            `json:"matched_tags"`
}

type InsightService struct {
	repository store.Repository
}

func NewInsightService(repository store.Repository) *InsightService {
	return &InsightService{repository: repository}
}

func (s *InsightService) Build(ctx context.Context, expeditionID string) (ExpeditionInsight, error) {
	if _, err := s.repository.GetExpedition(ctx, expeditionID); err != nil {
		return ExpeditionInsight{}, err
	}
	observations, err := s.repository.ListObservations(ctx, expeditionID)
	if err != nil {
		return ExpeditionInsight{}, err
	}
	specimens, err := s.repository.ListSpecimens(ctx, expeditionID)
	if err != nil {
		return ExpeditionInsight{}, err
	}
	first, last := analysis.ObservationWindow(observations)
	materials := ""
	if len(specimens) > 0 {
		materials = specimens[0].Material
	}
	tag := ""
	if len(observations) > 0 && len(observations[0].Tags) > 0 {
		tag = observations[0].Tags[0]
	}
	warnings := analysis.ReviewWarnings(observations, specimens)
	score := analysis.QualityScore(observations, specimens)
	return ExpeditionInsight{
		ExpeditionID:        expeditionID,
		ObservationDateList: sortedDays(observations),
		SiteCoverage:        analysis.SiteCoverage(observations),
		ConfidenceBands:     analysis.ConfidenceDistribution(observations),
		MaterialNames:       analysis.MaterialNames(observations, specimens),
		ReviewWarnings:      warnings,
		ReadyForClose:       analysis.IsReadyForClose(observations, specimens),
		ObservationWindow:   formatWindow(first, last),
		ObservationDays:     ObservationDayCount(observations),
		RepeatedSite:        HasRepeatedSite(observations),
		SiteNames:           SiteNames(observations),
		CoverageText:        DescribeCoverage(observations),
		ReviewStatus:        ReviewStatus(score, warnings),
		SortedSiteCounts:    SortedSiteCounts(observations),
		WindowHasData:       WindowContains(observations, first, last),
		MatchedMaterials:    len(MatchMaterial(specimens, materials)),
		MatchedTags:         len(MatchTag(observations, tag)),
	}, nil
}

func sortedDays(items []model.Observation) []string {
	set := make(map[string]struct{})
	for _, item := range items {
		set[item.RecordedAt.UTC().Format("2006-01-02")] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for day := range set {
		result = append(result, day)
	}
	sort.Strings(result)
	return result
}

func formatWindow(first, last time.Time) string {
	if first.IsZero() || last.IsZero() {
		return ""
	}
	return first.UTC().Format(time.RFC3339) + " to " + last.UTC().Format(time.RFC3339)
}

func ObservationDayCount(items []model.Observation) int {
	return len(sortedDays(items))
}

func HasRepeatedSite(items []model.Observation) bool {
	coverage := analysis.SiteCoverage(items)
	for _, count := range coverage {
		if count > 1 {
			return true
		}
	}
	return false
}

func SiteNames(items []model.Observation) []string {
	coverage := analysis.SiteCoverage(items)
	result := make([]string, 0, len(coverage))
	for site := range coverage {
		result = append(result, site)
	}
	sort.Strings(result)
	return result
}

func DescribeCoverage(items []model.Observation) []string {
	coverage := analysis.SiteCoverage(items)
	result := make([]string, 0, len(coverage))
	for site, count := range coverage {
		result = append(result, site+"="+formatCount(count))
	}
	sort.Strings(result)
	return result
}

func formatCount(value int) string {
	if value < 1 {
		return "0"
	}
	digits := []byte{}
	for value > 0 {
		digits = append([]byte{byte(value%10) + '0'}, digits...)
		value /= 10
	}
	return string(digits)
}

func MatchMaterial(items []model.Specimen, query string) []model.Specimen {
	target := strings.ToLower(strings.TrimSpace(query))
	result := make([]model.Specimen, 0)
	for _, item := range items {
		if strings.Contains(strings.ToLower(item.Material), target) {
			result = append(result, item)
		}
	}
	return result
}

func MatchTag(items []model.Observation, query string) []model.Observation {
	target := strings.ToLower(strings.TrimSpace(query))
	result := make([]model.Observation, 0)
	for _, item := range items {
		for _, tag := range item.Tags {
			if strings.EqualFold(strings.TrimSpace(tag), target) {
				result = append(result, item)
				break
			}
		}
	}
	return result
}

func ReviewStatus(score float64, warnings []string) string {
	switch {
	case score >= .9 && len(warnings) == 0:
		return "excellent"
	case score >= .7 && len(warnings) <= 1:
		return "good"
	case score > 0:
		return "needs-review"
	default:
		return "not-started"
	}
}

func SortedSiteCounts(items []model.Observation) []string {
	coverage := analysis.SiteCoverage(items)
	result := make([]string, 0, len(coverage))
	for site, count := range coverage {
		result = append(result, site+"="+formatCount(count))
	}
	sort.Slice(result, func(i, j int) bool {
		left, right := result[i], result[j]
		return left < right
	})
	return result
}

func WindowContains(items []model.Observation, start, end time.Time) bool {
	if start.IsZero() || end.IsZero() || end.Before(start) {
		return false
	}
	for _, item := range items {
		if !item.RecordedAt.Before(start) && !item.RecordedAt.After(end) {
			return true
		}
	}
	return false
}
