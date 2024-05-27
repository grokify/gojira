package gojira

import (
	"time"
)

type Config struct {
	ServerURL          string
	WorkingHoursPerDay float32
	WorkingDaysPerWeek float32
	StatusConfig       *StatusConfig
}

func NewConfigDefault() *Config {
	return &Config{
		WorkingHoursPerDay: WorkingHoursPerDayDefault,
		WorkingDaysPerWeek: WorkingDaysPerWeekDefault}
}

func (c *Config) SecondsToDays(sec int) float32 {
	return float32(sec) / 60 / 60 / c.WorkingHoursPerDay
}

func (c *Config) SecondsToWeeks(sec int) float32 {
	return c.SecondsToDays(sec) / c.WorkingDaysPerWeek
}

func (c *Config) CapacityForDaysPeople(days, people float32) time.Duration {
	return time.Duration(days) * time.Duration(c.WorkingHoursPerDay) * // hours
		60 * 60 * time.Second
}
