package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/field-expedition-ledger/internal/store"
)

func TestEmptySummaryUsesEmptyCollections(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	create := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/api/expeditions", strings.NewReader(`{"name":"Nil Ridge","region":"Basin","lead":"Ari","start_date":"2026-08-18T00:00:00Z"}`))
	request.Header.Set("Content-Type", "application/json")
	handler.ServeHTTP(create, request)
	var expedition struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(create.Body.Bytes(), &expedition); err != nil {
		t.Fatal(err)
	}
	summary := httptest.NewRecorder()
	summaryPath := strings.Replace("/api/expeditions/{id}/summary", "{id}", expedition.ID, 1)
	handler.ServeHTTP(summary, httptest.NewRequest(http.MethodGet, summaryPath, nil))
	if !strings.Contains(summary.Body.String(), `"specimen_materials":[]`) {
		t.Fatalf("empty collections were serialized as nil: %s", summary.Body.String())
	}
}
