package httpapi

import (
	"net/http"
	"time"

	"example.com/field-expedition-ledger/internal/service"
	"example.com/field-expedition-ledger/internal/store"
)

func NewHandler(repository store.Repository) http.Handler {
	expeditions := service.NewExpeditionService(repository)
	observations := service.NewObservationService(repository)
	specimens := service.NewSpecimenService(repository)
	summaries := service.NewSummaryService(repository)
	insights := service.NewInsightService(repository)
	query := service.NewQueryService(repository)
	expeditionRoutes := expeditionHandler{
		expeditions:  expeditions,
		observations: observations,
		specimens:    specimens,
		summaries:    summaries,
		insights:     insights,
		query:        query,
	}
	mux := http.NewServeMux()
	mux.Handle("/health", healthHandler{startedAt: time.Now().UTC()})
	mux.Handle("/api/expeditions", expeditionRoutes)
	mux.Handle("/api/expeditions/", expeditionRoutes)
	mux.Handle("/api/specimens/", specimenHandler{specimens: specimens})
	mux.Handle("/", staticHandler())
	return requestLogger(mux)
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}
