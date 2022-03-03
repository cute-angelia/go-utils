package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

var LogRequestAndData = true

func LogOn(on bool) {
	LogRequestAndData = on
}

func Success(w http.ResponseWriter, r *http.Request, data interface{}, msg string) {
	response := map[string]interface{}{
		"code":    0,
		"message": msg,
		"data":    data,
	}

	// 日志
	// FIXED 优化,协程处理
	// ApiMakeLog.createLog(r, response)
	if LogRequestAndData {
		logr(r, response, "[success]")
	}

	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 缓存返回
func SuccessCache(w http.ResponseWriter, code int, msg string, cacheData interface{}) {
	response := map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    cacheData,
	}
	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// 错误
func Error(w http.ResponseWriter, r *http.Request, data interface{}, msg string, code int32) {
	// 内部错误
	if code == 500 {
		http.Error(w, msg, 500)
		return
	}

	response := map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    data,
	}

	if LogRequestAndData {
		logr(r, response, "[error]")
	}

	// 日志
	// ApiMakeLog.createLog(r, response)

	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func logr(r *http.Request, response interface{}, msg string) {
	// 打印日志
	go func() {
		z, _ := json.Marshal(r.PostForm)
		z2, _ := json.Marshal(response)
		zuid := r.Header.Get("jwt_uid")

		log.Println("------------------------------------------------------------------------------")
		jwt_app_start_time := r.Header.Get("jwt_app_start_time")
		if len(jwt_app_start_time) > 0 {
			un, _ := strconv.Atoi(jwt_app_start_time)
			t2 := time.Unix(int64(un), 0)
			tc := time.Since(t2)

			flags := "Millisecond"
			if tc < time.Second {
				tc = tc / 10
			} else {
				flags = "Second"
			}

			log.Printf("用户: %s, TimeCost: %v %s", zuid, tc, flags)
		}

		log.Printf("%s 用户: %s, 请求地址: %s, 请求参数: %s, 请求数据: %s,", msg, zuid, r.URL.Path, r.URL.RawQuery, z)
		log.Printf("%s 用户: %s, 请求地址: %s, 响应数据: %s", msg, zuid, r.URL.Path, z2)
		log.Println("------------------------------------------------------------------------------")
	}()
}
