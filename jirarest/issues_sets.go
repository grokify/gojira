package jirarest

import (
	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/table"
	"github.com/grokify/mogo/type/maputil"
)

type IssuesSets struct {
	Order []string
	Data  map[string]IssuesSet
}

func NewIssuesSets() *IssuesSets {
	return &IssuesSets{
		Order: []string{},
		Data:  map[string]IssuesSet{}}
}

func (sets *IssuesSets) OrderOrDefault() []string {
	if len(sets.Order) > 0 {
		return sets.Order
	} else {
		return maputil.Keys(sets.Data)
	}
}

func (sets *IssuesSets) Upsert(setName string, set *IssuesSet) {
	sets.Data[setName] = *set
}

func (sets *IssuesSets) UpsertIssueKeys(jrClient *Client, setName string, issueKeys []string) error {
	if jrClient == nil {
		return ErrClientCannotBeNil
	}
	ii, err := jrClient.IssueAPI.Issues(issueKeys, false)
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
	for setName, set := range sets.Data {
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
