package itime

import (
	"fmt"
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

// ConvertVideoSecToStr 视频时长多少秒 =》 转 slice
// 100s => []string{00,01,40}
func ConvertVideoSecToStr(sec int64) []string {
	if sec <= 0 {
		return []string{}
	}
	lastsec := sec
	h := sec / 3600
	if h > 0 {
		lastsec = lastsec - h*3600
	}
	m := lastsec / 60
	if m > 0 {
		lastsec = lastsec - m*60
	}
	s := lastsec
	return []string{fmt.Sprintf("%02d", h), fmt.Sprintf("%02d", m), fmt.Sprintf("%02d", s)}
}
