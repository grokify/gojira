package jiraxml

import (
	"testing"
	"time"
)

var readFileTests = []struct {
	filename             string
	itemCount            int
	itemCountInProgress  int
	itemCountNeedsTriage int
	item0KeyDisplayName  string
	item0Title           string
	item0Created         string
	item0Updated         string
}{
	{"testdata/example_jira_mongodb.xml", 20, 3, 6, "SERVER-79445", "[SERVER-79445] Consider upgrading PCRE2",
		"2023-07-28T00:11:05Z", "2023-07-28T01:00:50Z"},
}

func TestReadFile(t *testing.T) {
	for _, tt := range readFileTests {
		j, err := ReadFile(tt.filename)
		if err != nil {
			t.Errorf("jiraxml.ReadFile(\"%s\") error: (%v)", tt.filename, err.Error())
		}
		// fmtutil.PrintJSON(j)
		stats := j.Channel.Items.Stats()
		if stats.ItemCount != tt.itemCount {
			t.Errorf("jiraxml.ReadFile(\"%s\") mismatch: want (%d), got (%d)", tt.filename, tt.itemCount, stats.ItemCount)
		}
		// fmtutil.PrintJSON(stats)
		if tt.itemCount <= 0 {
			continue
		}
		// fmtutil.PrintJSON(j.Channel.Items[0])
		item0 := j.Channel.Items[0]
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
	}
}
