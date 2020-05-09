package ip

import (
	"log"
	"strings"
	"time"
)

// https://best-ecology.oss-cn-hangzhou.aliyuncs.com/tool/ip/qqwry.dat

func GetInfo(ip string) map[string]ResultQQwry {
	IPData.FilePath = "/tmp/qqwry.dat"
	IPData.OnlineUrl = "https://best-ecology.oss-cn-hangzhou.aliyuncs.com/tool/ip/qqwry.dat"

	startTime := time.Now().UnixNano()
	res := IPData.InitIPData()

	if v, ok := res.(error); ok {
		log.Panic(v)
	}
	endTime := time.Now().UnixNano()

	log.Printf("IP 库加载完成 共加载:%d 条 IP 记录, 所花时间:%.1f ms\n", IPData.IPNum, float64(endTime-startTime)/1000000)

	ips := strings.Split(ip, ",")

	qqWry := NewQQwry()

	rs := map[string]ResultQQwry{}
	if len(ips) > 0 {
		for _, v := range ips {
			rs[v] = qqWry.Find(v)
		}
	}

	return rs
}
