package gojira

import (
	"errors"
	"net/url"

	"github.com/grokify/mogo/type/stringsutil"
)

type StatusConfig struct {
	Map         map[string]string
	StageConfig StageConfig
	// MetaStageOrder []string
}

func NewStatusConfig(stageConfig StageConfig) StatusConfig {
	return StatusConfig{
		Map:         map[string]string{}, // status to metastatus
		StageConfig: stageConfig,
	}
}

// AddMapSlice should be a map where the keys are meta statuses and the values are slices of Jira statuses.
func (ss *StatusConfig) AddMapSlice(m map[string][]string) error {
	for metaStatus, vals := range m {
		for _, status := range vals {
			if err := ss.Add(status, metaStatus); err != nil {
				return err
			}
		}
	}
	return nil
}

func (ss *StatusConfig) Add(status, metaStage string) error {
	if ss.Map == nil {
		ss.Map = map[string]string{}
	}
	if !ss.StageConfig.Exists(metaStage) {
		return errors.New("metaStage is not configured")
	}
	ss.Map[status] = metaStage
	return nil
}

func (ss *StatusConfig) MapMetaStageToStatuses() map[string][]string {
	out := url.Values{}
	for status, metaStage := range ss.Map {
		out.Add(metaStage, status)
	}
	return out
}

/*
func (ss *StatusConfig) DedupeMetaStageOrder() {
	if ss.MetaStageOrder == nil {
		ss.MetaStageOrder = []string{}
		return
	} else if len(ss.MetaStageOrder) == 0 || len(ss.MetaStageOrder) == 1 {
		return
	} else {
		ss.MetaStageOrder = stringsutil.SliceCondenseSpace(ss.MetaStageOrder, true, false)
	}
}
*/

// MetaStage returns the metastatus for a status. If there is no metastatus, an empty string is returned.
func (ss *StatusConfig) MetaStage(status string) string {
	if cat, ok := ss.Map[status]; ok {
		return cat
	} else {
		return ""
	}
}

/*
// MetaStageOrderMap returns a `map[string]uint` where the key is the meta status and the value is the index.
func (ss *StatusConfig) MetaStageOrderMap() map[string]uint {
	out := map[string]uint{}
	for i, ms := range ss.MetaStageOrder {
		out[ms] = uint(i)
	}
	return out
}
*/

func (ss *StatusConfig) StatusesReadyForPlanning() []string {
	if metaStageName := ss.StageConfig.ReadyforPlanningName(); metaStageName == "" {
		return []string{}
	} else {
		return ss.StatusesForMetaStage(metaStageName)
	}
}

func (ss *StatusConfig) StatusesInDevelopment() []string {
	if metaStageName := ss.StageConfig.InDevelopmentName(); metaStageName == "" {
		return []string{}
	} else {
		return ss.StatusesForMetaStage(metaStageName)
	}
}

func (ss *StatusConfig) StatusesDone() []string { // not backlog
	if metaStageName := ss.StageConfig.DoneName(); metaStageName == "" {
		return []string{}
	} else {
		return ss.StatusesForMetaStage(metaStageName)
	}
}

func (ss *StatusConfig) StatusesForMetaStage(metaStatus string) []string {
	var statuses []string
	for k, v := range ss.Map {
		if v == metaStatus {
			statuses = append(statuses, k)
		}
	}
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

func (ss *StatusConfig) StatusesInDevelopmentAndDone() []string { // not backlog
	var statuses []string
	statuses = append(statuses, ss.StatusesInDevelopment()...)
	statuses = append(statuses, ss.StatusesDone()...)
	return stringsutil.SliceCondenseSpace(statuses, true, true)
}

/*
func DefaultStatusesMapSlice() map[string][]string {
	return map[string][]string{
		StatusOpen:       {StatusOpen},
		StatusInProgress: {StatusInProgress},
		StatusDone:       {StatusDone},
	}
}
*/
