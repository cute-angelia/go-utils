package apiV3

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/cute-angelia/go-utils/utils/generator/random"
	"github.com/cute-angelia/go-utils/utils/iAes"
	"github.com/cute-angelia/go-utils/utils/iXor"
	"log"
	"net/http"
	"strconv"
	"time"
)

type render struct {
	w http.ResponseWriter
	r *http.Request

	isHasPage bool // 是否分页
	pager     Pagination

	cryptoType int32  // 加密方式：默认2
	cryptoKey  string // 是否加密：不为空为加密

	isLogOn bool // 打印日志

	reqStruct  any // 请求结构体
	respStruct Res // 返回结构体
}

// Res 标准JSON输出格式
type Res struct {
	// Code 响应的业务错误码。0表示业务执行成功，非0表示业务执行失败。
	Code int32 `json:"code"`
	// Msg 响应的参考消息。前端可使用msg来做提示
	Msg string `json:"msg"`
	// Data 响应的具体数据
	Data interface{} `json:"data"`

	Pagination *Pagination `json:"pagination,omitempty"`

	Ext *Ext `json:"ext"`
}

type Ext struct {
	ShowTips bool `json:"showTips"` // 弹消息提示
}

type Pagination struct {
	//  当前页
	PageNo int `json:"pageNo"`
	// PageSize 每页记录数
	PageSize int `json:"pageSize"`
	// PageTotal 总页数
	PageTotal int `json:"pageTotal"`
	// 总条数
	Count int `json:"count"`
}

// CalcTotal 计算总页数
func (p Pagination) CalcTotal(count int, pageSize int) int {
	var totalPages int

	if pageSize == 0 {
		pageSize = 1
	}

	if count%pageSize == 0 {
		totalPages = count / pageSize
	} else {
		totalPages = count/pageSize + 1
	}
	return totalPages
}

func NewPagination(count int, pageNo, pageSize int) Pagination {
	paginationor := Pagination{PageNo: pageNo, PageSize: pageSize, Count: count}
	paginationor.PageTotal = paginationor.CalcTotal(count, pageSize)
	return paginationor
}

func NewRender(w http.ResponseWriter, r *http.Request) *render {
	return &render{
		w:          w,
		r:          r,
		isLogOn:    true,                              // 默认：打开日志
		cryptoType: 2,                                 // 默认：加密方式
		cryptoKey:  CryptoEr.GetRequestContentType(r), // 默认：获取key
	}
}

// Decode request
func (that *render) Decode(v any) error {
	body, err := Decoder.Decode(that.r, v)
	that.reqStruct = body
	return err
}

// Success 成功返回
func (that *render) Success() {
	that.respStruct.Code = 0
	// 日志
	if that.isLogOn {
		that.logr("[success]")
	}

	// 加密
	that.cryptoData()

	// json
	that.w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(that.w).Encode(that.respStruct); err != nil {
		http.Error(that.w, err.Error(), 500)
		return
	}
}

func (that *render) ErrorCodeMsg(code int32, msg string) {
	err := NewApiError(code, msg)
	that.Error(err)
}

func (that *render) Error(err error) {
	that.respStruct.Code = -1

	if err != nil {
		var e *ApiError
		if errors.As(err, &e) {
			// 可以访问e.Code和e.Message
			that.respStruct.Code = e.Code
		}
		that.respStruct.Msg = err.Error()
	}

	if that.isLogOn {
		that.logr("[error]")
	}

	// 加密
	that.cryptoData()

	// json
	that.w.Header().Set("Content-Type", "application/json")

	if err = json.NewEncoder(that.w).Encode(that.respStruct); err != nil {
		http.Error(that.w, err.Error(), 500)
		return
	}
}

// SetData set data
func (that *render) SetData(data interface{}, msg string) *render {
	that.respStruct.Data = data
	that.respStruct.Msg = msg
	return that
}

func (that *render) SetPage(pager *Pagination) *render {
	that.respStruct.Pagination = pager
	return that
}

func (that *render) SetExt(ext *Ext) *render {
	that.respStruct.Ext = ext
	return that
}

func (that *render) SetLog(on bool) {
	that.isLogOn = on
}
func (that *render) SetCryptoType(cryptoType int32) *render {
	that.cryptoType = cryptoType
	return that
}
func (that *render) SetCryptoKey(cryptoKey string) *render {
	that.cryptoKey = cryptoKey
	return that
}

func (that *render) cryptoData() {
	crypto := that.r.URL.Query().Get("crypto")
	if len(crypto) > 0 {
		var randomKey = random.RandString(16, random.LetterAll)
		cryptoId := fmt.Sprintf("%s%s", that.cryptoKey, randomKey)
		datam, _ := json.Marshal(that.respStruct.Data)

		// Crypto 加密 Key：使用AES-GCM模式,处理密钥、认证、加密一次完成
		if that.cryptoType == 1 {
			encryptData, _ := iAes.EncryptCBCToBase64(datam, []byte(cryptoId))
			that.respStruct.Data = randomKey + encryptData
		}
		// xor
		if that.cryptoType == 2 {
			encryptData := iXor.XorEncrypt(datam, cryptoId)
			that.respStruct.Data = randomKey + encryptData
		}
	}
}

func (that *render) logr(msg string) {
	// 数据
	dataReq, _ := json.Marshal(that.reqStruct)
	dataResp, _ := json.Marshal(that.respStruct)

	// header
	uid := that.r.Header.Get("jwt_uid")
	appStartTime := that.r.Header.Get("jwt_app_start_time")

	log.Println("------------------------------------------------------------------------------")
	if len(appStartTime) > 0 {
		un, _ := strconv.Atoi(appStartTime)
		t2 := time.Unix(int64(un), 0)
		tc := time.Since(t2)

		flags := "Millisecond"
		if tc < time.Second {
			tc = tc / 10
		} else {
			flags = "Second"
		}

		log.Printf("用户: %s, TimeCost: %v %s", uid, tc, flags)
	}

	log.Printf("%s 用户: %s, 请求地址: %s, 请求参数: %s, 请求数据: %s,", msg, uid, that.r.URL.Path, that.r.URL.RawQuery, dataReq)
	log.Printf("%s 用户: %s, 请求地址: %s, 响应数据: %s", msg, uid, that.r.URL.Path, dataResp)
	log.Println("------------------------------------------------------------------------------")
}
