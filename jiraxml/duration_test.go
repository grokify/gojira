package jiraxml

/*
import (
	"testing"
)

var durationTests = []struct {
	whpd    float32
	wdpw    float32
	v       string
	seconds int64
}{
	{8.0, 5.0, "0 minutes", 0},
	{8.0, 5.0, "25 minutes", 1500},
	{8.0, 5.0, "30 minutes", 1800},
	{8, 5, "7 hours, 35 minutes", 27300},
	{8.0, 5.0, "1 day", 28800},
	{8.0, 5.0, "2 days", 57600},
	{8.0, 5.0, "3 days", 86400},
	{8.0, 5.0, "2 days, 3 hours, 25 minutes", 69900},
}

func TestDuration(t *testing.T) {
	for _, tt := range durationTests {
		di, err := ParseDurationInfo(tt.v)
		if err != nil {
			t.Errorf("jiraxml.ParseDurationInfo(\"%s\") error: (%s)", tt.v, err.Error())
		}
		d := di.Duration(uint(tt.whpd), uint(tt.wdpw))
		dursec := int64(d.Seconds())
		if dursec != tt.seconds {
			t.Errorf("DurationInfo.Duration(%d,%d) mismatch: want (%d), got (%d)", uint(tt.whpd), uint(tt.wdpw), tt.seconds, dursec)
		}
	}
}
*/
