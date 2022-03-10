package itime

import (
	"sort"
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
	// year
	t.Log(GetYearDay())

	// 获得当前月的初始和结束日期
	t.Log(GetMonthDay())

	// 周
	t.Log(GetWeekDay())

	// 季度
	t.Log(GetQuarterDay())

	// between
	dates := GetBetweenDates("2022-02-07", "2022-03-10")
	sort.Sort(sort.Reverse(sort.StringSlice(dates)))
	t.Log(dates)
}
