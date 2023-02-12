package apicache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/components/cache/ibunt"
	"github.com/gotomicro/ego/core/elog"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"
)

type Component struct {
	name   string
	config *config
	logger *elog.Component
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config, logger *elog.Component) *Component {
	return &Component{
		name:   compName,
		config: config,
		logger: logger,
	}
}

// 生成缓存 KEY generateCacheKey
func (e *Component) getSelfCacheKey() string {
	gCacheKey := fmt.Sprintf("%s%s", e.config.Prefix, e.config.CacheKey)
	if e.config.OnlyToday {
		gCacheKey = fmt.Sprintf("%s%s-%s", e.config.Prefix, e.config.CacheKey, time.Now().Format("20060102"))
	}
	return gCacheKey
}

// DEBUG
func (e *Component) debug(topic string, msg string) {
	if e.config.Debug {
		e.logger.Info(msg, elog.FieldKey(topic))
	}
}

type response struct {
	Code int             `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}

func (e Component) resp(w http.ResponseWriter, code int, msg string, cacheData string) {
	rdata := response{
		Code: code,
		Msg:  msg,
		Data: []byte(cacheData),
	}
	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(rdata); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// get cache
func (e *Component) GetCache() string {
	e.debug("===Get Cache===", e.getSelfCacheKey())
	data := ibunt.Get(e.config.DbName, e.getSelfCacheKey())
	if len(data) > 6 {
		return data
	} else {
		return ""
	}
}

// get cache and write
func (e *Component) GetCacheAndWriter(w http.ResponseWriter, msg string) (string, error) {
	e.debug(e.getSelfCacheKey()+"get cache", "start get cache")
	data := ibunt.Get(e.config.DbName, e.getSelfCacheKey())
	if len(data) > 6 {
		log.Println(e.getSelfCacheKey() + "get cache -> got")
		e.debug(e.getSelfCacheKey()+"get cache -> got", data)
		e.resp(w, 0, msg, data)
		return data, nil
	}
	return data, errors.New("读取缓存数据，数据不存在")
}

func (e *Component) SetCache(data interface{}) error {
	// prefix cache
	defer func() {
		if len(e.config.Prefix) > 0 {
			cacheData := ibunt.Get(e.config.DbName, e.config.Prefix)
			if len(cacheData) > 0 {
				cacheDatas := strings.Split(cacheData, "|")
				if len(cacheDatas) >= e.config.PrefixMaxNum {
					cacheDatas = cacheDatas[len(cacheDatas)-e.config.PrefixMaxNum : len(cacheDatas)]
				}
				cacheDatas = append(cacheDatas, e.getSelfCacheKey())
				ibunt.Set(e.config.DbName, e.config.Prefix, strings.Join(cacheDatas, "|"), e.config.Timeout)
			} else {
				ibunt.Set(e.config.DbName, e.config.Prefix, e.getSelfCacheKey(), e.config.Timeout)
			}
		}
	}()

	ds, _ := json.Marshal(data)
	e.debug("===Set Cache===", e.getSelfCacheKey()+" -> "+string(ds))
	return ibunt.Set(e.config.DbName, e.getSelfCacheKey(), string(ds), e.config.Timeout)
}

func (e *Component) DeleteCache() error {
	return ibunt.Delete(e.config.DbName, e.getSelfCacheKey())
}

func (e *Component) DeleteCacheAll() error {
	if len(e.config.Prefix) > 0 {
		cacheData := ibunt.Get(e.config.DbName, e.config.Prefix)
		if len(cacheData) > 0 {
			cacheDatas := strings.Split(cacheData, "|")
			for _, i2 := range cacheDatas {
				ibunt.Delete(e.config.DbName, i2)
			}
		}
		ibunt.Delete(e.config.DbName, e.config.Prefix)
	}
	return nil
}
