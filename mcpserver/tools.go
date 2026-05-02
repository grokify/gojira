package mcpserver

// GetTools returns all available Jira tools.
func GetTools() []Tool {
	return []Tool{
		{
			Name:        "jira_get_issue",
			Description: "Get a Jira issue by key with all fields including description, status, assignee, and custom fields",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
					"expand": map[string]any{
						"type":        "string",
						"description": "Comma-separated list of fields to expand (e.g., changelog,renderedFields)",
					},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "jira_search",
			Description: "Search Jira issues using JQL (Jira Query Language). Returns matching issues with key fields.",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"jql": map[string]any{
						"type":        "string",
						"description": "JQL query string (e.g., 'project = PROJ AND status = Open')",
					},
					"max_results": map[string]any{
						"type":        "integer",
						"description": "Maximum number of results to return (default: 50, max: 100)",
						"default":     50,
					},
					"fields": map[string]any{
						"type":        "string",
						"description": "Comma-separated list of fields to return (default: key,summary,status,assignee,created,updated)",
					},
				},
				"required": []string{"jql"},
			},
		},
		{
			Name:        "jira_update_issue",
			Description: "Update a Jira issue's fields such as summary, description, labels, or custom fields",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
					"summary": map[string]any{
						"type":        "string",
						"description": "New summary/title for the issue",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "New description for the issue",
					},
					"labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Labels to set on the issue (replaces existing labels)",
					},
					"add_labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Labels to add to the issue (preserves existing labels)",
					},
					"remove_labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Labels to remove from the issue",
					},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "jira_add_comment",
			Description: "Add a comment to a Jira issue",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
					"body": map[string]any{
						"type":        "string",
						"description": "Comment body text",
					},
				},
				"required": []string{"key", "body"},
			},
		},
		{
			Name:        "jira_get_transitions",
			Description: "Get available status transitions for a Jira issue",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "jira_transition_issue",
			Description: "Transition a Jira issue to a new status",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
					"transition_id": map[string]any{
						"type":        "string",
						"description": "Transition ID (get available transitions using jira_get_transitions)",
					},
					"comment": map[string]any{
						"type":        "string",
						"description": "Optional comment to add with the transition",
					},
				},
				"required": []string{"key", "transition_id"},
			},
		},
		{
			Name:        "jira_get_comments",
			Description: "Get comments on a Jira issue",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"key": map[string]any{
						"type":        "string",
						"description": "Issue key (e.g., PROJ-123)",
					},
					"max_results": map[string]any{
						"type":        "integer",
						"description": "Maximum number of comments to return (default: 50)",
						"default":     50,
					},
				},
				"required": []string{"key"},
			},
		},
		{
			Name:        "jira_get_projects",
			Description: "List available Jira projects",
			InputSchema: map[string]any{
				"type":       "object",
				"properties": map[string]any{},
			},
		},
		{
			Name:        "jira_create_issue",
			Description: "Create a new Jira issue (Story, Bug, Task, etc.) with support for custom fields",
			InputSchema: map[string]any{
				"type": "object",
				"properties": map[string]any{
					"project": map[string]any{
						"type":        "string",
						"description": "Project key (e.g., PROJ)",
					},
					"type": map[string]any{
						"type":        "string",
						"description": "Issue type (e.g., Story, Bug, Task, Epic)",
					},
					"summary": map[string]any{
						"type":        "string",
						"description": "Issue summary/title",
					},
					"description": map[string]any{
						"type":        "string",
						"description": "Issue description (supports Jira markdown)",
					},
					"parent": map[string]any{
						"type":        "string",
						"description": "Parent issue key for subtasks or stories under epics (e.g., PROJ-100)",
					},
					"labels": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Labels to apply to the issue",
					},
					"priority": map[string]any{
						"type":        "string",
						"description": "Priority name (e.g., High, Medium, Low)",
					},
					"assignee": map[string]any{
						"type":        "string",
						"description": "Assignee username or email",
					},
					"components": map[string]any{
						"type":        "array",
						"items":       map[string]any{"type": "string"},
						"description": "Component names",
					},
					"custom_fields": map[string]any{
						"type":        "object",
						"description": "Custom fields as key-value pairs (e.g., {\"customfield_12345\": \"value\"})",
					},
				},
				"required": []string{"project", "type", "summary"},
			},
		},
	}
}
