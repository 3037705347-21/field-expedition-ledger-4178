package analysis

import (
	"math"
	"sort"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/analytics"
	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/policy"
)

type Review struct {
	SiteCount          int      `json:"site_count"`
	QualityScore       float64  `json:"quality_score"`
	FollowUp           string   `json:"follow_up"`
	ObservedMaterials  []string `json:"observed_materials"`
	ActiveCustodians   []string `json:"active_custodians"`
	FirstObservation   string   `json:"first_observation,omitempty"`
	LastObservation    string   `json:"last_observation,omitempty"`
	ObservationDays    int      `json:"observation_days"`
	TaggedObservation  int      `json:"tagged_observation_count"`
	HighConfidence     int      `json:"high_confidence_count"`
	ElevationMatches   int      `json:"elevation_matches"`
	SpecimenWeightDays int      `json:"specimen_weight_days"`
	SortedObservation  []string `json:"sorted_observation_ids"`
	SortedSpecimens    []string `json:"sorted_specimen_ids"`
}

func BuildReview(observations []model.Observation, specimens []model.Specimen) Review {
	review := Review{ObservedMaterials: []string{}, ActiveCustodians: []string{}, FollowUp: "collect more observations"}
	if len(observations) == 0 && len(specimens) == 0 {
		return review
	}
	review.SiteCount = len(SiteCoverage(observations))
	review.QualityScore = QualityScore(observations, specimens)
	review.FollowUp = FollowUpAction(review.SiteCount, len(observations), len(specimens), review.QualityScore)
	review.ObservedMaterials = MaterialNames(observations, specimens)
	review.ActiveCustodians = CustodianNames(specimens)
	review.ObservationDays = len(DailyObservationCounts(observations))
	review.TaggedObservation = len(TagSet(observations))
	review.HighConfidence = len(FilterByConfidence(observations, .9))
	review.ElevationMatches = len(FilterByElevation(observations, 0, 3000))
	review.SpecimenWeightDays = len(DailySpecimenWeights(specimens))
	for _, item := range SortObservationsByConfidence(observations) {
		review.SortedObservation = append(review.SortedObservation, item.ID)
	}
	for _, item := range SortSpecimensByWeight(specimens) {
		review.SortedSpecimens = append(review.SortedSpecimens, item.ID)
	}
	if len(observations) > 0 {
		first, last := ObservationWindow(observations)
		review.FirstObservation = first.Format(time.RFC3339)
		review.LastObservation = last.Format(time.RFC3339)
	}
	return review
}

func SiteCoverage(items []model.Observation) map[string]int {
	coverage := make(map[string]int)
	for _, item := range items {
		site := strings.ToUpper(strings.TrimSpace(item.SiteCode))
		if site != "" {
			coverage[site]++
		}
	}
	return coverage
}

func MaterialNames(observations []model.Observation, specimens []model.Specimen) []string {
	set := make(map[string]struct{})
	for _, item := range observations {
		if value := strings.ToLower(strings.TrimSpace(item.RockType)); value != "" {
			set[value] = struct{}{}
		}
	}
	for _, item := range specimens {
		if value := strings.ToLower(strings.TrimSpace(item.Material)); value != "" {
			set[value] = struct{}{}
		}
	}
	return sortedSet(set)
}

func CustodianNames(items []model.Specimen) []string {
	set := make(map[string]struct{})
	for _, item := range items {
		if value := strings.TrimSpace(item.Custodian); value != "" {
			set[value] = struct{}{}
		}
	}
	return sortedSet(set)
}

func sortedSet(values map[string]struct{}) []string {
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}

func QualityScore(observations []model.Observation, specimens []model.Specimen) float64 {
	if len(observations) == 0 {
		if len(specimens) == 0 {
			return 0
		}
		return 0.25
	}
	confidence := 0.0
	descriptions := 0
	coordinates := 0
	for _, item := range observations {
		confidence += item.Confidence
		if strings.TrimSpace(item.Description) != "" {
			descriptions++
		}
		if item.Latitude != 0 || item.Longitude != 0 {
			coordinates++
		}
	}
	score := confidence / float64(len(observations))
	score += .15 * float64(descriptions) / float64(len(observations))
	score += .1 * float64(coordinates) / float64(len(observations))
	if len(specimens) > 0 {
		score += .05
	}
	return policy.Clamp(score, 0, 1)
}

func FollowUpAction(siteCount, observationCount, specimenCount int, score float64) string {
	switch {
	case observationCount == 0:
		return "record the first site observation"
	case siteCount < 2:
		return "visit a second site"
	case specimenCount == 0:
		return "register a specimen for custody"
	case score < .7:
		return "review low-confidence notes"
	default:
		return "ready for synthesis"
	}
}

func ObservationWindow(items []model.Observation) (time.Time, time.Time) {
	if len(items) == 0 {
		return time.Time{}, time.Time{}
	}
	first, last := items[0].RecordedAt, items[0].RecordedAt
	for _, item := range items[1:] {
		if item.RecordedAt.Before(first) {
			first = item.RecordedAt
		}
		if item.RecordedAt.After(last) {
			last = item.RecordedAt
		}
	}
	return first, last
}

func DailyObservationCounts(items []model.Observation) map[string]int {
	counts := make(map[string]int)
	for _, item := range items {
		counts[item.RecordedAt.UTC().Format("2006-01-02")]++
	}
	return counts
}

func DailySpecimenWeights(items []model.Specimen) map[string]float64 {
	weights := make(map[string]float64)
	for _, item := range items {
		weights[item.CollectedAt.UTC().Format("2006-01-02")] += item.WeightGrams
	}
	return weights
}

func ConfidenceDistribution(items []model.Observation) map[string]int {
	distribution := map[string]int{"high": 0, "medium": 0, "low": 0, "unrated": 0}
	for _, item := range items {
		distribution[analytics.ConfidenceBand(item.Confidence)]++
	}
	return distribution
}

func SortObservationsByConfidence(items []model.Observation) []model.Observation {
	result := append([]model.Observation(nil), items...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].Confidence == result[j].Confidence {
			return result[i].RecordedAt.Before(result[j].RecordedAt)
		}
		return result[i].Confidence > result[j].Confidence
	})
	return result
}

func SortSpecimensByWeight(items []model.Specimen) []model.Specimen {
	result := append([]model.Specimen(nil), items...)
	sort.SliceStable(result, func(i, j int) bool {
		if result[i].WeightGrams == result[j].WeightGrams {
			return result[i].CollectedAt.Before(result[j].CollectedAt)
		}
		return result[i].WeightGrams > result[j].WeightGrams
	})
	return result
}

func CoordinateBounds(items []model.Observation) (float64, float64, float64, float64, bool) {
	if len(items) == 0 {
		return 0, 0, 0, 0, false
	}
	minLat, maxLat := items[0].Latitude, items[0].Latitude
	minLon, maxLon := items[0].Longitude, items[0].Longitude
	for _, item := range items[1:] {
		minLat = math.Min(minLat, item.Latitude)
		maxLat = math.Max(maxLat, item.Latitude)
		minLon = math.Min(minLon, item.Longitude)
		maxLon = math.Max(maxLon, item.Longitude)
	}
	return minLat, maxLat, minLon, maxLon, true
}

func HasCoordinateDrift(items []model.Observation, threshold float64) bool {
	minLat, maxLat, minLon, maxLon, ok := CoordinateBounds(items)
	return ok && (maxLat-minLat > threshold || maxLon-minLon > threshold)
}

func TagSet(items []model.Observation) []string {
	set := make(map[string]struct{})
	for _, item := range items {
		for _, tag := range item.Tags {
			if value := strings.ToLower(strings.TrimSpace(tag)); value != "" {
				set[value] = struct{}{}
			}
		}
	}
	return sortedSet(set)
}

func FilterByConfidence(items []model.Observation, minimum float64) []model.Observation {
	result := make([]model.Observation, 0, len(items))
	for _, item := range items {
		if item.Confidence >= minimum {
			result = append(result, item)
		}
	}
	return result
}

func FilterByElevation(items []model.Observation, minimum, maximum float64) []model.Observation {
	result := make([]model.Observation, 0, len(items))
	for _, item := range items {
		if item.ElevationM >= minimum && item.ElevationM <= maximum {
			result = append(result, item)
		}
	}
	return result
}

func IsReadyForClose(observations []model.Observation, specimens []model.Specimen) bool {
	return len(observations) >= 2 && len(specimens) >= 1 && QualityScore(observations, specimens) >= .7
}

func ReviewWarnings(observations []model.Observation, specimens []model.Specimen) []string {
	warnings := make([]string, 0)
	if len(observations) == 0 {
		warnings = append(warnings, "no observations")
	}
	if len(specimens) == 0 {
		warnings = append(warnings, "no specimens")
	}
	if len(observations) > 0 && analytics.AverageConfidence(observations) < .7 {
		warnings = append(warnings, "low confidence")
	}
	if HasCoordinateDrift(observations, 5) {
		warnings = append(warnings, "wide coordinate spread")
	}
	return warnings
}
func (r Review) ApplyToSummary(summary *model.ExpeditionSummary) {
	summary.SiteCount = r.SiteCount
	summary.QualityScore = r.QualityScore
	summary.FollowUp = r.FollowUp
}
