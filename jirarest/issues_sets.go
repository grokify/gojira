package jirarest

import (
	"context"
	"slices"
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/charts/text/progressbarchart"
	"github.com/grokify/gocharts/v2/data/histogram"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/type/maputil"
	"github.com/olekukonko/errors"
)

type IssuesSets struct {
	Name  string
	Order []string
	Items map[string]IssuesSet
}

func NewIssuesSets() *IssuesSets {
	return &IssuesSets{
		Order: []string{},
		Items: map[string]IssuesSet{}}
}

func (sets *IssuesSets) AddIssuesSetFilterKeys(name string, iset *IssuesSet, keys []string, errOnUnfound bool) error {
	if iset == nil {
		return errors.New("issues set cannot be nil")
	} else if newIset, err := iset.FilterByKeys(keys, errOnUnfound); err != nil {
		return err
	} else {
		sets.Items[name] = *newIset
		return nil
	}
}

func (sets *IssuesSets) BarChartsText(inclProgress, inclFunnel bool, startNumber *int) (string, error) {
	var sb strings.Builder
	var useNumber bool
	nextNum := 0
	if startNumber != nil {
		useNumber = true
		nextNum = *startNumber
	}
	if inclProgress {
		var parts []string
		if useNumber {
			parts = append(parts, strconv.Itoa(nextNum)+".")
			nextNum++
		}
		if sets.Name != "" {
			parts = append(parts, sets.Name)
		}
		parts = append(parts, "Progress")
		if _, err := sb.WriteString(strings.Join(parts, " ") + "\n\n"); err != nil {
			return "", err
		}
		if cht := sets.BarChartTextProgress(); strings.TrimSpace(cht) != "" {
			if _, err := sb.WriteString(cht); err != nil {
				return "", err
			}
		}
	}
	if inclFunnel {
		if inclProgress {
			if _, err := sb.WriteString("\n\n"); err != nil {
				return "", err
			}
		}
		var parts []string

		if useNumber {
			parts = append(parts, strconv.Itoa(nextNum)+".")
		}
		if sets.Name != "" {
			parts = append(parts, sets.Name)
		}
		parts = append(parts, "Funnel")
		if _, err := sb.WriteString(strings.Join(parts, " ") + "\n\n"); err != nil {
			return "", err
		}
		if cht := sets.BarChartTextFunnel(); strings.TrimSpace(cht) != "" {
			if _, err := sb.WriteString(cht); err != nil {
				return "", err
			}
		}
	}

	return sb.String(), nil
}

func (sets *IssuesSets) BarChartTextProgress() string {
	h := sets.Histogram()
	tasks := progressbarchart.NewTasksFromHistogram(h)
	return tasks.ProgressBarChartText()
}

func (sets *IssuesSets) BarChartTextFunnel() string {
	h := sets.Histogram()
	tasks := progressbarchart.NewTasksFunnelFromHistogram(h)
	return tasks.ProgressBarChartText()
}

func (sets *IssuesSets) Histogram() *histogram.Histogram {
	h := histogram.NewHistogram("")
	h.Order = slices.Clone(sets.Order)
	for k, v := range sets.Items {
		h.Add(k, len(v.Items))
	}
	return h
}

func (sets *IssuesSets) OrderOrDefault() []string {
	if len(sets.Order) > 0 {
		return sets.Order
	} else {
		return maputil.Keys(sets.Items)
	}
}

func (sets *IssuesSets) Upsert(setName string, set *IssuesSet) {
	sets.Items[setName] = *set
}

func (sets *IssuesSets) UpsertIssueKeys(ctx context.Context, jrClient *Client, setName string, issueKeys []string) error {
	if jrClient == nil {
		return ErrClientCannotBeNil
	}
	ii, err := jrClient.IssueAPI.Issues(ctx, issueKeys, nil)
	if err != nil {
		return err
	}
	is, err := ii.IssuesSet(nil)
	if err != nil {
		return err
	}
	sets.Upsert(setName, is)
	return nil
}

func (sets *IssuesSets) TableSet(
	tblColsMapKeys []string,
	contColumnName string,
	fnIss func(iss *jira.Issue) (map[string]string, error),
	fnRowSort func(a, b []string) int,
) (*table.TableSet, error) {
	ts := table.NewTableSet("")
	for setName, set := range sets.Items {
		hmap, err := set.HistogramMapFunc(fnIss)
		if err != nil {
			return nil, err
		}
		tbl, err := hmap.TableMap(tblColsMapKeys, contColumnName, fnRowSort)
		tbl.Name = setName
		if err != nil {
			return nil, err
		}
		if err = ts.Add(tbl); err != nil {
			return nil, err
		}
	}
	return ts, nil
}
