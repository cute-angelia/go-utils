package apicache

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/cache/bunt"
	"github.com/gotomicro/ego/core/elog"
	"net/http"
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

// 获取 缓存 key
func (e *Component) getKey() string {
	if e.config.OnlyToday {
		return fmt.Sprintf("%s-%s", e.config.CacheKey, time.Now().Format("20060102"))
	} else {
		return e.config.CacheKey
	}
}

func (e *Component) debug(topic string, msg string) {
	if e.config.Debug {
		e.logger.Debug(msg, elog.FieldKey(topic))
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
func (e *Component) GetCache() (string, error) {
	e.debug(e.getKey()+"get cache", "start get cache")
	data := bunt.Get(e.config.DbName, e.getKey())
	if len(data) > 6 {
		e.debug(e.getKey()+"get cache -> got", data)
		return data, nil
	} else {
		e.debug(e.getKey()+"get cache -> gone", "数据不存在")
		return "", errors.New("数据不存在")
	}
}

// get cache and write
func (e *Component) GetCacheAndWriter(w http.ResponseWriter, msg string) (string, error) {
	e.debug(e.getKey()+"get cache", "start get cache")
	data := bunt.Get(e.config.DbName, e.getKey())
	if len(data) > 6 {
		e.debug(e.getKey()+"get cache -> got", data)
		e.resp(w, 0, msg, data)
		return data, nil
	} else {
		e.debug(e.getKey()+"get cache -> gone", "数据不存在")
		return "", errors.New("数据不存在")
	}
}

func (e *Component) SetCache(data interface{}) error {
	e.debug(e.getKey()+"set cache", "start set cache")
	ds, _ := json.Marshal(data)
	return bunt.Set(e.config.DbName, e.getKey(), string(ds), e.config.Timeout)
}
