package service

import (
	"context"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

type ObservationService struct {
	repository store.Repository
	clock      func() time.Time
}

func NewObservationService(repository store.Repository) *ObservationService {
	return &ObservationService{repository: repository, clock: time.Now}
}

func (s *ObservationService) Record(ctx context.Context, expeditionID, siteCode string, recordedAt time.Time, latitude, longitude, elevation float64, rockType, description string, tags []string, confidence float64) (model.Observation, error) {
	if err := ctx.Err(); err != nil {
		return model.Observation{}, err
	}
	expedition, err := s.repository.GetExpedition(ctx, strings.TrimSpace(expeditionID))
	if err != nil {
		return model.Observation{}, err
	}
	if expedition.Status == model.ExpeditionClosed {
		return model.Observation{}, model.ErrClosedExpedition
	}
	item := model.Observation{
		ID:           model.NewID("obs"),
		ExpeditionID: strings.TrimSpace(expeditionID),
		SiteCode:     siteCode,
		RecordedAt:   model.NormalizeTime(recordedAt),
		Latitude:     latitude,
		Longitude:    longitude,
		ElevationM:   elevation,
		RockType:     rockType,
		Description:  description,
		Tags:         append([]string(nil), tags...),
		Confidence:   confidence,
	}.Normalized()
	if err := item.Validate(); err != nil {
		return model.Observation{}, err
	}
	return s.repository.CreateObservation(ctx, item)
}

func (s *ObservationService) List(ctx context.Context, expeditionID string) ([]model.Observation, error) {
	if _, err := s.repository.GetExpedition(ctx, strings.TrimSpace(expeditionID)); err != nil {
		return nil, err
	}
	return s.repository.ListObservations(ctx, strings.TrimSpace(expeditionID))
}

func (s *ObservationService) Recent(ctx context.Context, expeditionID string, since time.Time) ([]model.Observation, error) {
	items, err := s.List(ctx, expeditionID)
	if err != nil {
		return nil, err
	}
	filtered := make([]model.Observation, 0, len(items))
	for _, item := range items {
		if !item.RecordedAt.Before(since) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}
