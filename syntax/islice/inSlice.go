package islice

func InSlice[T comparable](s []T, x T) bool {
	for _, v := range s {
		if v == x {
			return true
		}
	}
	return false
}

// InSliceBy returns true if predicate function return true.
func InSliceBy[T any](slice []T, predicate func(item T) bool) bool {
	for _, item := range slice {
		if predicate(item) {
			return true
		}
	}
	return false
}

// IndexOf returns the index of the first occurrence of v in s,
// or -1 if v is not present in s.
func IndexOf[T comparable](s []T, v T) int {
	for i, item := range s {
		if item == v {
			return i
		}
	}
	return -1
}
