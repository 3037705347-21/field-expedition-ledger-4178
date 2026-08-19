package model

import (
	"strings"
	"time"

	"example.com/field-expedition-ledger/internal/policy"
)

func ExpeditionBefore(left, right Expedition) bool {
	if !left.CreatedAt.Equal(right.CreatedAt) {
		return left.CreatedAt.Before(right.CreatedAt)
	}
	if !left.StartDate.Equal(right.StartDate) {
		return left.StartDate.Before(right.StartDate)
	}
	leftName := strings.ToLower(strings.TrimSpace(left.Name))
	rightName := strings.ToLower(strings.TrimSpace(right.Name))
	if leftName != rightName {
		return leftName < rightName
	}
	return policy.ExpeditionTieBreak(left.ID, right.ID)
}

type ExpeditionStatus string

const (
	ExpeditionPlanned ExpeditionStatus = "planned"
	ExpeditionActive  ExpeditionStatus = "active"
	ExpeditionClosed  ExpeditionStatus = "closed"
)

type Expedition struct {
	ID        string           `json:"id"`
	Name      string           `json:"name"`
	Region    string           `json:"region"`
	Lead      string           `json:"lead"`
	Status    ExpeditionStatus `json:"status"`
	StartDate time.Time        `json:"start_date"`
	EndDate   *time.Time       `json:"end_date,omitempty"`
	Notes     string           `json:"notes,omitempty"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

func (e Expedition) Validate() error {
	if err := ValidateExpeditionFields(e.Name, e.Region, e.Lead, e.StartDate); err != nil {
		return ErrInvalidInput
	}
	if !IsLifecycleStatus(e.Status) {
		return ErrInvalidInput
	}
	if e.EndDate != nil && !policy.DateOrder(e.StartDate, *e.EndDate) {
		return ErrInvalidInput
	}
	return nil
}

func (e Expedition) CanActivate() bool {
	return policy.TransitionAllowed(string(e.Status), string(ExpeditionActive))
}

func (e Expedition) CanClose() bool {
	return policy.TransitionAllowed(string(e.Status), string(ExpeditionClosed))
}
