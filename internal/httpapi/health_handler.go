package httpapi

import (
	"net/http"
	"time"
)

type healthHandler struct {
	startedAt time.Time
}

func (h healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"status":     "ok",
		"service":    "field-expedition-ledger",
		"started_at": h.startedAt,
	})
}
