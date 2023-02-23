package apicache

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Component struct {
	name   string
	config *config
	locker sync.Mutex
}

// newComponent ...
func newComponent(compName string, config *config) *Component {
	return &Component{
		name:   compName,
		config: config,
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

type response struct {
	Code int             `json:"code"`
	Msg  string          `json:"message"`
	Data json.RawMessage `json:"data"`
}

func (e *Component) resp(w http.ResponseWriter, code int, msg string, cacheData string) {
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

// Deprecated: GetCache 获取缓存
func (e *Component) GetCache() string {
	data, _ := e.config.Cache.Get(e.getSelfCacheKey())
	if len(data) > 6 {
		return data
	} else {
		return ""
	}
}

// Deprecated: GetCacheAndWriter get cache and write
func (e *Component) GetCacheAndWriter(w http.ResponseWriter, msg string) (string, error) {
	data, _ := e.config.Cache.Get(e.getSelfCacheKey())
	if len(data) > 6 {
		e.resp(w, 0, msg, data)
		return data, nil
	}
	return data, errors.New("读取缓存数据，数据不存在")
}

// Deprecated: SetCache
func (e *Component) SetCache(data interface{}) error {
	ds, _ := json.Marshal(data)
	return e.config.Cache.Set(e.getSelfCacheKey(), string(ds), e.config.Timeout)
}

// Deprecated: DeleteCache
func (e *Component) DeleteCache() error {
	return e.config.Cache.Delete(e.getSelfCacheKey())
}

// Deprecated: DeleteCacheAll
func (e *Component) DeleteCacheAll() error {
	return e.config.Cache.Flush()
}
