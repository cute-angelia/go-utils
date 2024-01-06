package islice

import (
	"math/rand"
	"time"
)

func GetRandomElement[T any](s []T) T {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return s[r.Intn(len(s))]
}

func GetRandomElementAndIndex[T any](s []T) (T, int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	index := r.Intn(len(s))
	return s[index], index
}

func RandomWeightIndex(weights []int) int {
	sum := 0
	for _, w := range weights {
		sum += w
	}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	v := r.Intn(sum + 1)
	for i, w := range weights {
		v -= w
		if v <= 0 {
			return i
		}
	}
	return 0
}
