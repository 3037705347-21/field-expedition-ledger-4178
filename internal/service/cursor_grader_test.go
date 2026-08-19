package service

import (
	"context"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestRecentExcludesObservationAtCursor(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	start := time.Date(2026, 8, 10, 0, 0, 0, 0, time.UTC)
	_, err := repository.CreateExpedition(ctx, model.Expedition{ID: "exp-1", Name: "Cursor Ridge", Region: "North", Lead: "Ari", Status: model.ExpeditionActive, StartDate: start, CreatedAt: start, UpdatedAt: start})
	if err != nil { t.Fatal(err) }
	service := NewObservationService(repository)
	for _, recordedAt := range []time.Time{start, start.Add(time.Hour)} {
		_, err = service.Record(ctx, "exp-1", "S1", recordedAt, 10, 20, 100, "granite", "cursor sample", nil, .9)
		if err != nil { t.Fatal(err) }
	}
	items, err := service.Recent(ctx, "exp-1", start)
	if err != nil { t.Fatal(err) }
	if len(items) != 1 || !items[0].RecordedAt.Equal(start.Add(time.Hour)) {
		t.Fatalf("recent observations=%v, want only the record after the cursor", items)
	}
}
