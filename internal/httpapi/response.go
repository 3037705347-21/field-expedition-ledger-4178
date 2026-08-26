package httpapi

import (
	"encoding/json"
	"errors"
	"net/http"

	"example.com/field-expedition-ledger/internal/model"
)

type errorPayload struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, err error) {
	status := http.StatusBadRequest
	code := "invalid_request"
	switch {
	case errors.Is(err, model.ErrNotFound):
		status = http.StatusNotFound
		code = "not_found"
	case errors.Is(err, model.ErrInvalidState), errors.Is(err, model.ErrClosedExpedition):
		status = http.StatusConflict
		code = "invalid_state"
	case errors.Is(err, model.ErrInactiveExpedition):
		status = http.StatusConflict
		code = "inactive_expedition"
	case errors.Is(err, model.ErrInvalidInput):
		status = http.StatusBadRequest
		code = "invalid_input"
	default:
		status = http.StatusInternalServerError
		code = "internal_error"
	}
	writeJSON(w, status, errorPayload{Code: code, Message: err.Error()})
}
