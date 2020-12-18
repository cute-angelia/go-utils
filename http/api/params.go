package api

import (
	"net/http"
	"strconv"
	"github.com/go-chi/chi"
)

// Query will get a query parameter by key.
func QueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// QueryInt will get a query parameter by key and convert it to an int or return an error.
// user?user_id=
func QueryStringInt32(r *http.Request, key string) int32 {
	v := r.URL.Query().Get(key)
	val, _ := strconv.Atoi(v)
	return int32(val)
}

// user/{user_id}
func Param(r *http.Request, name string) string {
	return chi.URLParam(r, name)
}

func ParamInt32(r *http.Request, name string) int32 {
	p, _ := strconv.Atoi(chi.URLParam(r, name))
	return int32(p)
}

func Post(r *http.Request, name string) string {
	//if len(requireds) > 0 {
	//	if vs := r.PostForm[name]; len(vs) > 0 {
	//		return vs[0], nil
	//	} else {
	//		return "", fmt.Errorf("缺少必传参数:" + name)
	//	}
	//}
	return r.PostFormValue(name)
}

func PostInt32(r *http.Request, name string) int32 {
	p, _ := strconv.Atoi(r.PostFormValue(name))
	return int32(p)
}

func PostInt64(r *http.Request, name string) int64 {
	p, _ := strconv.Atoi(r.PostFormValue(name))
	return int64(p)
}

func PostFloat64(r *http.Request, name string) float64 {
	p, _ := strconv.ParseFloat(r.PostFormValue(name), 64)
	return p
}

func GetUserUid(r *http.Request) int32 {
	uid, _ := strconv.Atoi(r.Header.Get("jwt_uid"))
	return int32(uid)
}
