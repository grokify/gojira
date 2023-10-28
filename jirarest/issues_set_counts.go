package jirarest

import (
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/time/timeutil"
)

func (is *IssuesSet) Counts() map[string]map[string]uint {
	mm := map[string]map[string]uint{
		"byProject":    is.CountsByProject(),
		"byProjectKey": is.CountsByProjectKey(),
		"byStatus":     is.CountsByStatus(),
		"byType":       is.CountsByType(),
		"byTime":       is.CountsByTime(),
	}
	return mm
}

func (is *IssuesSet) CountsByProject() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: &iss}
		m[im.Project()]++
	}
	return m
}

func (is *IssuesSet) CountsByProjectKey() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: &iss}
		m[im.ProjectKey()]++
	}
	return m
}

func (is *IssuesSet) CountsByStatus() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: &iss}
		//ifs := IssueFieldsSimple{Fields: iss.Fields}
		m[im.Status()]++
	}
	return m
}

func (is *IssuesSet) CountsByType() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		if iss.Fields != nil {
			m[iss.Fields.Type.Name]++
		}
	}
	return m
}

func (is *IssuesSet) CountsByTime() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		if iss.Fields.TimeEstimate <= 0 {
			m["TimeEstimateLTEZ"]++
		} else {
			m["TimeEstimateGTZ"]++
		}
		if iss.Fields.TimeOriginalEstimate <= 0 {
			m["TimeOriginalEstimateLTEZ"]++
		} else {
			m["TimeOriginalEstimateGTZ"]++
		}
		/*
					TimeTimeSpent                     = "Time Spent"
			TimeTimeEstimate                  = "Time Estimate"
			TimeTimeOriginalEstimate          = "Time Original Estimate"
			TimeAggregateTimeOriginalEstimate = "Aggregate Time Original Estimate"
			TimeAggregateTimeSpent            = "Aggregate Time Spent"
			TimeAggregateTimeEstimate         = "Aggregate Time Estimate"
			TimeTimeRemaining                 = "Time Remaining"
			TimeTimeRemainingOriginal         = "Time Remaining Original"
		*/
	}
	return m
}

func (is *IssuesSet) TimeStats() gojira.TimeStats {
	if is.Config == nil {
		is.Config = gojira.NewConfigDefault()
	}
	ts := gojira.TimeStats{
		TimeUnit:           timeutil.SecondString,
		ItemCount:          len(is.IssuesMap),
		WorkingDaysPerWeek: is.Config.WorkingDaysPerWeek,
		WorkingHoursPerDay: is.Config.WorkingHoursPerDay,
	}
	for _, iss := range is.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		timeRemainingOriginal, timeRemaining := gojira.TimeRemaining(iss.Fields.Status.Name, iss.Fields.TimeOriginalEstimate, iss.Fields.TimeEstimate, iss.Fields.TimeSpent)
		ts.TimeSpent += float32(iss.Fields.TimeSpent)
		ts.TimeEstimate += float32(iss.Fields.TimeEstimate)
		ts.TimeOriginalEstimate += float32(iss.Fields.TimeOriginalEstimate)
		ts.AggregateTimeOriginalEstimate += float32(iss.Fields.AggregateTimeOriginalEstimate)
		ts.AggregateTimeSpent += float32(iss.Fields.AggregateTimeSpent)
		ts.AggregateTimeEstimate += float32(iss.Fields.AggregateTimeEstimate)
		ts.TimeRemaining += float32(timeRemaining)
		ts.TimeRemainingOriginal += float32(timeRemainingOriginal)
	}
	return ts
}
