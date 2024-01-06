package go_oura

import (
	"fmt"
	"time"
)

type TimeStamp string

func (t TimeStamp) toTime() (time.Time, error) {
	formats := []string{
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
	}

	for _, format := range formats {
		parsedTime, err := time.Parse(format, string(t))
		if err == nil {
			return parsedTime, nil
		}
	}

	return time.Time{},
		fmt.Errorf("unable to parse timestamp: %s", string(t))
}
