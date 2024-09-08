package jirarest

import (
	"errors"
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type IssueService struct {
	Client *Client
}

func NewIssueService(client *Client) *IssueService {
	return &IssueService{Client: client}
}

func (c *IssueService) Issue(key string) (*jira.Issue, error) {
	key = strings.TrimSpace(key)
	jqlInfo := gojira.JQL{IssuesIncl: [][]string{{key}}}
	//jql := fmt.Sprintf("issue = %s", key)
	jql := jqlInfo.String()
	if key == "" {
		return nil, errors.New("issue key is required")
	} else if iss, err := c.SearchIssues(jql); err != nil {
		return nil, err
	} else if len(iss) == 0 {
		return nil, fmt.Errorf("key not found (%s)", key)
	} else if len(iss) > 1 {
		return nil, fmt.Errorf("too many issues (%d) found for (%s)", len(iss), key)
	} else {
		return &iss[0], nil
	}
}

func (c *IssueService) Issues(keys ...string) (Issues, error) {
	keys = stringsutil.SliceCondenseSpace(keys, true, true)
	iss := Issues{}
	if len(keys) == 0 {
		return iss, nil
	}
	for _, key := range keys {
		if is, err := c.Issue(key); err != nil {
			return iss, err
		} else {
			iss = append(iss, *is)
		}
	}
	return iss, nil
}

func (c *IssueService) GetIssuesSetForKeys(keys []string) (*IssuesSet, error) {
	is := NewIssuesSet(nil)

	keysSlice := slicesutil.SplitMaxLength(stringsutil.SliceCondenseSpace(keys, true, true), gojira.JQLMaxResults)

	for _, keysIter := range keysSlice {
		keysIter = stringsutil.SliceCondenseSpace(keysIter, true, true)
		if len(keysIter) == 0 {
			continue
		}
		jqlInfo := gojira.JQL{KeysIncl: [][]string{keysIter}}
		// jql := KeysJQL(keys)
		if jql := jqlInfo.String(); jql == "" {
			return is, nil
		} else if ii, err := c.SearchIssuesPages(jql, 0, 0, 0); err != nil {
			return nil, err
		} else if err = is.Add(ii...); err != nil {
			return nil, err
		}
	}

	return is, nil
}
