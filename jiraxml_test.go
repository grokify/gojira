package jiraxml

import (
	"strings"
	"testing"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

var readFileTests = []struct {
	filename                 string
	buildDateRFC3339FullDate string
	itemCount                int
	itemCountInProgress      int
	itemCountNeedsTriage     int
	item0KeyDisplayName      string
	item0Title               string
	item0Created             string
	item0Updated             string
	hoursPerDay              float32
	daysPerWeek              float32
}{
	{
		"testdata/example_jira_mongodb_new-issues.xml", "2022-07-20",
		20, 3, 6, "SERVER-79445", "[SERVER-79445] Consider upgrading PCRE2",
		"2023-07-28T00:11:05Z", "2023-07-28T01:00:50Z", 8, 5},
	{
		"testdata/example_jira_mongodb_resolved-recently.xml", "2022-07-20",
		20, 3, 6, "SERVER-79005", "[SERVER-79005] [SBE] Call registerSlot() lazily for certain variables",
		"2023-07-15T20:07:23Z", "2023-07-30T13:28:35Z", 8, 5},
}

func TestReadFile(t *testing.T) {
	for _, tt := range readFileTests {
		j, err := ReadFile(tt.filename)
		if err != nil {
			t.Errorf("jiraxml.ReadFile(\"%s\") error: (%s)", tt.filename, err.Error())
		}
		dt, err := j.Channel.BuildInfo.BuildDate.Time()
		if err != nil {
			t.Errorf("BuildDate.Time() error: (%s)", err.Error())
		}
		if dt.Format(timeutil.RFC3339FullDate) != tt.buildDateRFC3339FullDate {
			t.Errorf("file (%s) BuildDate.Time() mismatch: want (%s), got (%s)", tt.filename, tt.buildDateRFC3339FullDate, dt.Format(timeutil.RFC3339FullDate))
		}
		// fmtutil.PrintJSON(j)
		stats := j.Channel.Issues.Stats(tt.hoursPerDay, tt.daysPerWeek)
		if stats.ItemCount != tt.itemCount {
			t.Errorf("jiraxml.ReadFile(\"%s\") mismatch: want (%d), got (%d)", tt.filename, tt.itemCount, stats.ItemCount)
		}
		// fmtutil.PrintJSON(stats)
		if tt.itemCount <= 0 {
			continue
		}
		// fmtutil.PrintJSON(j.Channel.Items[0])
		item0 := j.Channel.Issues[0]
		if item0.Key.DisplayName != tt.item0KeyDisplayName {
			t.Errorf("item.Key.DisplayName mismatch: want (%s), got (%s)", tt.item0KeyDisplayName, item0.Key.DisplayName)
		}
		if item0.Title != tt.item0Title {
			t.Errorf("item.Title mismatch: want (%s), got (%s)", tt.item0Title, item0.Title)
		}
		item0Created, err := item0.Created.Time()
		if err != nil {
			if err != nil {
				t.Errorf("DateTime8.Time() error: (%s)", err.Error())
			}
		}
		if item0Created.Format(time.RFC3339) != tt.item0Created {
			t.Errorf("item.Key.DisplayName mismatch: want (%s), got (%s)", tt.item0Created, item0Created.Format(time.RFC3339))
		}
		item0updated, err := item0.Updated.Time()
		if err != nil {
			if err != nil {
				t.Errorf("DateTime8.Time() error: (%s)", err.Error())
			}
		}
		if item0updated.Format(time.RFC3339) != tt.item0Updated {
			t.Errorf("item.Key.DisplayName mismatch: want (%s), got (%s)", tt.item0Updated, item0updated.Format(time.RFC3339))
		}
		for _, itemx := range j.Channel.Issues {
			testReadFileDuration(t, itemx.TimeEstimate, tt.hoursPerDay, tt.daysPerWeek)
			testReadFileDuration(t, itemx.TimeOriginalEstimate, tt.hoursPerDay, tt.daysPerWeek)
			testReadFileDuration(t, itemx.TimeSpent, tt.hoursPerDay, tt.daysPerWeek)
			testReadFileDuration(t, itemx.AggregateTimeOriginalEstimate, tt.hoursPerDay, tt.daysPerWeek)
			testReadFileDuration(t, itemx.AggregateTimeRemainingEstimate, tt.hoursPerDay, tt.daysPerWeek)
			testReadFileDuration(t, itemx.AggregateTimeSpent, tt.hoursPerDay, tt.daysPerWeek)
			/*
				if len(strings.TrimSpace(itemx.TimeEstimate.Display)) != 0 {
					di, err := timeutil.ParseDurationInfo(itemx.TimeEstimate.Display)
					if err != nil {
						t.Errorf("timeutil.ParseDurationInfo(\"%s\") error: (%s)", itemx.TimeEstimate.Display, err.Error())
					}
					durSec := di.Duration(tt.hoursPerDay, tt.daysPerWeek).Seconds()
					if int64(durSec) != itemx.TimeEstimate.Seconds {
						t.Errorf("timeutil.ParseDurationInfo(\"%s\") mismatch: want (%d), got (%d)", itemx.TimeEstimate.Display, itemx.TimeEstimate.Seconds, int64(durSec))
					}
				}
			*/
		}
	}
}

func testReadFileDuration(t *testing.T, d Duration, hoursPerDay, daysPerWeek float32) {
	if len(strings.TrimSpace(d.Display)) != 0 {
		di, err := timeutil.ParseDurationInfo(d.Display)
		if err != nil {
			t.Errorf("timeutil.ParseDurationInfo(\"%s\") error: (%s)", d.Display, err.Error())
		}
		durSec := di.Duration(hoursPerDay, daysPerWeek).Seconds()
		if int64(durSec) != d.Seconds {
			t.Errorf("timeutil.ParseDurationInfo(\"%s\") mismatch: want (%d), got (%d)", d.Display, d.Seconds, int64(durSec))
		}
	}
}
