package iregexp

import "regexp"

func GetRegMatch(text string, regMatch string) []string {
	// `href="(.*)">立即下载`    <input type="hidden" name="formhash" value="(.*)" />
	reg, _ := regexp.Compile(regMatch)
	result := reg.FindStringSubmatch(text)
	return result
}
