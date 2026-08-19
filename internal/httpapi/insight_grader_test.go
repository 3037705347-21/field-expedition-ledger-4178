package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/store"
)

func TestInsightsSortSiteCountsNumerically(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	create := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/expeditions", strings.NewReader(`{"name":"Count Ridge","region":"Basin","lead":"Ari","start_date":"2026-08-18T00:00:00Z"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(create, request)
	if create.Code != http.StatusCreated {
		t.Fatalf("create status=%d body=%s", create.Code, create.Body.String())
	}
	var expedition struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(create.Body.Bytes(), &expedition); err != nil {
		t.Fatal(err)
	}

	for i := 0; i < 10; i++ {
		postObservation(t, handler, expedition.ID, "Ridge", time.Date(2026, 8, 18, i, 0, 0, 0, time.UTC))
	}
	for i := 0; i < 2; i++ {
		postObservation(t, handler, expedition.ID, "Basin", time.Date(2026, 8, 19, i, 0, 0, 0, time.UTC))
	}

	response := httptest.NewRecorder()
	insightsPath := strings.Replace("/api/expeditions/{id}/insights", "{id}", expedition.ID, 1)
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, insightsPath, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("insights status=%d body=%s", response.Code, response.Body.String())
	}
	var payload struct {
		SortedSiteCounts []string `json:"sorted_site_counts"`
	}
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if len(payload.SortedSiteCounts) != 2 || payload.SortedSiteCounts[0] != "RIDGE=10" || payload.SortedSiteCounts[1] != "BASIN=2" {
		t.Fatalf("unexpected sorted site counts: %v", payload.SortedSiteCounts)
	}
}

func postObservation(t *testing.T, handler http.Handler, expeditionID, site string, recordedAt time.Time) {
	t.Helper()
	body := `{"site_code":"` + site + `","recorded_at":"` + recordedAt.Format(time.RFC3339) + `","latitude":1,"longitude":2,"elevation_m":100,"rock_type":"granite","description":"outcrop","confidence":0.8}`
	request := httptest.NewRequest(http.MethodPost, "/api/expeditions/"+expeditionID+"/observations", strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("observation status=%d body=%s", response.Code, response.Body.String())
	}
}
