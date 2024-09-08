package jirarest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

type IssuePatchRequestBody struct {
	Fields map[string]any `json:"fields"`
}

// IssuePatch updates fields for an issue. See more here:
// https://community.developer.atlassian.com/t/update-issue-custom-field-value-via-api-without-going-forge/71161
func (c *IssueAPI) IssuePatch(issueKeyOrID string, issueUpdateRequestBody IssuePatchRequestBody) (*http.Response, error) {
	if issueKeyOrID = strings.TrimSpace(issueKeyOrID); issueKeyOrID == "" {
		return nil, errors.New("issue key or id must be provided")
	} else if len(issueUpdateRequestBody.Fields) == 0 {
		return nil, errors.New("issue fields must be provided")
	} else {
		return c.Client.simpleClient.Do(httpsimple.Request{
			Method:   http.MethodPut, // This only updates certain fields but uses a PUT http method.
			URL:      urlutil.JoinAbsolute(APIV3URLIssue, issueKeyOrID),
			Body:     issueUpdateRequestBody,
			BodyType: httpsimple.BodyTypeJSON})
	}
}
