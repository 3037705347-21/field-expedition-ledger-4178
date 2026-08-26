package analytics

import (
	"sort"
	"strings"

	"example.com/field-expedition-ledger/internal/model"
)

type ExpeditionStats struct {
	CountByStatus map[string]int `json:"count_by_status"`
	Regions       []string       `json:"regions"`
	Leads         []string       `json:"leads"`
	ActiveNames   []string       `json:"active_names"`
}

func SummarizeExpeditions(items []model.Expedition) ExpeditionStats {
	result := ExpeditionStats{
		CountByStatus: map[string]int{},
		Regions:       []string{},
		Leads:         []string{},
		ActiveNames:   []string{},
	}
	regions := make(map[string]struct{})
	leads := make(map[string]struct{})
	for _, item := range items {
		result.CountByStatus[string(item.Status)]++
		region := strings.TrimSpace(item.Region)
		if region != "" {
			regions[region] = struct{}{}
		}
		lead := strings.TrimSpace(item.Lead)
		if lead != "" {
			leads[lead] = struct{}{}
		}
		if item.Status == model.ExpeditionActive {
			result.ActiveNames = append(result.ActiveNames, item.Name)
		}
	}
	for region := range regions {
		result.Regions = append(result.Regions, region)
	}
	for lead := range leads {
		result.Leads = append(result.Leads, lead)
	}
	sort.Strings(result.Regions)
	sort.Strings(result.Leads)
	sort.Strings(result.ActiveNames)
	return result
}

func StatusNames() []string {
	return []string{
		string(model.ExpeditionPlanned),
		string(model.ExpeditionActive),
		string(model.ExpeditionClosed),
	}
}

func CountStatus(items []model.Expedition, status model.ExpeditionStatus) int {
	total := 0
	for _, item := range items {
		if item.Status == status {
			total++
		}
	}
	return total
}
