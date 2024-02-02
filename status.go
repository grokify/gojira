package gojira

import "github.com/grokify/mogo/type/stringsutil"

type StatusesSet struct {
	Map   map[string]string
	Order []string
}

func NewStatusesSet() StatusesSet {
	return StatusesSet{
		Map:   map[string]string{},
		Order: []string{},
	}
}

func (ss *StatusesSet) AddMapSlice(m map[string][]string) {
	for category, vals := range m {
		for _, status := range vals {
			ss.Add(status, category)
		}
	}
}

func (ss *StatusesSet) Add(status, statusCategory string) {
	if ss.Map == nil {
		ss.Map = map[string]string{}
	}
	ss.Map[status] = statusCategory
}

func (ss *StatusesSet) DedupeOrder() {
	if ss.Order == nil {
		ss.Order = []string{}
		return
	} else if len(ss.Order) == 0 || len(ss.Order) == 1 {
		return
	} else {
		ss.Order = stringsutil.SliceCondenseSpace(ss.Order, true, false)
	}
}

func (ss *StatusesSet) StatusCategory(status string) string {
	if cat, ok := ss.Map[status]; ok {
		return cat
	} else {
		return ""
	}
}

func (ss *StatusesSet) StatusesOpen() []string {
	return ss.StatusesForCategory(StatusOpen)
}

func (ss *StatusesSet) StatusesInProgress() []string {
	return ss.StatusesForCategory(StatusInProgress)
}

func (ss *StatusesSet) StatusesDone() []string { // not backlog
	return ss.StatusesForCategory(StatusDone)
}

func (ss *StatusesSet) StatusesForCategory(category string) []string {
	var statuses []string
	for k, v := range ss.Map {
		if v == category {
			statuses = append(statuses, k)
		}
	}
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

func (ss *StatusesSet) StatusesInProgressAndDone() []string { // not backlog
	var statuses []string
	statuses = append(statuses, ss.StatusesInProgress()...)
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
