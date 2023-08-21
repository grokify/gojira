package jiraxml

type Config struct {
	WorkingHoursPerDay float32
	WorkingDaysPerWeek float32
}

func NewConfigDefault() *Config {
	return &Config{
		WorkingHoursPerDay: WorkingHoursPerDayDefault,
		WorkingDaysPerWeek: WorkingDaysPerWeekDefault}
}
