package httpapi

import (
	"context"
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/store"
)

func TestSummaryUsesSpecimenOnlyReview(t *testing.T) {
	ctx := context.Background()
	repository := store.NewMemoryRepository()
	created := time.Date(2026, 8, 11, 0, 0, 0, 0, time.UTC)
	_, err := repository.CreateExpedition(ctx, model.Expedition{ID: "exp-summary", Name: "Inventory Ridge", Region: "North", Lead: "Ari", Status: model.ExpeditionActive, StartDate: created, CreatedAt: created, UpdatedAt: created})
	if err != nil {
		t.Fatal(err)
	}
	_, err = repository.CreateSpecimen(ctx, model.Specimen{ID: "spc-summary", ExpeditionID: "exp-summary", Label: "S-1", Material: "Granite", WeightGrams: 2.5, CollectedAt: created, Custodian: "Ari", Status: model.SpecimenCollected})
	if err != nil {
		t.Fatal(err)
	}
	request := httptest.NewRequest(http.MethodGet, "/api/expeditions/{id}/summary", nil)
	request.URL.Path = "/api/expeditions/exp-summary/summary"
	response := httptest.NewRecorder()
	NewHandler(repository).ServeHTTP(response, request)
	if response.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", response.Code, response.Body.String())
	}
	var summary model.ExpeditionSummary
	if err := json.NewDecoder(response.Body).Decode(&summary); err != nil {
		t.Fatal(err)
	}
	if math.Abs(summary.QualityScore-0.25) > 0.0001 || summary.FollowUp != "record the first site observation" {
		t.Fatalf("quality=%v follow_up=%q, want specimen-only review values", summary.QualityScore, summary.FollowUp)
	}
}
