package irandom

// RandInt generate random int between min and max, maybe min,  not be max
func RandInt(min, max int) int {
	if min == max {
		return min
	}
	if max < min {
		min, max = max, min
	}
	return r.Intn(max-min) + min
}

// 随机float64 [0.0,1.0)
func RandFloat64() float64 {
	return r.Float64()
}
