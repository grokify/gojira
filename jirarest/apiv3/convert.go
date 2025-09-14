package apiv3

import (
	"fmt"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

// ConvertToGoJiraIssue converts a V3 API Issue to a go-jira Issue
func (issue *Issue) ConvertToGoJiraIssue() (*jira.Issue, error) {
	goJiraIssue := &jira.Issue{
		Expand: issue.Expand,
		ID:     issue.ID,
		Key:    issue.Key,
		Self:   issue.Self,
		Fields: &jira.IssueFields{},
	}

	if issue.Fields != nil {
		if err := convertFields(issue.Fields, goJiraIssue.Fields); err != nil {
			return nil, fmt.Errorf("failed to convert fields: %w", err)
		}
	}

	return goJiraIssue, nil
}

// convertFields converts V3 Fields to go-jira IssueFields
func convertFields(v3Fields *Fields, goJiraFields *jira.IssueFields) error {
	// Summary
	goJiraFields.Summary = v3Fields.Summary

	// Description - handle ADF to string conversion
	if v3Fields.Description != nil {
		if desc, err := extractTextFromADF(v3Fields.Description); err == nil {
			goJiraFields.Description = desc
		}
	}

	// Issue Type
	if v3Fields.IssueType != nil {
		goJiraFields.Type = jira.IssueType{
			ID:          v3Fields.IssueType.ID,
			Name:        v3Fields.IssueType.Name,
			Description: v3Fields.IssueType.Description,
			IconURL:     v3Fields.IssueType.IconURL,
			Self:        v3Fields.IssueType.Self,
			Subtask:     v3Fields.IssueType.Subtask,
		}
	}

	// Status
	if v3Fields.Status != nil {
		goJiraFields.Status = &jira.Status{
			ID:          v3Fields.Status.ID,
			Name:        v3Fields.Status.Name,
			Description: v3Fields.Status.Description,
			IconURL:     v3Fields.Status.IconURL,
			Self:        v3Fields.Status.Self,
		}
	}

	// Priority
	if v3Fields.Priority != nil {
		goJiraFields.Priority = &jira.Priority{
			ID:      v3Fields.Priority.ID,
			Name:    v3Fields.Priority.Name,
			IconURL: v3Fields.Priority.IconURL,
			Self:    v3Fields.Priority.Self,
		}
	}

	// Project
	if v3Fields.Project != nil {
		goJiraFields.Project = jira.Project{
			ID:   v3Fields.Project.ID,
			Key:  v3Fields.Project.Key,
			Name: v3Fields.Project.Name,
			Self: v3Fields.Project.Self,
		}
	}

	// Assignee
	if v3Fields.Assignee != nil {
		goJiraFields.Assignee = convertUser(v3Fields.Assignee)
	}

	// Reporter
	if v3Fields.Reporter != nil {
		goJiraFields.Reporter = convertUser(v3Fields.Reporter)
	}

	// Creator
	if v3Fields.Creator != nil {
		goJiraFields.Creator = convertUser(v3Fields.Creator)
	}

	// Labels
	goJiraFields.Labels = v3Fields.Labels

	// Initialize Unknowns for custom fields handling
	if goJiraFields.Unknowns == nil {
		goJiraFields.Unknowns = make(map[string]interface{})
	}

	// Copy custom fields to Unknowns
	if len(v3Fields.CustomFields) > 0 {
		for key, value := range v3Fields.CustomFields {
			goJiraFields.Unknowns[key] = value
		}
	}

	return nil
}

// convertUser converts a V3 User to a go-jira User
func convertUser(v3User *User) *jira.User {
	user := &jira.User{
		AccountID:    v3User.AccountID,
		AccountType:  v3User.AccountType,
		Active:       v3User.Active,
		DisplayName:  v3User.DisplayName,
		EmailAddress: v3User.EmailAddress,
		Self:         v3User.Self,
		TimeZone:     v3User.TimeZone,
	}

	// Note: Skipping avatar URLs due to type conversion complexity

	return user
}

// extractTextFromADF extracts plain text from Atlassian Document Format (ADF) content
func extractTextFromADF(content interface{}) (string, error) {
	if content == nil {
		return "", nil
	}

	// If it's already a string, return as-is
	if str, ok := content.(string); ok {
		return str, nil
	}

	// If it's an ADF object, try to extract text
	if adfObj, ok := content.(map[string]interface{}); ok {
		return extractTextFromADFMap(adfObj), nil
	}

	// Try to convert to string as fallback
	return fmt.Sprintf("%v", content), nil
}

// extractTextFromADFMap recursively extracts text from an ADF map structure
func extractTextFromADFMap(adfMap map[string]interface{}) string {
	var textParts []string

	// Check for direct text content
	if text, exists := adfMap["text"]; exists {
		if textStr, ok := text.(string); ok {
			textParts = append(textParts, textStr)
		}
	}

	// Check for content array
	if content, exists := adfMap["content"]; exists {
		if contentArray, ok := content.([]interface{}); ok {
			for _, item := range contentArray {
				if itemMap, ok := item.(map[string]interface{}); ok {
					if extractedText := extractTextFromADFMap(itemMap); extractedText != "" {
						textParts = append(textParts, extractedText)
					}
				}
			}
		}
	}

	return strings.Join(textParts, " ")
}
