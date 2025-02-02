package jirarest

import (
	"errors"
	"fmt"
	"slices"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
)

// SearchIssues returns all issues for a JQL query, automatically handling API pagination.
func (svc *IssueService) SearchIssues(jql string) (Issues, error) {
	var issues Issues

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	// SearchPages will page through results and pass each issue to appendFunc
	// In this example, we'll search for all the issues in the target project
	err := svc.Client.JiraClient.Issue.SearchPages(jql, &jira.SearchOptions{Expand: "epic"}, appendFunc)
	return issues, err
}

func (svc *IssueService) SearchChildrenIssues(parentKeys []string) (Issues, error) {
	if parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true); len(parentKeys) == 0 {
		return Issues{}, errors.New("parentKeys cannot be empty")
	} else {
		jqlInfo := gojira.JQL{ParentsIncl: [][]string{parentKeys}}
		return svc.SearchIssues(jqlInfo.String())
	}
}

func (svc *IssueService) SearchChildrenIssuesSet(recursive, inclParents bool, parentKeys ...string) (*IssuesSet, error) {
	parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true)
	is := NewIssuesSet(svc.Client.Config)
	if len(parentKeys) == 0 {
		return is, nil
	}
	if inclParents {
		if iss, err := svc.Issues(parentKeys, false); err != nil {
			return nil, err
		} else if (len(iss)) == 0 {
			return nil, fmt.Errorf("no issues found for (%d) keys", len(parentKeys))
		} else if err := is.Add(iss...); err != nil {
			return nil, err
		}
	}
	seen := map[string]int{}
	seen, err := searchChildrenIssuesSetInternal(svc, is, parentKeys, seen)
	if err != nil {
		return nil, err
	}
	i := 0
	recurseLimit := 1000
	for {
		if slices.Equal(maputil.Keys(seen), maputil.Keys(is.IssuesMap)) {
			break
		}
		unseen := slicesutil.Sub(maputil.Keys(is.IssuesMap), maputil.Keys(seen))
		if len(unseen) == 0 {
			break
		}
		seen, err = searchChildrenIssuesSetInternal(svc, is, unseen, seen)
		if err != nil {
			return nil, err
		}
		i++
		if i >= recurseLimit {
			return is, fmt.Errorf("recurse limit of %d reached", recurseLimit)
		}
	}
	return is, nil
}

func searchChildrenIssuesSetInternal(svc *IssueService, set *IssuesSet, parentKeys []string, seen map[string]int) (map[string]int, error) {
	if ii, err := svc.SearchChildrenIssues(parentKeys); err != nil {
		return seen, err
	} else if len(ii) > 0 {
		if err := set.Add(ii...); err != nil {
			return seen, err
		}
	}
	for _, pk := range parentKeys {
		seen[pk]++
	}
	return seen, nil
}

func (svc *IssueService) SearchIssuesMulti(jqls ...string) (Issues, error) {
	var issues Issues
	for i, jql := range jqls {
		ii, err := svc.SearchIssues(jql)
		if err != nil {
			return issues, err
		}
		issues = append(issues, ii...)
		if svc.Client.LoggerZ != nil {
			svc.Client.LoggerZ.Info().
				Str("jql", jql).
				Int("index", i).
				Int("totalQueries", len(jqls)).
				Int("totalIssues", len(issues)).
				Msg("jira api iteration")
		}
	}
	return issues, nil
}

func (svc *IssueService) JQLResultsTotalCount(jql string) (int, error) {
	if _, resp, err := svc.Client.JiraClient.Issue.Search(jql, &jira.SearchOptions{
		MaxResults: 1,
		StartAt:    0,
	}); err != nil {
		return -1, err
	} else {
		return resp.Total, nil
	}
}

// SearchIssuesPage returns all issues for a JQL query, automatically handling API pagination.
// A `limit` value of `0` means the max results available. A `maxPages` of `0` means to retrieve
// all pages.
func (svc *IssueService) SearchIssuesPages(jql string, limit, offset, maxPages int) (Issues, error) {
	var issues Issues
	if limit < 0 || offset < 0 || maxPages < 0 {
		return issues, errors.New("limit, offset, and maxPages cannot be negative")
	}

	if limit == 0 {
		limit = gojira.JQLMaxResults
	}

	so := jira.SearchOptions{
		MaxResults: limit,
		StartAt:    offset,
	}

	i := 0
	for {
		if maxPages > 0 && i >= maxPages {
			break
		}
		ii, resp, err := svc.Client.JiraClient.Issue.Search(jql, &so)
		if err != nil {
			return issues, err
		} else if resp.Response.StatusCode >= 300 {
			return issues, fmt.Errorf("jira api status code (%d)", resp.Response.StatusCode)
		}
		if svc.Client.LoggerZ != nil {
			svc.Client.LoggerZ.Info().
				Int("iteration", i).
				Int("limit", resp.MaxResults).
				Int("offset", resp.StartAt).
				Int("total", resp.Total).
				Str("jql", jql).
				Msg("jira api iteration (SearchIssuesPages)")
		}
		if len(ii) > 0 {
			issues = append(issues, ii...)
		}
		if resp.StartAt+len(ii) >= resp.Total {
			break
		}
		so.StartAt += len(ii)
		i++
	}

	return issues, nil
}

func (svc *IssueService) SearchIssuesByMonth(jql gojira.JQL, createdGTE, createdLT time.Time, fnExec func(ii Issues, start time.Time) error) error {
	if createdGTE.IsZero() {
		createdGTE = month.MonthStart(time.Now(), 0)
	}
	if createdLT.IsZero() {
		createdLT = month.MonthStart(time.Now(), 1)
	}
	if createdGTE.Equal(createdLT) {
		createdLT = month.MonthStart(createdLT, 1)
	} else if createdLT.Before(createdGTE) {
		timeSwap := createdGTE
		createdGTE = createdLT
		createdLT = timeSwap
	}
	for createdGTE.Before(createdLT) {
		jql.CreatedGTE = createdGTE
		jql.CreatedLT = month.MonthStart(createdGTE, 1)
		// fmt.Printf("JQL [%s]\n", jql.String())
		if ii, err := svc.SearchIssues(jql.String()); err != nil {
			return err
		} else if err := fnExec(ii, createdGTE); err != nil {
			return err
		} else {
			createdGTE = month.MonthStart(createdGTE, 1)
		}
	}
	return nil
}

func (svc *IssueService) SearchIssuesSet(jql string) (*IssuesSet, error) {
	if ii, err := svc.SearchIssues(jql); err != nil {
		return nil, err
	} else {
		is := NewIssuesSet(svc.Client.Config)
		err = is.Add(ii...)
		return is, err
	}
}

func (svc *IssueService) SearchIssuesSetParents(set *IssuesSet) (*IssuesSet, error) {
	if set == nil {
		return nil, ErrIssuesSetCannotBeNil
	}
	// func (set *IssuesSet) RetrieveParentsIssuesSet(client *Client) (*IssuesSet, error) {
	parIssuesSet := NewIssuesSet(set.Config)
	parIDs := set.KeysParentsUnpopulated()

	iter := 0
	for {
		if len(parIDs) == 0 {
			break
		}

		if svc.Client.LoggerZ != nil {
			svc.Client.LoggerZ.Info().
				Int("iteration", iter).
				Msg("jira api populate parents (SearchIssuesSetParents)")
		}

		if parIssues, err := svc.GetIssuesSetForKeys(parIDs); err != nil {
			return nil, err
		} else if err := parIssuesSet.Add(parIssues.Issues()...); err != nil {
			return nil, err
		} else if parIDs, err = parIssuesSet.LineageTopKeysUnpopulated(); err != nil {
			// err = parIssuesSet.RetrieveParents(c) // don't use - use lineage instead.
			return nil, err
		}
		iter++
		if iter > 1000 {
			return nil, errors.New("search for parents over 1000 iterations")
		}
	}

	return parIssuesSet, nil
}
