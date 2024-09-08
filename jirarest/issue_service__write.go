package jirarest

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/net/http/httpsimple"
	"github.com/grokify/mogo/net/urlutil"
	"github.com/grokify/mogo/pointer"
	"github.com/grokify/mogo/type/stringsutil"
)

type IssuePatchRequestBody struct {
	Update *IssuePatchRequestBodyUpdate          `json:"update,omitempty"`
	Fields map[string]IssuePatchRequestBodyField `json:"fields,omitempty"`
}

func NewIssuePatchRequestBodyLabelAllRemove(label string, remove bool) IssuePatchRequestBody {
	labelUpdate := IssuePatchRequestBodyUpdateLabel{}
	if remove {
		labelUpdate.Remove = pointer.Pointer(label)
	} else {
		labelUpdate.Add = pointer.Pointer(label)
	}
	return IssuePatchRequestBody{
		Update: &IssuePatchRequestBodyUpdate{
			Labels: []IssuePatchRequestBodyUpdateLabel{
				labelUpdate,
			},
		},
	}
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

// IssuePatchLabelRecursive updates fields for an issue. See more here:
// https://community.developer.atlassian.com/t/update-issue-custom-field-value-via-api-without-going-forge/71161
func (c *IssueService) IssuePatchLabelRecursive(ctx context.Context, issueKeyOrID string, iss *jira.Issue, label string, removeLabel, processChildren bool, processChildrenTypes []string, skipUpdate bool) (int, error) {
	count := 0
	issueKeyOrID = strings.TrimSpace(issueKeyOrID)
	processChildrenTypes = stringsutil.SliceCondenseSpace(processChildrenTypes, true, true)
	var labelOperation string // for logging
	if removeLabel {
		labelOperation = OperationRemove
	} else {
		labelOperation = OperationAdd
	}
	if issueKeyOrID == "" {
		return count, errors.New("issue key or id must be supplied")
	}
	if iss == nil {
		if issGet, err := c.Issue(issueKeyOrID); err != nil {
			return 0, err
		} else {
			iss = issGet
		}
	}
	im := NewIssueMore(iss)

	labelExists := false
	if im.LabelExists(label) {
		labelExists = true
	}
	c.Client.LogOrNotAny(
		ctx,
		slog.LevelDebug,
		"validating issue labels",
		"issueKey", issueKeyOrID,
		"issueType", im.Type(),
		"labelAction", labelOperation,
		"label", label,
		"labelExists", labelExists)

	if !skipUpdate &&
		((removeLabel && im.LabelExists(label)) || (!removeLabel && !im.LabelExists(label))) {
		c.Client.LogOrNotAny(
			ctx,
			slog.LevelDebug,
			"updating issue labels",
			"issueKey", issueKeyOrID,
			"issueType", im.Type(),
			"labelAction", labelOperation,
			"label", label)
		reqBody := NewIssuePatchRequestBodyLabelAllRemove(label, removeLabel)
		if resp, err := c.IssuePatch(issueKeyOrID, reqBody); err != nil {
			return 0, nil
		} else if resp.StatusCode >= 300 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				body = []byte(err.Error())
			}
			return 0, fmt.Errorf("key (%s) status code (%d) bodyOrErr (%s)", im.Key(), resp.StatusCode, string(body))
		} else {
			count++
		}
	}
	if processChildren {
		ii, err := c.SearchChildrenIssues([]string{issueKeyOrID})
		if err != nil {
			return count, err
		}
		is := NewIssuesSet(nil)
		is.Add(ii...)
		if len(processChildrenTypes) > 0 {
			is, err = is.FilterType(processChildrenTypes...)
			if err != nil {
				return count, err
			}
		}
		if c.Client.LoggerZ != nil {
			c.Client.LoggerZ.Info().
				Str("issueKey", issueKeyOrID).
				Int("childrenCount", int(is.Len())).
				Msg("processing label update children count")
		}
		c.Client.LogOrNotAny(
			ctx,
			slog.LevelInfo,
			"processing label update children count",
			"issueKey", issueKeyOrID,
			"issueType", im.Type(),
			"labelAction", labelOperation,
			"label", label,
			"childrenCount", is.Len())
		issKeys := is.Keys()
		for _, issKey := range issKeys {
			cISS, err := is.Get(issKey)
			if err != nil {
				return count, err
			}
			// for _, cISS := range is.IssuesMap {
			im := NewIssueMore(&cISS)
			if countChildren, err := c.IssuePatchLabelRecursive(ctx, im.Key(), &cISS, label, removeLabel, processChildren, processChildrenTypes, skipUpdate); err != nil {
				return count, err
			} else {
				count += countChildren
			}
		}
	}
	return count, nil
}
