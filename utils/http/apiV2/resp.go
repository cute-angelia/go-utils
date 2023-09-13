package apiV2

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var LogRequestAndData = true

func LogOn(on bool) {
	LogRequestAndData = on
}

// Res 标准JSON输出格式
type Res struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int `json:"code"`
	// Msg 响应的参考消息。前端可使用msg来做提示
	Msg string `json:"msg"`
	// Data 响应的具体数据
	Data interface{} `json:"data"`
}

// ResPage 带分页的标准JSON输出格式
type ResPage struct {
	Res
	Pagination Pagination `json:"pagination"`
}

type Pagination struct {
	// Current 当前页
	Current int `json:"current"`
	// PageSize 每页记录数
	PageSize int `json:"pageSize"`
	// Total 总页数
	Total int64 `json:"total"`
}

// Success 成功返回
func Success(w http.ResponseWriter, r *http.Request, data interface{}, msg string) {
	response := Res{
		Code: 0,
		Msg:  msg,
		Data: data,
	}
	doResp(w, r, response)
}

// SuccessWithPage 成功分页返回
func SuccessWithPage(w http.ResponseWriter, r *http.Request, data interface{}, msg string, pager Pagination) {
	response := ResPage{
		Res: Res{
			Code: 0,
			Msg:  msg,
			Data: data,
		},
		Pagination: pager,
	}
	doResp(w, r, response)
}

func doResp(w http.ResponseWriter, r *http.Request, response interface{}) {
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
	response := Res{
		Code: 0,
		Msg:  msg,
		Data: cacheData,
	}
	// json
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

// Error 错误
func Error(w http.ResponseWriter, r *http.Request, err error) {
	response := Res{
		Code: -1,
	}

	if err != nil {
		if e, ok := err.(*ApiError); ok {
			// 可以访问e.Code和e.Message
			response.Code = e.Code
		}
		response.Msg = err.Error()
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
	// 从context获取body
	var z []byte
	if r.Header.Get("Content-Type") == ContentTypeWWWForm {
		if len(r.PostForm) > 0 {
			z, _ = json.Marshal(r.PostForm)
		} else {
			if err := r.ParseForm(); err != nil {
				log.Println(err)
			}
			z, _ = json.Marshal(r.PostForm)
		}
	}
	if r.Header.Get("Content-Type") == ContentTypeJson || strings.Contains(r.Header.Get("Content-Type"), ContentTypeJson) {
		z, _ = io.ReadAll(r.Body)
	}

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
}
