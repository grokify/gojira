package jiraxml

import (
	"strings"
	"time"

	jira "github.com/andygrunwald/go-jira"
	"github.com/grokify/mogo/encoding/xmlutil"
)

type XML struct {
	Channel Channel `xml:"channel"`
}

func ReadFile(name string) (XML, error) {
	x := XML{}
	err := xmlutil.UnmarshalFile(name, &x)
	if err != nil {
		return x, err
	}
	x.TrimSpace()
	return x, err
}

// TrimSpace removes leading and trailing space. It is useful when parsing XML that has been modified,
// such as by VS Code extensions.
func (x *XML) TrimSpace() {
	x.Channel.BuildInfo.BuildDate =
		DMYDateString(
			strings.TrimSpace(string(x.Channel.BuildInfo.BuildDate)))
	for i, ix := range x.Channel.Issues {
		ix.TrimSpace()
		x.Channel.Issues[i] = ix
	}
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	BuildInfo   BuildInfo `xml:"build-info"`
	Issues      Issues    `xml:"item"`
}

type BuildInfo struct {
	Version     string        `xml:"version"`
	BuildNumber int64         `xml:"build-number"`
	BuildDate   DMYDateString `xml:"build-date"`
}

type Issue struct {
	Type                           Simple         `xml:"type"`
	Title                          string         `xml:"title"`
	Description                    string         `xml:"description"`
	Link                           string         `xml:"link"`
	Key                            Simple         `xml:"key"`
	Project                        Project        `xml:"project"`
	Resolution                     Simple         `xml:"resolution"`
	Summary                        string         `xml:"summary"`
	Status                         Status         `xml:"status"`
	Assignee                       User           `xml:"assignee"`
	Reporter                       User           `xml:"reporter"`
	FixVersion                     string         `xml:"fixVersion"`
	TimeEstimate                   Duration       `xml:"timeestimate"`
	TimeOriginalEstimate           Duration       `xml:"timeoriginalestimate"`
	TimeSpent                      Duration       `xml:"timespent"`
	AggregateTimeEstimate          Duration       `xml:"aggregatetimeestimate"`
	AggregateTimeOriginalEstimate  Duration       `xml:"aggregatetimeoriginalestimate"`
	AggregateTimeRemainingEstimate Duration       `xml:"aggregatetimeremainingestimate"`
	AggregateTimeSpent             Duration       `xml:"aggregatetimespent"`
	Labels                         []Label        `xml:"labels"`
	Created                        RFC1123ZString `xml:"created"`  // RFC1123Z
	Updated                        RFC1123ZString `xml:"updated"`  // RFC1123Z
	Resolved                       RFC1123ZString `xml:"resolved"` // RFC1123Z
	Votes                          int            `json:"votes"`
	Watches                        int            `json:"watches"`
}

// TrimSpace removes leading and trailing space. It is useful when parsing XML that has been modified,
// such as by VS Code extensions.
func (i *Issue) TrimSpace() {
	i.Description = strings.TrimSpace(i.Description)
	i.FixVersion = strings.TrimSpace(i.FixVersion)
	i.Link = strings.TrimSpace(i.Link)
	i.Summary = strings.TrimSpace(i.Summary)
	i.Title = strings.TrimSpace(i.Title)
	i.AggregateTimeOriginalEstimate.TrimSpace()
	i.AggregateTimeRemainingEstimate.TrimSpace()
	i.AggregateTimeSpent.TrimSpace()
	i.Assignee.TrimSpace()
	i.Key.TrimSpace()
	i.Project.TrimSpace()
	i.Reporter.TrimSpace()
	i.Resolution.TrimSpace()
	i.TimeEstimate.TrimSpace()
	i.TimeOriginalEstimate.TrimSpace()
	i.TimeSpent.TrimSpace()
	i.Type.TrimSpace()
}

type RFC1123ZString string

func (s RFC1123ZString) Time() (time.Time, error) {
	return time.Parse(time.RFC1123Z, strings.TrimSpace(string(s)))
}

func RFC1123ZStringJiraTime(t jira.Time) RFC1123ZString {
	return RFC1123ZString(time.Time(t).Format(time.RFC1123Z))

}

const DMYDateFormat = "_2-01-2006"

type DMYDateString string

func (s DMYDateString) Time() (time.Time, error) {
	return time.Parse(DMYDateFormat, strings.TrimSpace(string(s)))
}

/*
type Type struct {
	DisplayName string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
}

type Key struct {
	DisplayName string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
}
*/

type Label struct {
	Label string `xml:"label"`
}

type Project struct {
	DisplayName string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
	Key         string `xml:"key,attr"`
}

func (p *Project) TrimSpace() {
	p.DisplayName = strings.TrimSpace(p.DisplayName)
	p.Key = strings.TrimSpace(p.Key)
}

type Simple struct {
	ID          int    `xml:"id,attr"`
	DisplayName string `xml:",chardata"`
}

func (s *Simple) TrimSpace() {
	s.DisplayName = strings.TrimSpace(s.DisplayName)
}

type Status struct {
	ID          int    `xml:"id,attr"`
	DisplayName string `xml:",chardata"`
	Description string `xml:"description,attr"`
	IconURL     string `xml:"iconUrl"`
}

type StatusCategory struct {
	ID          int    `xml:"id,attr"`
	DisplayName string `xml:",chardata"`
	Key         string `xml:"key,attr"`
	ColorName   string `xml:"colorName,attr"`
}

type Duration struct {
	Display string `xml:",chardata"`
	Seconds int64  `xml:"seconds,attr"`
}

func (d *Duration) Duration() time.Duration {
	return time.Duration(d.Seconds) * time.Second
}

func (d *Duration) TrimSpace() {
	d.Display = strings.TrimSpace(d.Display)
}

type User struct {
	Display  string `xml:",chardata"`
	Username string `xml:"username,attr"`
}

func (u *User) TrimSpace() {
	u.Display = strings.TrimSpace(u.Display)
	u.Username = strings.TrimSpace(u.Username)
}
