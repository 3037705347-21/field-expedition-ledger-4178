package httpapi

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/service"
	"example.com/field-expedition-ledger/internal/store"
)

func TestSpecimenLabelsAreUniquePerExpedition(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	expeditions := service.NewExpeditionService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	first, err := expeditions.Create(ctx, "North Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(ctx, first.ID); err != nil {
		t.Fatal(err)
	}
	second, err := expeditions.Create(ctx, "South Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := expeditions.Activate(ctx, second.ID); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	body := []byte(`{"label":"R-1","material":"granite","weight_grams":20,"collected_at":"2026-08-18T00:00:00Z","custodian":"Ari"}`)
	create := func(id string) *httptest.ResponseRecorder {
		response := httptest.NewRecorder()
		handler.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/api/expeditions/"+id+"/specimens", bytes.NewReader(body)))
		return response
	}
	if response := create(first.ID); response.Code != http.StatusCreated {
		t.Fatalf("first status=%d", response.Code)
	}
	duplicate := create(first.ID)
	if duplicate.Code != http.StatusConflict {
		t.Fatalf("duplicate status=%d", duplicate.Code)
	}
	otherExpedition := create(second.ID)
	if otherExpedition.Code != http.StatusCreated {
		t.Fatalf("other expedition status=%d", otherExpedition.Code)
	}
}
