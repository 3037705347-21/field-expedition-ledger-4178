package service

import (
	"context"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

type SpecimenService struct {
	repository store.Repository
	clock      func() time.Time
}

func NewSpecimenService(repository store.Repository) *SpecimenService {
	return &SpecimenService{repository: repository, clock: time.Now}
}

func (s *SpecimenService) Register(ctx context.Context, expeditionID, label, material string, weight float64, collectedAt time.Time, custodian, notes string) (model.Specimen, error) {
	expedition, err := s.repository.GetExpedition(ctx, strings.TrimSpace(expeditionID))
	if err != nil {
		return model.Specimen{}, err
	}
	if expedition.Status == model.ExpeditionClosed {
		return model.Specimen{}, model.ErrClosedExpedition
	}
	if !expedition.CanRecord() {
		return model.Specimen{}, model.ErrInactiveExpedition
	}
	item := model.Specimen{
		ID:           model.NewID("spc"),
		ExpeditionID: strings.TrimSpace(expeditionID),
		Label:        label,
		Material:     material,
		WeightGrams:  weight,
		CollectedAt:  model.NormalizeTime(collectedAt),
		Custodian:    custodian,
		Status:       model.SpecimenCollected,
		Notes:        strings.TrimSpace(notes),
	}.Normalized()
	if err := item.Validate(); err != nil {
		return model.Specimen{}, err
	}
	return s.repository.CreateSpecimen(ctx, item)
}

func (s *SpecimenService) Get(ctx context.Context, id string) (model.Specimen, error) {
	return s.repository.GetSpecimen(ctx, strings.TrimSpace(id))
}

func (s *SpecimenService) List(ctx context.Context, expeditionID string) ([]model.Specimen, error) {
	if _, err := s.repository.GetExpedition(ctx, strings.TrimSpace(expeditionID)); err != nil {
		return nil, err
	}
	return s.repository.ListSpecimens(ctx, strings.TrimSpace(expeditionID))
}
