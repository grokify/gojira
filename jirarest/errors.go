package jirarest

import "errors"

var (
	ErrCustomFieldLabelRequired         = errors.New("custom field label is required")
	ErrIssueKeyCannotBeEmpty            = errors.New("issue key cannot be empty")
	ErrIssueOrIssueKeyOrIssueIDRequired = errors.New("issue, issue id, or issue key required")
	ErrIssuesSetCannotBeNil             = errors.New("issuesSet cannot be nil")
)
