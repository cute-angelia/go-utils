package itime

import (
	"strings"
	"testing"
)

/*
=== RUN   TestCycle
    itime_test.go:7: 2021-07-01 00:00:00 2021-07-31 23:59:59
    itime_test.go:10: 2021-07-12 00:00:00 2021-07-18 23:59:59
    itime_test.go:13: 2021-07-01 00:00:00 2021-09-30 23:59:59
--- PASS: TestCycle (0.00s)
*/
func TestCycle(t *testing.T) {
	// 获得当前月的初始和结束日期
	t.Log(GetMonthDay())

	// 周
	t.Log(GetWeekDay())

	// 季度
	t.Log(GetQuarterDay())
}

func TestTime(t *testing.T) {
	t.Log( strings.Join( SecToStr(10), ":"))
	t.Log( strings.Join( SecToStr(60), ":"))
	t.Log( strings.Join( SecToStr(70), ":"))
	t.Log( strings.Join( SecToStr(100), ":"))
	t.Log( strings.Join( SecToStr(1000), ":"))
	t.Log( strings.Join( SecToStr(10000), ":"))
	t.Log( strings.Join( SecToStr(86400), ":"))
	t.Log( strings.Join( SecToStr(86401), ":"))
}