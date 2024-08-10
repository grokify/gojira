package gojira

import (
	"time"

	"golang.org/x/exp/slices"
)

const (
	// These are used by "GoJira" but not necessarily "Jira"
	FieldCreatedDate = "createddate"
	FieldFilter      = "filter"
	FieldIssue       = "issue" // issue keys
	FieldKey         = "key"
	FieldParent      = "parent"
	FieldProject     = "project" // project keys
	FieldProjectKey  = "projectkey"
	FieldResolution  = "resolution"
	FieldStatus      = "status"
	FieldSummary     = "summary"
	FieldType        = "type"

	CalcCreatedAgeDays = "createdagedays"
	CalcCreatedMonth   = "createdmonth"
	AliasIssueKey      = "issuekey"

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

	StagePlanning    = "Planning"
	StageDesign      = "Design"
	StageDevelopment = "Development"
	StageTesting     = "Testing"
	StageDeployment  = "Deployment"
	StageReview      = "Review"

	MetaStagePrefixReadyFor = "Ready for "
	MetaStagePrefixIn       = "In "

	MetaStageReadyForPlanning    = MetaStagePrefixReadyFor + StagePlanning
	MetaStageInPlanning          = MetaStagePrefixIn + StagePlanning
	MetaStageReadyForDesign      = MetaStagePrefixReadyFor + StageDesign
	MetaStageInDesign            = MetaStagePrefixIn + StageDesign
	MetaStageReadyForDevelopment = MetaStagePrefixReadyFor + StageDevelopment
	MetaStageInDevelopment       = MetaStagePrefixIn + StageDevelopment
	MetaStageReadyForTesting     = MetaStagePrefixReadyFor + StageTesting
	MetaStageInTesting           = MetaStagePrefixIn + StageTesting
	MetaStageReadyForDeployment  = MetaStagePrefixReadyFor + StageDeployment
	MetaStageInDeployment        = MetaStagePrefixIn + StageDeployment
	MetaStageReadyForReview      = MetaStagePrefixReadyFor + StageReview
	MetaStageInReview            = MetaStagePrefixIn + StageReview
	MetaStageDone                = StatusDone

	WorkingHoursPerDayDefault float32 = 8.0
	WorkingDaysPerWeekDefault float32 = 5.0

	JiraXMLGenerated = time.UnixDate // "Fri Jul 28 01:07:16 UTC 2023"

	JQLMaxResults = 100
	JQLMaxLength  = 6000 // https://jira.atlassian.com/browse/JRASERVER-41005
	JQLInSep      = ","
)

func MetaStageOrder() []string {
	return []string{
		MetaStageReadyForPlanning,
		MetaStageInPlanning,
		MetaStageReadyForDesign,
		MetaStageInDesign,
		MetaStageReadyForDevelopment,
		MetaStageInDevelopment,
		MetaStageReadyForTesting,
		MetaStageInTesting,
		MetaStageReadyForDeployment,
		MetaStageInDeployment,
		MetaStageReadyForReview,
		MetaStageInReview,
		MetaStageDone}
}

func IsMetaStage(status string) bool {
	return slices.Index(MetaStageOrder(), status) > -1
}

func StatusesInactive() []string {
	return []string{
		StatusDone,
	}
}
