package jiraxml

import (
	"time"

	"github.com/grokify/mogo/type/stringsutil"
)

type Items []Item

func (ii Items) FilterByStatus(statuses ...string) Items {
	new := Items{}
	mStatuses := map[string]int{}
	for _, s := range statuses {
		mStatuses[s]++
	}
	for _, ix := range ii {
		if _, ok := mStatuses[ix.Status.DisplayName]; ok {
			new = append(new, ix)
		}
	}
	return new
}

func (ii Items) Keys() []string {
	keys := []string{}
	for _, ix := range ii {
		keys = append(keys, ix.Key.DisplayName)
	}
	return stringsutil.SliceCondenseSpace(keys, true, true)
}

func (ii Items) Stats(workingHoursPerDay, workingDaysPerWeek float32) ItemsStats {
	if workingHoursPerDay == 0 {
		workingHoursPerDay = WorkingHoursPerDayDefault
	}
	if workingDaysPerWeek == 0 {
		workingDaysPerWeek = WorkingDaysPerWeekDefault
	}
	whpworkingHoursPerDay64 := float64(workingHoursPerDay)
	stats := ItemsStats{
		WorkingHoursPerDay:     workingHoursPerDay,
		WorkingDaysPerWeek:     workingDaysPerWeek,
		ItemCount:              len(ii),
		ItemCountByStatus:      map[string]int{},
		ItemCountByType:        map[string]int{},
		EstimateStatsByType:    map[string]EstimateStats{},
		ClosedEstimateVsActual: EstimateVsActual{},
	}
	for _, it := range ii {
		stats.TimeOriginalEstimate += it.TimeOriginalEstimate.Duration()
		stats.AggregateTimeSpent += it.AggregateTimeSpent.Duration()
		stats.ItemCountByStatus[it.Status.DisplayName]++
		stats.ItemCountByType[it.Type.DisplayName]++
		esStats, ok := stats.EstimateStatsByType[it.Type.DisplayName]
		if !ok {
			esStats = EstimateStats{}
		}
		if it.TimeOriginalEstimate.Seconds > 0 {
			esStats.WithEstimate++
		} else {
			esStats.WithoutEstimate++
		}
		stats.EstimateStatsByType[it.Type.DisplayName] = esStats
		if it.Status.DisplayName == StatusClosed {
			stats.ClosedEstimateVsActual.ClosedCount++
			if it.TimeOriginalEstimate.Seconds > 0 {
				stats.ClosedEstimateVsActual.ClosedCountWithEstimate++
				stats.ClosedEstimateVsActual.EstimateDays += it.TimeOriginalEstimate.Duration().Hours() / whpworkingHoursPerDay64
				stats.ClosedEstimateVsActual.ActualDays += it.AggregateTimeSpent.Duration().Hours() / whpworkingHoursPerDay64
			}
		}
	}
	stats.TimeOriginalEstimateDays = stats.TimeOriginalEstimate.Hours() / whpworkingHoursPerDay64
	stats.AggregateTimeSpentDays = stats.AggregateTimeSpent.Hours() / whpworkingHoursPerDay64
	stats.ClosedEstimateVsActual.Inflate()
	return stats
}

type ItemsStats struct {
	WorkingHoursPerDay       float32
	WorkingDaysPerWeek       float32
	ItemCount                int
	ItemCountByStatus        map[string]int
	ItemCountByType          map[string]int
	EstimateStatsByType      map[string]EstimateStats
	TimeOriginalEstimate     time.Duration
	TimeOriginalEstimateDays float64
	AggregateTimeSpent       time.Duration
	AggregateTimeSpentDays   float64
	ClosedEstimateVsActual   EstimateVsActual
}

type EstimateStats struct {
	WithEstimate    int
	WithoutEstimate int
}

type EstimateVsActual struct {
	ClosedCount             int
	ClosedCountWithEstimate int
	EstimateDays            float64
	ActualDays              float64
	EstimateRatio           float64
}

func (eva *EstimateVsActual) Inflate() {
	if eva.ActualDays > 0 {
		eva.EstimateRatio = eva.ActualDays / eva.EstimateDays
	}
}
