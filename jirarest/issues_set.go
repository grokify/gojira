package jirarest

import (
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"golang.org/x/exp/slices"
)

type IssuesSet struct {
	Config    *gojira.Config
	IssuesMap map[string]jira.Issue
	Parents   *IssuesSet
}

func NewIssuesSet(cfg *gojira.Config) *IssuesSet {
	if cfg == nil {
		cfg = gojira.NewConfigDefault()
	}
	return &IssuesSet{
		Config:    cfg,
		IssuesMap: map[string]jira.Issue{},
		Parents: &IssuesSet{
			Config:    cfg,
			IssuesMap: map[string]jira.Issue{},
		},
	}
}

func (set *IssuesSet) StatusesOrder() []string {
	if set.Config != nil && set.Config.StatusConfig != nil {
		// is.Config.StatusesSet.DedupeMetaStageOrder()
		return set.Config.StatusConfig.StageConfig.Order()
	} else {
		return []string{}
	}
}

func (set *IssuesSet) AddIssuesFile(filename string) error {
	if ii, err := IssuesReadFileJSON(filename); err != nil {
		return err
	} else {
		return set.Add(ii...)
	}
}

func (set *IssuesSet) Add(issues ...jira.Issue) error {
	if set.IssuesMap == nil {
		set.IssuesMap = map[string]jira.Issue{}
	}
	for _, iss := range issues {
		if key := strings.TrimSpace(iss.Key); key == "" {
			return errors.New("no key")
		} else {
			set.IssuesMap[key] = iss
		}
	}
	return nil
}

func (set *IssuesSet) IssueFirst() (jira.Issue, error) {
	keys := set.Keys()
	if len(keys) == 0 {
		return jira.Issue{}, errors.New("no issues present")
	} else if iss, ok := set.IssuesMap[keys[0]]; ok {
		return iss, nil
	} else {
		panic(fmt.Sprintf("issue key from map not found (%s)", keys[0]))
	}
}

// KeyExists returns a boolean representing the existence of an issue key.
func (set *IssuesSet) KeyExists(key string, inclParents bool) bool {
	if _, ok := set.IssuesMap[key]; ok {
		return true
	} else if !inclParents || set.Parents == nil {
		return false
	} else {
		return set.Parents.KeyExists(key, inclParents)
	}
}

func (set *IssuesSet) Keys() []string              { return maputil.Keys(set.IssuesMap) }
func (set *IssuesSet) Len() uint                   { return uint(len(set.IssuesMap)) }
func (set *IssuesSet) LenParents() uint            { return uint(len(set.KeysParents())) }
func (set *IssuesSet) LenParentsPopulated() uint   { return uint(len(set.KeysParentsPopulated())) }
func (set *IssuesSet) LenParentsUnpopulated() uint { return uint(len(set.KeysParentsUnpopulated())) }

func (set *IssuesSet) LenLineageTopKeysPopulated() uint {
	if linPopIDs, err := set.LineageTopKeysPopulated(); err != nil {
		panic(err)
	} else {
		return uint(len(linPopIDs))
	}
}

func (set *IssuesSet) LenLineageTopKeysUnpopulated() uint {
	if linUnpopIDs, err := set.LineageTopKeysUnpopulated(); err != nil {
		panic(err)
	} else {
		return uint(len(linUnpopIDs))
	}
}

// LenMap provides various metrics. It is useful for determining if all parents and lineages have been loaded.
func (set *IssuesSet) LenMap() map[string]uint {
	lenParentsSet := 0
	if set.Parents != nil {
		lenParentsSet = len(set.Parents.IssuesMap)
	}
	return map[string]uint{
		"len":                       set.Len(),
		"lineageTopKeysPopulated":   set.LenLineageTopKeysPopulated(),
		"lineageTopKeysUnpopulated": set.LenLineageTopKeysUnpopulated(),
		"parents":                   set.LenParents(),
		"parentsPopulated":          set.LenParentsPopulated(),
		"parentsUnpopulated":        set.LenParentsUnpopulated(),
		"parentsSetAll":             uint(lenParentsSet),
	}
}

func (set *IssuesSet) FilterByStatus(inclStatuses, exclStatuses []string) (*IssuesSet, error) {
	filteredIssuesSet := NewIssuesSet(set.Config)
	inclStatusesMap := map[string]int{}
	for _, s := range inclStatuses {
		inclStatusesMap[s]++
	}
	exclStatusesMap := map[string]int{}
	for _, s := range exclStatuses {
		exclStatusesMap[s]++
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		// ifs := IssueFieldsSimple{Fields: iss.Fields}
		statusName := im.Status()
		_, inclStatusOk := inclStatusesMap[statusName]
		_, exclStatusOk := exclStatusesMap[statusName]
		if len(inclStatusesMap) > 0 && !inclStatusOk {
			continue
		} else if len(exclStatuses) > 0 && exclStatusOk {
			continue
		}
		err := filteredIssuesSet.Add(iss)
		if err != nil {
			return nil, err
		}
	}
	return filteredIssuesSet, nil
}

func (set *IssuesSet) EpicKeys(customFieldID string) []string {
	keys := []string{}
	for _, iss := range set.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		if iss.Fields.Epic != nil {
			keys = append(keys, iss.Fields.Epic.Key)
		}
		epickey := IssueFieldsCustomFieldString(iss.Fields, customFieldID)
		if epickey != "" {
			keys = append(keys, epickey)
		}
	}
	keys = slicesutil.Dedupe(keys)
	sort.Strings(keys)
	return keys
}

func (set *IssuesSet) Get(key string) (jira.Issue, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return jira.Issue{}, errors.New("key not provided")
	}
	if iss, ok := set.IssuesMap[key]; ok {
		return iss, nil
	} else if set.Parents != nil {
		if iss, ok := set.Parents.IssuesMap[key]; ok {
			return iss, nil
		}
	}
	return jira.Issue{}, errors.New("key not found")
}

func (set *IssuesSet) InflateEpicKeys(customFieldEpicLinkID string) {
	for k, iss := range set.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		if iss.Fields.Epic != nil && strings.TrimSpace(iss.Fields.Epic.Key) != "" {
			continue
		}
		epicKey := IssueFieldsCustomFieldString(iss.Fields, customFieldEpicLinkID)
		if epicKey != "" {
			if iss.Fields.Epic == nil {
				iss.Fields.Epic = &jira.Epic{}
			}
			iss.Fields.Epic.Key = epicKey
		}
		set.IssuesMap[k] = iss
	}
}

// InflateEpics uses the Jira REST API to inflate the Issue struct with an Epic struct.
func (set *IssuesSet) InflateEpics(jclient *jira.Client, customFieldIDEpicLink string) error {
	epicKeys := set.EpicKeys(customFieldIDEpicLink)
	newEpicKeys := []string{}
	for _, key := range epicKeys {
		if _, ok := set.IssuesMap[key]; !ok {
			newEpicKeys = append(newEpicKeys, key)
		}
	}
	epicsSet := NewEpicsSet()
	err := epicsSet.GetKeys(jclient, newEpicKeys)
	if err != nil {
		return err
	}

	for k, iss := range set.IssuesMap {
		issEpicKey := strings.TrimSpace(IssueFieldsCustomFieldString(iss.Fields, customFieldIDEpicLink))
		if issEpicKey == "" {
			continue
		}
		epic, ok := epicsSet.EpicsMap[issEpicKey]
		if !ok {
			panic("not found")
		}
		iss.Fields.Epic = &epic
		set.IssuesMap[k] = iss
	}
	return nil
}

func (set *IssuesSet) FilterStatus(inclStatuses ...string) (*IssuesSet, error) {
	n := NewIssuesSet(set.Config)
	if len(inclStatuses) == 0 {
		return n, nil
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if slices.Index(inclStatuses, im.Status()) >= 0 {
			err := n.Add(iss)
			if err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}

func (set *IssuesSet) FilterType(inclTypes ...string) (*IssuesSet, error) {
	n := NewIssuesSet(set.Config)
	if len(inclTypes) == 0 {
		return n, nil
	}
	for _, iss := range set.IssuesMap {
		im := NewIssueMore(pointer.Pointer(iss))
		if slices.Index(inclTypes, im.Type()) >= 0 {
			err := n.Add(iss)
			if err != nil {
				return nil, err
			}
		}
	}
	return n, nil
}

func (set *IssuesSet) Issues() Issues {
	ii := Issues{}
	for _, iss := range set.IssuesMap {
		ii = append(ii, iss)
	}
	return ii
}

func (set *IssuesSet) IssueMetas(customFieldLabels []string) IssueMetas {
	var imetas IssueMetas
	for _, iss := range set.IssuesMap {
		iss := iss
		issMore := NewIssueMore(&iss)
		issMeta := issMore.Meta(set.Config.ServerURL, customFieldLabels)
		imetas = append(imetas, issMeta)
	}
	return imetas
}

func (set *IssuesSet) IssuesSetHighestType(issueType string) (*IssuesSet, error) {
	new := NewIssuesSet(set.Config)
	for _, iss := range set.IssuesMap {
		iss := iss
		issMore := NewIssueMore(&iss)
		issMeta := issMore.Meta(set.Config.ServerURL, []string{})
		issKey := strings.TrimSpace(issMeta.Key)
		if issKey != "" {
			lineage, err := set.Lineage(issKey, []string{})
			if err != nil {
				return nil, errorsutil.Wrapf(err, "error on `is.Lineage(%s)`", issKey)
			}
			if issMetaType := lineage.HighestType(issueType); issMetaType != nil && strings.TrimSpace(issMetaType.Key) != "" {
				if issType, err := set.Get(issMetaType.Key); err != nil {
					return nil, errorsutil.Wrapf(err, "error on `is.Get(%s)`", issMetaType.Key)
				} else {
					if err := new.Add(issType); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	new.Parents = set.Parents
	return new, nil
}

func (set *IssuesSet) WriteFileJSON(name, prefix, indent string) error {
	j, err := jsonutil.MarshalSimple(set, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(name, j, 0600)
}
