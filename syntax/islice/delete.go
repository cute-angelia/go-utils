package islice

// 介意顺序
func RemoveIntWithOrder(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
}

// 不介意顺序
func RemoveIntWithoutOrder(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func RemoveStringWithOrder(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func RemoveStringWithoutOrder(s []string, i int) []string {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}
