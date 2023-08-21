package jiraxml

import "time"

const (
	StatusClosed            = "Closed"
	StatusInProgress        = "In Progress"
	StatusPOReview          = "PO Review"
	StatusPendingValidation = "Pending Validation"
	StatusReady             = "Ready"

	TypeBug   = "Bug"
	TypeSpike = "Spike"
	TypeStory = "Story"

	WorkingHoursPerDayDefault float32 = 8.0
	WorkingDaysPerWeekDefault float32 = 5.0

	JiraXMLGenerated = time.UnixDate // "Fri Jul 28 01:07:16 UTC 2023"
)
