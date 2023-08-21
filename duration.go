package jiraxml

/*
import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grokify/mogo/time/timeutil"
)

// DurationInfo represents information for a human-readable duration as presented in Jira.
type DurationInfo struct {
	Weeks        int
	Days         int
	Hours        int
	Minutes      int
	Seconds      int
	Milliseconds int
	Microseconds int
	Nanoseconds  int
}

// Duration returns a `time.Duration` struct with an optional `workingHoursPerDay` and `workingDaysPerWeek`
// valuees. Default values are available as
func (di DurationInfo) Duration(workingHoursPerDay, workingDaysPerWeek uint) time.Duration {
	dur := time.Duration(di.Nanoseconds) +
		time.Duration(di.Microseconds)*time.Microsecond +
		time.Duration(di.Milliseconds)*time.Millisecond +
		time.Duration(di.Seconds)*time.Second +
		time.Duration(di.Minutes)*time.Minute +
		time.Duration(di.Hours)*time.Hour
	if di.Days != 0 {
		if workingHoursPerDay != 0 {
			dur += time.Duration(di.Days) * time.Duration(workingHoursPerDay) * time.Hour
		} else {
			dur += time.Duration(di.Days) * timeutil.DurationDay
		}
	}
	if di.Weeks != 0 {
		if workingDaysPerWeek != 0 {
			daysPerWeek := time.Duration(workingDaysPerWeek)
			if workingHoursPerDay != 0 {
				dur += time.Duration(di.Weeks) *
					time.Duration(workingDaysPerWeek) *
					time.Duration(workingHoursPerDay) *
					time.Hour
			} else {
				dur += time.Duration(di.Weeks) *
					daysPerWeek *
					timeutil.DurationDay
			}
		} else {
			dur += time.Duration(di.Weeks) * timeutil.DurationWeek
		}
	}
	return dur
}

// ParseDurationInfo converts a Jira human readable string into a `DurationInfo` struct.
func ParseDurationInfo(s string) (DurationInfo, error) {
	parts := strings.Split(strings.ToLower(s), ",")
	di := DurationInfo{}
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if len(p) == 0 {
			continue
		}
		ps := strings.Fields(p)
		if len(ps) != 2 {
			return di, fmt.Errorf("cannot parse (%s)", p)
		}
		v, err := strconv.Atoi(ps[0])
		if err != nil {
			return di, err
		}
		switch ps[1] {
		case "week", "weeks":
			di.Days = v
		case "day", "days":
			di.Days = v
		case "hour", "hours", "h":
			di.Hours = v
		case "minute", "minutes", "m":
			di.Minutes = v
		case "second", "seconds", "s":
			di.Seconds = v
		case "millisecond", "milliseconds", "ms":
			di.Milliseconds = v
		case "microsecond", "microseconds", "us", "µs":
			di.Microseconds = v
		case "nanosecond", "nanoseconds", "ns":
			di.Nanoseconds = v
		default:
			return di, fmt.Errorf("cannot parse (%s)", p)
		}
	}
	return di, nil
}
*/
