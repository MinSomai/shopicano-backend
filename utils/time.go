package utils

import "time"

const (
	DateTimeFormat                = time.RFC3339
	DateTimeFormatForDistribution = "Mon Jan/02/2006 15:04:05"
	DateTimeFormatForInput        = "02-01-2006T15:04"
	DateFormat                    = "02-01-2006"
)

func ParseDateTimeForInput(v string) (time.Time, error) {
	return time.Parse(DateTimeFormatForInput, v)
}
