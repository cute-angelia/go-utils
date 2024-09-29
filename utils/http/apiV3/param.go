package apiV3

import (
	"net/http"
	"strconv"
)

// from header
func GetHeaderValue(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func GetUid(r *http.Request) int64 {
	uid, _ := strconv.Atoi(GetHeaderValue(r, "jwt_uid"))
	return int64(uid)
}

func GetUidV2[T int | int32 | int64](r *http.Request) T {
	uid, _ := strconv.ParseInt(GetHeaderValue(r, "jwt_uid"), 10, 64)
	uidStr := GetHeaderValue(r, "jwt_uid")
	if uidStr == "" {
		return 0
	}

	uid, err := strconv.ParseInt(uidStr, 10, 64)
	if err != nil {
		return 0
	}

	switch any(T(0)).(type) {
	case int:
		return T(int(uid))
	case int32:
		return T(int32(uid))
	case int64:
		return T(uid)
	default:
		return 0
	}
}

func GetAppId(r *http.Request) string {
	return GetHeaderValue(r, "jwt_appid")
}

func GetCid(r *http.Request) int32 {
	cid, _ := strconv.Atoi(GetHeaderValue(r, "jwt_cid"))
	return int32(cid)
}

// from query

// QueryString Query will get a query parameter by key.
func QueryString(r *http.Request, key string) string {
	return r.URL.Query().Get(key)
}

// QueryStringInt32 QueryInt will get a query parameter by key and convert it to an int or return an error.
// user?user_id=
func QueryStringInt32(r *http.Request, key string) int32 {
	v := r.URL.Query().Get(key)
	val, _ := strconv.Atoi(v)
	return int32(val)
}
