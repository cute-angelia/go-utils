package itime

import (
	"strings"
	"time"
)

// 时间相关
const (
	TIME_FORMAT       = "2006-01-02 15:04:05"
	TIME_FORMAT_DATE  = "2006-01-02"
	TIME_FORMAT_TIME  = "15:04:05"
	TIME_FORMAT_MONTH = "2006-01"
)

func TimeZero() time.Time {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return t
}

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
