package iurl

import (
	"fmt"
	"github.com/cute-angelia/go-utils/syntax/istrings"
	"net/url"
	"strings"
)

// CleanUrlWithoutParm 获得一个干净的不带参数的地址
func CleanUrlWithoutParm(uri string) string {
	clearUrl, _ := url.QueryUnescape(uri)
	if uriInfo, err := url.Parse(clearUrl); err == nil {
		if len(uriInfo.Scheme) == 0 {
			uriInfo.Scheme = "https"
		}
		return fmt.Sprintf("%s://%s%s", uriInfo.Scheme, uriInfo.Host, uriInfo.Path)
	} else {
		if beforeUrl, _, ok := istrings.Cut(clearUrl, `?`); ok {
			if !strings.Contains(beforeUrl, "http") {
				beforeUrl = "https:" + beforeUrl
			}
			return beforeUrl
		} else {
			return clearUrl
		}
	}
}
