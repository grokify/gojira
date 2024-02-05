package gojira

import (
	"slices"
	"time"
)

const (
	FieldFilter  = "filter"
	FieldIssue   = "issue" // issue keys
	FieldKey     = "key"
	FieldParent  = "parent"
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

	StagePlanning    = "Planning"
	StageDesign      = "Design"
	StageDevelopment = "Development"
	StageTesting     = "Testing"
	StageDeployment  = "Deployment"
	StageReview      = "Review"

	metaStagePrefixReadyFor = "Ready for "
	metaStagePrefixIn       = "In "

	MetaStageReadyForPlanning    = metaStagePrefixReadyFor + StagePlanning
	MetaStageInPlanning          = metaStagePrefixIn + StagePlanning
	MetaStageReadyForDesign      = metaStagePrefixReadyFor + StageDesign
	MetaStageInDesign            = metaStagePrefixIn + StageDesign
	MetaStageReadyForDevelopment = metaStagePrefixReadyFor + StageDevelopment
	MetaStageInDevelopment       = metaStagePrefixIn + StageDevelopment
	MetaStageReadyForTesting     = metaStagePrefixReadyFor + StageTesting
	MetaStageInTesting           = metaStagePrefixIn + StageTesting
	MetaStageReadyForDeployment  = metaStagePrefixReadyFor + StageDeployment
	MetaStageInDeployment        = metaStagePrefixIn + StageDeployment
	MetaStageReadyForReview      = metaStagePrefixReadyFor + StageReview
	MetaStageInReview            = metaStagePrefixIn + StageReview
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
