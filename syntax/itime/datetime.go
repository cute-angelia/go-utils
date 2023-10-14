package itime

import (
	"log"
	"strings"
	"time"
)

// @referer https://github.com/duke-git/lancet/blob/main/datetime/conversion.go

type TheTime struct {
	unix int64
}

// NewUnixNow return unix timestamp of current time
func NewUnixNow() *TheTime {
	return &TheTime{unix: time.Now().Unix()}
}

func NewTime(t time.Time) *TheTime {
	return NewUnix(t.Unix())
}

// NewUnix return unix timestamp of specified time
func NewUnix(unix int64) *TheTime {
	return &TheTime{unix: unix}
}

// NewFormat return unix timestamp of specified time string, t should be "yyyy-mm-dd hh:mm:ss"
func NewFormat(t string) (*TheTime, error) {
	timeLayout := "2006-01-02 15:04:05"

	// 包含 T
	if strings.Contains(t, "T") {
		timeLayout = time.RFC3339
	}

	loc := time.FixedZone("CST", 8*3600)
	tt, err := time.ParseInLocation(timeLayout, t, loc)
	if err != nil {
		return nil, err
	}
	return &TheTime{unix: tt.Unix()}, nil
}

// NewFormatLayout 根据模板获取时间 强制+8时区
func NewFormatLayout(t string, timeLayout string) *TheTime {
	// 包含 T
	if strings.Contains(t, "T") {
		timeLayout = time.RFC3339
	}
	loc := time.FixedZone("CST", 8*3600)
	tt, err := time.ParseInLocation(timeLayout, t, loc)
	if err != nil {
		log.Println("NewFormatLayout error", t, timeLayout, err)
	}
	return &TheTime{unix: tt.Unix()}
}

// NewISO8601 return unix timestamp of specified iso8601 time string
func NewISO8601(iso8601 string) (*TheTime, error) {
	t, err := time.ParseInLocation(time.RFC3339, iso8601, time.UTC)
	if err != nil {
		return nil, err
	}
	return &TheTime{unix: t.Unix()}, nil
}

// GetTime return time.Time
func (t *TheTime) GetTime() time.Time {
	return time.Unix(t.unix, 0)
}

// GetUnix return unix timestamp
func (t *TheTime) GetUnix() int64 {
	return t.unix
}

// GetTimeZero 获取零点
func (t *TheTime) GetTimeZero() *TheTime {
	timeStr := t.FormatDate()
	t1, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	return NewTime(t1)
}

// GetMonthStartDay 获得当前月的初始天
func (t *TheTime) GetMonthStartDay() *TheTime {
	now := t.GetTime()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	return NewUnix(firstOfMonth.Unix())
}

// GetMonthEndDay 获得当前月的终
func (t *TheTime) GetMonthEndDay() *TheTime {
	now := t.GetTime()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return NewUnix(lastOfMonth.Unix())
}

// GetWeekDayChinese 中国的周码 （1，2，3，4，5，6，7）
func (t *TheTime) GetWeekDayChinese() int {
	now := t.GetTime()
	if now.Weekday() == 0 {
		return 7
	} else {
		return int(now.Weekday())
	}
}

// GetWeekStartDay 获得当前周的初始和结束日期
func (t *TheTime) GetWeekStartDay() *TheTime {
	now := t.GetTime()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}

	lastoffset := int(time.Saturday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}

	firstOfWeek := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, offset)
	return NewUnix(firstOfWeek.Unix())
}

// GetWeekEndDay 获得当前周的初始和结束日期
func (t *TheTime) GetWeekEndDay() *TheTime {
	now := t.GetTime()
	offset := int(time.Monday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if offset > 0 {
		offset = -6
	}
	lastoffset := int(time.Saturday - now.Weekday())
	//周日做特殊判断 因为time.Monday = 0
	if lastoffset == 6 {
		lastoffset = -1
	}
	lastOfWeeK := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, lastoffset+1)
	return NewUnix(lastOfWeeK.Unix())
}

// --------------------- format -----------------------

// Format return the time string 'yyyy-mm-dd hh:mm:ss' of unix time
func (t *TheTime) Format() string {
	return time.Unix(t.unix, 0).Format("2006-01-02 15:04:05")
}

func (t *TheTime) FormatDate() string {
	return time.Unix(t.unix, 0).Format("2006-01-02")
}

func (t *TheTime) FormatTime() string {
	return time.Unix(t.unix, 0).Format("15:04:05")
}

func (t *TheTime) FormatTimeZeroSec() string {
	return time.Unix(t.unix, 0).Format("2006-01-02 00:00:00")
}

func (t *TheTime) FormatTimeLastSec() string {
	return time.Unix(t.unix, 0).Format("2006-01-02 23:59:59")
}

// FormatForTpl return the time string which format is specified tpl
func (t *TheTime) FormatForTpl(tpl string) string {
	return time.Unix(t.unix, 0).Format(tpl)
}

// FormatIso8601 return iso8601 time string
func (t *TheTime) FormatIso8601() string {
	return time.Unix(t.unix, 0).Format(time.RFC3339)
}
