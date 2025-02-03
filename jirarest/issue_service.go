package jirarest

import (
	"errors"
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

type IssueService struct {
	Client *Client
}

func NewIssueService(client *Client) *IssueService {
	return &IssueService{Client: client}
}

func (svc *IssueService) Issue(key string) (*jira.Issue, error) {
	key = strings.TrimSpace(key)
	jqlInfo := gojira.JQL{IssuesIncl: [][]string{{key}}}
	//jql := fmt.Sprintf("issue = %s", key)
	jql := jqlInfo.String()
	if key == "" {
		return nil, errors.New("issue key is required")
	} else if iss, err := svc.SearchIssues(jql); err != nil {
		return nil, err
	} else if len(iss) == 0 {
		return nil, fmt.Errorf("key not found (%s)", key)
	} else if len(iss) > 1 {
		return nil, fmt.Errorf("too many issues (%d) found for (%s)", len(iss), key)
	} else {
		return &iss[0], nil
	}
}

// Issues returns a list of issues given a set of keys. If no keys are provided,
// any empty slice is returned. `skipNotFound` is useful if Jira ticket key no
// longer exists.
func (svc *IssueService) Issues(keys []string, skipNotFound bool) (Issues, error) {
	keys = stringsutil.SliceCondenseSpace(keys, true, true)
	iss := Issues{}
	if len(keys) == 0 {
		return iss, nil
	}
	for _, key := range keys {
		if is, err := svc.Issue(key); err != nil {
			if skipNotFound && err.Error() == ErrNotFound.Error() {
				continue
			} else {
				return iss, err
			}
		} else {
			iss = append(iss, *is)
		}
	}
	return iss, nil
}

// Issues returns an `IssuesSet{}` given a set of keys. If no keys are provided,
// any empty slice is returned.
func (svc *IssueService) GetIssuesSetForKeys(keys []string) (*IssuesSet, error) {
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
		} else if ii, err := svc.SearchIssuesPages(jql, 0, 0, 0); err != nil {
			return nil, err
		} else if err = is.Add(ii...); err != nil {
			return nil, err
		}
	}

	return is, nil
}

func (svc *IssueService) WebURL() string {
	issuesWebURL := ""
	if svc.Client != nil && svc.Client.Config != nil {
		if svrURL := strings.TrimSpace(svc.Client.Config.ServerURL); svrURL != "" {
			issuesWebURL = urlutil.JoinAbsolute(svrURL, "issues/?")
		}
	}
	return issuesWebURL
}
