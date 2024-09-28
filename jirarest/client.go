package jirarest

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth"
	"github.com/grokify/goauth/authutil"
	"github.com/grokify/gojira"
	"github.com/grokify/mogo/errors/errorsutil"
	"github.com/grokify/mogo/log/slogutil"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/rs/zerolog"
)

var (
	ErrClientCannotBeNil     = errors.New("client cannot be nil")
	ErrJiraClientCannotBeNil = errors.New("jira client cannot be nil")
)

type Client struct {
	Config         *gojira.Config
	HTTPClient     *http.Client
	JiraClient     *jira.Client
	simpleClient   *httpsimple.Client
	LoggerZ        *zerolog.Logger
	Logger         *slog.Logger
	BacklogAPI     *BacklogService
	CustomFieldAPI *CustomFieldService
	IssueAPI       *IssueService
	CustomFieldSet *CustomFieldSet
}

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

func NewClientFromGoauthCredentials(c *goauth.Credentials) (*Client, error) {
	if c == nil {
		return nil, errors.New("goauth.Credentials cannot be nil")
	} else if c.Type == goauth.TypeBasic && c.Basic != nil {
		return NewClientBasicAuth(c.Basic.ServerURL, c.Basic.Username, c.Basic.Password)
	} else {
		return nil, errors.New("auth method not supported or popuated")
	}
}

func NewClientGoauthBasicAuthFile(filename, credsKey string, addCustomFieldSet bool) (*Client, error) {
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
		if err := c.Inflate(addCustomFieldSet); err != nil {
			return nil, err
		}
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

func (c *Client) Inflate(addCustomFieldSet bool) error {
	c.BacklogAPI = NewBacklogService(c)
	c.CustomFieldAPI = NewCustomFieldService(c)
	c.IssueAPI = NewIssueService(c)
	if addCustomFieldSet {
		if err := c.LoadCustomFields(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) LoadCustomFields() error {
	if c.CustomFieldAPI == nil {
		c.CustomFieldAPI = NewCustomFieldService(c)
	}
	if fields, err := c.CustomFieldAPI.GetCustomFields(); err != nil {
		return err
	} else if len(fields) > 0 {
		c.CustomFieldSet = NewCustomFieldSet()
		if err := c.CustomFieldSet.Add(fields...); err != nil {
			return err
		}
	}
	return nil
}

func (c *Client) LogOrNotAny(ctx context.Context, level slog.Level, msg string, attrs ...any) {
	if c.Logger != nil {
		slogutil.LogOrNotAny(ctx, c.Logger, level, msg, attrs...)
	}
}
