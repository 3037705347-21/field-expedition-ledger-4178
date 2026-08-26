package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/store"
)

type validationErrorPayload struct {
	Code string `json:"code"`
}

func TestInvalidRecordReturnsClientError(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	healthRequest := httptest.NewRequest(http.MethodGet, "/health", nil)
	healthResponse := httptest.NewRecorder()
	handler.ServeHTTP(healthResponse, healthRequest)
	if healthResponse.Code != http.StatusOK {
		t.Fatalf("health status=%d", healthResponse.Code)
	}
	expeditionID := createActiveExpedition(t, handler)

	observationResponse := serveJSON(t, handler, http.MethodPost, "/api/expeditions/"+expeditionID+"/observations", map[string]any{
		"site_code":   "",
		"recorded_at": time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
		"latitude":    1,
		"longitude":   2,
		"elevation_m": 100,
		"rock_type":   "granite",
		"description": "clear outcrop",
		"confidence":  0.8,
	})
	assertInvalidInputResponse(t, "observation", observationResponse)

	specimenResponse := serveJSON(t, handler, http.MethodPost, "/api/expeditions/"+expeditionID+"/specimens", map[string]any{
		"label":        "A-1-01",
		"material":     "granite",
		"weight_grams": 0,
		"collected_at": time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
		"custodian":    "Ari",
	})
	assertInvalidInputResponse(t, "specimen", specimenResponse)

	validObservationResponse := serveJSON(t, handler, http.MethodPost, "/api/expeditions/"+expeditionID+"/observations", map[string]any{
		"site_code":   "A-1",
		"recorded_at": time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
		"latitude":    1,
		"longitude":   2,
		"elevation_m": 100,
		"rock_type":   "granite",
		"description": "clear outcrop",
		"confidence":  0.8,
	})
	if validObservationResponse.Code != http.StatusCreated {
		t.Fatalf("valid observation status=%d body=%s", validObservationResponse.Code, validObservationResponse.Body.String())
	}
}

func createActiveExpedition(t *testing.T, handler http.Handler) string {
	t.Helper()
	created := serveJSON(t, handler, http.MethodPost, "/api/expeditions", map[string]any{
		"name":       "Ridge",
		"region":     "Basin",
		"lead":       "Ari",
		"start_date": time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC),
	})
	if created.Code != http.StatusCreated {
		t.Fatalf("create expedition status=%d body=%s", created.Code, created.Body.String())
	}
	var expedition struct {
		ID string `json:"id"`
	}
	if err := json.NewDecoder(created.Body).Decode(&expedition); err != nil {
		t.Fatal(err)
	}
	activated := serveJSON(t, handler, http.MethodPost, "/api/expeditions/"+expedition.ID+"/activate", nil)
	if activated.Code != http.StatusOK {
		t.Fatalf("activate expedition status=%d body=%s", activated.Code, activated.Body.String())
	}
	return expedition.ID
}

func serveJSON(t *testing.T, handler http.Handler, method, path string, value any) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	if value != nil {
		if err := json.NewEncoder(&body).Encode(value); err != nil {
			t.Fatal(err)
		}
	}
	request := httptest.NewRequest(method, path, &body)
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func assertInvalidInputResponse(t *testing.T, name string, response *httptest.ResponseRecorder) {
	t.Helper()
	if response.Code != http.StatusBadRequest {
		t.Fatalf("%s status=%d body=%s", name, response.Code, response.Body.String())
	}
	var payload validationErrorPayload
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatalf("%s decode error: %v", name, err)
	}
	if payload.Code != "invalid_input" {
		t.Fatalf("%s code=%q", name, payload.Code)
	}
}
