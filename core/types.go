// Package core provides shared business logic for gojira CLI and MCP server.
package core

// IssueInput represents the input for creating or updating a Jira issue.
// Fields map directly to Jira API fields, with customfield_* supported.
type IssueInput struct {
	// Standard fields
	Project     string   `yaml:"project" json:"project"`
	Type        string   `yaml:"type" json:"type"`
	Summary     string   `yaml:"summary" json:"summary"`
	Description string   `yaml:"description" json:"description"`
	Parent      string   `yaml:"parent" json:"parent,omitempty"`
	Labels      []string `yaml:"labels" json:"labels,omitempty"`
	Priority    string   `yaml:"priority" json:"priority,omitempty"`
	Assignee    string   `yaml:"assignee" json:"assignee,omitempty"`
	Reporter    string   `yaml:"reporter" json:"reporter,omitempty"`
	Components  []string `yaml:"components" json:"components,omitempty"`
	FixVersions []string `yaml:"fix_versions" json:"fix_versions,omitempty"`

	// Custom fields - keys like customfield_12345
	CustomFields map[string]any `yaml:"-" json:"custom_fields,omitempty"`

	// RawFields captures all fields for custom field extraction
	RawFields map[string]any `yaml:",inline" json:"-"`
}

// IssueResult represents the result of a create/update operation.
type IssueResult struct {
	Key     string `json:"key"`
	ID      string `json:"id"`
	Self    string `json:"self"`
	Summary string `json:"summary,omitempty"`
}

// GetCustomFields extracts customfield_* entries from RawFields.
func (i *IssueInput) GetCustomFields() map[string]any {
	if i.CustomFields != nil {
		return i.CustomFields
	}

	result := make(map[string]any)
	for k, v := range i.RawFields {
		if len(k) > 12 && k[:12] == "customfield_" {
			result[k] = v
		}
	}
	return result
}
