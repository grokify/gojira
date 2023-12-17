package gojira

import "time"

const (
	FieldIssue   = "issue"   // issue keys
	FieldProject = "project" // project keys
	FieldStatus  = "status"
	FieldType    = "type"

	// Statuses: https://support.atlassian.com/jira-cloud-administration/docs/what-are-issue-statuses-priorities-and-resolutions/
	StatusClosed            = "Closed"
	StatusDone              = "Done"
	StatusInProgress        = "In Progress"
	StatusPOReview          = "PO Review"
	StatusPendingValidation = "Pending Validation"
	StatusReady             = "Ready"

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
)

func StatusesInactive() []string {
	return []string{
		StatusClosed,
		StatusDone,
	}
}
