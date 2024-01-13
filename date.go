package go_oura

/*
	Simply here to convert YYYY-MM-DD dates in API returns from string to time.Time
*/

import (
	"strings"
	"time"
)

type Date struct {
	time.Time
}

func (d *Date) UnmarshalJSON(input []byte) error {
	strInput := strings.Trim(string(input), `"`)

	if strInput == "null" || strInput == "" {
		return nil
	}

	newTime, err := time.Parse("2006-01-02", strInput)
	if err != nil {
		return err
	}

	d.Time = newTime
	return nil
}
