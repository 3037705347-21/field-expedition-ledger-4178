package httpapi

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestSummaryAccumulatesWeightsForRepeatedMaterial(t *testing.T) {
	repository := store.NewMemoryRepository()
	now := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	expedition := model.Expedition{
		ID: "exp-weight-total", Name: "Ridge", Region: "Basin", Lead: "Ari",
		Status: model.ExpeditionActive, StartDate: now, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := repository.CreateExpedition(context.Background(), expedition); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	for _, body := range []string{
		`{"label":"A-1","material":"Granite","weight_grams":12,"collected_at":"2026-08-18T10:00:00Z","custodian":"Ari"}`,
		`{"label":"A-2","material":"granite","weight_grams":18,"collected_at":"2026-08-18T11:00:00Z","custodian":"Ari"}`,
	} {
		request := httptest.NewRequest(http.MethodPost, "/api/expeditions/"+expedition.ID+"/specimens", bytes.NewBufferString(body))
		request.Header.Set("Content-Type", "application/json")
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusCreated {
			t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
		}
	}
	request := httptest.NewRequest(http.MethodGet, "/api/expeditions/"+expedition.ID+"/summary", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("summary status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var payload struct {
		MaterialWeights map[string]float64 `json:"material_weights"`
	}
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload.MaterialWeights["granite"] != 30 {
		t.Fatalf("material weights=%v", payload.MaterialWeights)
	}
}
