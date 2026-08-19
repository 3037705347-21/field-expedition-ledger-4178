package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestExpeditionListUsesStableIDTieBreakForEqualCreationTimes(t *testing.T) {
	repository := store.NewMemoryRepository()
	now := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	for _, item := range []model.Expedition{
		{ID: "exp-b", Name: "Bravo", Region: "Basin", Lead: "Ari", Status: model.ExpeditionPlanned, StartDate: now, CreatedAt: now, UpdatedAt: now},
		{ID: "exp-a", Name: "Alpha", Region: "Basin", Lead: "Ari", Status: model.ExpeditionPlanned, StartDate: now, CreatedAt: now, UpdatedAt: now},
	} {
		if _, err := repository.CreateExpedition(context.Background(), item); err != nil {
			t.Fatal(err)
		}
	}
	request := httptest.NewRequest(http.MethodGet, "/api/expeditions", nil)
	recorder := httptest.NewRecorder()
	NewHandler(repository).ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var payload struct {
		Items []model.Expedition `json:"items"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.Items) != 2 || payload.Items[0].ID != "exp-a" || payload.Items[1].ID != "exp-b" {
		t.Fatalf("items=%v", payload.Items)
	}
}
