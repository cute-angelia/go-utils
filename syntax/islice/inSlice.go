package islice

// golang 数组 Has

import (
	"sync"
)

type inSlice struct {
	data sync.Map
}

func NewInSlice(items ...interface{}) *inSlice {
	var islice = new(inSlice)
	islice.Add(items...)
	return islice
}

// Add 添加元素
func (s *inSlice) Add(items ...interface{}) {
	for _, v := range items {
		s.data.Store(v, true)
	}
}

// Remove 删除元素
func (s *inSlice) Remove(items ...interface{}) {
	for _, v := range items {
		s.data.Delete(v)
	}
}

// Has 判断元素是否存在
func (s *inSlice) Has(items ...interface{}) bool {
	for _, v := range items {
		if _, ok := s.data.Load(v); !ok {
			return false
		}
	}
	return true
}

func (s *inSlice) GetData() sync.Map {
	return s.data
}

// Count 元素个数
func (s *inSlice) Count() int {
	count := 0
	s.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// Clear 清空集合
func (s *inSlice) Clear() {
	s.data = sync.Map{}
}

// Empty 空集合判断
func (s *inSlice) Empty() bool {
	return s.Count() == 0
}

// List 无序列表
func (s *inSlice) List() []interface{} {
	var list []interface{}
	s.data.Range(func(key, value interface{}) bool {
		list = append(list, value)
		return true
	})
	return list
}
