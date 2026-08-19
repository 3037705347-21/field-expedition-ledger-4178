package policy

import (
	"sort"
	"strconv"
	"strings"
)

func SortCountLabelsDescending(labels []string) []string {
	result := append([]string(nil), labels...)
	sort.SliceStable(result, func(i, j int) bool {
		leftSite, leftCount, leftOK := splitCountLabel(result[i])
		rightSite, rightCount, rightOK := splitCountLabel(result[j])
		if !leftOK || !rightOK {
			return result[i] < result[j]
		}
		if leftCount == rightCount {
			return leftSite < rightSite
		}
		return leftCount > rightCount
	})
	return result
}

func splitCountLabel(value string) (string, int, bool) {
	site := CountLabelSite(value)
	count, ok := CountLabelValue(value)
	if !ok {
		return site, 0, false
	}
	return site, count, true
}

func CountLabelSite(value string) string {
	separator := strings.LastIndex(value, "=")
	if separator <= 0 {
		return value
	}
	return value[:separator]
}

func CountLabelValue(value string) (int, bool) {
	separator := strings.LastIndex(value, "=")
	if separator <= 0 || separator == len(value)-1 {
		return 0, false
	}
	count, err := strconv.Atoi(value[separator+1:])
	if err != nil {
		return 0, false
	}
	return count, true
}
