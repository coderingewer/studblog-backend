package utils

import (
	"time"
)

func CaluclateLastDay(tim time.Time) bool {
	now := time.Now()
	currentYear, currentMonth, currentDay := now.Date()
	day := currentDay - tim.Day()
	year := currentYear - tim.Year()
	month := currentMonth - tim.Month()
	total := day + year + int(month)
	if total < 1 {
		return true
	}
	return false
}
