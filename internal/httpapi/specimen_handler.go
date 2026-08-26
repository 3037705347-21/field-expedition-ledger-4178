package httpapi

import (
	"net/http"

	"example.com/field-expedition-ledger/internal/service"
)

type specimenHandler struct {
	specimens *service.SpecimenService
}

func (h specimenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := pathParts(r.URL.Path)
	if len(parts) != 3 || parts[0] != "api" || parts[1] != "specimens" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}
	item, err := h.specimens.Get(r.Context(), parts[2])
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, item)
}
