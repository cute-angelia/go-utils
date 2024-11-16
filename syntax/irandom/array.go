package irandom

import (
	"math"
)

// RandArray returns a random element from the given slice
func RandArray[T any](arr []T) T {
	index := math.Floor(r.Float64() * float64(len(arr)))
	return arr[int(index)]
}
