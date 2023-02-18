package islice

// RemoveDuplicateString 去重
// removeDuplicateElement函数总共初始化两个变量，一个长度为0的slice，一个空map。由于slice传参是按引用传递，没有创建占用额外的内存空间。
// map[string]struct{}{}创建了一个key类型为String值类型为空struct的map，等效于使用make(map[string]struct{})
// 空struct不占内存空间，使用它来实现我们的函数空间复杂度是最低的。
func RemoveDuplicateString(objs []string) []string {
	result := make([]string, 0, len(objs))
	temp := map[string]struct{}{}
	for _, item := range objs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}

func RemoveDuplicateInt32(objs []int32) []int32 {
	result := make([]int32, 0, len(objs))
	temp := map[int32]struct{}{}
	for _, item := range objs {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
