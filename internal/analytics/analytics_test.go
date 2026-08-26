package analytics

import (
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
)

func TestSummariesNormalizeAndSortCollections(t *testing.T) {
	now := time.Now().UTC()
	observations := []model.Observation{
		{SiteCode: "b-2", RockType: "Granite", ElevationM: 200, Confidence: .8, RecordedAt: now, Tags: []string{"Fresh"}},
		{SiteCode: "A-1", RockType: "basalt", ElevationM: 100, Confidence: .9, RecordedAt: now, Tags: []string{"fresh", "flow"}},
	}
	stats := SummarizeObservations(observations)
	if stats.Count != 2 || stats.MinimumElevation != 100 || stats.MaximumElevation != 200 {
		t.Fatalf("unexpected observation stats: %+v", stats)
	}
	if len(stats.Sites) != 2 || stats.Sites[0] != "A-1" {
		t.Fatalf("unexpected sites: %+v", stats.Sites)
	}
	specimens := []model.Specimen{
		{Material: "Basalt", WeightGrams: 10, Status: model.SpecimenStored, Custodian: "B"},
		{Material: "basalt", WeightGrams: 20, Status: model.SpecimenCollected, Custodian: "A"},
	}
	specStats := SummarizeSpecimens(specimens)
	if specStats.Count != 2 || specStats.TotalWeightGram != 30 {
		t.Fatalf("unexpected specimen stats: %+v", specStats)
	}
}
