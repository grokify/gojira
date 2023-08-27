package gojira

import (
	"time"

	"github.com/grokify/mogo/strconv/strconvutil"
)

type Config struct {
	BaseURL            string
	WorkingHoursPerDay float32
	WorkingDaysPerWeek float32
}

func NewConfigDefault() *Config {
	return &Config{
		WorkingHoursPerDay: WorkingHoursPerDayDefault,
		WorkingDaysPerWeek: WorkingDaysPerWeekDefault}
}

func (c *Config) SecondsToDays(sec int) float32 {
	return float32(sec) / 60 / 60 / c.WorkingHoursPerDay
}

func (c *Config) SecondsToDaysString(sec int) string {
	return strconvutil.FormatFloat64Simple(float64(c.SecondsToDays(sec)))
}

func (c *Config) CapacityForDaysPeople(days, people float32) time.Duration {
	return time.Duration(days) * time.Duration(c.WorkingHoursPerDay) * // hours
		60 * 60 * time.Second
}
