package utils

import (
	"fmt"
	"time"
)

func StrToTime(timestamp string) (time.Time, error) {
	layout := "2006-01-02T15:04:05.000Z"
	time, err := time.Parse(layout, timestamp)
	if err != nil {
		return time, fmt.Errorf("%w failed to parse time string to time.Time", err)
	}

	return time, nil
}
