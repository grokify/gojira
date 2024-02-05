package gojira

import "github.com/grokify/mogo/type/stringsutil"

type StatusesSet struct {
	Map            map[string]string
	MetaStageOrder []string
}

func NewStatusesSet() StatusesSet {
	return StatusesSet{
		Map:            map[string]string{}, // status to metastatus
		MetaStageOrder: MetaStageOrder(),
	}
}

// AddMapSlice should be a map where the keys are meta statuses and the values are slices of Jira statuses.
func (ss *StatusesSet) AddMapSlice(m map[string][]string) {
	for metaStatus, vals := range m {
		for _, status := range vals {
			ss.Add(status, metaStatus)
		}
	}
}

func (ss *StatusesSet) Add(status, metaStatus string) {
	if ss.Map == nil {
		ss.Map = map[string]string{}
	}
	ss.Map[status] = metaStatus
}

func (ss *StatusesSet) DedupeMetaStageOrder() {
	if ss.MetaStageOrder == nil {
		ss.MetaStageOrder = []string{}
		return
	} else if len(ss.MetaStageOrder) == 0 || len(ss.MetaStageOrder) == 1 {
		return
	} else {
		ss.MetaStageOrder = stringsutil.SliceCondenseSpace(ss.MetaStageOrder, true, false)
	}
}

// MetaStage returns the metastatus for a status. If there is no metastatus, an empty string is returned.
func (ss *StatusesSet) MetaStage(status string) string {
	if cat, ok := ss.Map[status]; ok {
		return cat
	} else {
		return ""
	}
}

// MetaStageOrderMap returns a `map[string]uint` where the key is the meta status and the value is the index.
func (ss *StatusesSet) MetaStageOrderMap() map[string]uint {
	out := map[string]uint{}
	for i, ms := range ss.MetaStageOrder {
		out[ms] = uint(i)
	}
	return out
}

func (ss *StatusesSet) StatusesReadyForPlanning() []string {
	return ss.StatusesForMetaStage(MetaStageReadyForPlanning)
}

func (ss *StatusesSet) StatusesInDevelopment() []string {
	return ss.StatusesForMetaStage(MetaStageInDevelopment)
}

func (ss *StatusesSet) StatusesDone() []string { // not backlog
	return ss.StatusesForMetaStage(StatusDone)
}

func (ss *StatusesSet) StatusesForMetaStage(metaStatus string) []string {
	var statuses []string
	for k, v := range ss.Map {
		if v == metaStatus {
			statuses = append(statuses, k)
		}
	}
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

func (ss *StatusesSet) StatusesInDevelopmentAndDone() []string { // not backlog
	var statuses []string
	statuses = append(statuses, ss.StatusesInDevelopment()...)
	statuses = append(statuses, ss.StatusesDone()...)
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

func DefaultStatusesMapSlice() map[string][]string {
	return map[string][]string{
		StatusOpen:       {StatusOpen},
		StatusInProgress: {StatusInProgress},
		StatusDone:       {StatusDone},
	}
}
