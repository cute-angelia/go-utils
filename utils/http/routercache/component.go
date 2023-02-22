package routercache

import (
	"bytes"
	"encoding/gob"
	"encoding/hex"
	"encoding/json"
	"github.com/schollz/progressbar/v3"
	"hash/fnv"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sort"
	"strings"
	"time"
)

type Component struct {
	config *config
	bar    *progressbar.ProgressBar
}

// newComponent ...
func newComponent(config *config) *Component {
	comp := &Component{}
	comp.config = config
	return comp
}

// NewMiddleware is the HTTP cache middleware handler.
func (c *Component) NewMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if c.cacheableMethod(r.Method) {
			c.sortURLParams(r.URL)
			key := c.generateKey(r)
			if r.Method == http.MethodPost && r.Body != nil {
				body, err := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				if err != nil {
					next.ServeHTTP(w, r)
					return
				}
				reader := ioutil.NopCloser(bytes.NewBuffer(body))
				key = c.generateKeyWithBody(r, body)
				r.Body = reader
			}

			params := r.URL.Query()
			if _, ok := params[c.config.RefreshKey]; ok {
				delete(params, c.config.RefreshKey)

				r.URL.RawQuery = params.Encode()
				key = c.generateKey(r)

				c.config.Store.Delete(key)
			} else {
				b, err := c.config.Store.Get(key)
				response := c.bytesToResponse([]byte(b))
				if err == nil {
					if response.Expiration.After(time.Now()) {
						response.LastAccess = time.Now()
						response.Frequency++
						c.config.Store.Set(key, response.String(), response.Expiration.Sub(time.Now()))

						for k, v := range response.Header {
							w.Header().Set(k, strings.Join(v, ","))
						}
						if c.config.WriteExpiresHeader {
							w.Header().Set("Expires", response.Expiration.UTC().Format(http.TimeFormat))
						}
						w.WriteHeader(response.StatusCode)
						w.Write(response.Value)

						if c.config.PrintLog {
							z, _ := json.Marshal(r.PostForm)
							z2, _ := json.Marshal(response)
							zuid := r.Header.Get("jwt_uid")
							log.Printf("%s 用户: %s, 请求地址: %s, 请求参数: %s, 请求数据: %s,", "[success cache]", zuid, r.URL.Path, r.URL.RawQuery, z)
							log.Printf("%s 用户: %s, 请求地址: %s, 响应数据: %s", "[success cache]", zuid, r.URL.Path, z2)
						}

						return
					}

					c.config.Store.Delete(key)
				}
			}

			rec := httptest.NewRecorder()
			next.ServeHTTP(rec, r)
			result := rec.Result()

			statusCode := result.StatusCode
			value := rec.Body.Bytes()
			now := time.Now()
			expires := now.Add(c.config.Ttl)

			if c.config.StatusCodeFilter(statusCode) {
				response := Response{
					Value:      value,
					Header:     result.Header,
					StatusCode: statusCode,
					Expiration: expires,
					LastAccess: now,
					Frequency:  1,
				}
				c.config.Store.Set(key, response.String(), response.Expiration.Sub(time.Now()))
			}
			for k, v := range result.Header {
				w.Header().Set(k, strings.Join(v, ","))
			}
			if c.config.WriteExpiresHeader {
				w.Header().Set("Expires", expires.UTC().Format(http.TimeFormat))
			}
			w.WriteHeader(statusCode)
			w.Write(value)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// DeleteCustomKey 删除自定义key缓存
func (c *Component) DeleteCustomKey(key string) {
	c.config.Store.Delete(key)
}

func (c *Component) cacheableMethod(method string) bool {
	for _, m := range c.config.Methods {
		if method == m {
			return true
		}
	}
	return false
}

func (c *Component) sortURLParams(URL *url.URL) {
	params := URL.Query()
	for _, param := range params {
		sort.Slice(param, func(i, j int) bool {
			return param[i] < param[j]
		})
	}
	URL.RawQuery = params.Encode()
}

func (c *Component) generateKey(r *http.Request) string {
	hash := fnv.New64a()
	if len(c.config.CustomKey) > 0 {
		hash.Write([]byte(c.config.CustomKey))
	} else {
		hash.Write([]byte(r.URL.String() + r.Header.Get("Authorization")))
	}

	return hex.EncodeToString(hash.Sum(nil))
}

func (c *Component) generateKeyWithBody(r *http.Request, body []byte) string {
	hash := fnv.New64a()
	if len(c.config.CustomKey) > 0 {
		body = []byte(c.config.CustomKey)
	} else {
		body = append([]byte(r.URL.String()+r.Header.Get("Authorization")), body...)
	}
	hash.Write(body)
	return hex.EncodeToString(hash.Sum(nil))
}

// BytesToResponse converts bytes array into Response data structure.
func (c *Component) bytesToResponse(b []byte) Response {
	var r Response
	dec := gob.NewDecoder(bytes.NewReader(b))
	dec.Decode(&r)

	return r
}
