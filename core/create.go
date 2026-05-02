package core

import (
	"context"
	"fmt"
	"os"
	"strings"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira/rest"
	"gopkg.in/yaml.v3"
)

// CreateIssueFromFile reads a YAML file and creates a Jira issue.
func CreateIssueFromFile(ctx context.Context, client *rest.Client, filename string) (*IssueResult, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("read file: %w", err)
	}

	return CreateIssueFromYAML(ctx, client, data)
}

// CreateIssueFromYAML parses YAML content and creates a Jira issue.
func CreateIssueFromYAML(ctx context.Context, client *rest.Client, data []byte) (*IssueResult, error) {
	input, err := ParseIssueYAML(data)
	if err != nil {
		return nil, fmt.Errorf("parse YAML: %w", err)
	}

	return CreateIssue(ctx, client, input)
}

// ParseIssueYAML parses YAML data into IssueInput.
func ParseIssueYAML(data []byte) (*IssueInput, error) {
	var input IssueInput

	// First unmarshal to get standard fields
	if err := yaml.Unmarshal(data, &input); err != nil {
		return nil, err
	}

	// Then unmarshal to raw map to capture custom fields
	var rawMap map[string]any
	if err := yaml.Unmarshal(data, &rawMap); err != nil {
		return nil, err
	}
	input.RawFields = rawMap

	return &input, nil
}

// CreateIssue creates a Jira issue from IssueInput.
func CreateIssue(ctx context.Context, client *rest.Client, input *IssueInput) (*IssueResult, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	// Build the Jira issue
	issue := &jira.Issue{
		Fields: &jira.IssueFields{
			Summary:     input.Summary,
			Description: input.Description,
			Project: jira.Project{
				Key: input.Project,
			},
			Type: jira.IssueType{
				Name: input.Type,
			},
			Labels: input.Labels,
		},
	}

	// Set parent if provided (for subtasks or stories under epics)
	if input.Parent != "" {
		issue.Fields.Parent = &jira.Parent{
			Key: input.Parent,
		}
	}

	// Set priority if provided
	if input.Priority != "" {
		issue.Fields.Priority = &jira.Priority{
			Name: input.Priority,
		}
	}

	// Set assignee if provided
	if input.Assignee != "" {
		issue.Fields.Assignee = &jira.User{
			Name: input.Assignee,
		}
	}

	// Set reporter if provided
	if input.Reporter != "" {
		issue.Fields.Reporter = &jira.User{
			Name: input.Reporter,
		}
	}

	// Set components if provided
	if len(input.Components) > 0 {
		for _, c := range input.Components {
			issue.Fields.Components = append(issue.Fields.Components, &jira.Component{
				Name: c,
			})
		}
	}

	// Set fix versions if provided
	if len(input.FixVersions) > 0 {
		for _, v := range input.FixVersions {
			issue.Fields.FixVersions = append(issue.Fields.FixVersions, &jira.FixVersion{
				Name: v,
			})
		}
	}

	// Handle custom fields
	customFields := input.GetCustomFields()
	if len(customFields) > 0 {
		if issue.Fields.Unknowns == nil {
			issue.Fields.Unknowns = make(map[string]any)
		}
		for k, v := range customFields {
			issue.Fields.Unknowns[k] = v
		}
	}

	// Create the issue
	created, resp, err := client.JiraClient.Issue.CreateWithContext(ctx, issue)
	if err != nil {
		if resp != nil && resp.StatusCode >= 400 {
			return nil, fmt.Errorf("create issue failed (status %d): %w", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("create issue: %w", err)
	}

	return &IssueResult{
		Key:     created.Key,
		ID:      created.ID,
		Self:    created.Self,
		Summary: input.Summary,
	}, nil
}

func validateInput(input *IssueInput) error {
	var missing []string

	if input.Project == "" {
		missing = append(missing, "project")
	}
	if input.Type == "" {
		missing = append(missing, "type")
	}
	if input.Summary == "" {
		missing = append(missing, "summary")
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required fields: %s", strings.Join(missing, ", "))
	}

	return nil
}

// DryRunCreate validates input and returns what would be created without actually creating.
func DryRunCreate(input *IssueInput) (*DryRunResult, error) {
	if err := validateInput(input); err != nil {
		return nil, err
	}

	return &DryRunResult{
		Project:      input.Project,
		Type:         input.Type,
		Summary:      input.Summary,
		Description:  input.Description,
		Parent:       input.Parent,
		Labels:       input.Labels,
		Priority:     input.Priority,
		Assignee:     input.Assignee,
		CustomFields: input.GetCustomFields(),
		Valid:        true,
	}, nil
}

// DryRunResult shows what would be created.
type DryRunResult struct {
	Valid        bool           `json:"valid"`
	Project      string         `json:"project"`
	Type         string         `json:"type"`
	Summary      string         `json:"summary"`
	Description  string         `json:"description,omitempty"`
	Parent       string         `json:"parent,omitempty"`
	Labels       []string       `json:"labels,omitempty"`
	Priority     string         `json:"priority,omitempty"`
	Assignee     string         `json:"assignee,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}
