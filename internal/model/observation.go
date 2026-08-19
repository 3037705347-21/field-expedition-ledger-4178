package model

import (
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/policy"
)

type Observation struct {
	ID           string    `json:"id"`
	ExpeditionID string    `json:"expedition_id"`
	SiteCode     string    `json:"site_code"`
	RecordedAt   time.Time `json:"recorded_at"`
	Latitude     float64   `json:"latitude"`
	Longitude    float64   `json:"longitude"`
	ElevationM   float64   `json:"elevation_m"`
	RockType     string    `json:"rock_type"`
	Description  string    `json:"description"`
	Tags         []string  `json:"tags,omitempty"`
	Confidence   float64   `json:"confidence"`
}

func (o Observation) Validate() error {
	if err := ValidateObservationFields(o.ExpeditionID, o.SiteCode, o.RockType, o.Description, o.RecordedAt, o.Latitude, o.Longitude, o.ElevationM, o.Confidence); err != nil {
		return ErrInvalidInput
	}
	return nil
}

func (o Observation) Normalized() Observation {
	o.SiteCode = strings.ToUpper(NormalizeName(o.SiteCode))
	o.RockType = NormalizeName(o.RockType)
	o.Description = NormalizeName(o.Description)
	o.Tags = policy.CleanList(policy.CloneStrings(o.Tags))
	o.RecordedAt = o.RecordedAt.UTC()
	return o
}
