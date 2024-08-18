package gojira

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/text/markdown"
	"github.com/grokify/mogo/type/stringsutil"
)

const (
	JiraIssueKeyURL = `%s/browse/%s`
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

func BuildIssueKeyURL(baseURL, key string) string {
	key = strings.TrimSpace(key)
	if urlutil.IsHTTP(key, true, true) {
		return key
	}
	return fmt.Sprintf(JiraIssueKeyURL, baseURL, key)
}

func BuildIssueKeyURLMarkdown(baseURL, key string) string {
	keyURL := BuildIssueKeyURL(baseURL, key)
	return markdown.Linkify(keyURL, key)
}
