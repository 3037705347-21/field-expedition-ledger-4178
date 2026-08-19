package analytics

import (
	"sort"
	"strings"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/policy"
)

type SpecimenStats struct {
	Count           int            `json:"count"`
	TotalWeightGram float64        `json:"total_weight_gram"`
	AverageWeight   float64        `json:"average_weight"`
	Materials       []string       `json:"materials"`
	StatusCounts    map[string]int `json:"status_counts"`
	Custodians      []string       `json:"custodians"`
	WeightBands     map[string]int `json:"weight_bands"`
}

func SummarizeSpecimens(items []model.Specimen) SpecimenStats {
	result := SpecimenStats{
		Materials:    []string{},
		StatusCounts: map[string]int{},
		Custodians:   []string{},
		WeightBands:  map[string]int{},
	}
	if len(items) == 0 {
		return result
	}
	materials := make(map[string]struct{})
	custodians := make(map[string]struct{})
	for _, item := range items {
		result.Count++
		result.TotalWeightGram += item.WeightGrams
		result.StatusCounts[string(item.Status)]++
		result.WeightBands[policy.WeightBand(item.WeightGrams)]++
		material := item.MaterialKey()
		if material != "" {
			materials[material] = struct{}{}
		}
		custodian := strings.TrimSpace(item.Custodian)
		if custodian != "" {
			custodians[custodian] = struct{}{}
		}
	}
	result.AverageWeight = result.TotalWeightGram / float64(result.Count)
	for key := range materials {
		result.Materials = append(result.Materials, key)
	}
	for key := range custodians {
		result.Custodians = append(result.Custodians, key)
	}
	sort.Strings(result.Materials)
	sort.Strings(result.Custodians)
	return result
}

func WeightByMaterial(items []model.Specimen) map[string]float64 {
	result := make(map[string]float64)
	for _, item := range items {
		key := item.MaterialKey()
		result[key] += item.WeightGrams
	}
	return result
}

func StatusBreakdown(items []model.Specimen) []string {
	counts := make(map[string]int)
	for _, item := range items {
		counts[string(item.Status)]++
	}
	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	result := make([]string, 0, len(keys))
	for _, key := range keys {
		result = append(result, key+":"+itoa(counts[key]))
	}
	return result
}

func itoa(value int) string {
	if value == 0 {
		return "0"
	}
	digits := make([]byte, 0, 4)
	for value > 0 {
		digits = append([]byte{byte('0' + value%10)}, digits...)
		value /= 10
	}
	return string(digits)
}
