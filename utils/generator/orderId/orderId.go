package orderId

import (
	"time"
	"fmt"
	"math/rand"
	"crypto/md5"
	"encoding/hex"
)

func GenerateOrderId() string {
	timenow := time.Now().Format("20060102150405")
	return fmt.Sprintf("%s%06v", timenow, rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
}

func GenerateTradeId() string {
	str := fmt.Sprintf("%d%06v", time.Now().Nanosecond(), rand.New(rand.NewSource(time.Now().UnixNano())).Int31n(1000000))
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}