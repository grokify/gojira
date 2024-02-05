package jirarest

import "errors"

var (
	ErrIssueKeyCannotBeEmpty = errors.New("issue key cannot be empty")
	ErrIssuesSetCannotBeNil  = errors.New("issuesSet cannot be nil")
)
