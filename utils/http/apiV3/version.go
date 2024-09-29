package apiV3

import (
	"strconv"
	"strings"
)

// CompareVersion 版本号比较 大于:1 等于:0 小于:-1
func CompareVersion(version1, version2 string) int {
	v1 := strings.Split(version1, ".")
	v2 := strings.Split(version2, ".")

	// 找出最长的版本号长度
	maxLen := len(v1)
	if len(v2) > maxLen {
		maxLen = len(v2)
	}

	for i := 0; i < maxLen; i++ {
		num1 := getVersionPart(v1, i)
		num2 := getVersionPart(v2, i)

		if num1 > num2 {
			return 1
		} else if num1 < num2 {
			return -1
		}
	}

	return 0
}

func getVersionPart(version []string, index int) int {
	if index >= len(version) {
		return 0
	}
	num, err := strconv.Atoi(version[index])
	if err != nil {
		return 0 // 或者可以选择在这里处理错误
	}
	return num
}
