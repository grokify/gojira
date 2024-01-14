package jirarest

const (
	APIURL2ListCustomFields = `/rest/api/2/field`

	StatusDone         = "Done"
	StatusOpen         = "Open"
	StatusCustomClosed = "Closed"

	MaxResults    = uint(1000)
	MetaParamRank = "_rank"

	TimeTimeSpent                     = "Time Spent"
	TimeTimeEstimate                  = "Time Estimate"
	TimeTimeOriginalEstimate          = "Time Original Estimate"
	TimeAggregateTimeOriginalEstimate = "Aggregate Time Original Estimate"
	TimeAggregateTimeSpent            = "Aggregate Time Spent"
	TimeAggregateTimeEstimate         = "Aggregate Time Estimate"
	TimeTimeRemaining                 = "Time Remaining"
	TimeTimeRemainingOriginal         = "Time Remaining Original"

	FieldSlugType       = "type"
	FieldSlugProjectkey = "projectkey"
)

type Status struct {
	Name        string
	Description string
}

func IssueStatuses() []Status {
	return []Status{
		{
			Name:        "Open",
			Description: "The issue is open and ready for the assignee to start work on it.",
		},
	}
}

/*

	In Progress

	This issue is being actively worked on at the moment by the assignee.

	Done

	Work has finished on the issue.

	To Do

	The issue has been reported and is waiting for the team to action it.

	In Review

	The assignee has carried out the work needed on the issue, and it needs peer review before being considered done.

	Under review

	A reviewer is currently assessing the work completed on the issue before considering it done.

	Approved

	A reviewer has approved the work completed on the issue and the issue is considered done.

	Cancelled

	Work has stopped on the issue and the issue is considered done.

	Rejected

	A reviewer has rejected the work completed on the issue and the issue is considered done.

	Draft

	For content management and document approval projects, the work described on the issue is being prepared for review and is considered in progress, in the draft stage of writing.

	Published

	For content management projects, the work described on the issue has been published and/or released for internal consumption. The issue is considered done.

	Interviewing

	For recruitment projects, this indicates that the candidate is currently in the interviewing stage of the hiring process.

	Interview Debrief

	For recruitment projects, this indicates that the candidate has completed interviewing and interviewers are discussing their next steps in the hiring process.

	Screening

	For recruitment projects, this indicates that the candidate has applied and is being considered for interviews.

	Offer Discussions

	For recruitment projects, this indicates that the candidate has been offered a position and recruiters are working to shore up the details.

	Accepted

	For recruitment projects, this indicates that the candidate has accepted the position. The issue is considered done.

	Applications

	For recruitment projects, this indicates that a candidate has applied and recruiters are waiting to screen them for future action in the hiring process.

	Second Review

	For document approval projects, the work described on the issue has passed its initial review and is being closely proofed for publication.

	Lost

	For lead tracking projects, this indicates that the lead was unsuccessful. The issue is considered done.

	Won

	For lead tracking projects, this indicates that the lead was successful. The issue is considered done.

	Contacted

	For lead tracking projects, this indicates that the sales representative has contacted their lead and the pitch is in progress.

	Opportunity

	For lead tracking projects, this indicates that the sales team has identified and opportunity they want to pursue.

	In Negotiation

	For lead tracking projects, this indicates that the sales team is adjusting their terms to make a sale.

	Purchased

	For procurement projects, this indicates that the service or item was purchased. The issue is considered done.

	Requested

	For procurement projects, this indicates that the service or item has been requested and is waiting for a procurement team member to action the request.
*/
