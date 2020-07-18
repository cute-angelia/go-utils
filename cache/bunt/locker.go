/*
* 限制类
  次数限制
  每日次数限制
*/
package bunt

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type LockerOpts struct {
	Limit int  // 限制次数
	Today bool // 是否限制每日
}

type LockerOpt func(opts *LockerOpts)

func NewLockerOpt(opts ...LockerOpt) LockerOpts {
	var sopt LockerOpts
	for _, opt := range opts {
		opt(&sopt)
	}
	if sopt.Limit == 0 {
		sopt.Limit = 1
	}
	return sopt
}

func WithLimit(limit int) LockerOpt {
	return func(options *LockerOpts) {
		options.Limit = limit
	}
}

func WithToday(today bool) LockerOpt {
	return func(options *LockerOpts) {
		options.Today = today
	}
}

/**
	没有锁住 处理业务逻辑
	true => do somethings
 */
func IsNotLockedInLimit(dbname string, key string, ttl time.Duration, opt LockerOpts) bool {
	if opt.Today {
		key = fmt.Sprintf("%s_%s", key, time.Now().Format("2006-01-02"))
	}
	value := Get(dbname, key)
	if len(value) > 0 {
		n, _ := strconv.Atoi(value)
		if n >= opt.Limit {
			return false
		} else {
			if err := Set(dbname, key, fmt.Sprintf("%d", n+1), ttl); err != nil {
				log.Println("IsLockedLimit error:", err.Error())
			}
			return true
		}
	} else {
		if err := Set(dbname, key, "1", ttl); err != nil {
			log.Println("IsLockedLimit error:", err.Error())
		}
		return true
	}
}
