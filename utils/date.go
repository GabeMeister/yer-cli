package utils

import (
	"time"
)

func FormatISODate(isoString string) (string, error) {
	t, err := time.Parse(time.RFC3339, isoString)
	if err != nil {
		return "", err
	}

	return t.Format("January 2, 2006"), nil
}

func GetYearFromDateStr(isoString string) int {
	// Parse the ISO string
	parsed, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", isoString)
	if err != nil {
		return -1
	}

	// Extract the year from the parsed time
	return parsed.Year()
}

// isoString: Thu Mar 21 08:28:07 2024 -0700
// year: 2023
func IsDateStrInYear(isoString string, year int) bool {
	parsedYear := GetYearFromDateStr(isoString)

	// Compare the extracted year with the given year
	return parsedYear == year
}

// isoString: Thu Mar 21 08:28:07 2024 -0700
// year: 2023
func IsDateStrBeforeYear(isoString string, year int) bool {
	// Parse the ISO string
	parsed, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", isoString)
	if err != nil {
		return false
	}

	// Extract the year from the parsed time
	parsedYear := parsed.Year()

	// Compare the extracted year with the given year
	return parsedYear < year
}
