package ilinked

import "fmt"

type SingleList struct {
	head *ListNode
	len  int
}

func (lst *SingleList) Init() {
	lst.head = new(ListNode)
	lst.head.lst = lst
	lst.len = 0
}

func New() *SingleList {
	lst := new(SingleList)
	lst.Init()
	return lst
}

// Clear ----------清空---------------
func (lst *SingleList) Clear() {
	lst.Init()
}

// Front -----------返回第一个结点--------
// 头结点的下一个，可能为nil
func (lst *SingleList) Front() *ListNode {
	return lst.head.Next()
}

// Back ------------返回尾结点---------
// 最后一个，可能为nil
func (lst *SingleList) Back() *ListNode {
	if lst.len == 0 {
		return nil
	}
	cur := lst.head.Next()
	for {
		if cur.Next() == nil {
			return cur
		}
		cur = cur.Next()
	}
}

// -----------插入---------
// insert at后插入e，返回e，仅在此函数增加长度
func (lst *SingleList) insert(e, at *ListNode) *ListNode {
	lst.len++
	e.next = at.next
	at.next = e
	e.lst = lst
	return e
}

// InsertAfter ------------指定元素后插--------
// mark后插入val
// 插入成功返回插入元素，否则返回nil
func (lst *SingleList) InsertAfter(val interface{}, mark *ListNode) *ListNode {
	if mark == nil {
		return nil
	}
	if mark.lst == lst {
		return lst.insert(&ListNode{Val: val}, mark)
	}
	return nil
}

// PushBack --------------尾插--------------
// 调用InsertAfter和Back进行val的尾插
// 成功返回插入的元素
func (lst *SingleList) PushBack(val interface{}) *ListNode {
	end := lst.Back()
	if end == nil {
		return lst.InsertAfter(val, lst.head)
	} else {
		return lst.InsertAfter(val, end)
	}
}

// PushFront --------------头插------------
// 调用InsertAfter进行val的头插
// 成功返回插入的元素，否则返回nil
func (lst *SingleList) PushFront(val interface{}) *ListNode {
	return lst.InsertAfter(val, lst.head)
}

// remove -------------------------删除-----------------------
// 仅在此函数减少长度，返回删除的元素或nil
func (lst *SingleList) remove(e *ListNode) *ListNode {
	nextOne := e.next
	if nextOne == nil {
		return nil
	}
	lst.len--
	e.next = e.next.next
	nextOne.lst = nil
	return nextOne
}

// RemoveAfter -------指定元素后删---------
func (lst *SingleList) RemoveAfter(e *ListNode) *ListNode {
	if e == nil {
		return nil
	}
	if e.lst == lst {
		return lst.remove(e)
	}
	return nil
}

// RemoveFront --------------前删---------------
func (lst *SingleList) RemoveFront() *ListNode {
	return lst.remove(lst.head)
}

// ShowList ------------显示单链表----------
// head:头结点
func (lst *SingleList) ShowList() {
	for cur := lst.Front(); cur != nil; cur = cur.Next() {
		fmt.Print(cur.Val, " ")
		cur = cur.Next()
	}
	fmt.Println()
}

// Len ------------单链表长度-----------
func (lst *SingleList) Len() int {
	return lst.len
}

// ListNode ---------单链表节点-----------
type ListNode struct {
	Val  interface{}
	next *ListNode
	lst  *SingleList
}

// Next ----------获得下一个元素----------
// 用于遍历
func (node *ListNode) Next() *ListNode {
	return node.next
}
