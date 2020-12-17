package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

var LogRequestAndData = false

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
		logr(r, response)
	}

	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

/**
错误
*/
func Error(w http.ResponseWriter, r *http.Request, data interface{}, msg string, code int32) {
	// 内部错误
	if code == 500 {
		http.Error(w, msg, 500)
		return
	}

	response := map[string]interface{}{
		"code":    code,
		"message": msg,
		"data":    nil,
	}

	if LogRequestAndData {
		logr(r, response)
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

func logr(r *http.Request, response interface{}) {
	// 打印日志
	go func() {
		z, _ := json.Marshal(r.PostForm)
		z2, _ := json.Marshal(response)
		zuid := r.Header.Get("jwt_uid")

		log.Println("------------------------------------------------------------------------------")
		log.Printf("用户: %s, 请求地址: %s?%s", zuid, r.URL.Path, r.URL.RawQuery)

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

		log.Printf("用户: %s, 请求数据: %s,", zuid, z)
		log.Printf("用户: %s, 响应数据: %s", zuid, z2)
		log.Println("------------------------------------------------------------------------------")
	}()
}
