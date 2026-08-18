package httpapi

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/field-expedition-ledger/internal/store"
)

func TestHealthAndListEndpoints(t *testing.T) {
	handler := NewHandler(store.NewMemoryRepository())
	server := httptest.NewServer(handler)
	defer server.Close()
	response, err := http.Get(server.URL + "/health")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("health status=%d", response.StatusCode)
	}
	response, err = http.Get(server.URL + "/api/expeditions")
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		t.Fatalf("list status=%d", response.StatusCode)
	}
	if !strings.Contains(response.Header.Get("Content-Type"), "application/json") {
		t.Fatalf("unexpected content type: %s", response.Header.Get("Content-Type"))
	}
}
