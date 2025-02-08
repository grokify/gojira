package jirarest

import (
	"github.com/grokify/gojira"
)

// JQLsReportMarkdownLines provides Markdownlines for a set of JQLs, including querying the number
// of results for each JQL via the Jira API.
func (svc *IssueService) JQLsReportMarkdownLines(jqls gojira.JQLs, opts *gojira.JQLsReportMarkdownOpts) ([]string, error) {
	if jqls, err := svc.JQLsAddMetadata(jqls); err != nil {
		return []string{}, err
	} else {
		if opts == nil {
			opts = &gojira.JQLsReportMarkdownOpts{}
		}
		opts.IssuesWebURL = svc.WebURL()
		return jqls.ReportMarkdownLines(opts)
	}
}
