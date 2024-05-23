package apiV3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/utils/generator/random"
	"github.com/cute-angelia/go-utils/utils/iAes"
	"github.com/cute-angelia/go-utils/utils/iXor"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type response struct {
	w http.ResponseWriter
	r *http.Request

	isHasPage bool // 是否分页
	pager     Pagination

	CryptoType int32  // 加密方式：默认2
	CryptoKey  string // 是否加密：不为空为加密

	isLogOn bool // 打印日志

	res Res // 返回结构体
}

// Res 标准JSON输出格式
type Res struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int32 `json:"code"`
	// Msg 响应的参考消息。前端可使用msg来做提示
	Msg string `json:"msg"`
	// Data 响应的具体数据
	Data interface{} `json:"data"`

	Pagination Pagination `json:"pagination,omitempty"`
}

type Pagination struct {
	// Current 当前页
	Current int `json:"current"`
	// PageSize 每页记录数
	PageSize int `json:"pageSize"`
	// Total 总页数
	Total int64 `json:"total"`
}

func NewResponse(w http.ResponseWriter, r *http.Request) *response {
	return &response{
		w: w,
		r: r,
	}
}

func (that *response) SetData(data interface{}, msg string) *response {
	that.res.Data = data
	that.res.Msg = msg
	return that
}

// Success 成功返回
func (that *response) Success() {
	that.res.Code = 0
	// 日志
	if that.isLogOn {
		that.logr("[success]")
	}

	// 加密
	that.cryptoData()

	// json
	that.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(that.w).Encode(that.res); err != nil {
		http.Error(that.w, err.Error(), 500)
		return
	}
}

func (that *response) Error(err error) {
	res := Res{
		Code: -1,
	}
	if err != nil {
		var e *ApiError
		if errors.As(err, &e) {
			// 可以访问e.Code和e.Message
			res.Code = e.Code
		}
		res.Msg = err.Error()
	}

	if that.isLogOn {
		that.logr("[error]")
	}

	// 加密
	that.cryptoData()

	// json
	that.w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(that.w).Encode(res); err != nil {
		http.Error(that.w, err.Error(), 500)
		return
	}
}

func (that *response) SetPage(pager Pagination) {
	that.res.Pagination = pager
}
func (that *response) SetLog(on bool) {
	that.isLogOn = on
}

func (that *response) cryptoData() {
	// crypto := r.URL.Query().Get("crypto")
	if len(that.CryptoKey) > 0 {
		var randomKey = random.RandString(16, random.LetterAll)
		cryptoId := fmt.Sprintf("%s%s", that.CryptoKey, randomKey)
		datam, _ := json.Marshal(that.res.Data)

		// Crypto 加密 Key：使用AES-GCM模式,处理密钥、认证、加密一次完成
		if that.CryptoType == 1 {
			encryptData, _ := iAes.EncryptCBCToBase64(datam, []byte(cryptoId))
			that.res.Data = randomKey + encryptData
		}
		// xor
		if that.CryptoType == 2 {
			encryptData := iXor.XorEncrypt(datam, cryptoId)
			that.res.Data = randomKey + encryptData
		}
	}
}

func (that *response) logr(msg string) {
	// 从context获取body
	var z []byte

	contentType := ContentTyper.GetRequestContentType(that.r)
	if contentType == ContentTypeForm {
		if len(that.r.PostForm) > 0 {
			z, _ = json.Marshal(that.r.PostForm)
		} else {
			if err := that.r.ParseForm(); err != nil {
				log.Println(err)
			}
			z, _ = json.Marshal(that.r.PostForm)
		}
	}
	if contentType == ContentTypeJSON {
		z, _ = io.ReadAll(that.r.Body)
	}

	z2, _ := json.Marshal(that.res)
	zuid := that.r.Header.Get("jwt_uid")

	log.Println("------------------------------------------------------------------------------")
	jwt_app_start_time := that.r.Header.Get("jwt_app_start_time")
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

	log.Printf("%s 用户: %s, 请求地址: %s, 请求参数: %s, 请求数据: %s,", msg, zuid, that.r.URL.Path, that.r.URL.RawQuery, z)
	log.Printf("%s 用户: %s, 请求地址: %s, 响应数据: %s", msg, zuid, that.r.URL.Path, z2)
	log.Println("------------------------------------------------------------------------------")
}
