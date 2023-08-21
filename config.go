package gojira

import "github.com/grokify/mogo/strconv/strconvutil"

type Config struct {
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
