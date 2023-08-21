package jirarest

import (
	"encoding/json"
	"io"
	"os"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gocharts/v2/data/histogram"
)

// IssuesResponse is only a small wrapper around the Search (with JQL) method to be able to parse the results
type IssuesResponse struct {
	Issues     Issues `json:"issues" structs:"issues"`
	Expand     string `json:"expand"`
	StartAt    int    `json:"startAt" structs:"startAt"`
	MaxResults int    `json:"maxResults" structs:"maxResults"`
	Total      int    `json:"total" structs:"total"`
}

func IssuesResponseReadFile(filename string) (*IssuesResponse, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return ParseIssuesResponseReader(f)
}

func ParseIssuesResponseReader(r io.Reader) (*IssuesResponse, error) {
	b, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return ParseIssuesResponseBytes(b)
}

func ParseIssuesResponseBytes(b []byte) (*IssuesResponse, error) {
	ir := IssuesResponse{}
	return &ir, json.Unmarshal(b, &ir)
}

type Issues []jira.Issue

func (ii Issues) CountsByType() map[string]int {
	counts := map[string]int{}
	for _, iss := range ii {
		name := iss.Fields.Type.Name
		counts[name]++
		counts["_total"]++
	}
	return counts
}

// CountsByProjectTypeStatus returns a `*histogram.Histogram` with issue counts
// by project, type, and status. This can be used to export CSV and XLSX sheets
// for analysis.
func (ii Issues) CountsByProjectTypeStatus() *histogram.HistogramSets {
	hsets := histogram.NewHistogramSets("")
	for _, iss := range ii {
		hsets.Add(
			iss.Fields.Project.Key,
			iss.Fields.Type.Name,
			iss.Fields.Status.Name,
			1,
			true)
	}
	return hsets
}

func GetIssueCustomValueStruct(iss jira.Issue) (*IssueCustomField, error) {
	if iss.Fields == nil {
		return nil, nil
	}
	unv, ok := iss.Fields.Unknowns["customfield_12461"]
	if !ok {
		return nil, nil
	}
	icf := &IssueCustomField{}
	err := UnmarshalAny(unv, icf)
	return icf, err
}
