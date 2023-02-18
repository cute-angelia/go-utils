package itime

// from https://github.com/gookit/goutil/blob/master/fmtutil/time.go

import (
	"fmt"
)

type Language int32

const (
	LanguageEn  = Language(0)
	LanguageChs = Language(1)
)

var timeFormats = [][]int{
	{0},
	{1},
	{2, 1},
	{60},
	{120, 60},
	{3600},
	{7200, 3600},
	{86400},
	{172800, 86400},
}

var timeMessagesEn = []string{
	"< 1 sec", "1 sec", "secs", "1 min", "mins", "1 hr", "hrs", "1 day", "days",
}
var timeMessagesChs = []string{
	"< 1 秒",
	"1 秒",
	"秒",
	"1 分钟",
	"分钟",
	"1 小时",
	"小时",
	"1 天",
	"天",
}

// HowLongAgo format a seconds, get how lang ago
func HowLongAgo(sec int64, lang Language) string {
	timeMessages := timeMessagesEn
	if lang == LanguageChs {
		timeMessages = timeMessagesChs
	}

	intVal := int(sec)
	length := len(timeFormats)

	for i, item := range timeFormats {
		if intVal >= item[0] {
			ni := i + 1
			match := false

			if ni < length { // next exists
				next := timeFormats[ni]
				if intVal < next[0] { // current <= intVal < next
					match = true
				}
			} else if ni == length { // current is last
				match = true
			}

			if match { // match success
				if len(item) == 1 {
					return timeMessages[i]
				}

				// len is 2
				return fmt.Sprintf("%d %s", intVal/item[1], timeMessages[i])
			}
		}
	}

	return "unknown" // He should never happen
}
