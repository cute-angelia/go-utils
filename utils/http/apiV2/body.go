package apiV2

import (
	"bytes"
	"encoding/json"
	"github.com/gorilla/schema"
	jsoniter "github.com/json-iterator/go"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const (
	ContentTypeWWWForm = "application/x-www-form-urlencoded"
	ContentTypeJson    = "application/json"
)

type body struct {
	r           *http.Request
	isJson      bool
	dataJson    map[string]interface{}
	dataWWWForm url.Values
	dataBody    []byte
}

var jsonLib = jsoniter.ConfigCompatibleWithStandardLibrary

// SetMemory 设置内存
func SetMemory(r *http.Request, max int64) {
	r.ParseMultipartForm(max)
}

func NewBody(thatReq *http.Request) *body {
	b := body{
		r: thatReq,
	}
	if thatReq.Header.Get("Content-Type") == ContentTypeWWWForm {
		if thatReq.PostForm == nil {
			thatReq.ParseMultipartForm(32 << 20) // 32MB
		}
		b.dataWWWForm = thatReq.Form
	}
	if thatReq.Header.Get("Content-Type") == ContentTypeJson || strings.Contains(thatReq.Header.Get("Content-Type"), ContentTypeJson) {
		b.isJson = true

		buf, _ := io.ReadAll(thatReq.Body)
		b.dataBody = buf

		thatReq.Body = io.NopCloser(bytes.NewReader(buf))
	}
	return &b
}

// Decode 解码器
func (b *body) Decode(dst interface{}) error {
	if b.isJson {
		d := jsonLib.NewDecoder(bytes.NewReader(b.dataBody))
		d.UseNumber()
		return d.Decode(dst)
	} else {
		var decoder = schema.NewDecoder()
		return decoder.Decode(dst, b.r.PostForm)
	}
}

func (b *body) EncodeJson() (val string) {
	if b.isJson {
		val1, _ := jsonLib.Marshal(b.dataJson)
		return string(val1)
	} else {
		val2, _ := jsonLib.Marshal(b.dataWWWForm)
		return string(val2)
	}
}

// Get Value
func (b *body) GetString(key string) (val string) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(string); ok {
				val = strValue
			}
		}
	} else {
		val = b.dataWWWForm.Get(key)
	}
	return val
}

// PostSlice 数组，组合可能很多，有map ，strings ，int
func (b *body) GetSlice(key string) (val []interface{}) {
	if b.isJson {
		if _, ok := b.dataJson[key]; !ok {
			return []interface{}{}
		} else {
			val = b.dataJson[key].([]interface{})
		}
		log.Println(val)
	} else {
		val = append(val, b.dataWWWForm.Get(key))
	}
	return val
}
func (b *body) GetBool(key string) (val bool) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(bool); ok {
				val = strValue
			}
		}
	} else {
		val, _ = strconv.ParseBool(b.dataWWWForm.Get(key))
	}
	return val
}
func (b *body) GetInt(key string) (val int) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, err := value.(json.Number).Int64(); err == nil {
				val = int(strValue)
			} else {
				log.Println(err)
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = ageInt
	}
	return val
}
func (b *body) GetInt32(key string) (val int32) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, err := value.(json.Number).Int64(); err == nil {
				val = int32(strValue)
			} else {
				log.Println(err)
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = int32(ageInt)
	}
	return val
}
func (b *body) GetInt64(key string) (val int64) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, err := value.(json.Number).Int64(); err == nil {
				val = strValue
			} else {
				log.Println(err)
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = int64(ageInt)
	}
	return val
}
func (b *body) GetFloat64(key string) (val float64) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, err := value.(json.Number).Float64(); err == nil {
				val = strValue
			}
		}
	} else {
		ageInt, _ := strconv.ParseFloat(b.dataWWWForm.Get(key), 10)
		val = ageInt
	}
	return val
}
func (b *body) GetUploadString(key string) string {
	return b.r.FormValue(key)
}
func (b *body) GetUploadInt32(key string) int32 {
	p, _ := strconv.Atoi(b.r.FormValue(key))
	return int32(p)
}
