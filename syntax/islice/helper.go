package islice

import "errors"

// GetOneSafeInSlice 安全取回一个元素
func GetOneSafeInSlice[T comparable](arr []T, randomIndex int) (T, error) {
	length := len(arr)
	if randomIndex >= 0 && randomIndex < length {
		return arr[randomIndex], nil
	} else {
		// 处理索引越界的情况, 默认返回最后一个
		return arr[length-1], errors.New("索引越界")
	}
}
