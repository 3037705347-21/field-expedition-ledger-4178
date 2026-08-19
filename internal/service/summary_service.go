package service

import (
	"context"

	"example.com/field-expedition-ledger/internal/analysis"
	"example.com/field-expedition-ledger/internal/analytics"
	"example.com/field-expedition-ledger/internal/model"
	"example.com/field-expedition-ledger/internal/policy"
	"example.com/field-expedition-ledger/internal/store"
)

type SummaryService struct {
	repository store.Repository
}

func NewSummaryService(repository store.Repository) *SummaryService {
	return &SummaryService{repository: repository}
}

func (s *SummaryService) Build(ctx context.Context, expeditionID string) (model.ExpeditionSummary, error) {
	expedition, err := s.repository.GetExpedition(ctx, expeditionID)
	if err != nil {
		return model.ExpeditionSummary{}, err
	}
	observations, err := s.repository.ListObservations(ctx, expeditionID)
	if err != nil {
		return model.ExpeditionSummary{}, err
	}
	specimens, err := s.repository.ListSpecimens(ctx, expeditionID)
	if err != nil {
		return model.ExpeditionSummary{}, err
	}
	summary := model.NewExpeditionSummary(expedition)
	summary.ObservationCount = len(observations)
	summary.SpecimenCount = len(specimens)
	expeditionStats := analytics.SummarizeExpeditions([]model.Expedition{expedition})
	for _, status := range analytics.StatusNames() {
		summary.StatusCounts[status] = expeditionStats.CountByStatus[status]
	}
	summary.ActiveCount = analytics.CountStatus([]model.Expedition{expedition}, model.ExpeditionActive)
	summary.RockTypes = analytics.UniqueRockTypes(observations)
	summary.SiteGroupCount = len(analytics.GroupBySite(observations))
	specimenStats := analytics.SummarizeSpecimens(specimens)
	summary.SpecimenMaterials = specimenStats.Materials
	summary.MaterialWeights = analytics.WeightByMaterial(specimens)
	summary.SpecimenStatuses = analytics.StatusBreakdown(specimens)
	review := analysis.BuildReview(observations, specimens)
	if len(observations) == 0 && !specimenStats.HasInventory() {
		return summary, nil
	}
	review.ApplyToSummary(&summary)
	if len(observations) == 0 {
		return summary, nil
	}
	derived := analytics.SummarizeObservations(observations)
	minimum, maximum, hasRange := analytics.ElevationRange(observations)
	if hasRange {
		summary.MinimumElevationM = minimum
		summary.MaximumElevationM = maximum
		summary.ElevationRange = policy.DescribeRange(minimum, maximum, "m")
	}
	summary.MinimumElevationM = derived.MinimumElevation
	summary.MaximumElevationM = derived.MaximumElevation
	summary.AverageConfidence = derived.AverageConfidence
	return summary, nil
}
