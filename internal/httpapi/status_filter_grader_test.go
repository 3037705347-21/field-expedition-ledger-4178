package httpapi

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/service"
	"example.com/field-expedition-ledger/internal/store"
)

func TestExpeditionStatusFilterReturnsOnlyMatchingItems(t *testing.T) {
	handler, expeditions := newExpeditionStatusFilterHandler(t)
	for _, testCase := range []struct {
		name       string
		path       string
		status     model.ExpeditionStatus
		expectedID []string
	}{
		{name: "planned", path: "/api/expeditions?status=planned", status: model.ExpeditionPlanned, expectedID: []string{expeditions["planned"].ID}},
		{name: "active", path: "/api/expeditions?status=active", status: model.ExpeditionActive, expectedID: []string{expeditions["activeNorth"].ID, expeditions["activeSouth"].ID}},
		{name: "closed", path: "/api/expeditions?status=closed", status: model.ExpeditionClosed, expectedID: []string{expeditions["closed"].ID}},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			items := listExpeditions(t, handler, testCase.path)
			assertMatchingExpeditions(t, items, testCase.status, testCase.expectedID)
		})
	}
}

func TestExpeditionStatusFilterCombinesWithSearchAndRegion(t *testing.T) {
	handler, expeditions := newExpeditionStatusFilterHandler(t)
	items := listExpeditions(t, handler, "/api/expeditions?status=active&region=North%20Basin&q=Ridge")
	assertMatchingExpeditions(t, items, model.ExpeditionActive, []string{expeditions["activeNorth"].ID})
}

func newExpeditionStatusFilterHandler(t *testing.T) (http.Handler, map[string]model.Expedition) {
	t.Helper()
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	expeditions := service.NewExpeditionService(repository)
	start := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	items := map[string]model.Expedition{}
	items["planned"] = createExpedition(t, expeditions, ctx, "Planned Ridge", "North Basin", start)
	items["activeNorth"] = createExpedition(t, expeditions, ctx, "Ridge Survey", "North Basin", start)
	items["activeSouth"] = createExpedition(t, expeditions, ctx, "Ridge Survey", "South Basin", start)
	items["closed"] = createExpedition(t, expeditions, ctx, "Ridge Archive", "North Basin", start)
	for _, name := range []string{"activeNorth", "activeSouth", "closed"} {
		active, err := expeditions.Activate(ctx, items[name].ID)
		if err != nil {
			t.Fatal(err)
		}
		items[name] = active
	}
	closed := items["closed"]
	closed.Status = model.ExpeditionClosed
	updated, err := repository.UpdateExpedition(ctx, closed)
	if err != nil {
		t.Fatal(err)
	}
	items["closed"] = updated
	return NewHandler(repository), items
}

func createExpedition(t *testing.T, expeditions *service.ExpeditionService, ctx context.Context, name, region string, start time.Time) model.Expedition {
	t.Helper()
	item, err := expeditions.Create(ctx, name, region, "Ari", start, "")
	if err != nil {
		t.Fatal(err)
	}
	return item
}

func listExpeditions(t *testing.T, handler http.Handler, path string) []model.Expedition {
	t.Helper()
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d", response.Code)
	}
	var payload struct {
		Items []model.Expedition `json:"items"`
	}
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	return payload.Items
}

func assertMatchingExpeditions(t *testing.T, items []model.Expedition, status model.ExpeditionStatus, expectedIDs []string) {
	t.Helper()
	if len(items) != len(expectedIDs) {
		t.Fatalf("items=%+v want %d items", items, len(expectedIDs))
	}
	expected := make(map[string]struct{}, len(expectedIDs))
	for _, id := range expectedIDs {
		expected[id] = struct{}{}
	}
	for _, item := range items {
		if item.Status != status {
			t.Fatalf("item %s has status %q, want %q", item.ID, item.Status, status)
		}
		if _, ok := expected[item.ID]; !ok {
			t.Fatalf("unexpected item: %+v", item)
		}
	}
}
