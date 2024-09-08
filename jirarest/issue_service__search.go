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
func (c *IssueService) SearchIssues(jql string) (Issues, error) {
	var issues Issues

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	// SearchPages will page through results and pass each issue to appendFunc
	// In this example, we'll search for all the issues in the target project
	err := c.Client.JiraClient.Issue.SearchPages(jql, &jira.SearchOptions{Expand: "epic"}, appendFunc)
	return issues, err
}

func (c *IssueService) SearchChildrenIssues(parentKeys []string) (Issues, error) {
	if parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true); len(parentKeys) == 0 {
		return Issues{}, errors.New("parentKeys cannot be empty")
	} else {
		jqlInfo := gojira.JQL{ParentsIncl: [][]string{parentKeys}}
		return c.SearchIssues(jqlInfo.String())
	}
}

func (c *IssueService) SearchChildrenIssuesSet(recursive bool, parentKeys ...string) (*IssuesSet, error) {
	is := NewIssuesSet(c.Client.Config)
	seen := map[string]int{}
	seen, err := searchChildrenIssuesSetInternal(c, is, parentKeys, seen)
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
		seen, err = searchChildrenIssuesSetInternal(c, is, unseen, seen)
		if err != nil {
			return nil, err
		}
		i++
		if i >= recurseLimit {
			break
		}
	}

	return is, nil
}

func searchChildrenIssuesSetInternal(c *IssueService, is *IssuesSet, parentKeys []string, seen map[string]int) (map[string]int, error) {
	if ii, err := c.SearchChildrenIssues(parentKeys); err != nil {
		return seen, err
	} else if len(ii) > 0 {
		if err := is.Add(ii...); err != nil {
			return seen, err
		}
	}
	for _, pk := range parentKeys {
		seen[pk]++
	}
	return seen, nil
}

func (c *IssueService) SearchIssuesMulti(jqls ...string) (Issues, error) {
	var issues Issues
	for i, jql := range jqls {
		ii, err := c.SearchIssues(jql)
		if err != nil {
			return issues, err
		}
		issues = append(issues, ii...)
		if c.Client.LoggerZ != nil {
			c.Client.LoggerZ.Info().
				Str("jql", jql).
				Int("index", i).
				Int("totalQueries", len(jqls)).
				Int("totalIssues", len(issues)).
				Msg("jira api iteration")
		}
	}
	return issues, nil
}

// SearchIssuesPage returns all issues for a JQL query, automatically handling API pagination.
// A `limit` value of `0` means the max results available. A `maxPages` of `0` means to retrieve
// all pages.
func (c *IssueService) SearchIssuesPages(jql string, limit, offset, maxPages uint) (Issues, error) {
	var issues Issues

	if limit == 0 {
		limit = gojira.JQLMaxResults
	}

	so := jira.SearchOptions{
		MaxResults: int(limit),
		StartAt:    int(offset),
	}

	i := uint(0)
	for {
		if maxPages > 0 && i >= maxPages {
			break
		}
		ii, resp, err := c.Client.JiraClient.Issue.Search(jql, &so)
		if err != nil {
			return issues, err
		} else if resp.Response.StatusCode >= 300 {
			return issues, fmt.Errorf("jira api status code (%d)", resp.Response.StatusCode)
		}
		if c.Client.LoggerZ != nil {
			c.Client.LoggerZ.Info().
				Int("iteration", int(i)).
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

func (c *IssueService) SearchIssuesByMonth(jql gojira.JQL, createdGTE, createdLT time.Time, fnExec func(ii Issues, start time.Time) error) error {
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
		if ii, err := c.SearchIssues(jql.String()); err != nil {
			return err
		} else if err := fnExec(ii, createdGTE); err != nil {
			return err
		} else {
			createdGTE = month.MonthStart(createdGTE, 1)
		}
	}
	return nil
}

func (c *IssueService) SearchIssuesSet(jql string) (*IssuesSet, error) {
	if ii, err := c.SearchIssues(jql); err != nil {
		return nil, err
	} else {
		is := NewIssuesSet(c.Client.Config)
		err = is.Add(ii...)
		return is, err
	}
}

func (c *IssueService) SearchIssuesSetParents(is *IssuesSet) (*IssuesSet, error) {
	if is == nil {
		return nil, errors.New("issues set must be set")
	}
	// func (is *IssuesSet) RetrieveParentsIssuesSet(client *Client) (*IssuesSet, error) {
	parIssuesSet := NewIssuesSet(is.Config)
	parIDs := is.KeysParentsUnpopulated()

	iter := 0
	for {
		if len(parIDs) == 0 {
			break
		}

		if c.Client.LoggerZ != nil {
			c.Client.LoggerZ.Info().
				Int("iteration", iter).
				Msg("jira api populate parents (SearchIssuesSetParents)")
		}

		if parIssues, err := c.GetIssuesSetForKeys(parIDs); err != nil {
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
