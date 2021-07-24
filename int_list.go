package main

import "sync/atomic"

type SyncIntList struct {
	head   *intNode
	length int64
}

func NewInt() IntList {
	return &SyncIntList{
		head:   newNode(0),
		length: 0,
	}
}

func (s *SyncIntList) Contains(value int) bool {
	curr := s.head.getNextNode()
	for curr != nil && curr.getValue() < value {
		curr = curr.getNextNode()
	}
	if curr == nil || curr.isDeleted() {
		return false
	}
	return curr.getValue() == value
}

func (s *SyncIntList) Insert(value int) bool {
	var pre, curr *intNode
	for {
		pre = s.head
		curr = pre.getNextNode()
		for curr != nil && curr.getValue() < value {
			pre = curr
			curr = curr.getNextNode()
		}
		if curr != nil && curr.getValue() == value {
			return false
		}
		pre.mu.Lock()
		if pre.getNextNode() != curr || (pre != nil && pre.isDeleted()) {
			pre.mu.Unlock()
			continue
		} else {
			break
		}
	}
	nextNode := newNode(value)
	nextNode.setNextNode(curr)
	pre.setNextNode(nextNode)
	atomic.AddInt64(&s.length, 1)
	pre.mu.Unlock()
	return true
}

func (s *SyncIntList) Delete(value int) bool {
	var pre, curr *intNode
	for {
		pre = s.head
		curr = pre.getNextNode()
		for curr != nil && curr.getValue() < value {
			pre = curr
			curr = curr.getNextNode()
		}
		if curr == nil || curr.getValue() != value {
			return false
		}
		curr.mu.Lock()
		if curr.isDeleted() {
			curr.mu.Unlock()
			continue
		}
		pre.mu.Lock()
		if pre.getNextNode() != curr || pre.isDeleted() {
			curr.mu.Unlock()
			pre.mu.Unlock()
			continue
		} else {
			break
		}
	}
	curr.setDeleted()
	pre.setNextNode(curr.getNextNode())
	atomic.AddInt64(&s.length, -1)
	pre.mu.Unlock()
	curr.mu.Unlock()
	return true
}

func (s *SyncIntList) Range(f func(value int) bool) {
	curr := s.head.getNextNode()
	for curr != nil {
		if !f(curr.getValue()) {
			break
		}
		curr = curr.getNextNode()
	}
}

func (s *SyncIntList) Len() int {
	return int(atomic.LoadInt64(&s.length))
}

type IntList interface {
	// Contains 检查一个元素是否存在，如果存在则返回 true，否则返回 false
	Contains(value int) bool

	// Insert 插入一个元素，如果此操作成功插入一个元素，则返回 true，否则返回 false
	Insert(value int) bool

	// Delete 删除一个元素，如果此操作成功删除一个元素，则返回 true，否则返回 false
	Delete(value int) bool

	// Range 遍历此有序链表的所有元素，如果 f 返回 false，则停止遍历
	Range(f func(value int) bool)

	// Len 返回有序链表的元素个数
	Len() int
}
