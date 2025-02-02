package gojira

import (
	"fmt"
	"strings"
	"time"

	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/time/timeutil"
)

type JQLsReportMarkdownOpts struct {
	TimeZone       string
	AddCount       bool
	AddExplicitURL bool
}

// JQLsReportMarkdownLines provides Markdownlines for a set of JQLs.
// The `JQLsReportMarkdownOpts.AddCount` option adds a static count to the report. This is useful when the
// report isn't auto-updating, such as on a code repo.
// The `JQLsReportMarkdownOpts.AddExplicitURL` option adds a URL which can be pasted into Confluence. This
// will be interpreted to load a dynamic table in the Confluence page. This is not needed for git repo
// pages.
func JQLsReportMarkdownLines(issuesWebURL, headerPrefix string, jqls []JQL, opts JQLsReportMarkdownOpts) ([]string, error) {
	var lines []string
	issuesWebURL = strings.TrimSpace(issuesWebURL) // should end with `issues/?`
	timeZone := strings.TrimSpace(opts.TimeZone)
	for i, j := range jqls {
		dtStr := ""
		count := -1
		dt := j.Meta.QueryTime
		if !dt.IsZero() {
			if timeZone != "" {
				if dtTry, err := timeutil.TimeUpdateLocation(dt, timeZone); err != nil {
					return lines, err
				} else {
					dt = dtTry
				}
			}
			dtStr = fmt.Sprintf(" at %s", dt.Format(time.RFC1123))
			count = j.Meta.QueryTotalCount
		}

		name := strings.TrimSpace(j.Meta.Name)
		if name == "" {
			name = fmt.Sprintf("JQL #%d", i+1)
		}
		lines = append(lines, "", fmt.Sprintf("%s%s", headerPrefix, name))

		if key := strings.TrimSpace(j.Meta.Key); key != "" {
			lines = append(lines, fmt.Sprintf("* Key: `%s`", key))
		}
		if desc := strings.TrimSpace(j.Meta.Description); desc != "" {
			lines = append(lines, fmt.Sprintf("* Description: %s", desc))
		}
		if opts.AddCount && count > -1 {
			lines = append(lines, fmt.Sprintf("* Issue Count: **%d**%s", count, dtStr))
		}
		js := strings.TrimSpace(j.String())
		jq := strings.TrimSpace(j.QueryString())
		if jq != "" {
			jq = issuesWebURL + jq
		}
		mk := markdown.Linkify(jq, "`"+js+"`")
		if mk != "" {
			lines = append(lines, fmt.Sprintf("* JQL: %s", mk))
		}
		if opts.AddExplicitURL && jq != "" {
			lines = append(lines, "", markdown.Linkify(jq, jq))
		}
	}
	return lines, nil
}
