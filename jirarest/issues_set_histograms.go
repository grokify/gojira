package jirarest

import (
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

// Workstream | Story|Bug|Spike | Status | Team

type IssueCalcField struct {
	Key     string
	ValFunc func(iss *jira.Issue) (string, error)
}

type CustomJiraProcessor struct {
	*IssuesSet
}
