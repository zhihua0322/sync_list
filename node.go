package main

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

const isDeleted = 1

type intNode struct {
	value  int
	marked uint32
	next   *intNode
	mu     sync.Mutex
}

func newNode(value int) *intNode {
	return &intNode{
		value:  value,
		marked: 0,
		next:   nil,
		mu:     sync.Mutex{},
	}
}

// util
func (n *intNode) getValue() int {
	return int(atomic.LoadInt64((*int64)(unsafe.Pointer(&n.value))))
}

func (n *intNode) setValue(value int) {
	atomic.StoreInt64((*int64)(unsafe.Pointer(&n.value)), int64(value))
}

func (n *intNode) getNextNode() *intNode {
	return (*intNode)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&n.next))))
}

func (n *intNode) setNextNode(next *intNode) {
	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&n.next)), unsafe.Pointer(next))
}

func (n *intNode) isDeleted() bool {
	return int(atomic.LoadUint32(&n.marked)) == isDeleted
}

func (n *intNode) setDeleted() {
	atomic.StoreUint32(&n.marked, uint32(isDeleted))
}
