package helper

import (
	"strings"
	"time"
)

func ConverRFC3339(layout string, inTime string) string {
	t := ConverRFC3339ToTime(layout, inTime)
	return t.Format(layout)
}

func ConverRFC3339ToTime(layout string, inTime string) time.Time {
	if len(inTime) == 0 {
		t, _ := time.ParseInLocation(layout, "1970-01-01 00:00:00", time.Local)
		return t
	}

	if !strings.Contains(inTime, "T") {
		t, _ := time.ParseInLocation(layout, inTime, time.Local)
		return t
	}

	t, _ := time.ParseInLocation(time.RFC3339, inTime, time.Local)
	if t.Unix() < 0 {
		t, _ = time.ParseInLocation(layout, inTime, time.Local)
	}
	return t
}
