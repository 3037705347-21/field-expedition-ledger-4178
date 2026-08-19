package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/field-expedition-ledger/internal/store"
)

func TestInvalidExpeditionKeepsClientErrorStatus(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	request := httptest.NewRequest(http.MethodPost, "/api/expeditions", strings.NewReader(`{"name":"","region":"Basin","lead":"Ari","start_date":"2026-08-18T00:00:00Z"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("invalid input status=%d body=%s", response.Code, response.Body.String())
	}
}
