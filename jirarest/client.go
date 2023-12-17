package jirarest

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/rs/zerolog"
)

func NewClientGoauthBasicAuthFile(filename, credsKey string) (*Client, error) {
	c := &Client{}
	hclient, serverURL, err := NewClientHTTPBasicAuthFile(filename, credsKey)
	if err != nil {
		return nil, errorsutil.Wrapf(err, `jirarest.ClientsBasicAuthFile() (%s)`, filename)
	}
	c.HTTPClient = hclient
	cfg := gojira.NewConfigDefault()
	cfg.ServerURL = serverURL
	c.Config = *cfg
	jclient, err := NewClientJiraBasicAuthFile(filename, credsKey)
	if err != nil {
		return c, errorsutil.Wrap(err, `jirarest.ClientsBasicAuthFile()..JiraClientBasicAuthFile()`)
	}
	c.JiraClient = jclient
	return c, nil
}

func NewCredentialsBasicAuthGoauthFile(filename, credsKey string) (*goauth.CredentialsBasicAuth, error) {
	// func UserPassCredsBasic(filename, credsKey string) (*goauth.CredentialsBasicAuth, error) {
	cs, err := goauth.ReadFileCredentialsSet(filename, true)
	if err != nil {
		return nil, err
	}

	creds, err := cs.Get(credsKey)
	if err != nil {
		return nil, err
	}

	return creds.Basic, nil
}

func NewClientHTTPBasicAuthFile(filename, credsKey string) (hclient *http.Client, serverURL string, err error) {
	creds, err := NewCredentialsBasicAuthGoauthFile(filename, credsKey)
	if err != nil {
		return nil, "", err
	}
	hclient, err = creds.NewClient()
	if err != nil {
		return hclient, "", err
	}
	serverURL = creds.ServerURL
	return
}

func NewClientJiraBasicAuthFile(filename, credsKey string) (*jira.Client, error) {
	creds, err := NewCredentialsBasicAuthGoauthFile(filename, credsKey)
	if err != nil {
		return nil, err
	}
	return JiraClientBasicAuth(creds)
}

func JiraClientBasicAuth(creds *goauth.CredentialsBasicAuth) (*jira.Client, error) {
	if creds == nil {
		return nil, errors.New("goauth.CredentialsBasicAuth cannot be nil")
	}
	tp := jira.BasicAuthTransport{
		Username: creds.Username,
		Password: creds.Password,
	}
	return jira.NewClient(tp.Client(), creds.ServerURL)
}

type Client struct {
	Config     gojira.Config
	HTTPClient *http.Client
	JiraClient *jira.Client
	Logger     *zerolog.Logger
}

func (c *Client) Issue(key string) (*jira.Issue, error) {
	key = strings.TrimSpace(key)
	jqlInfo := gojira.JQL{IssuesIncl: []string{key}}
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

func (c *Client) SearchIssuesMulti(jqls ...string) (Issues, error) {
	var issues Issues
	for i, jql := range jqls {
		ii, err := c.SearchIssues(jql)
		if err != nil {
			return issues, err
		}
		issues = append(issues, ii...)
		if c.Logger != nil {
			c.Logger.Info().
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
func (c *Client) SearchIssuesPages(jql string, limit, offset, maxPages uint) (Issues, error) {
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
		ii, resp, err := c.JiraClient.Issue.Search(jql, &so)
		if err != nil {
			return issues, err
		} else if resp.Response.StatusCode >= 300 {
			return issues, fmt.Errorf("jira api status code (%d)", resp.Response.StatusCode)
		}
		if c.Logger != nil {
			c.Logger.Info().
				Int("limit", resp.MaxResults).
				Int("offset", resp.StartAt).
				Int("total", resp.Total).
				Msg("jira api iteration")
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

// SearchIssues returns all issues for a JQL query, automatically handling API pagination.
func (c *Client) SearchIssues(jql string) (Issues, error) {
	var issues Issues

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	// SearchPages will page through results and pass each issue to appendFunc
	// In this example, we'll search for all the issues in the target project
	err := c.JiraClient.Issue.SearchPages(jql, &jira.SearchOptions{Expand: "epic"}, appendFunc)
	return issues, err
}

func (c *Client) SearchIssuesSet(jql string) (*IssuesSet, error) {
	ii, err := c.SearchIssues(jql)
	if err != nil {
		return nil, err
	}
	is := NewIssuesSet(&c.Config)
	err = is.Add(ii...)
	return is, err
}

func (c *Client) GetIssuesSetForKeys(keys []string) (*IssuesSet, error) {
	is := NewIssuesSet(nil)
	jql := KeysJQL(keys)
	if jql == "" {
		return is, nil
	}
	ii, err := c.SearchIssues(jql)
	if err != nil {
		return is, nil
	}

	err = is.Add(ii...)
	return is, err
}

func (c *Client) SearchIssuesSetParents(is *IssuesSet) (*IssuesSet, error) {
	// func (is *IssuesSet) RetrieveParentsIssuesSet(client *Client) (*IssuesSet, error) {
	parIssuesSet := NewIssuesSet(is.Config)
	parIDs := is.UnknownParents()
	if len(parIDs) == 0 {
		return parIssuesSet, nil
	}

	err := parIssuesSet.RetrieveIssues(c, parIDs)
	if err != nil {
		return nil, err
	}

	err = parIssuesSet.RetrieveParents(c)

	return parIssuesSet, err
}

func (c *Client) IssuesSetAddParents(is *IssuesSet) error {
	if is == nil {
		return errors.New("issues set is nil")
	}
	parents, err := c.SearchIssuesSetParents(is)
	if err != nil {
		return err
	}
	is.Parents = parents
	return nil
}
