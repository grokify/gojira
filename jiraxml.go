package jiraxml

import (
	"strings"
	"time"

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
	for i, ix := range x.Channel.Items {
		ix.TrimSpace()
		x.Channel.Items[i] = ix
	}
}

type Channel struct {
	Title       string    `xml:"title"`
	Link        string    `xml:"link"`
	Description string    `xml:"description"`
	Language    string    `xml:"language"`
	BuildInfo   BuildInfo `xml:"build-info"`
	Items       Items     `xml:"item"`
}

type BuildInfo struct {
	Version     string        `xml:"version"`
	BuildNumber int64         `xml:"build-number"`
	BuildDate   DMYDateString `xml:"build-date"`
}

type Item struct {
	Type                           Type           `xml:"type"`
	Title                          string         `xml:"title"`
	Link                           string         `xml:"link"`
	Key                            Key            `xml:"key"`
	Project                        Project        `xml:"project"`
	Summary                        string         `xml:"summary"`
	Status                         Status         `xml:"status"`
	FixVersion                     string         `xml:"fixVersion"`
	TimeEstimate                   Duration       `xml:"timeestimate"`
	TimeOriginalEstimate           Duration       `xml:"timeoriginalestimate"`
	TimeSpent                      Duration       `xml:"timespent"`
	AggregateTimeOriginalEstimate  Duration       `xml:"aggregatetimeoriginalestimate"`
	AggregateTimeRemainingEstimate Duration       `xml:"aggregatetimeremainingestimate"`
	AggregateTimeSpent             Duration       `xml:"aggregatetimespent"`
	Labels                         []Label        `xml:"labels"`
	Created                        RFC1123ZString `xml:"created"` // RFC1123Z
	Updated                        RFC1123ZString `xml:"updated"` // RFC1123Z
	Votes                          int            `json:"votes"`
	Watches                        int            `json:"watches"`
}

// TrimSpace removes leading and trailing space. It is useful when parsing XML that has been modified,
// such as by VS Code extensions.
func (i *Item) TrimSpace() {
	i.Type.DisplayName = strings.TrimSpace(i.Type.DisplayName)
	i.Title = strings.TrimSpace(i.Title)
	i.Link = strings.TrimSpace(i.Link)
	i.FixVersion = strings.TrimSpace(i.FixVersion)
	i.Key.DisplayName = strings.TrimSpace(i.Key.DisplayName)
	i.Project.DisplayName = strings.TrimSpace(i.Project.DisplayName)
	i.Status.ID = strings.TrimSpace(i.Status.ID)
	i.Summary = strings.TrimSpace(i.Summary)
	i.TimeEstimate.Display = strings.TrimSpace(i.TimeEstimate.Display)
}

type RFC1123ZString string

func (s RFC1123ZString) Time() (time.Time, error) {
	return time.Parse(time.RFC1123Z, strings.TrimSpace(string(s)))
}

const DMYDateFormat = "_2-01-2006"

type DMYDateString string

func (s DMYDateString) Time() (time.Time, error) {
	return time.Parse(DMYDateFormat, string(s))
}

type Type struct {
	DisplayName string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
}

type Key struct {
	DisplayName string `xml:",chardata"`
	ID          int    `xml:"id,attr"`
}

type Label struct {
	Label string `xml:"label"`
}

type Project struct {
	DisplayName string `xml:",chardata"`
	ID          string `xml:"id,attr"`
	Key         string `xml:"key,attr"`
}

type Status struct {
	ID          string `xml:"id,attr"`
	DisplayName string `xml:",chardata"`
	Description string `xml:"description,attr"`
}

type Duration struct {
	Display string `xml:",chardata"`
	Seconds int64  `xml:"seconds,attr"`
}

func (d Duration) Duration() time.Duration {
	return time.Duration(d.Seconds) * time.Second
}
