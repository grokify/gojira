package jiraxml

import (
	"strconv"
	"strings"

	jira "github.com/andygrunwald/go-jira"
)

func IssueFromAPI(iss jira.Issue) Issue {
	if iss.Fields == nil {
		return Issue{}
	}
	x := Issue{
		Key:                           Simple{DisplayName: iss.Key},
		Summary:                       strings.TrimSpace(iss.Fields.Summary),
		Description:                   strings.TrimSpace(iss.Fields.Description),
		TimeEstimate:                  Duration{Seconds: int64(iss.Fields.TimeEstimate)},
		TimeOriginalEstimate:          Duration{Seconds: int64(iss.Fields.TimeOriginalEstimate)},
		TimeSpent:                     Duration{Seconds: int64(iss.Fields.TimeSpent)},
		AggregateTimeEstimate:         Duration{Seconds: int64(iss.Fields.AggregateTimeEstimate)},
		AggregateTimeOriginalEstimate: Duration{Seconds: int64(iss.Fields.AggregateTimeOriginalEstimate)},
		AggregateTimeSpent:            Duration{Seconds: int64(iss.Fields.AggregateTimeSpent)},
		Created:                       RFC1123ZStringJiraTime(iss.Fields.Created),
		Updated:                       RFC1123ZStringJiraTime(iss.Fields.Updated),
	}
	x.Type = Simple{
		DisplayName: iss.Fields.Type.Name,
	}
	if r := iss.Fields.Resolution; r != nil && len(strings.TrimSpace(r.ID)) > 0 {
		id, err := strconv.Atoi(r.ID)
		if err != nil {
			panic(err)
		}
		x.Resolution = Simple{
			ID:          id,
			DisplayName: strings.TrimSpace(r.Name)}
	}
	if s := iss.Fields.Status; s != nil && len(strings.TrimSpace(s.ID)) > 0 {
		id, err := strconv.Atoi(s.ID)
		if err != nil {
			panic(err)
		}
		x.Status = Status{
			ID:          id,
			Description: strings.TrimSpace(s.Description),
			DisplayName: strings.TrimSpace(s.Name),
			IconURL:     strings.TrimSpace(s.IconURL),
		}
	}
	return x
}
