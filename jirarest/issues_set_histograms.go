package jirarest

import (
	"errors"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/gojira"
)

// HistogramMapProjectTypeStatus provides issue counts by: Project, Type, and Status.
func (is *IssuesSet) HistogramMapProjectTypeStatus() *histogram.Histogram {
	h := histogram.NewHistogram(gojira.FieldIssuePlural)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := IssueMore{Issue: &iss}
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
		im := IssueMore{Issue: &iss}
		hset.Add(im.ProjectKey(), im.Type(), 1)
	}
	return hset
}

// HistogramSetsProjectTypeStatus provides issue counts by: Project, Type, and Status.
func (is *IssuesSet) HistogramSetsProjectTypeStatus() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets(gojira.FieldIssuePlural)
	for _, iss := range is.IssuesMap {
		iss := iss
		im := IssueMore{Issue: &iss}
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
		} else if is.Config.StatusesSet == nil {
			return nil, errors.New("statusesSet not set")
		} else {
			statusCategoryFunc = is.Config.StatusesSet.StatusCategory
		}
	}
	for _, iss := range is.IssuesMap {
		iss := iss
		im := IssueMore{Issue: &iss}
		key := im.Key()
		if key == "" {
			return nil, errors.New("issue key cannot be empty")
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

type (
	WorkstreamFuncMake func(issueKey string) (string, error)
	WorkstreamFuncIncl func(ws string) bool
)

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
