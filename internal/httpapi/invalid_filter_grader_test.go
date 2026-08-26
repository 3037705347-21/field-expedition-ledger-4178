package httpapi

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/field-expedition-ledger/internal/store"
)

func TestInvalidBeforeFilterReturnsClientError(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	request := httptest.NewRequest(http.MethodGet, "/api/expeditions?before=not-a-date", nil)
	recorder := httptest.NewRecorder()

	handler.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("status=%d body=%s", recorder.Code, recorder.Body.String())
	}
	var payload map[string]string
	if err := json.Unmarshal(recorder.Body.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["code"] != "invalid_filter" {
		t.Fatalf("code=%q body=%s", payload["code"], recorder.Body.String())
	}
}
