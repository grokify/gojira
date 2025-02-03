package jirarest

import (
	"github.com/grokify/gojira"
)

// JQLsReportMarkdownLines provides Markdownlines for a set of JQLs, including querying the number
// of results for each JQL via the Jira API.
func (svc *IssueService) JQLsReportMarkdownLines(headerPrefix string, jqls gojira.JQLs, opts gojira.JQLsReportMarkdownOpts) ([]string, error) {
	if jqls, err := svc.JQLsAddMetadata(jqls); err != nil {
		return []string{}, err
	} else {
		return jqls.ReportMarkdownLines(svc.WebURL(), headerPrefix, opts)
	}
}
