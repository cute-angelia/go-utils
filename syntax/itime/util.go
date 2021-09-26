package itime

import (
	"fmt"
	"time"
)

//@brief：耗时统计函数
func TimeCost(startTime time.Time) (string, time.Duration) {
	tc := time.Since(startTime)
	return fmt.Sprintf("time cost = %v", tc), tc
}

// ElapsedTime calc elapsed time 计算运行时间消耗 单位 ms(毫秒)
func ElapsedTime(startTime time.Time) string {
	return fmt.Sprintf("time cost = %.3f", time.Since(startTime).Seconds()*1000)
}

// 视频时长多少秒 =》 转 slice
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
