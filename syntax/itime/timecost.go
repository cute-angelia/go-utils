package itime

import (
	"fmt"
	"time"
)

// TimeCost @brief：耗时统计函数
func TimeCost(startTime time.Time) (string, time.Duration) {
	tc := time.Since(startTime)
	return fmt.Sprintf("time cost = %v", tc), tc
}

// ElapsedTime calc elapsed time 计算运行时间消耗 单位 ms(毫秒)
func ElapsedTime(startTime time.Time) string {
	return fmt.Sprintf("time cost = %.3f", time.Since(startTime).Seconds()*1000)
}
