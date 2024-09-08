package jirarest

import (
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/timeutil"
	"github.com/grokify/mogo/type/stringsutil"
)

func (set *IssuesSet) Counts() map[string]map[string]uint {
	mm := map[string]map[string]uint{
		"byProject":    set.CountsByProject(),
		"byProjectKey": set.CountsByProjectKey(),
		"byStatus":     set.CountsByStatus(),
		"byType":       set.CountsByType(true, false),
		"byTime":       set.CountsByTime(),
	}
	return mm
}

/*
func (is *IssuesSet) TimeSeriesCreated() (timeseries.TimeSeries, error) {
	ts := timeseries.NewTimeSeries("")
	for _, iss := range is.IssuesMap {
		iss := iss
		im := IssueMore{Issue: &iss}
		ts.AddInt64(im.CreateTime().UTC(), 1)
	}
	return ts, nil
}
*/

// CountsByCustomFieldValues returns a list of custom field value counts where `customField` is in
// the format `customfield_12345`.
func (set *IssuesSet) CountsByCustomFieldValues(customField string) (map[string]uint, error) {
	out := map[string]uint{}
	for _, iss := range set.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		cfInfo, err := im.CustomField(customField)
		if err != nil {
			return out, err
		}
		out[cfInfo.Value]++
	}
	return out, nil
}

func (set *IssuesSet) CountsByProject() map[string]uint {
	m := map[string]uint{}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		m[im.Project()]++
	}
	return m
}

func (set *IssuesSet) CountsByProjectKey() map[string]uint {
	m := map[string]uint{}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		m[im.ProjectKey()]++
	}
	return m
}

func (set *IssuesSet) CountsByStatus() map[string]uint {
	m := map[string]uint{}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		//ifs := IssueFieldsSimple{Fields: iss.Fields}
		m[im.Status()]++
	}
	return m
}

func (set *IssuesSet) CountsByMetaStage(inclTypeFilter []string) map[string]uint {
	inclTypeFilter = stringsutil.SliceCondenseSpace(inclTypeFilter, true, true)
	inclTypeFilterMap := map[string]int{}
	for _, filter := range inclTypeFilter {
		inclTypeFilterMap[filter]++
	}
	out := map[string]uint{}
	count := uint(0)
	unknownStatus := map[string]uint{}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if len(inclTypeFilterMap) > 0 {
			if _, ok := inclTypeFilterMap[im.Type()]; !ok {
				continue
			}
		}
		metaStage := ""
		if set.Config != nil && set.Config.StatusConfig != nil {
			metaStage = set.Config.StatusConfig.MetaStage(im.Status())
		}
		if metaStage == "" {
			unknownStatus[im.Status()]++
		}
		out[metaStage]++
		count++
	}
	if msuCount(out) != count {
		panic("count mismatch")
	}
	return out
}

func (is *IssuesSet) CountsByProjectAndMetaStage(inclTypeFilter []string) *histogram.HistogramSet {
	out := histogram.NewHistogramSet("")
	inclTypeFilter = stringsutil.SliceCondenseSpace(inclTypeFilter, true, true)
	inclTypeFilterMap := map[string]int{}
	for _, filter := range inclTypeFilter {
		inclTypeFilterMap[filter]++
	}
	count := int(0)
	for _, iss := range is.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if len(inclTypeFilterMap) > 0 {
			if _, ok := inclTypeFilterMap[im.Type()]; !ok {
				continue
			}
		}
		projectName := im.Project()
		metaStage := ""
		if is.Config != nil && is.Config.StatusConfig != nil {
			metaStage = is.Config.StatusConfig.MetaStage(im.Status())
		}
		out.Add(projectName, metaStage, 1)
		count++
	}
	if count != out.Sum() {
		panic("count mismatch")
	}
	return out
}

func msuCount(m map[string]uint) uint {
	c := uint(0)
	for _, v := range m {
		c += v
	}
	return c
}

func (is *IssuesSet) CountWithTypeFilter(inclTypeFilter []string) uint {
	inclTypeFilter = stringsutil.SliceCondenseSpace(inclTypeFilter, true, true)
	inclTypeFilterMap := map[string]int{}
	for _, filter := range inclTypeFilter {
		inclTypeFilterMap[filter]++
	}
	count := uint(0)
	for _, iss := range is.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if len(inclTypeFilterMap) > 0 {
			if _, ok := inclTypeFilterMap[im.Type()]; !ok {
				continue
			}
		}
		count++
	}
	return count
}

func (is *IssuesSet) CountsByType(inclLeafs, inclParents bool) map[string]uint {
	m := map[string]uint{}
	if inclLeafs {
		for _, iss := range is.IssuesMap {
			iss := iss
			im := NewIssueMore(&iss)
			m[im.Type()]++
		}
	}
	if inclParents && is.Parents != nil {
		for _, iss := range is.Parents.IssuesMap {
			iss := iss
			im := NewIssueMore(&iss)
			m[im.Type()]++
		}
	}
	return m
}

func (is *IssuesSet) CountsByTime() map[string]uint {
	out := map[string]uint{}
	for _, iss := range is.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		if iss.Fields.TimeEstimate <= 0 {
			out["TimeEstimateLTEZ"]++
		} else {
			out["TimeEstimateGTZ"]++
		}
		if iss.Fields.TimeOriginalEstimate <= 0 {
			out["TimeOriginalEstimateLTEZ"]++
		} else {
			out["TimeOriginalEstimateGTZ"]++
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
	return out
}

func (is *IssuesSet) CountsByWorkstream(wsFuncMake WorkstreamFuncMake, inclTypeFilter []string) (map[string]uint, error) {
	inclTypeFilter = stringsutil.SliceCondenseSpace(inclTypeFilter, true, true)
	inclTypeFilterMap := map[string]int{}
	for _, filter := range inclTypeFilter {
		inclTypeFilterMap[filter]++
	}
	out := map[string]uint{}
	for _, iss := range is.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if len(inclTypeFilterMap) > 0 {
			if _, ok := inclTypeFilterMap[im.Type()]; !ok {
				continue
			}
		}
		if ws, err := wsFuncMake(im.Key()); err != nil {
			return nil, err
		} else {
			out[ws]++
		}
	}
	return out, nil
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
