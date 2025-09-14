package jirarest

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/time/month"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"

	"github.com/grokify/gojira"
	"github.com/grokify/gojira/jirarest/apiv3"
)

/*
// SearchIssues returns all issues for a JQL query, automatically handling API pagination.
func (svc *IssueService) SearchIssuesDeprecated(jql string) (Issues, error) {
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
*/

// SearchIssues returns all issues for a JQL query using the V3 API endpoint /rest/api/3/search/jql.
// If retrieveAll is true, it will paginate through all results until no more issues are available.
func (svc *IssueService) SearchIssues(jql string, retrieveAll bool) (Issues, error) {
	if svc.Client == nil {
		return nil, ErrClientCannotBeNil
	}
	if svc.Client.simpleClient == nil {
		return nil, ErrSimpleClientCannotBeNil
	}
	if strings.TrimSpace(jql) == "" {
		return Issues{}, nil
	}

	var allIssues Issues
	//startAt := 0
	maxResults := MaxResults

	nextPageToken := ""

	for {
		query := map[string][]string{
			"jql":        {jql},
			"maxResults": {fmt.Sprintf("%d", maxResults)},
			"fields":     {"*all"},
			// "startAt":    {fmt.Sprintf("%d", startAt)},
		}
		if nextPageToken != "" {
			query["nextPageToken"] = []string{nextPageToken}
		}

		resp, err := svc.Client.simpleClient.Do(context.Background(), httpsimple.Request{
			Method: http.MethodGet,
			URL:    APIV3URLSearchJQL,
			Query:  query,
		})
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 300 {
			return nil, fmt.Errorf("jira api status code (%d)", resp.StatusCode)
		}

		// Read response body
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		/*
			fmt.Println(string(body))
			err = os.WriteFile("issues_response.json", body, 0600)
			logutil.FatalErr(err)
			panic("Z")
		*/

		// Parse the response using the V3 typed structs
		var v3Response apiv3.IssuesResponse
		if err := json.Unmarshal(body, &v3Response); err != nil {
			return nil, err
		}

		// Convert V3 issues to go-jira Issues
		for _, v3Issue := range v3Response.Issues {
			goJiraIssue, err := v3Issue.ConvertToGoJiraIssue()
			if err != nil {
				return nil, fmt.Errorf("failed to convert V3 issue %s: %w", v3Issue.Key, err)
			}
			allIssues = append(allIssues, *goJiraIssue)
		}

		// If not retrieving all results, or we've got all the results, break
		/*
			if !retrieveAll || len(v3Response.Issues) == 0 || v3Response.StartAt+len(v3Response.Issues) >= v3Response.Total {
				break
			}
		*/

		nextPageToken = strings.TrimSpace(v3Response.NextPageToken)
		if v3Response.IsLast || nextPageToken == "" {
			break
		}
	}

	return allIssues, nil
}

func (svc *IssueService) SearchChildrenIssues(parentKeys []string) (Issues, error) {
	if parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true); len(parentKeys) == 0 {
		return Issues{}, errors.New("parentKeys cannot be empty")
	} else {
		jqlInfo := gojira.JQL{ParentsIncl: [][]string{parentKeys}}
		return svc.SearchIssues(jqlInfo.String(), true)
	}
}

func (svc *IssueService) SearchChildrenIssuesSet(ctx context.Context, parentKeys []string, opts *GetQueryOptions) (*IssuesSet, error) {
	// func (svc *IssueService) SearchChildrenIssuesSet(recursive, inclParents bool, parentKeys ...string) (*IssuesSet, error) {
	parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true)
	is := NewIssuesSet(svc.Client.Config)
	if len(parentKeys) == 0 {
		return is, nil
	}

	if opts.XIncludeParents {
		if iss, err := svc.Issues(ctx, parentKeys, opts); err != nil {
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
	for !slices.Equal(maputil.Keys(seen), maputil.Keys(is.IssuesMap)) {
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
		ii, err := svc.SearchIssues(jql, true)
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
	if jql = strings.TrimSpace(jql); jql == "" {
		return 0, nil
	} else if svc.Client == nil {
		return -1, ErrClientCannotBeNil
	} else if svc.Client.JiraClient == nil {
		return -1, ErrJiraClientCannotBeNil
	} else if _, resp, err := svc.Client.JiraClient.Issue.Search(jql, &jira.SearchOptions{
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
		StartAt:    offset}

	i := 0
	for maxPages == 0 || i < maxPages {
		ii, resp, err := svc.Client.JiraClient.Issue.Search(jql, &so)
		if err != nil {
			return issues, err
		} else if resp.StatusCode >= 300 {
			return issues, fmt.Errorf("jira api status code (%d)", resp.StatusCode)
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
		jql.CreatedGTE = &createdGTE
		jql.CreatedLT = pointer.Pointer(month.MonthStart(createdGTE, 1))
		if ii, err := svc.SearchIssues(jql.String(), true); err != nil {
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
	if ii, err := svc.SearchIssues(jql, true); err != nil {
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

	parIssuesSet := NewIssuesSet(set.Config)
	parIDs := set.KeysParentsUnpopulated()

	iter := 0
	for len(parIDs) > 0 {
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
			return nil, err
		}
		iter++
		if iter > 1000 {
			return nil, errors.New("search for parents over 1000 iterations")
		}
	}

	return parIssuesSet, nil
}

// JQLsAddMetadata returns all issues for a JQL query, automatically handling API pagination.
func (svc *IssueService) JQLsAddMetadata(jqls gojira.JQLs) (gojira.JQLs, error) {
	if len(jqls) == 0 {
		return jqls, nil
	}
	for i, j := range jqls {
		if count, err := svc.JQLResultsTotalCount(j.String()); err != nil {
			return jqls, err
		} else {
			j.Meta.QueryTime = time.Now()
			j.Meta.QueryTotalCount = count
			jqls[i] = j
		}
	}
	return jqls, nil
}
