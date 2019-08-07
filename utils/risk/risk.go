package risk

import (
	"time"
	"github.com/patrickmn/go-cache"
	"sync"
	"fmt"
)

type riskRule struct {
	Key      string
	TTL      time.Duration
	MaxCount int
}

type riskCache struct {
	Count int
}

// 风控
type Risk struct {
	Rules map[string]riskRule
	Cache *cache.Cache
	Lock  sync.Mutex
}

func (self *Risk) getCacheKey(key string) string {
	return "risk_" + key
}

func (self *Risk) getRule(key string) riskRule {
	if v, ok := self.Rules[key]; ok {
		return v
	} else {
		return riskRule{}
	}
}

func (self *Risk) Increase(key string) error {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	cacheKey := self.getCacheKey(key)
	rule := self.getRule(key)

	if &rule != nil {
		if c, found := self.Cache.Get(cacheKey); found {
			d := c.(riskCache)
			d.Count += 1
			self.Cache.Set(self.getCacheKey(key), d, rule.TTL)
		} else {
			c := riskCache{
				Count: 1,
			}
			self.Cache.Set(self.getCacheKey(key), c, rule.TTL)
		}
		return nil
	} else {
		return fmt.Errorf("未发现规则: %s" , key)
	}
}

func (self *Risk) Check(key string) error {
	self.Lock.Lock()
	defer self.Lock.Unlock()

	cacheKey := self.getCacheKey(key)
	rule := self.getRule(key)

	if rule.MaxCount > 0 {
		if c, found := self.Cache.Get(cacheKey); found {
			d := c.(riskCache)
			if d.Count > rule.MaxCount {
				return fmt.Errorf("数量超过限制 now:%d => max:%d", d.Count, rule.MaxCount)
			}
		}
		return nil
	} else {
		return fmt.Errorf("未发现规则: %s" , key)
	}
}

type RiskOption func(*Risk)

func NewRisk(opts ...RiskOption) *Risk {
	var sopt Risk

	sopt.Rules = map[string]riskRule{}
	sopt.Cache = cache.New(5*time.Minute, 10*time.Minute)

	for _, opt := range opts {
		opt(&sopt)
	}

	return &sopt
}

func Rules(key string, ttl time.Duration, maxCount int) RiskOption {
	return func(options *Risk) {
		options.Rules[key] = riskRule{
			Key:      key,
			TTL:      ttl,
			MaxCount: maxCount,
		}
	}
}
