package jiraweb

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/text/markdown"
)

const (
	WebSlugBrowse = "/browse"
)

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
