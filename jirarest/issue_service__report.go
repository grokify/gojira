package jirarest

import (
	"strings"

	"github.com/grokify/gojira"
	"github.com/grokify/mogo/net/urlutil"
)

// JQLsReportMarkdownLines provides Markdownlines for a set of JQLs, including querying the number
// of results for each JQL via the Jira API.
func (svc *IssueService) JQLsReportMarkdownLines(headerPrefix string, jqls gojira.JQLs, opts gojira.JQLsReportMarkdownOpts) ([]string, error) {
	jqls, err := svc.JQLsAddMetadata(jqls)
	if err != nil {
		return []string{}, err
	}
	issuesWebURL := ""
	if svc.Client != nil && svc.Client.Config != nil {
		if svrURL := strings.TrimSpace(svc.Client.Config.ServerURL); svrURL != "" {
			issuesWebURL = urlutil.JoinAbsolute(svrURL, "issues/?")
		}
	}
	return jqls.ReportMarkdownLines(issuesWebURL, headerPrefix, opts)
}
