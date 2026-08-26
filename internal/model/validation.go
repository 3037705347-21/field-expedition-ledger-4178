package model

import (
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/policy"
)

func ValidateExpeditionFields(name, region, lead string, startDate time.Time) error {
	if !policy.Required(name) || !policy.Required(region) || !policy.Required(lead) || !policy.ValidDate(startDate) {
		return ErrInvalidInput
	}
	return nil
}

func ValidateObservationFields(expeditionID, siteCode, rockType, description string, recordedAt time.Time, latitude, longitude, elevation, confidence float64) error {
	if !policy.Required(expeditionID) || !policy.Required(siteCode) || !policy.Required(rockType) || !policy.Required(description) {
		return ErrInvalidInput
	}
	if !policy.ValidDate(recordedAt) || !policy.Coordinates(latitude, longitude) || !policy.Elevation(elevation) || !policy.Confidence(confidence) {
		return ErrInvalidInput
	}
	return nil
}

func ValidateSpecimenFields(expeditionID, label, material, custodian string, weight float64, collectedAt time.Time) error {
	if !policy.Required(expeditionID) || !policy.Required(label) || !policy.Required(material) || !policy.Required(custodian) {
		return ErrInvalidInput
	}
	if !policy.PositiveWeight(weight) || !policy.ValidDate(collectedAt) {
		return ErrInvalidInput
	}
	return nil
}

func NormalizeName(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func IsLifecycleStatus(value ExpeditionStatus) bool {
	return policy.StatusAllowed(string(value), string(ExpeditionPlanned), string(ExpeditionActive), string(ExpeditionClosed))
}

func IsSpecimenStatus(value SpecimenStatus) bool {
	return policy.StatusAllowed(string(value), string(SpecimenCollected), string(SpecimenStored), string(SpecimenReleased))
}
