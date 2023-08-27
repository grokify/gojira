package jirarest

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/grokify/gojira"
	"github.com/grokify/mogo/net/http/httpsimple"
)

// project = foundation AND resolution = Unresolved AND status!=Closed AND (Sprint not in openSprints() OR Sprint is EMPTY) AND type not in (Epic, Sub-Task) ORDER BY Rank ASC
/*
func BacklogJQL(projectName string) string {
	parts := []string{
		"resolution = Unresolved",
		"status!=Closed",
		"(Sprint not in openSprints() OR Sprint is EMPTY)",
		"type not in (Epic, Sub-Task)",
	}
	projectName = strings.TrimSpace(projectName)
	if projectName != "" {
		// parts = []string{projectName, parts... }
	}
}
*/

func Shift[S ~[]E, E any](s S) (E, S) {
	if len(s) == 0 {
		return *new(E), []E{}
	}
	return s[0], s[1:]
}

/*
func Prepend[S ~[]E, E any](s []E        , e E) []E   {
	return append     {[]E{e }  ,   s...  }
}*/

const (
	ParamFields        = "fields"
	ParamJQL           = "jql"
	ParamMaxResults    = "maxResults"
	ParamStartAt       = "startAt"
	ParamValidateQuery = "validateQuery"
)

type BoardBacklogParams struct {
	StartAt       uint   `url:"startAt"`
	MaxResults    uint   `url:"maxResults"`
	JQL           string `url:"jql"`
	ValidateQuery bool   `url:"validateQuery"`
	Fields        string `url:"fields"`
	Expand        string `url:"expand"`
}

func (p BoardBacklogParams) URLValues() url.Values {
	u := url.Values{}
	if p.StartAt > 0 {
		u.Add(ParamStartAt, strconv.Itoa(int(p.StartAt)))
	}
	if p.MaxResults > 0 {
		u.Add(ParamMaxResults, strconv.Itoa(int(p.MaxResults)))
	}
	jql := strings.TrimSpace(p.JQL)
	if jql != "" {
		u.Add(ParamJQL, jql)
		if p.ValidateQuery {
			u.Add(ParamValidateQuery, "true")
		} else {
			u.Add(ParamValidateQuery, "false) ")
		}
	}
	fields := strings.TrimSpace(p.Fields)
	if fields != "" {
		u.Add(ParamFields, fields)
	}
	return u
}

// BacklogAPIURL returns a backlog issues API URL described at https://docs.atlassian.com/jira-software/REST/7.3.1/ .
// The description is here: Returns all issues from the board's backlog, for the given board Id. This only includes issues that the user has permission to view. The backlog contains incomplete issues that are not assigned to any future or active sprint. Note, if the user does not have permission to view the board, no issues will be returned at all. Issues returned from this resource include Agile fields, like sprint, closedSprints, flagged, and epic. By default, the returned issues are ordered by rank.
// Reference: https://docs.atlassian.com/jira-software/REST/7.3.1/#agile/1.0/board-getIssuesForBacklog
func BacklogAPIURL(baseURL string, boardID uint, qry *BoardBacklogParams) string {
	apiURL := baseURL + fmt.Sprintf(`/rest/agile/1.0/board/%d/backlog`, boardID)
	if qry != nil {
		v := qry.URLValues()
		q := strings.TrimSpace(v.Encode())
		if q != "" {
			apiURL += "?" + q
		}
	}
	return apiURL
}

type BacklogService struct {
	config  gojira.Config
	sclient httpsimple.SimpleClient
	// client    *http.Client
	// serverURL string
}

func NewBacklogService(client *http.Client, cfg *gojira.Config) *BacklogService {
	return &BacklogService{
		config: *cfg,
		sclient: httpsimple.SimpleClient{
			HTTPClient: client,
			BaseURL:    cfg.BaseURL}}
}

func (s *BacklogService) GetBacklogIssuesResponse(boardID uint, qry *BoardBacklogParams) (*IssuesResponse, []byte, error) {
	if s.sclient.HTTPClient == nil {
		return nil, []byte{}, errors.New("client not set")
	}
	sreq := httpsimple.SimpleRequest{
		Method: http.MethodGet,
		URL:    BacklogAPIURL("", boardID, nil),
		Query:  qry.URLValues(),
	}
	resp, err := s.sclient.Do(sreq)
	if err != nil {
		return nil, []byte{}, err
	} else if resp.StatusCode >= 300 {
		return nil, []byte{}, fmt.Errorf("statusCode (%d)", resp.StatusCode)
	}
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, []byte{}, err
	}
	ir, err := ParseIssuesResponseBytes(b)
	return ir, b, err
}

func (s *BacklogService) GetBacklogIssuesAll(boardID uint, jql string) (*IssuesResponse, [][]byte, error) {
	iragg := &IssuesResponse{}
	bb := [][]byte{}
	issues := Issues{}

	opts := &BoardBacklogParams{
		StartAt:    uint(0),
		MaxResults: MaxResults,
		JQL:        jql,
	}
	for {
		ir, b, err := s.GetBacklogIssuesResponse(boardID, opts)
		if err != nil {
			return iragg, bb, err
		} else {
			iragg.Total = ir.Total
			iragg.MaxResults = ir.Total
		}
		bb = append(bb, b)
		if len(ir.Issues) > 0 {
			issues = append(issues, ir.Issues...)
		} else {
			break
		}
		if opts.StartAt+opts.MaxResults > uint(ir.Total) {
			break
		} else {
			opts.StartAt += opts.MaxResults
		}
	}
	iragg.Issues = issues.AddRank()
	return iragg, bb, nil
}

func (s *BacklogService) GetBacklogIssuesSetAll(boardID uint, jql string) (*IssuesSet, [][]byte, error) {
	iir, b, err := s.GetBacklogIssuesAll(boardID, jql)
	if err != nil {
		return nil, b, err
	}
	is := NewIssuesSet(&s.config)
	is.Add(iir.Issues...)
	return is, b, err
}
