package apiV2

import (
	"net/http"
	"strconv"
	"strings"
)

// from header

func GetHeaderValue(r *http.Request, key string) string {
	return r.Header.Get(key)
}

func GetLoginUid(r *http.Request) int32 {
	uid, _ := strconv.Atoi(GetHeaderValue(r, "jwt_uid"))
	return int32(uid)
}

func GetLoginUidInt64(r *http.Request) int64 {
	uid, _ := strconv.Atoi(GetHeaderValue(r, "jwt_uid"))
	return int64(uid)
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

func CompareVersion(version1 string, version2 string) int {
	versionA := strings.Split(version1, ".")
	versionB := strings.Split(version2, ".")

	for i := len(versionA); i < 4; i++ {
		versionA = append(versionA, "0")
	}
	for i := len(versionB); i < 4; i++ {
		versionB = append(versionB, "0")
	}
	for i := 0; i < 4; i++ {
		version1, _ := strconv.Atoi(versionA[i])
		version2, _ := strconv.Atoi(versionB[i])
		if version1 == version2 {
			continue
		} else if version1 > version2 {
			return 1
		} else {
			return -1
		}
	}
	return 0
}
