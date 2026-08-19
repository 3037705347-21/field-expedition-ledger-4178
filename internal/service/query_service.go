package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/policy"
	"example.com/field-expedition-ledger/internal/store"
)

type QueryService struct {
	repository store.Repository
}

func NewQueryService(repository store.Repository) *QueryService {
	return &QueryService{repository: repository}
}

func (s *QueryService) Page(ctx context.Context, filter model.ExpeditionFilter, offset, limit int) (model.Page[model.Expedition], error) {
	items, err := s.repository.ListExpeditions(ctx, filter)
	if err != nil {
		return model.Page[model.Expedition]{}, err
	}
	start, end := policy.PageBounds(offset, limit, len(items))
	return model.NewPage(items[start:end], len(items), start, end-start), nil
}

func (s *QueryService) Active(ctx context.Context) ([]model.Expedition, error) {
	return s.repository.ListExpeditions(ctx, model.ExpeditionFilter{Status: string(model.ExpeditionActive)})
}

func (s *QueryService) ByStatus(ctx context.Context, status string) ([]model.Expedition, error) {
	return s.repository.ListExpeditions(ctx, model.ExpeditionFilter{Status: policy.CanonicalStatus(status)})
}

func (s *QueryService) Search(ctx context.Context, query string) ([]model.Expedition, error) {
	clean := strings.TrimSpace(query)
	if clean == "" {
		return []model.Expedition{}, nil
	}
	return s.repository.ListExpeditions(ctx, model.ExpeditionFilter{Search: clean})
}

func (s *QueryService) ByRegion(ctx context.Context, region string) ([]model.Expedition, error) {
	return s.repository.ListExpeditions(ctx, model.ExpeditionFilter{Region: strings.TrimSpace(region)})
}

func SortByStartDate(items []model.Expedition, newestFirst bool) []model.Expedition {
	result := append([]model.Expedition(nil), items...)
	sort.SliceStable(result, func(i, j int) bool {
		if newestFirst {
			return result[i].StartDate.After(result[j].StartDate)
		}
		return result[i].StartDate.Before(result[j].StartDate)
	})
	return result
}

func FilterStartedBefore(items []model.Expedition, moment time.Time) []model.Expedition {
	result := make([]model.Expedition, 0, len(items))
	for _, item := range items {
		if item.StartDate.Before(moment) {
			result = append(result, item)
		}
	}
	return result
}

func FilterWithNotes(items []model.Expedition) []model.Expedition {
	result := make([]model.Expedition, 0, len(items))
	for _, item := range items {
		if strings.TrimSpace(item.Notes) != "" {
			result = append(result, item)
		}
	}
	return result
}

func RegionNames(items []model.Expedition) []string {
	return policy.UniqueSorted(regions(items))
}

func LeadNames(items []model.Expedition) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Lead)
	}
	return policy.UniqueSorted(values)
}

func regions(items []model.Expedition) []string {
	values := make([]string, 0, len(items))
	for _, item := range items {
		values = append(values, item.Region)
	}
	return values
}
