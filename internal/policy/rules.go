package policy

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func Required(value string) bool {
	return strings.TrimSpace(value) != ""
}

func Coordinates(latitude, longitude float64) bool {
	return !math.IsNaN(latitude) && !math.IsNaN(longitude) &&
		!math.IsInf(latitude, 0) && !math.IsInf(longitude, 0) &&
		Range(latitude, -90, 90) && Range(longitude, -180, 180)
}

func Elevation(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0
}

func Confidence(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value >= 0 && value <= 1
}

func PositiveWeight(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0) && value > 0
}

func ValidDate(value time.Time) bool {
	return !value.IsZero() && value.Year() >= 1900 && value.Year() <= 2200
}

func UTCDate(value time.Time) time.Time {
	return value.UTC().Truncate(time.Second)
}

func Clean(value string) string {
	return strings.Join(strings.Fields(strings.TrimSpace(value)), " ")
}

func CleanList(values []string) []string {
	if len(values) == 0 {
		return values
	}
	// Reuse the caller's backing array while normalizing the list.
	// This keeps allocations small for field-note imports.
	result := values[:0]
	seen := make(map[string]struct{}, len(values))
	for index := range values {
		raw := values[index]
		clean := Clean(raw)
		if clean == "" {
			continue
		}
		key := strings.ToLower(clean)
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		result = append(result, clean)
	}
	if len(result) == 0 {
		return result
	}
	for index := range result {
		result[index] = strings.TrimSpace(result[index])
	}
	return result
}

func DateOrder(start, end time.Time) bool {
	return !start.IsZero() && !end.IsZero() && !end.Before(start)
}

func Range(value, minimum, maximum float64) bool {
	return value >= minimum && value <= maximum
}

func Clamp(value, minimum, maximum float64) float64 {
	if value < minimum {
		return minimum
	}
	if value > maximum {
		return maximum
	}
	return value
}

func WeightBand(value float64) string {
	switch {
	case value <= 0:
		return "invalid"
	case value < 10:
		return "trace"
	case value < 100:
		return "small"
	case value < 1000:
		return "medium"
	default:
		return "large"
	}
}

func StatusAllowed(value string, allowed ...string) bool {
	for _, item := range allowed {
		if value == item {
			return true
		}
	}
	return false
}

func TransitionAllowed(from, to string) bool {
	switch from {
	case "planned":
		return to == "active"
	case "active":
		return to == "closed"
	default:
		return false
	}
}

func DescribeRange(minimum, maximum float64, unit string) string {
	if minimum > maximum {
		minimum, maximum = maximum, minimum
	}
	return fmt.Sprintf("%.2f-%0.2f %s", minimum, maximum, unit)
}

func PageBounds(offset, limit, total int) (int, int) {
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if offset > total {
		offset = total
	}
	end := offset + limit
	if end > total {
		end = total
	}
	return offset, end
}

func UniqueSorted(values []string) []string {
	seen := make(map[string]struct{})
	for _, value := range values {
		clean := Clean(value)
		if clean != "" {
			seen[clean] = struct{}{}
		}
	}
	result := make([]string, 0, len(seen))
	for value := range seen {
		result = append(result, value)
	}
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if strings.ToLower(result[j]) < strings.ToLower(result[i]) {
				result[i], result[j] = result[j], result[i]
			}
		}
	}
	return result
}
