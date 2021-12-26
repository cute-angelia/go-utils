package islice

// golang 数组 Has

import (
	"sync"
)

type SetSyncMap struct {
	data sync.Map
}

// 新建集合对象
func NewSetSyncMap(items ...interface{}) *Set {
	s := &Set{}
	s.Add(items...)
	return s
}

// 添加元素
func (s *SetSyncMap) Add(items ...interface{}) {
	for _, v := range items {
		s.data.Store(v, true)
	}
}

// 删除元素
func (s *SetSyncMap) Remove(items ...interface{}) {
	for _, v := range items {
		s.data.Delete(v)
	}
}

// 判断元素是否存在
func (s *SetSyncMap) Has(items ...interface{}) bool {
	for _, v := range items {
		if _, ok := s.data.Load(v); !ok {
			return false
		}
	}
	return true
}

func (s *SetSyncMap) GetData() sync.Map {
	return s.data
}

// 元素个数
func (s *SetSyncMap) Count() int {
	count := 0
	s.data.Range(func(key, value interface{}) bool {
		count++
		return true
	})
	return count
}

// 清空集合
func (s *SetSyncMap) Clear() {
	s.data = sync.Map{}
}

// 空集合判断
func (s *SetSyncMap) Empty() bool {
	return s.Count() == 0
}

// 无序列表
func (s *SetSyncMap) List() []interface{} {
	var list []interface{}
	s.data.Range(func(key, value interface{}) bool {
		list = append(list, value)
		return true
	})
	return list
}
