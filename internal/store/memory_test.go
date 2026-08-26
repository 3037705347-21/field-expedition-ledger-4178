package store

import (
	"context"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
)

func TestMemoryRepositoryStoresAndFiltersExpeditions(t *testing.T) {
	repository := NewMemoryRepository()
	now := time.Now().UTC()
	_, err := repository.CreateExpedition(context.Background(), model.Expedition{
		ID: "exp-1", Name: "North", Region: "Basin", Lead: "Ari",
		Status: model.ExpeditionPlanned, StartDate: now, CreatedAt: now, UpdatedAt: now,
	})
	if err != nil {
		t.Fatalf("create expedition: %v", err)
	}
	items, err := repository.ListExpeditions(context.Background(), model.ExpeditionFilter{Region: "basin"})
	if err != nil || len(items) != 1 {
		t.Fatalf("list filtered expeditions: len=%d err=%v", len(items), err)
	}
}
