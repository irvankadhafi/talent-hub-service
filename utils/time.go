package utils

import "time"

// FormatTimeRFC3339 Format time according to RFC3339Nano
func FormatTimeRFC3339(t *time.Time) (s string) {
	if t == nil {
		return
	}

	if t.Nanosecond() == 0 {
		return t.Format("2006-01-02T15:04:05.000000000Z07:00")
	}

	return t.Format(time.RFC3339Nano)
}

// ParseDurationWithDefault parses a duration string and returns the parsed duration.
// If the parsing fails, it returns the default duration provided.
func ParseDurationWithDefault(input string, defaultDuration time.Duration) time.Duration {
	parsedDuration, err := time.ParseDuration(input)
	if err != nil {
		return defaultDuration
	}

	return parsedDuration
}

// ParseDate parses a date string using the specified layout and returns the corresponding time.Time value.
// If the parsing fails, it returns the zero time value.
func ParseDate(layout, value string) time.Time {
	date, _ := time.Parse(layout, value)

	return date
}

// ParseDatetimeToRFC3339 formats the provided time value to the RFC3339 format.
func ParseDatetimeToRFC3339(inputTime *time.Time) string {
	return inputTime.Format(time.RFC3339)
}
