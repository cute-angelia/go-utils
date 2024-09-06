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

// Chunk creates a slice of elements split into groups the length of size.
// Play: https://go.dev/play/p/b4Pou5j2L_C
func Chunk[T any](slice []T, size int) [][]T {
	result := [][]T{}

	if len(slice) == 0 || size <= 0 {
		return result
	}

	for _, item := range slice {
		l := len(result)
		if l == 0 || len(result[l-1]) == size {
			result = append(result, []T{})
			l++
		}

		result[l-1] = append(result[l-1], item)
	}

	return result
}

// Compact creates a slice with all falsey values removed. The values false, nil, 0, and "" are falsey.
// Play: https://go.dev/play/p/pO5AnxEr3TK
func Compact[T comparable](slice []T) []T {
	var zero T

	result := make([]T, 0, len(slice))

	for _, v := range slice {
		if v != zero {
			result = append(result, v)
		}
	}
	return result[:len(result):len(result)]
}

// Concat creates a new slice concatenating slice with any additional slices.
// Play: https://go.dev/play/p/gPt-q7zr5mk
func Concat[T any](slices ...[]T) []T {
	totalLen := 0
	for _, v := range slices {
		totalLen += len(v)
		if totalLen < 0 {
			panic("len out of range")
		}
	}
	result := make([]T, 0, totalLen)

	for _, v := range slices {
		result = append(result, v...)
	}

	return result
}
