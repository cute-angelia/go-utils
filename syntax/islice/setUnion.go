package islice

// golang 数组之间的交集，差集，并集，补集

import (
	"sync"
)

type SetUnion struct {
	sync.RWMutex
	m map[interface{}]bool
}

// 新建集合对象
func New(items ...interface{}) *SetUnion {
	s := &SetUnion{
		m: make(map[interface{}]bool, len(items)),
	}
	s.Add(items...)
	return s
}

// 添加元素
func (s *SetUnion) Add(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// 删除元素
func (s *SetUnion) Remove(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// 判断元素是否存在
func (s *SetUnion) Has(items ...interface{}) bool {
	s.RLock()
	defer s.RUnlock()
	for _, v := range items {
		if _, ok := s.m[v]; !ok {
			return false
		}
	}
	return true
}

// 元素个数
func (s *SetUnion) Count() interface{} {
	return len(s.m)
}

// 清空集合
func (s *SetUnion) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

// 空集合判断
func (s *SetUnion) Empty() bool {
	return len(s.m) == 0
}

// 无序列表
func (s *SetUnion) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := make([]interface{}, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// 排序列表
//func (s *SetUnion) SortList() []interface{} {
//	s.RLock()
//	defer s.RUnlock()
//	list := make([]interface{}, 0, len(s.m))
//	for item := range s.m {
//		list = append(list, item)
//	}
//	sort.Ints(list)
//	return list
//}

// 并集
func (s *SetUnion) Union(sets ...*SetUnion) *SetUnion {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// 差集
func (s *SetUnion) Minus(sets ...*SetUnion) *SetUnion {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			if _, ok := s.m[e]; ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// 交集
func (s *SetUnion) Intersect(sets ...*SetUnion) *SetUnion {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range s.m {
			if _, ok := set.m[e]; !ok {
				delete(r.m, e)
			}
		}
	}
	return r
}

// 补集
func (s *SetUnion) Complement(full *SetUnion) *SetUnion {
	r := New()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.Add(e)
		}
	}
	return r
}
