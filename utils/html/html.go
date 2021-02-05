package html

import (
	"regexp"
	"strings"
)

// 获取第一张图片
func GetFirstImg(content string) string {
	reg, _ := regexp.Compile(`(?U)src="(.*)"`)
	imgs := reg.FindStringIndex(content)
	if len(imgs) > 0 {
		img_str := content[(imgs[0]):(imgs[1])]
		return strings.TrimLeft(strings.TrimRight(img_str, `"`), `imgsrc="`)
	} else {
		return ""
	}
}

// golang 去除html标签
func TrimHtml(src string, length int) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	returns := strings.TrimSpace(src)
	if length > 0 {
		nameRune := []rune(returns)

		if length > len(nameRune) {
			length = len(nameRune)
		}

		return string(nameRune[0:length])
	} else {
		return returns
	}
}
