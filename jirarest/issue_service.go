package jirarest

import (
	"context"
	"errors"
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gojira"
)

type IssueService struct {
	Client *Client
}

func NewIssueService(client *Client) *IssueService {
	return &IssueService{Client: client}
}

type GetQueryOptions struct {
	ExpandChangelog    bool // sent to andygrunwald SDK
	XMultiSkipNotFound bool // not sent to andygrunwald SDK; used for getting multiple issues
	XIncludeParents    bool
	// XMultiRecursive    bool
}

// Build returns a `*jira.GetQueryOptions` for the andygrunwald SDK.
func (opts GetQueryOptions) Build() *jira.GetQueryOptions {
	out := &jira.GetQueryOptions{}
	if opts.ExpandChangelog {
		out.Expand = "changelog"
	}
	return out
}

func (svc *IssueService) Issue(ctx context.Context, issueIDOrKey string, opts *GetQueryOptions) (*jira.Issue, error) {
	issueIDOrKey = strings.TrimSpace(issueIDOrKey)

	var opts2 *jira.GetQueryOptions
	if opts != nil {
		opts2 = opts.Build()
	}
	if issueIDOrKey == "" {
		return nil, errors.New("issue key cannot be empty")
	} else if svc.Client == nil {
		return nil, errors.New("gojira.Client cannot be nil")
	} else if svc.Client.JiraClient == nil {
		return nil, errors.New("gojira.Client.JiraClient cannot be nil")
	} else if svc.Client.JiraClient.Issue == nil {
		return nil, errors.New("gojira.Client.JiraClient.issue cannot be nil")
	} else if iss, resp, err := svc.Client.JiraClient.Issue.GetWithContext(ctx, issueIDOrKey, opts2); err != nil {
		return nil, err
	} else if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("unsuccessful jira api http status code (%d)", resp.StatusCode)
	} else {
		return iss, nil
	}
	/*
		jql := gojira.JQL{IssuesIncl: [][]string{{key}}}
		if key == "" {
			return nil, errors.New("issue key is required")
		} else if iss, err := svc.SearchIssues(jql.String(), false); err != nil {
			return nil, err
		} else if len(iss) == 0 {
			return nil, fmt.Errorf("key not found (%s)", key)
		} else if len(iss) > 1 {
			return nil, fmt.Errorf("too many issues (%d) found for (%s)", len(iss), key)
		} else {
			return &iss[0], nil
		}
	*/
}

// Issues returns a list of issues given a set of keys. If no keys are provided,
// any empty slice is returned. `skipNotFound` is useful if Jira ticket key no
// longer exists.
func (svc *IssueService) Issues(ctx context.Context, keys []string, opts *GetQueryOptions) (Issues, error) {
	keys = stringsutil.SliceCondenseSpace(keys, true, true)
	iss := Issues{}
	if len(keys) == 0 {
		return iss, nil
	}
	if opts == nil {
		opts = &GetQueryOptions{}
	}
	for _, key := range keys {
		if is, err := svc.Issue(ctx, key, opts); err != nil {
			if err.Error() == ErrNotFound.Error() && opts.XMultiSkipNotFound {
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
		if jql := jqlInfo.String(); jql == "" {
			continue
		} else if ii, err := svc.SearchIssuesPages(jql, 0, 0, 0); err != nil {
			return nil, err
		} else if err = is.Add(ii...); err != nil {
			return nil, err
		}
	}

	return is, nil
}
