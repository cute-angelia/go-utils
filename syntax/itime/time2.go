package itime

// from https://github.com/gookit/goutil/blob/master/fmtutil/time.go

import (
	"fmt"
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

var timeMessages = []string{
	"< 1 sec", "1 sec", "secs", "1 min", "mins", "1 hr", "hrs", "1 day", "days",
}

// HowLongAgo format a seconds, get how lang ago
func HowLongAgo(sec int64) string {
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

// 100s => []string{00,01,40}
func SecToStr(sec int64) []string {
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
