package apiV2

import (
	jsoniter "github.com/json-iterator/go"
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
}

func NewBody(r *http.Request) *body {
	b := body{
		r: r,
	}

	if r.Header.Get("Content-Type") == ContentTypeWWWForm {
		if r.PostForm == nil {
			r.ParseMultipartForm(32 << 20) // 32MB
		}
		b.dataWWWForm = r.Form
	}

	if r.Header.Get("Content-Type") == ContentTypeJson || strings.Contains(r.Header.Get("Content-Type"), ContentTypeJson) {
		b.isJson = true
		json := jsoniter.ConfigCompatibleWithStandardLibrary
		// 2. 解析json
		json.NewDecoder(r.Body).Decode(&b.dataJson)
	}

	return &b
}

// from post

func (b *body) PostString(key string) (val string) {
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
func (b *body) PostInt(key string) (val int) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(int); ok {
				val = strValue
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = ageInt
	}
	return val
}
func (b *body) PostInt32(key string) (val int32) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(int32); ok {
				val = strValue
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = int32(ageInt)
	}
	return val
}
func (b *body) PostInt64(key string) (val int64) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(int64); ok {
				val = strValue
			}
		}
	} else {
		ageInt, _ := strconv.Atoi(b.dataWWWForm.Get(key))
		val = int64(ageInt)
	}
	return val
}
func (b *body) PostFloat64(key string) (val float64) {
	if b.isJson {
		if value, ok := b.dataJson[key]; ok {
			if strValue, ok := value.(float64); ok {
				val = strValue
			}
		}
	} else {
		ageInt, _ := strconv.ParseFloat(b.dataWWWForm.Get(key), 10)
		val = ageInt
	}
	return val
}

// from header

func (b *body) GetUserUid() int32 {
	uid, _ := strconv.Atoi(b.r.Header.Get("jwt_uid"))
	return int32(uid)
}
func (b *body) GetUserUidInt64() int64 {
	uid, _ := strconv.Atoi(b.r.Header.Get("jwt_uid"))
	return int64(uid)
}
func (b *body) GetAppId() string {
	appid := b.r.Header.Get("jwt_appid")
	return appid
}
func (b *body) GetCid() int32 {
	cid, _ := strconv.Atoi(b.r.Header.Get("jwt_cid"))
	return int32(cid)
}

// from query

// QueryString Query will get a query parameter by key.
func (b *body) QueryString(key string) string {
	return b.r.URL.Query().Get(key)
}

// QueryStringInt32 QueryInt will get a query parameter by key and convert it to an int or return an error.
// user?user_id=
func (b *body) QueryStringInt32(key string) int32 {
	v := b.r.URL.Query().Get(key)
	val, _ := strconv.Atoi(v)
	return int32(val)
}

// from upload

func (b *body) UploadParseMultipartForm(maxMemory int64) error {
	return b.r.ParseMultipartForm(maxMemory)
}

func (b *body) Upload(key string) string {
	return b.r.FormValue(key)
}

func (b *body) Upload32(key string) int32 {
	p, _ := strconv.Atoi(b.r.FormValue(key))
	return int32(p)
}

func (b *body) PostArray(key string) []string {
	return b.r.PostForm[key]
}
