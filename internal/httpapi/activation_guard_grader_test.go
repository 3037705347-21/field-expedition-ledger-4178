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

func TestInactiveExpeditionRejectsRecords(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	expeditions := service.NewExpeditionService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	item, err := expeditions.Create(ctx, "North Ridge", "Basin", "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	observationBody := []byte(`{"site_code":"NR-1","recorded_at":"2026-08-18T00:00:00Z","latitude":1,"longitude":2,"elevation_m":100,"rock_type":"granite","description":"outcrop","confidence":0.9}`)
	observationBlocked := httptest.NewRecorder()
	handler.ServeHTTP(observationBlocked, httptest.NewRequest(http.MethodPost, "/api/expeditions/"+item.ID+"/observations", bytes.NewReader(observationBody)))
	if observationBlocked.Code != http.StatusConflict {
		t.Fatalf("inactive observation status=%d", observationBlocked.Code)
	}
	specimenBody := []byte(`{"label":"NR-1-01","material":"granite","weight_grams":20,"collected_at":"2026-08-18T00:00:00Z","custodian":"Ari"}`)
	specimenBlocked := httptest.NewRecorder()
	handler.ServeHTTP(specimenBlocked, httptest.NewRequest(http.MethodPost, "/api/expeditions/"+item.ID+"/specimens", bytes.NewReader(specimenBody)))
	if specimenBlocked.Code != http.StatusConflict {
		t.Fatalf("inactive specimen status=%d", specimenBlocked.Code)
	}
	if _, err := expeditions.Activate(ctx, item.ID); err != nil {
		t.Fatal(err)
	}
	observationAllowed := httptest.NewRecorder()
	handler.ServeHTTP(observationAllowed, httptest.NewRequest(http.MethodPost, "/api/expeditions/"+item.ID+"/observations", bytes.NewReader(observationBody)))
	if observationAllowed.Code != http.StatusCreated {
		t.Fatalf("active observation status=%d", observationAllowed.Code)
	}
	specimenAllowed := httptest.NewRecorder()
	handler.ServeHTTP(specimenAllowed, httptest.NewRequest(http.MethodPost, "/api/expeditions/"+item.ID+"/specimens", bytes.NewReader(specimenBody)))
	if specimenAllowed.Code != http.StatusCreated {
		t.Fatalf("active specimen status=%d", specimenAllowed.Code)
	}
}
