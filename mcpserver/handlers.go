package mcpserver

import (
	"context"
	"fmt"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/gojira/core"
	"github.com/grokify/gojira/rest"
)

// CallTool dispatches a tool call to the appropriate handler.
func (s *Server) CallTool(ctx context.Context, name string, args map[string]any) (any, error) {
	switch name {
	case "jira_get_issue":
		return s.handleGetIssue(ctx, args)
	case "jira_search":
		return s.handleSearch(ctx, args)
	case "jira_update_issue":
		return s.handleUpdateIssue(ctx, args)
	case "jira_add_comment":
		return s.handleAddComment(ctx, args)
	case "jira_get_transitions":
		return s.handleGetTransitions(ctx, args)
	case "jira_transition_issue":
		return s.handleTransitionIssue(ctx, args)
	case "jira_get_comments":
		return s.handleGetComments(ctx, args)
	case "jira_get_projects":
		return s.handleGetProjects(ctx, args)
	case "jira_create_issue":
		return s.handleCreateIssue(ctx, args)
	default:
		return nil, fmt.Errorf("unknown tool: %s", name)
	}
}

// IssueResult is a simplified issue representation for tool output.
type IssueResult struct {
	Key          string         `json:"key"`
	Summary      string         `json:"summary"`
	Description  string         `json:"description,omitempty"`
	Status       string         `json:"status"`
	Type         string         `json:"type"`
	Priority     string         `json:"priority,omitempty"`
	Assignee     string         `json:"assignee,omitempty"`
	Reporter     string         `json:"reporter,omitempty"`
	Labels       []string       `json:"labels,omitempty"`
	Created      string         `json:"created,omitempty"`
	Updated      string         `json:"updated,omitempty"`
	Project      string         `json:"project,omitempty"`
	Parent       string         `json:"parent,omitempty"`
	CustomFields map[string]any `json:"custom_fields,omitempty"`
}

func issueToResult(issue *jira.Issue) IssueResult {
	result := IssueResult{
		Key:    issue.Key,
		Labels: issue.Fields.Labels,
	}

	if issue.Fields != nil {
		result.Summary = issue.Fields.Summary
		result.Description = issue.Fields.Description

		if issue.Fields.Status != nil {
			result.Status = issue.Fields.Status.Name
		}
		if issue.Fields.Type.Name != "" {
			result.Type = issue.Fields.Type.Name
		}
		if issue.Fields.Priority != nil {
			result.Priority = issue.Fields.Priority.Name
		}
		if issue.Fields.Assignee != nil {
			result.Assignee = issue.Fields.Assignee.DisplayName
		}
		if issue.Fields.Reporter != nil {
			result.Reporter = issue.Fields.Reporter.DisplayName
		}
		if issue.Fields.Project.Key != "" {
			result.Project = issue.Fields.Project.Key
		}
		if issue.Fields.Parent != nil {
			result.Parent = issue.Fields.Parent.Key
		}
		if !time.Time(issue.Fields.Created).IsZero() {
			result.Created = time.Time(issue.Fields.Created).Format("2006-01-02T15:04:05Z")
		}
		if !time.Time(issue.Fields.Updated).IsZero() {
			result.Updated = time.Time(issue.Fields.Updated).Format("2006-01-02T15:04:05Z")
		}
	}

	return result
}

func (s *Server) handleGetIssue(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	var opts *rest.GetQueryOptions
	if expand, ok := args["expand"].(string); ok && expand != "" {
		opts = &rest.GetQueryOptions{
			ExpandChangelog: expand == "changelog" || expand == "all",
		}
	}

	issue, err := s.client.IssueAPI.Issue(ctx, key, opts)
	if err != nil {
		return nil, fmt.Errorf("get issue %s: %w", key, err)
	}

	return issueToResult(issue), nil
}

func (s *Server) handleSearch(ctx context.Context, args map[string]any) (any, error) {
	jql, ok := args["jql"].(string)
	if !ok || jql == "" {
		return nil, fmt.Errorf("jql is required")
	}

	maxResults := 50
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
		if maxResults > 100 {
			maxResults = 100
		}
	}

	// Use the context-aware V3 API for search
	issues, err := s.client.IssueAPI.SearchIssuesAPIV3(ctx, jql, false)
	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	// Limit results
	if len(issues) > maxResults {
		issues = issues[:maxResults]
	}

	results := make([]IssueResult, 0, len(issues))
	for _, issue := range issues {
		results = append(results, issueToResult(&issue))
	}

	return map[string]any{
		"total":  len(results),
		"issues": results,
	}, nil
}

func (s *Server) handleUpdateIssue(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	// Build update request
	updateBody := rest.IssuePatchRequestBody{}
	hasUpdate := false

	// Handle label operations
	if addLabels, ok := args["add_labels"].([]any); ok && len(addLabels) > 0 {
		if updateBody.Update == nil {
			updateBody.Update = &rest.IssuePatchRequestBodyUpdate{}
		}
		for _, label := range addLabels {
			if labelStr, ok := label.(string); ok {
				labelCopy := labelStr
				updateBody.Update.Labels = append(updateBody.Update.Labels, rest.IssuePatchRequestBodyUpdateLabel{
					Add: &labelCopy,
				})
			}
		}
		hasUpdate = true
	}

	if removeLabels, ok := args["remove_labels"].([]any); ok && len(removeLabels) > 0 {
		if updateBody.Update == nil {
			updateBody.Update = &rest.IssuePatchRequestBodyUpdate{}
		}
		for _, label := range removeLabels {
			if labelStr, ok := label.(string); ok {
				labelCopy := labelStr
				updateBody.Update.Labels = append(updateBody.Update.Labels, rest.IssuePatchRequestBodyUpdateLabel{
					Remove: &labelCopy,
				})
			}
		}
		hasUpdate = true
	}

	// Handle summary and description through the SDK's update method
	if summary, ok := args["summary"].(string); ok && summary != "" {
		if updateBody.Fields == nil {
			updateBody.Fields = make(map[string]rest.IssuePatchRequestBodyField)
		}
		updateBody.Fields["summary"] = rest.IssuePatchRequestBodyField{Value: summary}
		hasUpdate = true
	}

	if description, ok := args["description"].(string); ok && description != "" {
		if updateBody.Fields == nil {
			updateBody.Fields = make(map[string]rest.IssuePatchRequestBodyField)
		}
		updateBody.Fields["description"] = rest.IssuePatchRequestBodyField{Value: description}
		hasUpdate = true
	}

	if !hasUpdate {
		return nil, fmt.Errorf("no update fields provided")
	}

	_, err := s.client.IssueAPI.IssuePatch(ctx, key, updateBody)
	if err != nil {
		return nil, fmt.Errorf("update issue %s: %w", key, err)
	}

	return map[string]any{
		"success": true,
		"key":     key,
		"message": "Issue updated successfully",
	}, nil
}

func (s *Server) handleAddComment(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	body, ok := args["body"].(string)
	if !ok || body == "" {
		return nil, fmt.Errorf("body is required")
	}

	comment := &jira.Comment{
		Body: body,
	}

	addedComment, _, err := s.client.JiraClient.Issue.AddCommentWithContext(ctx, key, comment)
	if err != nil {
		return nil, fmt.Errorf("add comment to %s: %w", key, err)
	}

	return map[string]any{
		"success":    true,
		"key":        key,
		"comment_id": addedComment.ID,
		"message":    "Comment added successfully",
	}, nil
}

func (s *Server) handleGetTransitions(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	transitions, _, err := s.client.IssueAPI.GetTransitions(ctx, key, false)
	if err != nil {
		return nil, fmt.Errorf("get transitions for %s: %w", key, err)
	}

	results := make([]map[string]any, 0, len(transitions))
	for _, t := range transitions {
		results = append(results, map[string]any{
			"id":   t.ID,
			"name": t.Name,
			"to":   t.To.Name,
		})
	}

	return map[string]any{
		"key":         key,
		"transitions": results,
	}, nil
}

func (s *Server) handleTransitionIssue(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	transitionID, ok := args["transition_id"].(string)
	if !ok || transitionID == "" {
		return nil, fmt.Errorf("transition_id is required")
	}

	_, err := s.client.JiraClient.Issue.DoTransitionWithContext(ctx, key, transitionID)
	if err != nil {
		return nil, fmt.Errorf("transition issue %s: %w", key, err)
	}

	// Add comment if provided
	if comment, ok := args["comment"].(string); ok && comment != "" {
		_, _, err := s.client.JiraClient.Issue.AddCommentWithContext(ctx, key, &jira.Comment{Body: comment})
		if err != nil {
			return map[string]any{
				"success":       true,
				"key":           key,
				"transition_id": transitionID,
				"comment_error": err.Error(),
				"message":       "Issue transitioned but comment failed",
			}, nil
		}
	}

	return map[string]any{
		"success":       true,
		"key":           key,
		"transition_id": transitionID,
		"message":       "Issue transitioned successfully",
	}, nil
}

func (s *Server) handleGetComments(ctx context.Context, args map[string]any) (any, error) {
	key, ok := args["key"].(string)
	if !ok || key == "" {
		return nil, fmt.Errorf("key is required")
	}

	// Get issue with comments expanded
	issue, _, err := s.client.JiraClient.Issue.GetWithContext(ctx, key, &jira.GetQueryOptions{
		Expand: "renderedFields",
	})
	if err != nil {
		return nil, fmt.Errorf("get issue %s: %w", key, err)
	}

	if issue.Fields.Comments == nil {
		return map[string]any{
			"key":      key,
			"total":    0,
			"comments": []any{},
		}, nil
	}

	maxResults := 50
	if mr, ok := args["max_results"].(float64); ok {
		maxResults = int(mr)
	}

	comments := issue.Fields.Comments.Comments
	if len(comments) > maxResults {
		comments = comments[:maxResults]
	}

	results := make([]map[string]any, 0, len(comments))
	for _, c := range comments {
		author := ""
		if c.Author.DisplayName != "" {
			author = c.Author.DisplayName
		}
		results = append(results, map[string]any{
			"id":      c.ID,
			"author":  author,
			"body":    c.Body,
			"created": c.Created,
			"updated": c.Updated,
		})
	}

	return map[string]any{
		"key":      key,
		"total":    len(results),
		"comments": results,
	}, nil
}

func (s *Server) handleGetProjects(ctx context.Context, _ map[string]any) (any, error) {
	projects, _, err := s.client.JiraClient.Project.GetListWithContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("get projects: %w", err)
	}

	results := make([]map[string]any, 0, len(*projects))
	for _, p := range *projects {
		results = append(results, map[string]any{
			"key":  p.Key,
			"name": p.Name,
			"id":   p.ID,
		})
	}

	return map[string]any{
		"total":    len(results),
		"projects": results,
	}, nil
}

func (s *Server) handleCreateIssue(ctx context.Context, args map[string]any) (any, error) {
	// Build IssueInput from args
	input := &core.IssueInput{}

	// Required fields
	project, ok := args["project"].(string)
	if !ok || project == "" {
		return nil, fmt.Errorf("project is required")
	}
	input.Project = project

	issueType, ok := args["type"].(string)
	if !ok || issueType == "" {
		return nil, fmt.Errorf("type is required")
	}
	input.Type = issueType

	summary, ok := args["summary"].(string)
	if !ok || summary == "" {
		return nil, fmt.Errorf("summary is required")
	}
	input.Summary = summary

	// Optional fields
	if description, ok := args["description"].(string); ok {
		input.Description = description
	}

	if parent, ok := args["parent"].(string); ok {
		input.Parent = parent
	}

	if priority, ok := args["priority"].(string); ok {
		input.Priority = priority
	}

	if assignee, ok := args["assignee"].(string); ok {
		input.Assignee = assignee
	}

	// Handle labels array
	if labels, ok := args["labels"].([]any); ok {
		for _, label := range labels {
			if labelStr, ok := label.(string); ok {
				input.Labels = append(input.Labels, labelStr)
			}
		}
	}

	// Handle components array
	if components, ok := args["components"].([]any); ok {
		for _, comp := range components {
			if compStr, ok := comp.(string); ok {
				input.Components = append(input.Components, compStr)
			}
		}
	}

	// Handle custom fields
	if customFields, ok := args["custom_fields"].(map[string]any); ok {
		input.CustomFields = customFields
	}

	// Create the issue using core package
	result, err := core.CreateIssue(ctx, s.client, input)
	if err != nil {
		return nil, fmt.Errorf("create issue: %w", err)
	}

	return map[string]any{
		"success": true,
		"key":     result.Key,
		"id":      result.ID,
		"self":    result.Self,
		"summary": result.Summary,
		"message": fmt.Sprintf("Issue %s created successfully", result.Key),
	}, nil
}
