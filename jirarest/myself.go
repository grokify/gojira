package jirarest

// accountId=712020:c043668b-904d-4ecc-bde1-990ce20bd437

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/goauth/oidc"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

const APIURLMyself = "/rest/api/3/myself"

func (c *Client) Myself() (*jira.User, *http.Response, error) {
	usr := jira.User{}
	if c.simpleClient == nil {
		return nil, nil, ErrSimpleClientCannotBeNil
	} else if resp, err := c.simpleClient.Do(httpsimple.Request{
		Method: http.MethodGet,
		URL:    urlutil.JoinAbsolute(c.Config.ServerURL, APIURLMyself)}); err != nil {
		return nil, nil, err
	} else if b, err := io.ReadAll(resp.Body); err != nil {
		return nil, nil, err
	} else if err = json.Unmarshal(b, &usr); err != nil {
		return nil, nil, err
	} else {
		return &usr, resp, err
	}
}

func (c *Client) MyselfUserInfo(ctx context.Context) (*oidc.UserInfo, *jira.User, *jira.Response, error) {
	if c.JiraClient == nil {
		return nil, nil, nil, ErrJiraClientCannotBeNil
	} else if u, resp, err := c.JiraClient.User.GetSelfWithContext(ctx); err != nil {
		return nil, nil, nil, err
	} else if resp.StatusCode > 299 {
		return nil, nil, nil, fmt.Errorf("bad status code (%d)", resp.StatusCode)
	} else {
		return UserJiraToOIDC(u, c.Config.ServerURL), u, resp, nil
	}
}

func UserJiraToOIDC(u *jira.User, serverURL string) *oidc.UserInfo {
	if u == nil {
		return nil
	} else {
		ui := &oidc.UserInfo{
			Issuer:  serverURL,
			Picture: u.AvatarUrls.Four8X48}
		ui.AddEmail(u.EmailAddress, true)
		ui.AddName(u.DisplayName, true)
		return ui
	}
}
