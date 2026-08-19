package store

import (
	"context"
	"sort"
	"sync"

	"example.com/field-expedition-ledger/internal/model"
)

type MemoryRepository struct {
	mu           sync.RWMutex
	expeditions  map[string]model.Expedition
	observations map[string]model.Observation
	specimens    map[string]model.Specimen
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{
		expeditions:  make(map[string]model.Expedition),
		observations: make(map[string]model.Observation),
		specimens:    make(map[string]model.Specimen),
	}
}

func (r *MemoryRepository) CreateExpedition(ctx context.Context, expedition model.Expedition) (model.Expedition, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Expedition{}, err
	}
	if err := expedition.Validate(); err != nil {
		return model.Expedition{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.expeditions[expedition.ID] = expedition
	return expedition, nil
}

func (r *MemoryRepository) GetExpedition(ctx context.Context, id string) (model.Expedition, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Expedition{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.expeditions[id]
	if !ok {
		return model.Expedition{}, model.ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) UpdateExpedition(ctx context.Context, expedition model.Expedition) (model.Expedition, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Expedition{}, err
	}
	if err := expedition.Validate(); err != nil {
		return model.Expedition{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.expeditions[expedition.ID]; !ok {
		return model.Expedition{}, model.ErrNotFound
	}
	r.expeditions[expedition.ID] = expedition
	return expedition, nil
}

func (r *MemoryRepository) ListExpeditions(ctx context.Context, filter model.ExpeditionFilter) ([]model.Expedition, error) {
	if err := model.ContextReady(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	items := make([]model.Expedition, 0, len(r.expeditions))
	for _, item := range r.expeditions {
		if filter.Matches(item) {
			items = append(items, item)
		}
	}
	r.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.Before(items[j].CreatedAt)
	})
	return items, nil
}

func (r *MemoryRepository) CreateObservation(ctx context.Context, observation model.Observation) (model.Observation, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Observation{}, err
	}
	if err := observation.Validate(); err != nil {
		return model.Observation{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.expeditions[observation.ExpeditionID]; !ok {
		return model.Observation{}, model.ErrNotFound
	}
	r.observations[observation.ID] = observation
	return observation, nil
}

func (r *MemoryRepository) ListObservations(ctx context.Context, expeditionID string) ([]model.Observation, error) {
	if err := model.ContextReady(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	items := make([]model.Observation, 0)
	for _, item := range r.observations {
		if item.ExpeditionID == expeditionID {
			items = append(items, item)
		}
	}
	r.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].RecordedAt.Before(items[j].RecordedAt)
	})
	return items, nil
}

func (r *MemoryRepository) CreateSpecimen(ctx context.Context, specimen model.Specimen) (model.Specimen, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Specimen{}, err
	}
	if err := specimen.Validate(); err != nil {
		return model.Specimen{}, err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.expeditions[specimen.ExpeditionID]; !ok {
		return model.Specimen{}, model.ErrNotFound
	}
	r.specimens[specimen.ID] = specimen
	return specimen, nil
}

func (r *MemoryRepository) GetSpecimen(ctx context.Context, id string) (model.Specimen, error) {
	if err := model.ContextReady(ctx); err != nil {
		return model.Specimen{}, err
	}
	r.mu.RLock()
	defer r.mu.RUnlock()
	item, ok := r.specimens[id]
	if !ok {
		return model.Specimen{}, model.ErrNotFound
	}
	return item, nil
}

func (r *MemoryRepository) ListSpecimens(ctx context.Context, expeditionID string) ([]model.Specimen, error) {
	if err := model.ContextReady(ctx); err != nil {
		return nil, err
	}
	r.mu.RLock()
	items := make([]model.Specimen, 0)
	for _, item := range r.specimens {
		if item.ExpeditionID == expeditionID {
			items = append(items, item)
		}
	}
	r.mu.RUnlock()
	sort.Slice(items, func(i, j int) bool {
		return items[i].CollectedAt.Before(items[j].CollectedAt)
	})
	return items, nil
}
