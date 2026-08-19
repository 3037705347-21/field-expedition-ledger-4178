package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

type ExpeditionService struct {
	repository store.Repository
	clock      func() time.Time
}

func NewExpeditionService(repository store.Repository) *ExpeditionService {
	return &ExpeditionService{repository: repository, clock: time.Now}
}

func (s *ExpeditionService) Create(ctx context.Context, name, region, lead string, startDate time.Time, notes string) (model.Expedition, error) {
	now := model.NormalizeTime(s.clock())
	item := model.Expedition{
		ID:        model.NewID("exp"),
		Name:      strings.TrimSpace(name),
		Region:    strings.TrimSpace(region),
		Lead:      strings.TrimSpace(lead),
		Status:    model.ExpeditionPlanned,
		StartDate: model.NormalizeTime(startDate),
		Notes:     strings.TrimSpace(notes),
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := item.Validate(); err != nil {
		return model.Expedition{}, fmt.Errorf("expedition input: %v", err)
	}
	return s.repository.CreateExpedition(ctx, item)
}

func (s *ExpeditionService) List(ctx context.Context, filter model.ExpeditionFilter) ([]model.Expedition, error) {
	return s.repository.ListExpeditions(ctx, filter)
}

func (s *ExpeditionService) Get(ctx context.Context, id string) (model.Expedition, error) {
	return s.repository.GetExpedition(ctx, strings.TrimSpace(id))
}

func (s *ExpeditionService) Activate(ctx context.Context, id string) (model.Expedition, error) {
	item, err := s.repository.GetExpedition(ctx, strings.TrimSpace(id))
	if err != nil {
		return model.Expedition{}, err
	}
	if !item.CanActivate() {
		return model.Expedition{}, model.ErrInvalidState
	}
	item.Status = model.ExpeditionActive
	item.UpdatedAt = model.NormalizeTime(s.clock())
	return s.repository.UpdateExpedition(ctx, item)
}

func (s *ExpeditionService) Close(ctx context.Context, id string) (model.Expedition, error) {
	item, err := s.repository.GetExpedition(ctx, strings.TrimSpace(id))
	if err != nil {
		return model.Expedition{}, err
	}
	if !item.CanClose() {
		return model.Expedition{}, model.ErrInvalidState
	}
	now := model.NormalizeTime(s.clock())
	item.Status = model.ExpeditionClosed
	item.EndDate = &now
	item.UpdatedAt = now
	return s.repository.UpdateExpedition(ctx, item)
}
