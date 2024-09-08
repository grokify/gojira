package jirarest

import (
	"errors"
	"fmt"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gocharts/v2/data/timeseries"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/time/timeutil"
)

// TimeSeriesCreatedMonth provides issue counts by month by create date
func (is *IssuesSet) TimeSeriesCreatedMonth() *timeseries.TimeSeries {
	ts := timeseries.NewTimeSeries("by month")
	ts.Interval = timeutil.IntervalMonth
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		ts.AddInt64(im.CreateTime(), 1)
	}
	ts2 := ts.ToMonth(true)
	return &ts2
}

// TimeSeriesSetCreatedMonthByCustomField provides issue counts by custom field and month by create date.
// `customFieldID` is aunit for the integer part of `customfield_12345` or `cf[12345]`.
func (is *IssuesSet) TimeSeriesSetCreatedMonthByCustomField(cumulative, inflate, popLast bool, monthsFilter []time.Month, customFieldID uint) (*timeseries.TimeSeriesSet, error) {
	customFieldLabel := fmt.Sprintf("customfield_%d", customFieldID)
	return is.TimeSeriesSetCreatedMonthByKey(
		cumulative, inflate, popLast, monthsFilter,
		func(iss jira.Issue) (string, error) {
			im := NewIssueMore(&iss)
			icf, err := im.CustomField(customFieldLabel)
			if err != nil {
				return "", err
			}
			return icf.Value, nil
		},
	)
}

// TimeSeriesSetCreatedMonthByProject provides issue counts by project and month by create date
func (is *IssuesSet) TimeSeriesSetCreatedMonthByProject(cumulative, inflate, popLast bool, monthsFilter []time.Month) (*timeseries.TimeSeriesSet, error) {
	return is.TimeSeriesSetCreatedMonthByKey(
		cumulative, inflate, popLast, monthsFilter,
		func(iss jira.Issue) (string, error) {
			im := NewIssueMore(&iss)
			return im.ProjectKey(), nil
		},
	)
}

// TimeSeriesSetCreatedMonthByResolution provides issue counts by resolution and month by create date
func (is *IssuesSet) TimeSeriesSetCreatedMonthByResolution(cumulative, inflate, popLast bool, monthsFilter []time.Month) (*timeseries.TimeSeriesSet, error) {
	return is.TimeSeriesSetCreatedMonthByKey(
		cumulative, inflate, popLast, monthsFilter,
		func(iss jira.Issue) (string, error) {
			im := NewIssueMore(&iss)
			return im.Resolution(), nil
		},
	)
}

// TimeSeriesSetCreatedMonthByStatus provides issue counts by status and month by create date
func (is *IssuesSet) TimeSeriesSetCreatedMonthByStatus(cumulative, inflate, popLast bool, monthsFilter []time.Month) (*timeseries.TimeSeriesSet, error) {
	return is.TimeSeriesSetCreatedMonthByKey(
		cumulative, inflate, popLast, monthsFilter,
		func(iss jira.Issue) (string, error) {
			im := NewIssueMore(&iss)
			return im.Status(), nil
		},
	)
}

// TimeSeriesCreatedMonth provides issue counts by month by create date
func (is *IssuesSet) TimeSeriesSetCreatedMonthByKey(cumulative, inflate, popLast bool, monthsFilter []time.Month, fnKey func(iss jira.Issue) (string, error)) (*timeseries.TimeSeriesSet, error) {
	if fnKey == nil {
		return nil, errors.New("fnKey cannot be nil")
	}
	tss := timeseries.NewTimeSeriesSet("By Project By Month")
	tss.Interval = timeutil.IntervalMonth
	for _, iss := range is.IssuesMap {
		iss := iss
		tssKey, err := fnKey(iss)
		if err != nil {
			return nil, err
		}
		im := NewIssueMore(&iss)
		tss.AddInt64(tssKey, im.CreateTime(), 1)
	}
	tssm, err := tss.ToMonth(cumulative, inflate, popLast, monthsFilter)
	return &tssm, err
}

// HistogramMapProjectTypeStatus provides issue counts by: Project, Type, and Status.
func (is *IssuesSet) HistogramMapProjectTypeStatus() *histogram.Histogram {
	h := histogram.NewHistogram(gojira.FieldIssuePlural)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		h.AddMap(map[string]string{
			gojira.FieldProject: im.ProjectKey(),
			gojira.FieldType:    im.Type(),
			gojira.FieldStatus:  im.Status(),
		}, 1)
	}
	return h
}

func (is *IssuesSet) TableSetProjectTypeStatus(tsConfig *histogram.HistogramMapTableSetConfig) (*table.TableSet, error) {
	hist := is.HistogramMapProjectTypeStatus()
	if tsConfig == nil {
		tsConfig = DefaultHistogramMapTableConfig([]string{})
	}
	return hist.TableSetMap(tsConfig.Configs)
}

func DefaultHistogramMapTableConfig(projectKeys []string) *histogram.HistogramMapTableSetConfig {
	colNameIssueCount := "Issue Count"
	return &histogram.HistogramMapTableSetConfig{
		Configs: []histogram.HistogramMapTableConfig{
			{
				TableName: "Project Type Status",
				ColumnKeys: []string{
					gojira.FieldProject,
					gojira.FieldType,
					gojira.FieldStatus},
				ColNameCount: colNameIssueCount,
			},
			{
				TableName: "Meta Type",
				ColumnKeys: []string{
					gojira.FieldProject,
					gojira.FieldType},
				ColNameCount: colNameIssueCount,
			},
			{
				TableName: "Meta Status",
				ColumnKeys: []string{
					gojira.FieldProject,
					gojira.FieldStatus},
				ColNameCount: colNameIssueCount,
			},
			{
				TableNamePrefix:    "Project - ",
				SplitKey:           gojira.FieldProject,
				SplitValFilterIncl: projectKeys,
				ColumnKeys: []string{
					gojira.FieldType,
					gojira.FieldStatus},
				ColNameCount: colNameIssueCount,
			},
		},
	}
}

// HistogramSetProjectType returns a list of histograms by Project and Type.
func (is *IssuesSet) HistogramSetProjectType() *histogram.HistogramSet {
	hset := histogram.NewHistogramSet(gojira.FieldIssuePlural)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		hset.Add(im.ProjectKey(), im.Type(), 1)
	}
	return hset
}

// HistogramSetsProjectTypeStatus provides issue counts by: Project, Type, and Status.
func (is *IssuesSet) HistogramSetsProjectTypeStatus() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets(gojira.FieldIssuePlural)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		hsets.Add(
			im.ProjectKey(),
			im.Type(),
			im.Status(),
			1,
			true)
	}
	return hsets
}

func (is *IssuesSet) HistogramMap(stdKeys []string, calcFields []IssueCalcField) (*histogram.Histogram, error) {
	h := histogram.NewHistogram("")
	return h, nil
}

func (is *IssuesSet) ExportWorkstremaFilter(wsFuncMake WorkstreamFuncMake, wsFuncIncl WorkstreamFuncIncl, customFieldLabels []string) (*IssuesSet, error) {
	out := NewIssuesSet(is.Config)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		key := im.Key()
		if key == "" {
			return nil, ErrIssueKeyCannotBeEmpty
		} else if ws, err := wsFuncMake(key); err != nil {
			return nil, err
		} else if wsFuncIncl != nil && !wsFuncIncl(ws) {
			continue
		} else if err = out.Add(iss); err != nil {
			return nil, err
		} else if lineages, err := is.Lineage(key, customFieldLabels); err != nil {
			return nil, err
		} else {
			for _, im := range lineages {
				if im.Key == key {
					continue
				} else if iss, err := is.Get(im.Key); err != nil {
					return nil, err
				} else if err = out.Parents.Add(iss); err != nil {
					return nil, err
				}
			}
		}
	}
	return out, nil
}

type (
	WorkstreamFuncMake func(issueKey string) (string, error)
	WorkstreamFuncIncl func(ws string) bool
)

func (is *IssuesSet) ExportWorkstreamXfieldStatusHistogramSets(
	wsFuncMake WorkstreamFuncMake,
	wsFuncIncl WorkstreamFuncIncl,
	xfieldSlug string,
	useStatusCategory bool) (*histogram.HistogramSets, error) {
	if wsFuncMake == nil {
		return nil, errors.New("workstream func not supplied")
	}
	if wsFuncIncl == nil {
		wsFuncIncl = func(ws string) bool { return true }
	}
	xfieldSlugs := map[string]int{
		FieldSlugProjectkey: 1,
		FieldSlugType:       1,
	}
	if _, ok := xfieldSlugs[xfieldSlug]; !ok {
		return nil, errors.New("xfieldSlug not known")
	}
	hss := histogram.NewHistogramSets("issues")
	statusCategoryFunc := func(s string) string { return s }
	if useStatusCategory {
		if is.Config == nil {
			return nil, errors.New("config not set")
		} else if is.Config.StatusConfig == nil {
			return nil, errors.New("statusesSet not set")
		} else {
			statusCategoryFunc = is.Config.StatusConfig.MetaStage
		}
	}
	for _, iss := range is.IssuesMap {
		iss := iss
		im := NewIssueMore(&iss)
		key := im.Key()
		if key == "" {
			return nil, ErrIssueKeyCannotBeEmpty
		}
		ws, err := wsFuncMake(key)
		if err != nil {
			return nil, err
		}
		if wsFuncIncl != nil && !wsFuncIncl(ws) {
			continue
		}
		status := im.Status()
		if useStatusCategory {
			statusCategory := statusCategoryFunc(status)
			if statusCategory != "" {
				status = statusCategory
			}
		}
		xfield := ""
		switch xfieldSlug {
		case FieldSlugProjectkey:
			xfield = im.ProjectKey()
		case FieldSlugType:
			xfield = im.Type()
		}

		hss.Add(ws, xfield, status, 1, true)
	}
	return hss, nil
}

func (is *IssuesSet) ExportWorkstreamXfieldStatusTablePivot(wsFuncMake WorkstreamFuncMake, wsFuncIncl WorkstreamFuncIncl, xfieldSlug, xfieldName string, useStatusCategory bool) (*table.Table, error) {
	hss, err := is.ExportWorkstreamXfieldStatusHistogramSets(wsFuncMake, wsFuncIncl, xfieldSlug, useStatusCategory)
	if err != nil {
		return nil, err
	}
	// tbl := hss.TablePivot("issues", "Workstream", xfieldName, "Status: ", "", is.StatusesOrder(), true)
	tbl := hss.TablePivot(histogram.TablePivotOpts{
		TableName:           "issues",
		ColNameHistogramSet: "Workstream",
		ColNameHistogram:    xfieldName,
		ColNameBinPrefix:    "Status: ",
		BinNamesOrder:       is.StatusesOrder(),
		InclBinsUnordered:   true,
		InclBinCounts:       true,
		InclBinCountsSum:    true,
		InclBinPercentages:  true,
	})
	return &tbl, nil
}

func (is *IssuesSet) ExportWorkstreamProjectkeyStatusTablePivot(wsFuncMake WorkstreamFuncMake, wsFuncIncl WorkstreamFuncIncl, useStatusCategory bool) (*table.Table, error) {
	hss, err := is.ExportWorkstreamXfieldStatusHistogramSets(wsFuncMake, wsFuncIncl, FieldSlugProjectkey, useStatusCategory)
	if err != nil {
		return nil, err
	}
	// tbl := hss.TablePivot("issues", "Workstream", "Project Key", "Status: ", "", is.StatusesOrder(), true)
	tbl := hss.TablePivot(histogram.TablePivotOpts{
		TableName:           "issues",
		ColNameHistogramSet: "Workstream",
		ColNameHistogram:    "Project Key",
		ColNameBinPrefix:    "Status: ",
		BinNamesOrder:       is.StatusesOrder(),
		InclBinsUnordered:   true,
		InclBinCounts:       true,
		InclBinCountsSum:    true,
		InclBinPercentages:  true,
	})

	return &tbl, nil
}

func (is *IssuesSet) ExportWorkstreamTypeStatusTablePivot(wsFuncMake WorkstreamFuncMake, wsFuncIncl WorkstreamFuncIncl, useStatusCategory bool) (*table.Table, error) {
	hss, err := is.ExportWorkstreamXfieldStatusHistogramSets(wsFuncMake, wsFuncIncl, FieldSlugType, useStatusCategory)
	if err != nil {
		return nil, err
	}
	// tbl := hss.TablePivot("issues", "Workstream", "Type", "Status: ", "", is.StatusesOrder(), true)
	tbl := hss.TablePivot(histogram.TablePivotOpts{
		TableName:           "issues",
		ColNameHistogramSet: "Workstream",
		ColNameHistogram:    "Type",
		ColNameBinPrefix:    "Status: ",
		BinNamesOrder:       is.StatusesOrder(),
		InclBinsUnordered:   true,
		InclBinCounts:       true,
		InclBinCountsSum:    true,
		InclBinPercentages:  true,
	})
	return &tbl, nil
}

// Workstream | Story|Bug|Spike | Status | Team

type IssueCalcField struct {
	Key     string
	ValFunc func(iss *jira.Issue) (string, error)
}

type CustomJiraProcessor struct {
	*IssuesSet
}
