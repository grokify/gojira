package gojira

import (
	jira "github.com/andygrunwald/go-jira"
)

type TimeStatsSets struct {
	Map map[string]TimeStatsSet
}

func (tss *TimeStatsSets) AddIssue(iss jira.Issue) {
	if iss.Fields == nil {
		return
	}
}

type TimeStatsSet struct {
	Map map[string]TimeStats
}

type TimeStats struct {
	WorkingHoursPerDay            float32
	WorkingDaysPerWeek            float32
	ItemCount                     int
	TimeSpent                     int
	TimeEstimate                  int
	TimeOriginalEstimate          int
	AggregateTimeOriginalEstimate int
	AggregateTimeSpent            int
	AggregateTimeEstimate         int
	TimeRemaining                 int
	TimeRemainingOriginal         int
	/*
		ItemCountByStatus  map[string]int
		ItemCountByType    map[string]int
		//EstimateStatsByType      map[string]EstimateStats
		TimeOriginalEstimate     time.Duration
		TimeOriginalEstimateDays float64
		AggregateTimeSpent       time.Duration
		AggregateTimeSpentDays   float64
		//ClosedEstimateVsActual   EstimateVsActual
	*/
}
