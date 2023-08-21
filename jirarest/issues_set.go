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
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/type/slicesutil"
)

/*
const (
	WorkingHoursPerDayDefault float32 = 8.0
	WorkingDaysPerWeekDefault float32 = 5.0
)

type Config struct {
	WorkingHoursPerDay float32
	WorkingDaysPerWeek float32
}

func (c *Config) SecondsToDays(sec int) float32 {
	return float32(sec) / 60 / 60 / c.WorkingHoursPerDay
}

func (c *Config) SecondsToDaysString(sec int) string {
	return strconvutil.FormatFloat64Simple(float64(c.SecondsToDays(sec)))
}

func NewConfigDefault() *Config {
	return &Config{
		WorkingHoursPerDay: WorkingHoursPerDayDefault,
		WorkingDaysPerWeek: WorkingDaysPerWeekDefault}
}
*/

type IssuesSet struct {
	Config    *gojira.Config
	IssuesMap map[string]jira.Issue
}

func NewIssuesSet(cfg *gojira.Config) IssuesSet {
	if cfg == nil {
		cfg = gojira.NewConfigDefault()
	}
	return IssuesSet{
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

func (is *IssuesSet) Issues() Issues {
	ii := Issues{}
	for _, iss := range is.IssuesMap {
		ii = append(ii, iss)
	}
	return ii
}

func (is *IssuesSet) TimeStats() gojira.TimeStats {
	ts := gojira.TimeStats{
		WorkingDaysPerWeek: is.Config.WorkingDaysPerWeek,
		WorkingHoursPerDay: is.Config.WorkingHoursPerDay,
	}
	for _, iss := range is.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		timeRemainingOriginal, timeRemaining := gojira.TimeRemaining(iss.Fields.Status.Name, iss.Fields.TimeOriginalEstimate, iss.Fields.TimeEstimate, iss.Fields.TimeSpent)
		ts.TimeSpent += iss.Fields.TimeSpent
		ts.TimeEstimate += iss.Fields.TimeEstimate
		ts.TimeOriginalEstimate += iss.Fields.TimeOriginalEstimate
		ts.AggregateTimeOriginalEstimate += iss.Fields.AggregateTimeOriginalEstimate
		ts.AggregateTimeSpent += iss.Fields.AggregateTimeSpent
		ts.AggregateTimeEstimate += iss.Fields.AggregateTimeEstimate
		ts.TimeRemaining += timeRemaining
		ts.TimeRemainingOriginal += timeRemainingOriginal
	}
	return ts
}

func (is *IssuesSet) WriteJSONFile(filename string) error {
	b, err := json.Marshal(is)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0600)
}

func (is *IssuesSet) HistogramSets(baseURL string) *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets("issues")

	for _, iss := range is.IssuesMap {
		if iss.Fields == nil {
			continue
		}
		sev, ok := iss.Fields.Unknowns["AB"]
		if !ok {
			continue
		}
		fmt.Printf("%v\n", sev)
	}

	return hsets
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

func DefaultIssuesSetTableColumns() *table.ColumnDefinitions {
	return &table.ColumnDefinitions{
		Definitions: []table.ColumnDefinition{
			{Name: "Epic Key", Format: table.FormatURL},
			{Name: "Epic Name"},
			{Name: "Project"},
			{Name: "Type"},
			{Name: "Key", Format: table.FormatURL},
			{Name: "Summary"},
			{Name: "Status"},
			{Name: "Resolution"},
			// {Name: "Aggregate Original Time Estimate Seconds", Format: table.FormatInt},
			// {Name: "Original Estimate Seconds", Format: table.FormatInt},
			{Name: "Original Estimate Days", Format: table.FormatFloat},
			{Name: "Estimate Days", Format: table.FormatFloat},
			{Name: "Time Spent Days", Format: table.FormatFloat},
			{Name: "Time Remaining Days", Format: table.FormatFloat},
		},
	}
}

func BuildJiraIssueURL(baseURL, issueKey string) string {
	return urlutil.JoinAbsolute(baseURL, "/browse/", issueKey)
}

func (is *IssuesSet) Table(baseURL string, customCols *CustomTableCols) (table.Table, error) {
	baseURL = strings.TrimSpace(baseURL)

	if is.Config == nil {
		is.Config = gojira.NewConfigDefault()
	}
	tbl := table.NewTable("issues")
	tbl.LoadColumnDefinitions(DefaultIssuesSetTableColumns())

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
			name := strings.TrimSpace(customCol.Name)
			if name != "" {
				tbl.Columns = append(tbl.Columns, name)
			} else {
				tbl.Columns = append(tbl.Columns, fmt.Sprintf("Column %d", j+1))
			}
		}
	}

	for key, iss := range is.IssuesMap {
		ifs := IssueFieldsSimple{Fields: iss.Fields}

		keyDisplay := key
		epicKeyDisplay := ifs.EpicKey()
		if len(baseURL) > 0 {
			keyURL := BuildJiraIssueURL(baseURL, key)
			keyDisplay = markdown.Linkify(keyURL, key)

			if len(epicKeyDisplay) > 0 {
				epicKeyURL := BuildJiraIssueURL(baseURL, ifs.EpicKey())
				epicKeyDisplay = markdown.Linkify(epicKeyURL, ifs.EpicKey())
			}
		}

		timeRemainingSecs := iss.Fields.TimeEstimate - iss.Fields.TimeSpent
		if timeRemainingSecs < 0 ||
			ifs.StatusName() == "Closed" ||
			ifs.StatusName() == "Done" {
			timeRemainingSecs = 0
		}

		row := []string{
			epicKeyDisplay,
			ifs.EpicName(),
			iss.Fields.Project.Name,
			iss.Fields.Type.Name,
			keyDisplay,
			iss.Fields.Summary,
			ifs.StatusName(),
			ifs.ResolutionName(),
			// strconv.Itoa(iss.Fields.AggregateTimeOriginalEstimate),
			// strconv.Itoa(iss.Fields.TimeOriginalEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeOriginalEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeEstimate),
			is.Config.SecondsToDaysString(iss.Fields.TimeSpent),
			is.Config.SecondsToDaysString(timeRemainingSecs),
			// time.Time(iss.Fields.Created).Format(time.RFC3339),
			// strconvutil.FormatFloat64Simple(float64(ix.TimeRemainingEstimate.Days(is.Config.WorkingHoursPerDay))),
		}

		if customCols != nil {
			for _, cc := range customCols.Cols {
				if cc.RenderSkip {
					continue
				}
				if cc.Func == nil {
					row = append(row, "")
				} else {
					val, err := cc.Func(iss)
					if err != nil {
						return tbl, err
					}
					row = append(row, val)
				}
			}
		}

		tbl.Rows = append(tbl.Rows, row)
	}
	return tbl, nil
}

type IssueFieldsSimple struct {
	Fields *jira.IssueFields
}

func (ifs IssueFieldsSimple) EpicKey() string {
	if ifs.Fields == nil || ifs.Fields.Epic == nil {
		return ""
	} else {
		return ifs.Fields.Epic.Key
	}
}

func (ifs IssueFieldsSimple) EpicName() string {
	if ifs.Fields == nil || ifs.Fields.Epic == nil {
		return ""
	} else {
		if ifs.Fields.Epic.Name != "" {
			return ifs.Fields.Epic.Name
		}
		if ifs.Fields.Epic.Summary != "" {
			return ifs.Fields.Epic.Summary
		}
		return " "
	}
}

func (ifs IssueFieldsSimple) ResolutionName() string {
	if ifs.Fields == nil || ifs.Fields.Resolution == nil {
		return ""
	} else {
		return ifs.Fields.Resolution.Name
	}
}

func (ifs IssueFieldsSimple) StatusName() string {
	if ifs.Fields == nil || ifs.Fields.Status == nil {
		return ""
	} else {
		return ifs.Fields.Status.Name
	}
}
