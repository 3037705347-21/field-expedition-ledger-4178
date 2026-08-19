package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"example.com/field-expedition-ledger/internal/store"
)

func TestCanceledRecentObservationRequestStopsEarly(t *testing.T) {
	repository := store.NewMemoryRepository()
	if err := store.SeedDemo(context.Background(), repository); err != nil {
		t.Fatal(err)
	}
	handler := NewHandler(repository)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	request := httptest.NewRequest(http.MethodGet, "/api/expeditions/exp-demo-ridge/observations?since=2026-08-18T00:00:00Z", nil).WithContext(ctx)
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != http.StatusInternalServerError {
		t.Fatalf("canceled request status=%d", response.Code)
	}
}
