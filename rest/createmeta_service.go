package rest

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grokify/mogo/encoding/jsonutil"
	"github.com/grokify/mogo/net/urlutil"
)

// CreateMetaService provides access to Jira's createmeta API for discovering
// which fields are available for issue creation in specific projects.
type CreateMetaService struct {
	JRClient *Client
}

// NewCreateMetaService creates a new CreateMetaService.
func NewCreateMetaService(client *Client) *CreateMetaService {
	return &CreateMetaService{JRClient: client}
}

// GetIssueTypes returns available issue types for a project.
// Uses GET /rest/api/3/issue/createmeta/{projectKey}/issuetypes
func (svc *CreateMetaService) GetIssueTypes(ctx context.Context, projectKey string) ([]CreateMetaIssueType, error) {
	if svc.JRClient == nil {
		return nil, ErrJiraRESTClientCannotBeNil
	}

	apiURL := urlutil.JoinAbsolute(
		svc.JRClient.Config.ServerURL,
		APIV3URLCreateMeta,
		projectKey,
		"issuetypes",
	)

	hclient := svc.JRClient.HTTPClient
	if hclient == nil {
		hclient = &http.Client{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := hclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error status code (%d) for project %q", resp.StatusCode, projectKey)
	}

	var result CreateMetaIssueTypesResponse
	if _, err := jsonutil.UnmarshalReader(resp.Body, &result); err != nil {
		return nil, err
	}

	return result.IssueTypes, nil
}

// GetFields returns available fields for a project/issue-type combination.
// Uses GET /rest/api/3/issue/createmeta/{projectKey}/issuetypes/{issueTypeId}
func (svc *CreateMetaService) GetFields(ctx context.Context, projectKey, issueTypeID string) (CreateMetaFields, error) {
	if svc.JRClient == nil {
		return nil, ErrJiraRESTClientCannotBeNil
	}

	apiURL := urlutil.JoinAbsolute(
		svc.JRClient.Config.ServerURL,
		APIV3URLCreateMeta,
		projectKey,
		"issuetypes",
		issueTypeID,
	)

	hclient := svc.JRClient.HTTPClient
	if hclient == nil {
		hclient = &http.Client{}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	resp, err := hclient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return nil, fmt.Errorf("error status code (%d) for project %q, issue type %q", resp.StatusCode, projectKey, issueTypeID)
	}

	var result CreateMetaFieldsResponse
	if _, err := jsonutil.UnmarshalReader(resp.Body, &result); err != nil {
		return nil, err
	}

	return result.Values, nil
}

// GetAllFieldsForProject returns all custom fields available across all issue types
// in a project. It queries each issue type and aggregates the unique fields.
func (svc *CreateMetaService) GetAllFieldsForProject(ctx context.Context, projectKey string) (CreateMetaFields, error) {
	issueTypes, err := svc.GetIssueTypes(ctx, projectKey)
	if err != nil {
		return nil, fmt.Errorf("getting issue types for project %q: %w", projectKey, err)
	}

	seen := make(map[string]CreateMetaField)
	for _, it := range issueTypes {
		fields, err := svc.GetFields(ctx, projectKey, it.ID)
		if err != nil {
			return nil, fmt.Errorf("getting fields for project %q, issue type %q: %w", projectKey, it.ID, err)
		}
		for _, f := range fields {
			if _, exists := seen[f.Key]; !exists {
				seen[f.Key] = f
			}
		}
	}

	result := make(CreateMetaFields, 0, len(seen))
	for _, f := range seen {
		result = append(result, f)
	}
	return result, nil
}
