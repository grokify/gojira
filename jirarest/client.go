package jirarest

import (
	"errors"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/net/http/httpsimple"
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
		c.Inflate()
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
	IssueAPI     *IssueAPI
}

func (c *Client) Inflate() {
	c.IssueAPI = &IssueAPI{Client: c}
}
