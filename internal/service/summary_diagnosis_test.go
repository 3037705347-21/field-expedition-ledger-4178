package service

import (
	"context"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestEmptyExpeditionSummaryProvidesFirstObservationFollowUp(t *testing.T) {
	repository := store.NewMemoryRepository()
	now := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	if _, err := repository.CreateExpedition(context.Background(), model.Expedition{
		ID: "exp-empty-review", Name: "Quiet Ridge", Region: "Basin", Lead: "Ari",
		Status: model.ExpeditionActive, StartDate: now, CreatedAt: now, UpdatedAt: now,
	}); err != nil {
		t.Fatal(err)
	}

	summary, err := NewSummaryService(repository).Build(context.Background(), "exp-empty-review")
	if err != nil {
		t.Fatal(err)
	}
	if summary.FollowUp != "record the first site observation" {
		t.Fatalf("follow_up=%q", summary.FollowUp)
	}
}
