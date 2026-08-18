package service

import (
	"context"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestExpeditionLifecycleAndSummary(t *testing.T) {
	repository := store.NewMemoryRepository()
	expeditions := NewExpeditionService(repository)
	observations := NewObservationService(repository)
	specimens := NewSpecimenService(repository)
	summaries := NewSummaryService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	item, err := expeditions.Create(context.Background(), "Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(context.Background(), item.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := observations.Record(context.Background(), item.ID, "A-1", start, 1, 2, 100, "granite", "A clear outcrop.", nil, .8); err != nil {
		t.Fatal(err)
	}
	if _, err := specimens.Register(context.Background(), item.ID, "A-1-01", "granite", 20, start, "Ari", ""); err != nil {
		t.Fatal(err)
	}
	summary, err := summaries.Build(context.Background(), item.ID)
	if err != nil || summary.ObservationCount != 1 || summary.SpecimenCount != 1 {
		t.Fatalf("summary=%+v err=%v", summary, err)
	}
	closed, err := expeditions.Close(context.Background(), item.ID)
	if err != nil || closed.Status != model.ExpeditionClosed {
		t.Fatalf("close=%+v err=%v", closed, err)
	}
}
