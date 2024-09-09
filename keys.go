package gojira

import (
	"regexp"

	"github.com/grokify/mogo/type/stringsutil"
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
