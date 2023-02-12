package ibunt

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

// LockerOpts
// 次数限制
// 每日次数限制
type LockerOpts struct {
	Limit int   // 限制次数
	Today bool  // 是否限制每日
	Uid   int32 // 针对某人
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

func WithUid(uid int32) LockerOpt {
	return func(options *LockerOpts) {
		options.Uid = uid
	}
}

// IsNotLockedInLimit
// true => 非锁定状态，处理正常逻辑
// false => 锁定状态，处理错误逻辑
// if !bunt.IsNotLockedInLimit("cache", "SIGNPRE_REPEAT_"+nonce, time.Minute*60, bunt.NewLockerOpt(bunt.WithToday(true))) {
func IsNotLockedInLimit(dbname string, key string, ttl time.Duration, opt LockerOpts) (bool, error) {
	if _, ok := BuntCaches.Load(dbname); !ok {
		return true, fmt.Errorf("[%s] 缓存未初始化 %s", PackageName, dbname)
	}

	if opt.Today {
		key = fmt.Sprintf("%s_%s", key, time.Now().Format("2006-01-02"))
	}
	if opt.Uid > 0 {
		key = fmt.Sprintf("%s_%d", key, opt.Uid)
	}
	value := Get(dbname, key)
	if len(value) > 0 {
		n, _ := strconv.Atoi(value)
		if n >= opt.Limit {
			return false, nil
		} else {
			if err := Set(dbname, key, fmt.Sprintf("%d", n+1), ttl); err != nil {
				log.Println("IsLockedLimit error:", err.Error())
				return true, err
			}
			return true, nil
		}
	} else {
		if err := Set(dbname, key, "1", ttl); err != nil {
			log.Println("IsLockedLimit error:", err.Error())
			return true, err
		}
		return true, nil
	}
}

// IsLockedInLimit 被锁定
func IsLockedInLimit(dbname string, key string, ttl time.Duration, opt LockerOpts) bool {
	if b, e := IsNotLockedInLimit(dbname, key, ttl, opt); e != nil {
		return false
	} else {
		return !b
	}
}
