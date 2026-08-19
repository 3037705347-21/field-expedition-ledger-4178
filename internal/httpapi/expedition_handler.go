package httpapi

import (
	"net/http"
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/service"
)

type expeditionHandler struct {
	expeditions  *service.ExpeditionService
	observations *service.ObservationService
	specimens    *service.SpecimenService
	summaries    *service.SummaryService
	insights     *service.InsightService
	query        *service.QueryService
}

func (h expeditionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)
	if len(parts) == 2 && parts[0] == "api" && parts[1] == "expeditions" && r.Method == http.MethodPost {
		h.create(w, r)
		return
	}
	if len(parts) == 2 && parts[0] == "api" && parts[1] == "expeditions" && r.Method == http.MethodGet {
		h.list(w, r)
		return
	}
	if len(parts) < 4 || parts[0] != "api" || parts[1] != "expeditions" {
		http.NotFound(w, r)
		return
	}
	id := parts[2]
	switch {
	case len(parts) == 4 && parts[3] == "activate" && r.Method == http.MethodPost:
		h.transition(w, r, id, true)
	case len(parts) == 4 && parts[3] == "close" && r.Method == http.MethodPost:
		h.transition(w, r, id, false)
	case len(parts) == 4 && parts[3] == "observations" && r.Method == http.MethodPost:
		h.createObservation(w, r, id)
	case len(parts) == 4 && parts[3] == "observations" && r.Method == http.MethodGet:
		h.listObservations(w, r, id)
	case len(parts) == 4 && parts[3] == "specimens" && r.Method == http.MethodPost:
		h.createSpecimen(w, r, id)
	case len(parts) == 4 && parts[3] == "specimens" && r.Method == http.MethodGet:
		h.listSpecimens(w, r, id)
	case len(parts) == 4 && parts[3] == "summary" && r.Method == http.MethodGet:
		h.summary(w, r, id)
	case len(parts) == 4 && parts[3] == "insights" && r.Method == http.MethodGet:
		h.insightsView(w, r, id)
	default:
		http.NotFound(w, r)
	}
}

func (h expeditionHandler) create(w http.ResponseWriter, r *http.Request) {
	var request createExpeditionRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, model.ErrInvalidInput)
		return
	}
	item, err := h.expeditions.Create(r.Context(), request.Name, request.Region, request.Lead, request.StartDate, request.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h expeditionHandler) list(w http.ResponseWriter, r *http.Request) {
	status := queryValue(r, "status")
	region := queryValue(r, "region")
	search := queryValue(r, "q")
	var (
		items []model.Expedition
		err   error
	)
	switch {
	case search != "":
		items, err = h.query.Search(r.Context(), search)
	case status == string(model.ExpeditionActive):
		items, err = h.query.Active(r.Context())
	case region != "":
		items, err = h.query.ByRegion(r.Context(), region)
	default:
		page, pageErr := h.query.Page(r.Context(), model.ExpeditionFilter{}, 0, 50)
		items, err = page.Items, pageErr
	}
	if err != nil {
		writeError(w, err)
		return
	}
	if queryValue(r, "newest") == "true" {
		items = service.SortByStartDate(items, true)
	}
	if before := queryValue(r, "before"); before != "" {
		moment, parseErr := model.ParseBeforeFilter(before)
		if parseErr != nil {
			writeError(w, parseErr)
			return
		}
		items = service.FilterStartedBefore(items, moment)
	}
	if queryValue(r, "notes") == "true" {
		items = service.FilterWithNotes(items)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"items":   items,
		"total":   len(items),
		"offset":  0,
		"limit":   len(items),
		"regions": service.RegionNames(items),
		"leads":   service.LeadNames(items),
	})
}

func (h expeditionHandler) transition(w http.ResponseWriter, r *http.Request, id string, activate bool) {
	var (
		item model.Expedition
		err  error
	)
	if activate {
		item, err = h.expeditions.Activate(r.Context(), id)
	} else {
		item, err = h.expeditions.Close(r.Context(), id)
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h expeditionHandler) createObservation(w http.ResponseWriter, r *http.Request, expeditionID string) {
	var request createObservationRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, model.ErrInvalidInput)
		return
	}
	item, err := h.observations.Record(r.Context(), expeditionID, request.SiteCode, request.RecordedAt, request.Latitude, request.Longitude, request.ElevationM, request.RockType, request.Description, request.Tags, request.Confidence)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h expeditionHandler) listObservations(w http.ResponseWriter, r *http.Request, expeditionID string) {
	var (
		items []model.Observation
		err   error
	)
	if since := queryValue(r, "since"); since != "" {
		moment, parseErr := time.Parse(time.RFC3339, since)
		if parseErr != nil {
			writeError(w, model.ErrInvalidInput)
			return
		}
		items, err = h.observations.Recent(r.Context(), expeditionID, moment)
	} else {
		items, err = h.observations.List(r.Context(), expeditionID)
	}
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.NewPage(items, len(items), 0, len(items)))
}

func (h expeditionHandler) createSpecimen(w http.ResponseWriter, r *http.Request, expeditionID string) {
	var request createSpecimenRequest
	if err := decodeJSON(r, &request); err != nil {
		writeError(w, model.ErrInvalidInput)
		return
	}
	item, err := h.specimens.Register(r.Context(), expeditionID, request.Label, request.Material, request.WeightGrams, request.CollectedAt, request.Custodian, request.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, item)
}

func (h expeditionHandler) listSpecimens(w http.ResponseWriter, r *http.Request, expeditionID string) {
	items, err := h.specimens.List(r.Context(), expeditionID)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, model.NewPage(items, len(items), 0, len(items)))
}

func (h expeditionHandler) summary(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.summaries.Build(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (h expeditionHandler) insightsView(w http.ResponseWriter, r *http.Request, id string) {
	item, err := h.insights.Build(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func pathParts(path string) []string {
	trimmed := strings.Trim(path, "/")
	if trimmed == "" {
		return nil
	}
	return strings.Split(trimmed, "/")
}
