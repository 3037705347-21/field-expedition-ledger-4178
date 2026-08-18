package store

import (
	"context"
	"time"

	"example.com/field-expedition-ledger/internal/model"
)

func SeedDemo(ctx context.Context, repository Repository) error {
	now := time.Date(2026, 8, 18, 0, 0, 0, 0, time.UTC)
	expedition := model.Expedition{
		ID:        "exp-demo-ridge",
		Name:      "Ridge of Quiet Basalt",
		Region:    "North Basin",
		Lead:      "Mara Chen",
		Status:    model.ExpeditionActive,
		StartDate: now,
		CreatedAt: now,
		UpdatedAt: now,
		Notes:     "Demonstration record for the local notebook.",
	}
	if _, err := repository.CreateExpedition(ctx, expedition); err != nil {
		return err
	}
	observation := model.Observation{
		ID:           "obs-demo-001",
		ExpeditionID: expedition.ID,
		SiteCode:     "NB-01",
		RecordedAt:   now.Add(2 * time.Hour),
		Latitude:     41.1,
		Longitude:    -110.2,
		ElevationM:   1820,
		RockType:     "basalt",
		Description:  "Dark vesicular flow with a narrow weathered edge.",
		Tags:         []string{"vesicular", "flow"},
		Confidence:   0.92,
	}
	if _, err := repository.CreateObservation(ctx, observation); err != nil {
		return err
	}
	specimen := model.Specimen{
		ID:           "spc-demo-001",
		ExpeditionID: expedition.ID,
		Label:        "NB-01-A",
		Material:     "basalt",
		WeightGrams:  184.5,
		CollectedAt:  now.Add(3 * time.Hour),
		Custodian:    "Mara Chen",
		Status:       model.SpecimenStored,
		Notes:        "Wrapped in archival paper.",
	}
	_, err := repository.CreateSpecimen(ctx, specimen)
	return err
}
