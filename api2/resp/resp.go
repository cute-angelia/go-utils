package api

import (
	"net/http"
	"encoding/json"
)

func Success(w http.ResponseWriter, r *http.Request, data interface{}, msg string) {
	response := map[string]interface{}{
		"code": 0,
		"msg":  msg,
		"data": data,
	}

	// 日志
	// FIXED 优化,协程处理
	// ApiMakeLog.createLog(r, response)

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
		"code": code,
		"msg":  msg,
		"data": nil,
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
