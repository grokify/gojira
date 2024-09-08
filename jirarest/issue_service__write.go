package jirarest

import (
	"errors"
	"net/http"
	"strings"

	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
)

type IssuePatchRequestBody struct {
	Update *IssuePatchRequestBodyUpdate          `json:"update,omitempty"`
	Fields map[string]IssuePatchRequestBodyField `json:"fields,omitempty"`
}

type IssuePatchRequestBodyUpdate struct {
	Labels []IssuePatchRequestBodyUpdateLabel `json:"labels,omitempty"`
}

type IssuePatchRequestBodyUpdateLabel struct {
	// cannot have both
	Add    *string `json:"add,omitempty"`
	Remove *string `json:"remove,omitempty"`
}

func (body IssuePatchRequestBody) Validate() error {
	if body.Update != nil {
		if len(body.Update.Labels) > 0 {
			for _, l := range body.Update.Labels {
				if l.Add != nil && l.Remove != nil {
					return errors.New("label update cannot have both add and remove")
				}
			}
		}
	}
	return nil
}

// FieldPatchRequestObject can be used IssuePatchRequestBody.Fields
type IssuePatchRequestBodyField struct {
	Value string                      `json:"value"`
	Child *IssuePatchRequestBodyField `json:"child,omitempty"`
}

// IssuePatch updates fields for an issue. See more here:
// https://community.developer.atlassian.com/t/update-issue-custom-field-value-via-api-without-going-forge/71161
func (c *IssueService) IssuePatch(issueKeyOrID string, issueUpdateRequestBody IssuePatchRequestBody) (*http.Response, error) {
	if err := issueUpdateRequestBody.Validate(); err != nil {
		return nil, err
	} else if issueKeyOrID = strings.TrimSpace(issueKeyOrID); issueKeyOrID == "" {
		return nil, errors.New("issue key or id must be provided")
	} else if issueUpdateRequestBody.Update == nil && len(issueUpdateRequestBody.Fields) == 0 {
		return nil, errors.New("issue `update` or `fields` must be provided")
	} else {
		return c.Client.simpleClient.Do(httpsimple.Request{
			Method:   http.MethodPut, // This only updates certain fields but uses a PUT http method.
			URL:      urlutil.JoinAbsolute(APIV3URLIssue, issueKeyOrID),
			Body:     issueUpdateRequestBody,
			BodyType: httpsimple.BodyTypeJSON})
	}
}
