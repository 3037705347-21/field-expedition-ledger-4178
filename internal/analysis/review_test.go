package analysis

import (
	"testing"
	"time"

	"example.com/field-expedition-ledger/internal/model"
)

func TestBuildReviewSuggestsSynthesis(t *testing.T) {
	now := time.Now().UTC()
	observations := []model.Observation{
		{SiteCode: "A", RecordedAt: now, Latitude: 1, Longitude: 1, RockType: "granite", Description: "clear", Confidence: .95},
		{SiteCode: "B", RecordedAt: now.Add(time.Hour), Latitude: 1.1, Longitude: 1.1, RockType: "basalt", Description: "clear", Confidence: .95},
	}
	specimens := []model.Specimen{{Material: "granite", WeightGrams: 10, Custodian: "A"}}
	review := BuildReview(observations, specimens)
	if review.SiteCount != 2 || review.FollowUp != "ready for synthesis" {
		t.Fatalf("unexpected review: %+v", review)
	}
	if len(ReviewWarnings(observations, specimens)) != 0 {
		t.Fatalf("unexpected warnings")
	}
}
