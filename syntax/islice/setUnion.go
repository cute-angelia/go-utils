package islice

// golang 数组之间的交集，差集，并集，补集

import (
	"sync"
)

type Set struct {
	sync.RWMutex
	m map[interface{}]bool
}

// 新建集合对象
func New(items ...interface{}) *Set {
	s := &Set{
		m: make(map[interface{}]bool, len(items)),
	}
	s.Add(items...)
	return s
}

// 添加元素
func (s *Set) Add(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		s.m[v] = true
	}
}

// 删除元素
func (s *Set) Remove(items ...interface{}) {
	s.Lock()
	defer s.Unlock()
	for _, v := range items {
		delete(s.m, v)
	}
}

// 判断元素是否存在
func (s *Set) Has(items ...interface{}) bool {
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
func (s *Set) Count() interface{} {
	return len(s.m)
}

// 清空集合
func (s *Set) Clear() {
	s.Lock()
	defer s.Unlock()
	s.m = map[interface{}]bool{}
}

// 空集合判断
func (s *Set) Empty() bool {
	return len(s.m) == 0
}

// 无序列表
func (s *Set) List() []interface{} {
	s.RLock()
	defer s.RUnlock()
	list := make([]interface{}, 0, len(s.m))
	for item := range s.m {
		list = append(list, item)
	}
	return list
}

// 排序列表
//func (s *Set) SortList() []interface{} {
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
func (s *Set) Union(sets ...*Set) *Set {
	r := New(s.List()...)
	for _, set := range sets {
		for e := range set.m {
			r.m[e] = true
		}
	}
	return r
}

// 差集
func (s *Set) Minus(sets ...*Set) *Set {
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
func (s *Set) Intersect(sets ...*Set) *Set {
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
func (s *Set) Complement(full *Set) *Set {
	r := New()
	for e := range full.m {
		if _, ok := s.m[e]; !ok {
			r.Add(e)
		}
	}
	return r
}
