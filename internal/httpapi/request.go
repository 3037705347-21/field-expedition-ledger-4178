package httpapi

import (
	"encoding/json"
	"net/http"
	"time"
)

type createExpeditionRequest struct {
	Name      string    `json:"name"`
	Region    string    `json:"region"`
	Lead      string    `json:"lead"`
	StartDate time.Time `json:"start_date"`
	Notes     string    `json:"notes"`
}

type createObservationRequest struct {
	SiteCode    string    `json:"site_code"`
	RecordedAt  time.Time `json:"recorded_at"`
	Latitude    float64   `json:"latitude"`
	Longitude   float64   `json:"longitude"`
	ElevationM  float64   `json:"elevation_m"`
	RockType    string    `json:"rock_type"`
	Description string    `json:"description"`
	Tags        []string  `json:"tags"`
	Confidence  float64   `json:"confidence"`
}

type createSpecimenRequest struct {
	Label       string    `json:"label"`
	Material    string    `json:"material"`
	WeightGrams float64   `json:"weight_grams"`
	CollectedAt time.Time `json:"collected_at"`
	Custodian   string    `json:"custodian"`
	Notes       string    `json:"notes"`
}

func decodeJSON(r *http.Request, destination any) error {
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	return decoder.Decode(destination)
}

func queryValue(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}
