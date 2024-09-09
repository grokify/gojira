package jiraweb

import (
	"errors"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	WebSlugBrowse  = "/browse"
	IssueURLFormat = `%s/browse/%s`
)

var rxJiraTicket = regexp.MustCompile(`([A-Z]+\-[0-9]+)`)

func ParseKeys(s string, unique, asc bool) []string {
	var keys []string
	m := rxJiraTicket.FindAllStringSubmatch(s, -1)
	for _, n := range m {
		if len(n) > 1 {
			keys = append(keys, n[1])
		}
	}
	return stringsutil.SliceCondenseSpace(keys, unique, asc)
}

func IssueLinkWebMarkdownOrEmptyFromIssueKey(serverURL, issueKey string) string {
	if issueKey = strings.TrimSpace(issueKey); issueKey == "" {
		return ""
	} else {
		if u, err := IssueURLWebFromIssueKey(serverURL, issueKey); err != nil {
			return ""
		} else {
			return markdown.Linkify(u, issueKey)
		}
	}
}

func IssueURLWebOrEmptyFromIssueKey(serverURL, issueKey string) string {
	if issueKey = strings.TrimSpace(issueKey); issueKey == "" {
		return ""
	} else {
		if u, err := IssueURLWebFromIssueKey(serverURL, issueKey); err != nil {
			return ""
		} else {
			return u
		}
	}
}

func IssueURLWebFromIssueKey(serverURL, issueKey string) (string, error) {
	var parts = []string{}
	if issueKey = strings.TrimSpace(issueKey); issueKey == "" {
		return "", errors.New("issue key cannot be empty")
	}
	if strings.TrimSpace(serverURL) == "" {
		parts = []string{WebSlugBrowse, issueKey}
	} else {
		parts = []string{serverURL, WebSlugBrowse, issueKey}
	}
	return urlutil.JoinAbsolute(parts...), nil
}
