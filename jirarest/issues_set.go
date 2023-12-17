package jirarest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/urlutil"
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

func (is *IssuesSet) Keys() []string {
	return maputil.Keys(is.IssuesMap)
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
		im := IssueMore{Issue: pointer.Pointer(iss)}
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
		im := IssueMore{Issue: pointer.Pointer(iss)}
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
		im := IssueMore{Issue: pointer.Pointer(iss)}
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

func (is *IssuesSet) IssueMetas() IssueMetas {
	var imetas IssueMetas
	for _, iss := range is.IssuesMap {
		iss := iss
		issMore := IssueMore{Issue: &iss}
		issMeta := issMore.Meta(is.Config.ServerURL)
		imetas = append(imetas, issMeta)
	}
	return imetas
}

// HistogramSets returns a list of histograms by Project and Type.
func (is *IssuesSet) HistogramSetProjectType() *histogram.HistogramSet {
	hset := histogram.NewHistogramSet("issues")

	serverURL := ""
	if is.Config != nil {
		serverURL = is.Config.ServerURL
	}

	for _, iss := range is.IssuesMap {
		iss := iss
		issMore := IssueMore{Issue: &iss}
		issMeta := issMore.Meta(serverURL)
		projKey := issMeta.ProjectKey
		issType := issMeta.Type
		hset.Add(projKey, issType, 1)
	}

	return hset
}

type CustomTableCols struct {
	Cols []CustomCol
}

type CustomCol struct {
	Name       string
	Type       string
	Func       func(iss jira.Issue) (string, error)
	RenderSkip bool
}

func DefaultIssuesSetTableColumns(inclInitiative, inclEpic bool) *table.ColumnDefinitions {
	var defs []table.ColumnDefinition
	if inclInitiative {
		initiativeCols := []table.ColumnDefinition{
			{Name: "Initiative Key", Format: table.FormatURL},
			{Name: "Initiative Name"}}
		defs = append(defs, initiativeCols...)
	}
	if inclEpic {
		epicCols := []table.ColumnDefinition{
			{Name: "Epic Key", Format: table.FormatURL},
			{Name: "Epic Name"}}
		defs = append(defs, epicCols...)
	}
	stdCols := []table.ColumnDefinition{
		{Name: "Issue Key", Format: table.FormatURL},
		{Name: "Issue Type"},
		{Name: "Project"},
		{Name: "Summary"},
		{Name: "Status"},
		{Name: "Resolution"},
		// {Name: "Aggregate Original Time Estimate Seconds", Format: table.FormatInt},
		// {Name: "Original Estimate Seconds", Format: table.FormatInt},
		{Name: "Original Estimate Days", Format: table.FormatFloat},
		{Name: "Estimate Days", Format: table.FormatFloat},
		{Name: "Time Spent Days", Format: table.FormatFloat},
		{Name: "Time Remaining Days", Format: table.FormatFloat},
	}
	defs = append(defs, stdCols...)

	// defs = append(defs, {Name: "Epic Key", Format: table.FormatURL},

	return &table.ColumnDefinitions{
		Definitions: defs,
	}
}

func BuildJiraIssueURL(baseURL, issueKey string) string {
	issueKey = strings.TrimSpace(issueKey)
	return urlutil.JoinAbsolute(baseURL, "/browse/", issueKey)
}

func (is *IssuesSet) IssuesSetHighestType(issueType string) (*IssuesSet, error) {
	new := NewIssuesSet(is.Config)
	for _, iss := range is.IssuesMap {
		iss := iss
		issMore := IssueMore{Issue: &iss}
		issMeta := issMore.Meta(is.Config.ServerURL)
		issKey := strings.TrimSpace(issMeta.Key)
		if issKey != "" {
			lineage, err := is.Lineage(issKey)
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

// TableSet is designed to return a `table.TableSet` where the tables include a list of issues and optionally, epics, and/or initiatives.
func (is *IssuesSet) TableSet(customCols *CustomTableCols, inclEpic bool, initiativeType string) (*table.TableSet, error) {
	ts := table.NewTableSet("Jira Issues")
	tbl1Issues, err := is.Table(customCols, inclEpic, initiativeType)
	if err != nil {
		return nil, err
	}
	tbl1Issues.Name = gojira.TypeIssue
	ts.TableMap[tbl1Issues.Name] = tbl1Issues
	ts.Order = append(ts.Order, tbl1Issues.Name)
	if inclEpic {
		isEpic, err := is.IssuesSetHighestType(gojira.TypeEpic)
		if err != nil {
			return nil, errorsutil.Wrapf(err, "error on `is.IssuesSetHighestType(%s)`", gojira.TypeEpic)
		}
		tbl2Epics, err := isEpic.Table(customCols, false, initiativeType)
		if err != nil {
			return nil, err
		}
		tbl2Epics.Name = gojira.TypeEpic
		ts.TableMap[tbl2Epics.Name] = tbl2Epics
		ts.Order = append(ts.Order, tbl2Epics.Name)
	}

	if initiativeType != "" {
		isInit, err := is.IssuesSetHighestType(initiativeType)
		if err != nil {
			return nil, err
		}
		tbl3Initiatives, err := isInit.Table(customCols, false, "")
		if err != nil {
			return nil, err
		}
		tbl3Initiatives.Name = initiativeType
		ts.TableMap[tbl3Initiatives.Name] = tbl3Initiatives
		ts.Order = append(ts.Order, tbl3Initiatives.Name)
	}
	return ts, nil
}

func (is *IssuesSet) Table(customCols *CustomTableCols, inclEpic bool, initiativeType string) (*table.Table, error) {
	if is.Config == nil {
		is.Config = gojira.NewConfigDefault()
	}
	initiativeType = strings.TrimSpace(initiativeType)
	inclInitiative := false
	if initiativeType != "" {
		inclInitiative = true
	}
	baseURL := strings.TrimSpace(is.Config.ServerURL)

	tbl := table.NewTable("issues")

	tbl.LoadColumnDefinitions(DefaultIssuesSetTableColumns(inclInitiative, inclEpic))

	if customCols != nil {
		lenCols := len(tbl.Columns)
		for i, customCol := range customCols.Cols {
			if customCol.RenderSkip {
				continue
			}
			j := lenCols + i
			customCol.Type = strings.TrimSpace(customCol.Type)
			if customCol.Type != "" {
				tbl.FormatMap[j] = customCol.Type
			}
			if name := strings.TrimSpace(customCol.Name); name != "" {
				tbl.Columns = append(tbl.Columns, name)
			} else {
				tbl.Columns = append(tbl.Columns, fmt.Sprintf("Column %d", j+1))
			}
		}
	}

	for key, iss := range is.IssuesMap {
		issMore := IssueMore{Issue: pointer.Pointer(iss)}
		issMeta := issMore.Meta(baseURL)

		lineage, err := is.Lineage(key)
		if err != nil {
			return nil, err
		}

		timeRemainingSecs := iss.Fields.TimeEstimate - iss.Fields.TimeSpent
		if timeRemainingSecs < 0 ||
			issMore.Status() == gojira.StatusClosed ||
			issMore.Status() == gojira.StatusDone {
			timeRemainingSecs = 0
		}

		row := []string{}

		if inclInitiative {
			initKeyDispplay := ""
			initName := ""
			if initiative := lineage.HighestType(initiativeType); initiative != nil {
				initiative.BuildKeyURL(baseURL) // should not be needed.
				initKeyDispplay = initiative.KeyLinkMarkdown()
				initName = initiative.Summary
			}
			row = append(row, initKeyDispplay, initName)
		}

		if inclEpic {
			epicKeyDisplay := issMore.EpicKey()
			epicName := ""
			epic := lineage.HighestEpic()
			if epic != nil {
				epic.BuildKeyURL(baseURL) // should not be needed.
				epicKeyDisplay = epic.KeyLinkMarkdown()
				epicName = epic.Summary
			}
			row = append(row, epicKeyDisplay, epicName)
		}

		stdCells := []string{
			issMeta.KeyLinkMarkdown(),
			issMore.Type(),
			issMore.Project(),
			issMore.Summary(),
			issMore.Status(),
			issMore.Resolution(),
			// strconv.Itoa(iss.Fields.AggregateTimeOriginalEstimate),
			// strconv.Itoa(iss.Fields.TimeOriginalEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeOriginalEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeSpent),
			is.Config.SecondsToDaysString(timeRemainingSecs),
			// time.Time(iss.Fields.Created).Format(time.RFC3339),
			// strconvutil.FormatFloat64Simple(float64(ix.TimeRemainingEstimate.Days(is.Config.WorkingHoursPerDay))),
		}
		row = append(row, stdCells...)

		if customCols != nil {
			for _, cc := range customCols.Cols {
				if cc.RenderSkip {
					continue
				} else if cc.Func == nil {
					row = append(row, "")
				} else if val, err := cc.Func(iss); err != nil {
					return nil, err
				} else {
					row = append(row, val)
				}
			}
		}

		tbl.Rows = append(tbl.Rows, row)
	}
	return &tbl, nil
}

func IssuesSetReadFileJSON(filename string) (*IssuesSet, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	is := &IssuesSet{}
	err = json.Unmarshal(b, is)
	return is, err
}

func (is *IssuesSet) WriteFileJSON(name, prefix, indent string) error {
	j, err := jsonutil.MarshalSimple(is, prefix, indent)
	if err != nil {
		return err
	}
	return os.WriteFile(name, j, 0600)
}
