package jirarest

// accountId=712020:c043668b-904d-4ecc-bde1-990ce20bd437

import (
	"encoding/json"
	"io"
	"net/http"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

const APIURLMyself = "/rest/api/3/myself"

func (c *Client) Myself() (*jira.User, *http.Response, error) {
	usr := jira.User{}
	if resp, err := c.simpleClient.Do(httpsimple.Request{
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
