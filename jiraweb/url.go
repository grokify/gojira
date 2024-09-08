package jiraweb

import (
	"errors"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
)

const (
	WebSlugBrowse = "/browse"
)

func IssueIDToURL(serverURL, issueKey string) (string, error) {
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
