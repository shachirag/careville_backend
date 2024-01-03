package utils

import "time"

func ParseDate(dateStr string) (time.Time, error) {
	// Try parsing the date with the expected format "date/month/year"
	date, err := time.Parse("01-02-2006", dateStr)
	if err != nil {
		// If parsing with the expected format fails, try parsing with RFC3339 format
		date, err = time.Parse(time.RFC3339, dateStr)
		if err != nil {
			// If parsing with RFC3339 format also fails, return an error
			return time.Time{}, err
		}
	}

	return date, nil
}
