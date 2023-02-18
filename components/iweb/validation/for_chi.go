package validation

import (
	"net/http"
	"strconv"
)

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

func PostInt(r *http.Request, name string) int {
	p, _ := strconv.Atoi(r.PostFormValue(name))
	return p
}

func PostFloat64(r *http.Request, name string) float64 {
	p, _ := strconv.ParseFloat(r.PostFormValue(name), 64)
	return p
}
