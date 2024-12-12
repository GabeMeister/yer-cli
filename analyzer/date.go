package analyzer

import (
	"time"
)

func GetNumWorkDaysInCurrYear() int {
	total := 0
	firstDayOfYear := time.Date(CURR_YEAR, time.January, 1, 0, 0, 0, 0, time.UTC)
	currDay := time.Now()

	tmpDay := firstDayOfYear
	for tmpDay.Before(currDay) {
		weekDay := tmpDay.Weekday()

		if weekDay != time.Saturday && weekDay != time.Sunday {
			total += 1
		}
		tmpDay = tmpDay.AddDate(0, 0, 1)
	}

	return total
}
