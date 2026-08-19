package httpapi

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
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

func observationCursor(r *http.Request) (time.Time, bool, error) {
	raw := strings.TrimSpace(queryValue(r, "since"))
	if raw == "" {
		return time.Time{}, false, nil
	}
	moment, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return time.Time{}, true, model.ErrInvalidInput
	}
	return model.NormalizeTime(moment), true, nil
}
