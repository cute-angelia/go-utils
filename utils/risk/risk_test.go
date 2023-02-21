package risk

import (
	"testing"
	"time"
	"log"
)

func TestRisk(t *testing.T) {
	risk := NewRisk(
		Rules("reg", time.Minute*10, 10),
		Rules("login", time.Minute*10, 10),
	)

	// 循环20次
	key := "reg"
	i := 0
	for ; i < 20; i++ {
		//log.Println("increase", risk.Increase("reg"))
		risk.Increase(key)
	}
	log.Println("check", risk.Check(key))

	// 循环10次
	key = "login"
	i = 0
	for ; i < 10; i++ {
		risk.Increase(key)
	}
	log.Println("check", risk.Check(key))


	key = "other"
	i = 0
	for ; i < 10; i++ {
		risk.Increase(key)
	}
	log.Println("check", risk.Check(key))
}

//2019/08/07 12:59:01 check 数量超过限制 now:20 => max:10
//2019/08/07 12:59:01 check <nil>
//2019/08/07 12:59:01 check 未发现规则: other