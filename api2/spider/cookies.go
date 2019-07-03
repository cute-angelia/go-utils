package spider

import (
	"net/http"
	"fmt"
	"log"
)

/**
	打印日志
 */
func LogCookie(postion string, cookies []*http.Cookie) {
	str := "\n"
	for _, v := range cookies {
		str += fmt.Sprintf("%s=%s;\n", v.Name, v.Value)
	}
	log.Println(postion, str)
}

/*
	Make Request Set Cookie
 */
func SetCookie(cookies []*http.Cookie, req *http.Request) {
	for _, c := range cookies {
		req.AddCookie(c)
	}
}

/**
	Merge Cookie
 */
func MergeCookie(cookies []*http.Cookie, cookiesNew []*http.Cookie) []*http.Cookie {
	// 合并新的 cookies
	cookies = append(cookies, cookiesNew...)

	// 倒序过滤
	keys := make(map[string]bool)
	for i := len(cookies) - 1; i >= 0; i-- {
		//log.Println("name",i,keys[cookies[i].Name],  cookies[i].Name)
		if _, value := keys[cookies[i].Name]; !value {
			keys[cookies[i].Name] = true
			continue
		} else {
			cookies = append(cookies[:i], cookies[i+1:]...)
		}
	}
	return cookies
}
