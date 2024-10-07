package utils

import "time"

func FormatISODate(isoString string) (string, error) {
	t, err := time.Parse(time.RFC3339, isoString)
	if err != nil {
		return "", err
	}

	return t.Format("January 2, 2006"), nil
}
