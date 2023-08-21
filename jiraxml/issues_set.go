package jiraxml

import (
	"fmt"
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/strconv/strconvutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/type/maputil"
)

type IssuesSet struct {
	Config    *gojira.Config
	IssuesMap map[string]Issue
}

func NewIssuesSet(cfg *gojira.Config) IssuesSet {
	if cfg == nil {
		cfg = gojira.NewConfigDefault()
	}
	return IssuesSet{
		Config:    cfg,
		IssuesMap: map[string]Issue{}}
}

func (is *IssuesSet) AddFromAPI(issues ...jira.Issue) error {
	for _, iss := range issues {
		err := is.Add(IssueFromAPI(iss))
		if err != nil {
			return err
		}
	}
	return nil
}

func (is *IssuesSet) Add(issues ...Issue) error {
	if is.IssuesMap == nil {
		is.IssuesMap = map[string]Issue{}
	}
	missingKeyIndexes := []string{}
	for i, ix := range issues {
		k := ix.GetKey()
		if k == "" {
			missingKeyIndexes = append(missingKeyIndexes, strconv.Itoa(i))
			continue
		}
		ix.TrimSpace()
		is.IssuesMap[k] = ix
	}
	if len(missingKeyIndexes) > 0 {
		return fmt.Errorf("missingkeyindexes (%s)", strings.Join(missingKeyIndexes, ","))
	}
	return nil
}

func (is *IssuesSet) ReadFile(filename string) error {
	x, err := ReadFile(filename)
	if err != nil {
		return err
	}
	return is.Add(x.Channel.Issues...)
}

func (is *IssuesSet) Keys() []string {
	return maputil.Keys(is.IssuesMap)
}

func (is *IssuesSet) Table(baseURL string) table.Table {
	baseURL = strings.TrimSpace(baseURL)

	if is.Config == nil {
		is.Config = gojira.NewConfigDefault()
	}
	tbl := table.NewTable("issues")
	tbl.Columns = []string{"Type", "Key", "Summary", "Status", "Resolution", "Aggregate Original Time Estimate Seconds", "Original Estimate Seconds", "Original Estimate Days", "Estimate Days", "Time Spent", "Time Remaining"}
	tbl.FormatMap = map[int]string{
		1: table.FormatURL,
		5: table.FormatInt,
		6: table.FormatInt,
		7: table.FormatFloat,
		8: table.FormatFloat,
		9: table.FormatFloat,
	}
	for k, ix := range is.IssuesMap {
		keyDisplay := k
		if len(baseURL) > 0 {
			keyURL := urlutil.JoinAbsolute(baseURL, "/browse/", k)
			keyDisplay = markdown.Linkify(keyURL, k)
		}
		tbl.Rows = append(tbl.Rows, []string{
			ix.Type.DisplayName,
			keyDisplay,
			// ix.Title,
			ix.Summary,
			ix.Status.DisplayName,
			ix.Resolution.DisplayName,
			strconv.Itoa(int(ix.AggregateTimeOriginalEstimate.Seconds)),
			strconv.Itoa(int(ix.TimeOriginalEstimate.Seconds)),
			strconvutil.FormatFloat64Simple(float64(ix.TimeOriginalEstimate.Days(is.Config.WorkingHoursPerDay))),
			strconvutil.FormatFloat64Simple(float64(ix.TimeEstimate.Days(is.Config.WorkingHoursPerDay))),
			strconvutil.FormatFloat64Simple(float64(ix.TimeSpent.Days(is.Config.WorkingHoursPerDay))),
			// strconvutil.FormatFloat64Simple(float64(ix.TimeRemainingEstimate.Days(is.Config.WorkingHoursPerDay))),
		})
	}
	return tbl
}

func (is *IssuesSet) AggregateTimeEstimate() int64 {
	var agg int64
	for _, ix := range is.IssuesMap {
		agg += ix.TimeEstimate.Seconds
	}
	return agg
}

func (is *IssuesSet) AggregateTimeOriginalEstimate() int64 {
	var agg int64
	for _, ix := range is.IssuesMap {
		agg += ix.TimeOriginalEstimate.Seconds
	}
	return agg
}

func (is *IssuesSet) AggregateTimeRemainingEstimate() int64 {
	var agg int64
	for _, ix := range is.IssuesMap {
		agg += ix.TimeSpent.Seconds
	}
	return agg
}

func (is *IssuesSet) AggregateTimeSpent() int64 {
	var agg int64
	for _, ix := range is.IssuesMap {
		agg += ix.TimeSpent.Seconds
	}
	return agg
}

// TSRHistogramSets returns a `*histogram.HistogramSets` for Type, Status and Resolution.
func (is *IssuesSet) TSRHistogramSets(name string) *histogram.HistogramSets {
	if strings.TrimSpace(name) == "" {
		name = "TSR"
	}
	hset := histogram.NewHistogramSets(name)
	for _, iss := range is.IssuesMap {
		hset.Add(
			iss.Type.DisplayName,
			iss.Status.DisplayName,
			iss.Resolution.DisplayName,
			1, true)
	}
	return hset
}

// TSRTable returns a `table.Table` for Type, Status and Resolution.
func (is *IssuesSet) TSRTable(name string) table.Table {
	hset := is.TSRHistogramSets(name)
	return hset.Table("Jira Issues", "Type", "Status", "Resolution", "Count")
}

// TSRWriteCSV writes a CSV file for Type, Status and Resolution.
func (is *IssuesSet) TSRWriteCSV(filename string) error {
	tbl := is.TSRTable("")
	return tbl.WriteCSV(filename)
}

// TSRWriteCSV writes a CSV file for Type, Status and Resolution.
func (is *IssuesSet) TSRWriteXLSX(filename, sheetname string) error {
	tbl := is.TSRTable("")
	return tbl.WriteXLSX(filename, sheetname)
}
