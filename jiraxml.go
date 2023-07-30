package jiraxml

import (
	"time"

	"github.com/grokify/mogo/encoding/xmlutil"
)

type XML struct {
	Channel Channel `xml:"channel"`
}

func ReadFile(name string) (XML, error) {
	x := XML{}
	err := xmlutil.UnmarshalFile(name, &x)
	return x, err
}

type Channel struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	Items Items  `xml:"item"`
}

type Item struct {
	Type                 Type           `xml:"type"`
	Title                string         `xml:"title"`
	Link                 string         `xml:"link"`
	Key                  Key            `xml:"key"`
	Project              Project        `xml:"project"`
	Summary              string         `xml:"summary"`
	Status               Status         `xml:"status"`
	FixVersion           string         `xml:"fixVersion"`
	TimeOriginalEstimate Duration       `xml:"timeoriginalestimate"`
	AggregateTimeSpent   Duration       `xml:"aggregatetimespent"`
	Labels               []Label        `xml:"labels"`
	Created              RFC1123ZString `xml:"created"` // RFC1123Z
	Updated              RFC1123ZString `xml:"updated"` // RFC1123Z
}

type RFC1123ZString string

func (s RFC1123ZString) Time() (time.Time, error) {
	return time.Parse(time.RFC1123Z, string(s))
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
	Seconds int    `xml:"seconds,attr"`
}

func (d Duration) Duration() time.Duration {
	return time.Duration(int64(d.Seconds * int(time.Second)))
}
