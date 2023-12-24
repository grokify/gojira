package jirarest

import (
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
)

func (is *IssuesSet) Counts() map[string]map[string]uint {
	mm := map[string]map[string]uint{
		"byProject":    is.CountsByProject(),
		"byProjectKey": is.CountsByProjectKey(),
		"byStatus":     is.CountsByStatus(),
		"byType":       is.CountsByType(true, false),
		"byTime":       is.CountsByTime(),
	}
	return mm
}

// CountsByCustomFieldValues returns a list of custom field value counts where `customField` is in
// the format `customfield_12345`.
func (is *IssuesSet) CountsByCustomFieldValues(customField string) (map[string]uint, error) {
	out := map[string]uint{}
	for _, iss := range is.IssuesMap {
		iss := iss
		im := IssueMore{Issue: &iss}
		cfInfo, err := im.CustomField(customField)
		if err != nil {
			return out, err
		}
		out[cfInfo.Value]++

	}
	return out, nil
}

func (is *IssuesSet) CountsByProject() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: pointer.Pointer(iss)}
		m[im.Project()]++
	}
	return m
}

func (is *IssuesSet) CountsByProjectKey() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: pointer.Pointer(iss)}
		m[im.ProjectKey()]++
	}
	return m
}

func (is *IssuesSet) CountsByStatus() map[string]uint {
	m := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := IssueMore{Issue: pointer.Pointer(iss)}
		//ifs := IssueFieldsSimple{Fields: iss.Fields}
		m[im.Status()]++
	}
	return m
}

func (is *IssuesSet) CountsByType(inclLeafs, inclParents bool) map[string]uint {
	m := map[string]uint{}
	if inclLeafs {
		for _, iss := range is.IssuesMap {
			iss := iss
			im := IssueMore{Issue: &iss}
			m[im.Type()]++
		}
	}
	if inclParents && is.Parents != nil {
		for _, iss := range is.Parents.IssuesMap {
			iss := iss
			im := IssueMore{Issue: &iss}
			m[im.Type()]++
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
