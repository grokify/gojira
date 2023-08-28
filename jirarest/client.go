package jirarest

import (
	"errors"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
)

func UserPassCredsBasic(filename, credsKey string) (*goauth.CredentialsBasicAuth, error) {
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

func ClientsBasicAuthFile(filename, credsKey string) (*http.Client, *jira.Client, string, error) {
	hclient, serverURL, err := HTTPClientBasicAuthFile(filename, credsKey)
	if err != nil {
		return nil, nil, "", errorsutil.Wrapf(err, `jirarest.ClientsBasicAuthFile() (%s)`, filename)
	}
	jclient, err := JiraClientBasicAuthFile(filename, credsKey)
	if err != nil {
		return hclient, jclient, serverURL, errorsutil.Wrap(err, `jirarest.ClientsBasicAuthFile()..JiraClientBasicAuthFile()`)
	}
	return hclient, jclient, serverURL, nil
}

func HTTPClientBasicAuthFile(filename, credsKey string) (hclient *http.Client, serverURL string, err error) {
	creds, err := UserPassCredsBasic(filename, credsKey)
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

func JiraClientBasicAuthFile(filename, credsKey string) (*jira.Client, error) {
	creds, err := UserPassCredsBasic(filename, credsKey)
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

// SearchIssues returns all issues for a JQL query, automatically handling API pagination.
func SearchIssues(client *jira.Client, jql string) (Issues, error) {
	var issues Issues

	// appendFunc will append jira issues to []jira.Issue
	appendFunc := func(i jira.Issue) (err error) {
		issues = append(issues, i)
		return err
	}

	// SearchPages will page through results and pass each issue to appendFunc
	// In this example, we'll search for all the issues in the target project
	err := client.Issue.SearchPages(jql, &jira.SearchOptions{Expand: "epic"}, appendFunc)
	return issues, err
}

func SearchIssuesSetForJQL(client *jira.Client, jql string, cfg *gojira.Config) (*IssuesSet, error) {
	ii, err := SearchIssues(client, jql)
	if err != nil {
		return nil, err
	}
	is := NewIssuesSet(cfg)
	err = is.Add(ii...)
	return is, err
}

func GetIssuesSetForKeys(client *jira.Client, keys []string) (*IssuesSet, error) {
	is := NewIssuesSet(nil)
	jql := KeysJQL(keys)
	if jql == "" {
		return is, nil
	}
	ii, err := SearchIssues(client, jql)
	if err != nil {
		return is, nil
	}

	err = is.Add(ii...)
	return is, err
}
