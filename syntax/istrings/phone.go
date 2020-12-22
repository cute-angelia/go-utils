package istrings

import "strings"

// 隐藏手机号码
func FormatHiddenMobileText(str string) string {
	if len(str) > 5 {
		i := strings.Index(str, "+")
		if i > -1 {
			return str[:6] + "****" + str[len(str)-4:]
		} else {
			return str[:3] + "****" + str[len(str)-4:]
		}
	} else {
		return str
	}
}
