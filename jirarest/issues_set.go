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

func (is *IssuesSet) StatusesOrder() []string {
	if is.Config != nil && is.Config.StatusConfig != nil {
		// is.Config.StatusesSet.DedupeMetaStageOrder()
		return is.Config.StatusConfig.StageConfig.Order()
	} else {
		return []string{}
	}
}

func (is *IssuesSet) AddIssuesFile(filename string) error {
	if ii, err := IssuesReadFileJSON(filename); err != nil {
		return err
	} else {
		return is.Add(ii...)
	}
}

func (is *IssuesSet) Add(issues ...jira.Issue) error {
	if is.IssuesMap == nil {
		is.IssuesMap = map[string]jira.Issue{}
	}
	for _, iss := range issues {
		if key := strings.TrimSpace(iss.Key); key == "" {
			return errors.New("no key")
		} else {
			is.IssuesMap[key] = iss
		}
	}
	return nil
}

func (is *IssuesSet) IssueFirst() (jira.Issue, error) {
	keys := is.Keys()
	if len(keys) == 0 {
		return jira.Issue{}, errors.New("no issues present")
	} else if iss, ok := is.IssuesMap[keys[0]]; ok {
		return iss, nil
	} else {
		panic(fmt.Sprintf("issue key from map not found (%s)", keys[0]))
	}
}

// KeyExists returns a boolean representing the existence of an issue key.
func (is *IssuesSet) KeyExists(key string, inclParents bool) bool {
	if _, ok := is.IssuesMap[key]; ok {
		return true
	} else if !inclParents || is.Parents == nil {
		return false
	} else {
		return is.Parents.KeyExists(key, inclParents)
	}
}

func (is *IssuesSet) Keys() []string              { return maputil.Keys(is.IssuesMap) }
func (is *IssuesSet) Len() uint                   { return uint(len(is.IssuesMap)) }
func (is *IssuesSet) LenParents() uint            { return uint(len(is.KeysParents())) }
func (is *IssuesSet) LenParentsPopulated() uint   { return uint(len(is.KeysParentsPopulated())) }
func (is *IssuesSet) LenParentsUnpopulated() uint { return uint(len(is.KeysParentsUnpopulated())) }

func (is *IssuesSet) LenLineageTopKeysPopulated() uint {
	if linPopIDs, err := is.LineageTopKeysPopulated(); err != nil {
		panic(err)
	} else {
		return uint(len(linPopIDs))
	}
}

func (is *IssuesSet) LenLineageTopKeysUnpopulated() uint {
	if linUnpopIDs, err := is.LineageTopKeysUnpopulated(); err != nil {
		panic(err)
	} else {
		return uint(len(linUnpopIDs))
	}
}

// LenMap provides various metrics. It is useful for determining if all parents and lineages have been loaded.
func (is *IssuesSet) LenMap() map[string]uint {
	lenParentsSet := 0
	if is.Parents != nil {
		lenParentsSet = len(is.Parents.IssuesMap)
	}
	return map[string]uint{
		"len":                       is.Len(),
		"lineageTopKeysPopulated":   is.LenLineageTopKeysPopulated(),
		"lineageTopKeysUnpopulated": is.LenLineageTopKeysUnpopulated(),
		"parents":                   is.LenParents(),
		"parentsPopulated":          is.LenParentsPopulated(),
		"parentsUnpopulated":        is.LenParentsUnpopulated(),
		"parentsSetAll":             uint(lenParentsSet),
	}
}

func (is *IssuesSet) FilterByStatus(inclStatuses, exclStatuses []string) (*IssuesSet, error) {
	filteredIssuesSet := NewIssuesSet(is.Config)
	inclStatusesMap := map[string]int{}
	for _, s := range inclStatuses {
		inclStatusesMap[s]++
	}
	exclStatusesMap := map[string]int{}
	for _, s := range exclStatuses {
		exclStatusesMap[s]++
	}
	for _, iss := range is.IssuesMap {
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

func (is *IssuesSet) EpicKeys(customFieldID string) []string {
	keys := []string{}
	for _, iss := range is.IssuesMap {
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

func (is *IssuesSet) Get(key string) (jira.Issue, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return jira.Issue{}, errors.New("key not provided")
	}
	if iss, ok := is.IssuesMap[key]; ok {
		return iss, nil
	} else if is.Parents != nil {
		if iss, ok := is.Parents.IssuesMap[key]; ok {
			return iss, nil
		}
	}
	return jira.Issue{}, errors.New("key not found")
}

func (is *IssuesSet) InflateEpicKeys(customFieldEpicLinkID string) {
	for k, iss := range is.IssuesMap {
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
		is.IssuesMap[k] = iss
	}
}

// InflateEpics uses the Jira REST API to inflate the Issue struct with an Epic struct.
func (is *IssuesSet) InflateEpics(jclient *jira.Client, customFieldIDEpicLink string) error {
	epicKeys := is.EpicKeys(customFieldIDEpicLink)
	newEpicKeys := []string{}
	for _, key := range epicKeys {
		if _, ok := is.IssuesMap[key]; !ok {
			newEpicKeys = append(newEpicKeys, key)
		}
	}
	epicsSet := NewEpicsSet()
	err := epicsSet.GetKeys(jclient, newEpicKeys)
	if err != nil {
		return err
	}

	for k, iss := range is.IssuesMap {
		issEpicKey := strings.TrimSpace(IssueFieldsCustomFieldString(iss.Fields, customFieldIDEpicLink))
		if issEpicKey == "" {
			continue
		}
		epic, ok := epicsSet.EpicsMap[issEpicKey]
		if !ok {
			panic("not found")
		}
		iss.Fields.Epic = &epic
		is.IssuesMap[k] = iss
	}
	return nil
}

func (is *IssuesSet) FilterStatus(inclStatuses ...string) (*IssuesSet, error) {
	n := NewIssuesSet(is.Config)
	if len(inclStatuses) == 0 {
		return n, nil
	}
	for _, iss := range is.IssuesMap {
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

func (is *IssuesSet) FilterType(inclTypes ...string) (*IssuesSet, error) {
	n := NewIssuesSet(is.Config)
	if len(inclTypes) == 0 {
		return n, nil
	}
	for _, iss := range is.IssuesMap {
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

func (is *IssuesSet) Issues() Issues {
	ii := Issues{}
	for _, iss := range is.IssuesMap {
		ii = append(ii, iss)
	}
	return ii
}

func (is *IssuesSet) IssueMetas(customFieldLabels []string) IssueMetas {
	var imetas IssueMetas
	for _, iss := range is.IssuesMap {
		iss := iss
		issMore := NewIssueMore(&iss)
		issMeta := issMore.Meta(is.Config.ServerURL, customFieldLabels)
		imetas = append(imetas, issMeta)
	}
	return imetas
}

func (is *IssuesSet) IssuesSetHighestType(issueType string) (*IssuesSet, error) {
	new := NewIssuesSet(is.Config)
	for _, iss := range is.IssuesMap {
		iss := iss
		issMore := NewIssueMore(&iss)
		issMeta := issMore.Meta(is.Config.ServerURL, []string{})
		issKey := strings.TrimSpace(issMeta.Key)
		if issKey != "" {
			lineage, err := is.Lineage(issKey, []string{})
			if err != nil {
				return nil, errorsutil.Wrapf(err, "error on `is.Lineage(%s)`", issKey)
			}
			if issMetaType := lineage.HighestType(issueType); issMetaType != nil && strings.TrimSpace(issMetaType.Key) != "" {
				if issType, err := is.Get(issMetaType.Key); err != nil {
					return nil, errorsutil.Wrapf(err, "error on `is.Get(%s)`", issMetaType.Key)
				} else {
					if err := new.Add(issType); err != nil {
						return nil, err
					}
				}
			}
		}
	}
	new.Parents = is.Parents
	return new, nil
}

func (is *IssuesSet) WriteFileJSON(name, prefix, indent string) error {
	j, err := jsonutil.MarshalSimple(is, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(name, j, 0600)
}
