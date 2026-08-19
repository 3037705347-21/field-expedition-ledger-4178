package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/store"
)

func TestRecordStopsWhenContextIsCanceled(t *testing.T) {
	repository := store.NewMemoryRepository()
	expeditions := NewExpeditionService(repository)
	observations := NewObservationService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	expedition, err := expeditions.Create(context.Background(), "Context Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(context.Background(), expedition.ID); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = observations.Record(ctx, expedition.ID, "C-1", start, 1, 2, 100, "granite", "note", nil, .8)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context cancellation, got %v", err)
	}
}
