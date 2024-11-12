package utils

import (
	"time"
)

var COMMIT_DATE_FORMAT = "Mon Jan 2 15:04:05 2006 -0700"

func FormatISODate(isoString string) (string, error) {
	t, err := time.Parse(time.RFC3339, isoString)
	if err != nil {
		return "", err
	}

	return t.Format("January 2, 2006"), nil
}

func GetYearFromDateStr(isoString string) int {
	// Parse the ISO string
	parsed, err := time.Parse(COMMIT_DATE_FORMAT, isoString)
	if err != nil {
		return -1
	}

	// Extract the year from the parsed time
	return parsed.Year()
}

// Returns "1/2/06"
func GetSimpleDateStr(isoString string) string {
	// Parse the ISO string
	parsed, err := time.Parse(COMMIT_DATE_FORMAT, isoString)
	if err != nil {
		return ""
	}

	// Extract the year from the parsed time
	return parsed.Format("1/2/06")
}

// Returns "January 2, 2006"
func GetHumanReadableDateStr(isoString string) string {
	// Parse the ISO string
	parsed, err := time.Parse(COMMIT_DATE_FORMAT, isoString)
	if err != nil {
		return ""
	}

	// Extract the year from the parsed time
	return parsed.Format("January 2, 2006")
}

// Returns "2006-01-02"
func GetMachineReadableDateStr(isoString string) string {
	// Parse the ISO string
	parsed, err := time.Parse(COMMIT_DATE_FORMAT, isoString)
	if err != nil {
		return ""
	}

	// Extract the year from the parsed time
	return parsed.Format("2006-01-02")
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

func GetDaysOfYear(year int) []string {
	dates := []string{}

	currDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	for currDate.Year() == year {
		dates = append(dates, currDate.Format("2006-01-02"))
		currDate = currDate.AddDate(0, 0, 1)
	}

	return dates
}

func GetWeeksOfYear() []int {
	weeks := []int{}
	for i := 1; i <= 52; i++ {
		weeks = append(weeks, i)
	}

	return weeks
}

func GetDateFromISOString(isoString string) time.Time {
	d, err := time.Parse(COMMIT_DATE_FORMAT, isoString)
	if err != nil {
		panic(err)
	}

	return d
}
