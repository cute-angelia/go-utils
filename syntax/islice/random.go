package islice

import (
	"errors"
	"math/rand"
	"time"
)

// GetRandomElement 随机获取一个元素
func GetRandomElement[T any](s []T) T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return s[r.Intn(len(s))]
}

// GetRandomElementAndIndex 随机获取一个元素并获取索引
func GetRandomElementAndIndex[T any](s []T) (T, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(s))
	return s[index], index
}

// RandomWeightIndex 权重随机获取索引
func RandomWeightIndex[T int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64](weights []T) (int, error) {
	if len(weights) == 0 {
		return 0, errors.New("传入数组为空")
	}

	var sum T
	for _, w := range weights {
		sum += w
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v := r.Int63n(int64(sum) + 1)
	for i, w := range weights {
		v -= int64(w)
		if v <= 0 {
			return i, nil
		}
	}
	return 0, nil
}
