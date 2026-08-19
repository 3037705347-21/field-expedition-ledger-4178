package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/service"
	"example.com/field-expedition-ledger/internal/store"
)

func TestInsightSiteCountsSortNumerically(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	expeditions := service.NewExpeditionService(repository)
	observations := service.NewObservationService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	item, err := expeditions.Create(ctx, "North Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(ctx, item.ID); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 10; i++ {
		if _, err := observations.Record(ctx, item.ID, "A", start.Add(time.Duration(i)*time.Hour), 1, 2, 100, "granite", "outcrop", nil, .9); err != nil {
			t.Fatal(err)
		}
	}
	for i := 0; i < 2; i++ {
		if _, err := observations.Record(ctx, item.ID, "B", start.Add(time.Duration(i+10)*time.Hour), 1, 2, 100, "granite", "outcrop", nil, .9); err != nil {
			t.Fatal(err)
		}
	}
	handler := NewHandler(repository)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/expeditions/"+item.ID+"/insights", bytes.NewReader(nil)))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d", response.Code)
	}
	var payload struct {
		SortedSiteCounts []string `json:"sorted_site_counts"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.SortedSiteCounts) < 2 || payload.SortedSiteCounts[0] != "A=10" {
		t.Fatalf("sorted site counts=%v", payload.SortedSiteCounts)
	}
}
