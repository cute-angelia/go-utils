package islice

// EqualOrdered 顺序敏感的 slice 比较（元素顺序必须一致）
func EqualOrdered[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// EqualUnordered 顺序无关的 slice 比较（忽略元素顺序）
func EqualUnordered[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}
	// 使用 map 计数
	count := make(map[T]int)
	for _, v := range a {
		count[v]++
	}
	for _, v := range b {
		count[v]--
		if count[v] < 0 {
			return false
		}
	}
	return true
}

// EqualSliceWithFunc 使用自定义比较函数
// 示例使用
//
//	func main() {
//		type Person struct{ Name string; Age int }
//
//		p1 := []Person{{"Alice", 30}, {"Bob", 25}}
//		p2 := []Person{{"Alice", 30}, {"Bob", 26}}
//
//		// 只比较名字
//		nameEqual := EqualSliceWithFunc(p1, p2,
//			func(a, b Person) bool { return a.Name == b.Name })
//
//		// 完整比较
//		fullEqual := EqualSliceWithFunc(p1, p2,
//			func(a, b Person) bool { return a == b })
//
//		fmt.Println("只比较名字:", nameEqual) // true
//		fmt.Println("完整比较:", fullEqual)   // false
//	}
func EqualSliceWithFunc[T any](a, b []T, equal func(T, T) bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if !equal(a[i], b[i]) {
			return false
		}
	}
	return true
}
