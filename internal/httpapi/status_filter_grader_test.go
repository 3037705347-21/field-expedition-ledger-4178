package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/service"
	"example.com/field-expedition-ledger/internal/store"
)

func TestExpeditionStatusFilterReturnsOnlyMatchingItems(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	expeditions := service.NewExpeditionService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	planned, err := expeditions.Create(ctx, "Planned Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	active, err := expeditions.Create(ctx, "Active Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(ctx, active.ID); err != nil {
		t.Fatal(err)
	}
	closed, err := expeditions.Create(ctx, "Closed Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(ctx, closed.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := repository.UpdateExpedition(ctx, model.Expedition{
		ID: closed.ID, Name: closed.Name, Region: closed.Region, Lead: closed.Lead,
		Status: model.ExpeditionClosed, StartDate: closed.StartDate,
		CreatedAt: closed.CreatedAt, UpdatedAt: closed.UpdatedAt,
	}); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/expeditions?status=closed", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d", response.Code)
	}
	var payload struct {
		Items []model.Expedition `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Items) != 1 || payload.Items[0].ID != closed.ID || payload.Items[0].Status != model.ExpeditionClosed {
		t.Fatalf("filtered items=%+v planned=%s active=%s", payload.Items, planned.ID, active.ID)
	}
}
