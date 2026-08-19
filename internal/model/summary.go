package model

type ExpeditionSummary struct {
	ExpeditionID       string             `json:"expedition_id"`
	ExpeditionName     string             `json:"expedition_name"`
	Status             string             `json:"status"`
	ObservationCount   int                `json:"observation_count"`
	SpecimenCount      int                `json:"specimen_count"`
	MinimumElevationM  float64            `json:"minimum_elevation_m"`
	MaximumElevationM  float64            `json:"maximum_elevation_m"`
	AverageConfidence  float64            `json:"average_confidence"`
	SiteCount           int               `json:"site_count"`
	QualityScore        float64            `json:"quality_score"`
	FollowUp            string             `json:"follow_up"`
	RockTypes           []string           `json:"rock_types"`
	SiteGroupCount      int                `json:"site_group_count"`
	SpecimenMaterials   []string           `json:"specimen_materials"`
	MaterialWeights     map[string]float64 `json:"material_weights"`
	SpecimenStatuses    []string           `json:"specimen_statuses"`
	StatusCounts        map[string]int     `json:"status_counts"`
	ActiveCount         int                `json:"active_count"`
	ElevationRange      string             `json:"elevation_range"`
}

func (s ExpeditionSummary) NormalizeCollections() ExpeditionSummary {
	if s.RockTypes == nil {
		s.RockTypes = []string{}
	}
	if s.SpecimenMaterials == nil {
		s.SpecimenMaterials = []string{}
	}
	if s.MaterialWeights == nil {
		s.MaterialWeights = map[string]float64{}
	}
	if s.SpecimenStatuses == nil {
		s.SpecimenStatuses = []string{}
	}
	if s.StatusCounts == nil {
		s.StatusCounts = map[string]int{}
	}
	return s
}
