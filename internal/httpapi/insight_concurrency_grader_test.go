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

func TestInsightRefreshesAfterObservationIsRecorded(t *testing.T) {
	repository := store.NewMemoryRepository()
	now := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	expedition := model.Expedition{
		ID: "exp-concurrent-insight", Name: "Ridge", Region: "Basin", Lead: "Ari",
		Status: model.ExpeditionActive, StartDate: now, CreatedAt: now, UpdatedAt: now,
	}
	if _, err := repository.CreateExpedition(context.Background(), expedition); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	getInsight := func() map[string]any {
		request := httptest.NewRequest(http.MethodGet, "/api/expeditions/"+expedition.ID+"/insights", nil)
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("insight status=%d body=%s", recorder.Code, recorder.Body.String())
		}
		var payload map[string]any
		if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
			t.Fatal(err)
		}
		return payload
	}
	if days := getInsight()["observation_day_count"].(float64); days != 0 {
		t.Fatalf("initial day count=%v", days)
	}
	body := `{"site_code":"A-1","recorded_at":"2026-08-18T10:00:00Z","latitude":1,"longitude":2,"elevation_m":100,"rock_type":"granite","description":"clear","confidence":0.8}`
	request := httptest.NewRequest(http.MethodPost, "/api/expeditions/"+expedition.ID+"/observations", bytes.NewBufferString(body))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	if days := getInsight()["observation_day_count"].(float64); days != 1 {
		t.Fatalf("refreshed day count=%v", days)
	}
}
