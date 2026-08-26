package policy

import (
	"testing"
	"time"
)

func TestRulesCoverMeasurementAndLifecycleBoundaries(t *testing.T) {
	if !Coordinates(0, 0) || Coordinates(91, 0) || Coordinates(0, 181) {
		t.Fatal("coordinate rule failed")
	}
	if !Confidence(.7) || Confidence(-.1) || Confidence(1.1) {
		t.Fatal("confidence rule failed")
	}
	if !PositiveWeight(.1) || PositiveWeight(0) {
		t.Fatal("weight rule failed")
	}
	if !TransitionAllowed("planned", "active") || TransitionAllowed("closed", "active") {
		t.Fatal("transition rule failed")
	}
	if !ValidDate(time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)) {
		t.Fatal("date rule failed")
	}
	if WeightBand(100) != "medium" || !Confidence(.95) {
		t.Fatal("band rule failed")
	}
}
