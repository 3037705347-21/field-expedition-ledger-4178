package store

import (
	"context"

	"example.com/field-expedition-ledger/internal/model"
)

type Repository interface {
	CreateExpedition(context.Context, model.Expedition) (model.Expedition, error)
	GetExpedition(context.Context, string) (model.Expedition, error)
	UpdateExpedition(context.Context, model.Expedition) (model.Expedition, error)
	ListExpeditions(context.Context, model.ExpeditionFilter) ([]model.Expedition, error)
	CreateObservation(context.Context, model.Observation) (model.Observation, error)
	ListObservations(context.Context, string) ([]model.Observation, error)
	CreateSpecimen(context.Context, model.Specimen) (model.Specimen, error)
	GetSpecimen(context.Context, string) (model.Specimen, error)
	ListSpecimens(context.Context, string) ([]model.Specimen, error)
}
