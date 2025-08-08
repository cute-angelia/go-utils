package islice

// Unique 唯一
func Unique[T comparable](objs []T) []T {
	result := make([]T, 0, len(objs))
	temp := make(map[T]struct{})

	for _, item := range objs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
