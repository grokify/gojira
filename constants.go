package gojira

import "time"

const (
	FieldFilter  = "filter"
	FieldIssue   = "issue" // issue keys
	FieldKey     = "key"
	FieldProject = "project" // project keys
	FieldStatus  = "status"
	FieldType    = "type"

	FieldIssuePlural = "issues"

	// Statuses: https://support.atlassian.com/jira-cloud-administration/docs/what-are-issue-statuses-priorities-and-resolutions/
	StatusOpen        = "Open"
	StatusInProgress  = "In Progress"
	StatusDone        = "Done"
	StatusToDo        = "To Do"
	StatusInReview    = "In Review"
	StatusUnderReview = "Under review"
	StatusApproved    = "Approved" // Done

	TypeIssue           = "Issue"
	TypeIssuePlural     = "Issues"
	TypeBug             = "Bug"
	TypeBugPlural       = "Bugs"
	TypeEpic            = "Epic"
	TypeEpicPlural      = "Epics"
	TypeSpike           = "Spike"
	TypeSpikePlural     = "Spikes"
	TypeStory           = "Story"
	TypeStoryPlural     = "Stories"
	TypeInitiative      = "Initiative"
	TypeInitiativePlura = "Initiatives"

	WorkingHoursPerDayDefault float32 = 8.0
	WorkingDaysPerWeekDefault float32 = 5.0

	JiraXMLGenerated = time.UnixDate // "Fri Jul 28 01:07:16 UTC 2023"

	JQLMaxResults = 100
	JQLMaxLength  = 6000 // https://jira.atlassian.com/browse/JRASERVER-41005
	JQLInSep      = ","
)

func StatusesInactive() []string {
	return []string{
		StatusDone,
	}
}
