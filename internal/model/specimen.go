package model

import (
	"strings"
	"time"
)

type SpecimenStatus string

const (
	SpecimenCollected SpecimenStatus = "collected"
	SpecimenStored    SpecimenStatus = "stored"
	SpecimenReleased  SpecimenStatus = "released"
)

type Specimen struct {
	ID           string         `json:"id"`
	ExpeditionID string         `json:"expedition_id"`
	Label        string         `json:"label"`
	Material     string         `json:"material"`
	WeightGrams  float64        `json:"weight_grams"`
	CollectedAt  time.Time      `json:"collected_at"`
	Custodian    string         `json:"custodian"`
	Status       SpecimenStatus `json:"status"`
	Notes        string         `json:"notes,omitempty"`
}

func (s Specimen) MaterialKey() string {
	return strings.ToLower(strings.Join(strings.Fields(strings.TrimSpace(s.Material)), " "))
}

func (s Specimen) Validate() error {
	if err := ValidateSpecimenFields(s.ExpeditionID, s.Label, s.Material, s.Custodian, s.WeightGrams, s.CollectedAt); err != nil {
		return ErrInvalidInput
	}
	if !IsSpecimenStatus(s.Status) {
		return ErrInvalidInput
	}
	return nil
}

func (s Specimen) Normalized() Specimen {
	s.Label = NormalizeName(s.Label)
	s.Material = NormalizeName(s.Material)
	s.Custodian = NormalizeName(s.Custodian)
	s.CollectedAt = s.CollectedAt.UTC()
	return s
}
