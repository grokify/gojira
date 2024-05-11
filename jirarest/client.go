package jirarest

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/type/maputil"
	"github.com/grokify/mogo/type/slicesutil"
	"github.com/grokify/mogo/type/stringsutil"
	"github.com/rs/zerolog"
)

var (
	ErrClientCannotBeNil     = errors.New("client cannot be nil")
	ErrJiraClientCannotBeNil = errors.New("jira client cannot be nil")
)

func NewClientBasicAuth(serverURL, username, password string) (*Client, error) {
	if hclient, err := authutil.NewClientBasicAuth(username, password, false); err != nil {
		return nil, err
	} else if jclient, err := JiraClientBasicAuth(serverURL, username, password); err != nil {
		return nil, err
	} else {
		c := &Client{
			HTTPClient: hclient,
			JiraClient: jclient}
		cfg := gojira.NewConfigDefault()
		cfg.ServerURL = serverURL
		c.Config = cfg
		return c, nil
	}
}

func NewClientGoauthBasicAuthFile(filename, credsKey string) (*Client, error) {
	if hclient, serverURL, err := NewClientHTTPBasicAuthFile(filename, credsKey); err != nil {
		return nil, errorsutil.Wrapf(err, `jirarest.ClientsBasicAuthFile() (%s)`, filename)
	} else if jclient, err := NewClientJiraBasicAuthFile(filename, credsKey); err != nil {
		return nil, errorsutil.Wrap(err, `jirarest.ClientsBasicAuthFile()..JiraClientBasicAuthFile()`)
	} else {
		c := &Client{
			HTTPClient: hclient,
			JiraClient: jclient}
		cfg := gojira.NewConfigDefault()
		cfg.ServerURL = serverURL
		sc := httpsimple.NewClient(hclient, serverURL)
		c.simpleClient = &sc
		c.Config = cfg
		return c, nil
	}
}

func NewCredentialsBasicAuthGoauthFile(filename, credsKey string) (*goauth.CredentialsBasicAuth, error) {
	// func UserPassCredsBasic(filename, credsKey string) (*goauth.CredentialsBasicAuth, error) {
	if cs, err := goauth.ReadFileCredentialsSet(filename, true); err != nil {
		return nil, err
	} else if creds, err := cs.Get(credsKey); err != nil {
		return nil, err
	} else {
		return creds.Basic, nil
	}
}

func NewClientHTTPBasicAuthFile(filename, credsKey string) (hclient *http.Client, serverURL string, err error) {
	if creds, err := NewCredentialsBasicAuthGoauthFile(filename, credsKey); err != nil {
		return nil, "", err
	} else if hclient, err = creds.NewClient(); err != nil {
		return hclient, "", err
	} else {
		serverURL = creds.ServerURL
	}
	return
}

func NewClientJiraBasicAuthFile(filename, credsKey string) (*jira.Client, error) {
	if creds, err := NewCredentialsBasicAuthGoauthFile(filename, credsKey); err != nil {
		return nil, err
	} else {
		return JiraClientBasicAuthGoauth(creds)
	}
}

func JiraClientBasicAuth(serverURL, username, password string) (*jira.Client, error) {
	tp := jira.BasicAuthTransport{
		Username: username,
		Password: password}
	return jira.NewClient(tp.Client(), serverURL)
}

func JiraClientBasicAuthGoauth(creds *goauth.CredentialsBasicAuth) (*jira.Client, error) {
	if creds == nil {
		return nil, errors.New("goauth.CredentialsBasicAuth cannot be nil")
	}
	return JiraClientBasicAuth(creds.ServerURL, creds.Username, creds.Password)
}

type Client struct {
	Config       *gojira.Config
	HTTPClient   *http.Client
	JiraClient   *jira.Client
	simpleClient *httpsimple.Client
	Logger       *zerolog.Logger
}

func (c *Client) Issue(key string) (*jira.Issue, error) {
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

func (c *Client) Issues(keys ...string) (Issues, error) {
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

func (c *Client) SearchChildrenIssues(parentKeys ...string) (Issues, error) {
	if parentKeys = stringsutil.SliceCondenseSpace(parentKeys, true, true); len(parentKeys) == 0 {
		return Issues{}, errors.New("parentKeys cannot be empty")
	} else {
		jqlInfo := gojira.JQL{ParentsIncl: [][]string{parentKeys}}
		return c.SearchIssues(jqlInfo.String())
	}
}

func (c *Client) SearchChildrenIssuesSet(recursive bool, parentKeys ...string) (*IssuesSet, error) {
	is := NewIssuesSet(c.Config)
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

func searchChildrenIssuesSetInternal(c *Client, is *IssuesSet, parentKeys []string, seen map[string]int) (map[string]int, error) {
	if ii, err := c.SearchChildrenIssues(parentKeys...); err != nil {
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
// A `limit` value of `0` means the max results available. A `maxPages` of `0` means to retrieve
// all pages.
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
	if ii, err := c.SearchIssues(jql); err != nil {
		return nil, err
	} else {
		is := NewIssuesSet(c.Config)
		err = is.Add(ii...)
		return is, err
	}
}

func (c *Client) GetIssuesSetForKeys(keys []string) (*IssuesSet, error) {
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

func (c *Client) SearchIssuesSetParents(is *IssuesSet) (*IssuesSet, error) {
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

		if c.Logger != nil {
			c.Logger.Info().
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

func (c *Client) IssuesSetAddParents(is *IssuesSet) error {
	if is == nil {
		return errors.New("issues set is nil")
	} else if parents, err := c.SearchIssuesSetParents(is); err != nil {
		return err
	} else {
		is.Parents = parents
		return nil
	}
}
