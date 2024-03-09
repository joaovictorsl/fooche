package evict

import (
	"fmt"
	"strings"
)

type (
	// LRU is a simple implementation of a Least Recently Used cache eviction policy.
	//
	// Not concurrent-safe.
	LRU[T comparable] struct {
		head       *lruNode[T]
		tail       *lruNode[T]
		keyNodeMap map[T]*lruNode[T]
		size       int
		cap        int
	}

	// Node used in LRU
	lruNode[T comparable] struct {
		key  T
		next *lruNode[T]
		prev *lruNode[T]
	}
)

// Indicates wether the LRU is empty or not.
//
// # Returns
//
// True if the node's value is equal to the zero value of type T, false otherwise.
func (n *lruNode[T]) isEmpty() bool {
	var empty T
	return n.key == empty
}

func NewLRU[T comparable](cap int) *LRU[T] {
	firstNode := &lruNode[T]{}

	return &LRU[T]{
		head:       firstNode,
		tail:       firstNode,
		keyNodeMap: make(map[T]*lruNode[T]),
		size:       0,
		cap:        cap,
	}
}

func (l *LRU[T]) RecordAccess(key T) (evictedKey T, eviction bool) {
	if node, ok := l.keyNodeMap[key]; ok {
		l.moveToFront(node)
		return evictedKey, false
	}

	eviction = false
	if l.size == l.cap {
		evictedKey = l.tail.key
		l.remove(l.tail)
		delete(l.keyNodeMap, evictedKey)
		eviction = true
	}

	l.addToFront(key)
	l.keyNodeMap[key] = l.head
	if !eviction {
		l.size++
	}
	return evictedKey, eviction
}

func (l *LRU[T]) Remove(key T) {
	if node, ok := l.keyNodeMap[key]; ok {
		l.remove(node)
		delete(l.keyNodeMap, key)
		l.size--
	}
}

// Adds a key to MRU position.
//
// This method does not update the size variable.
//
// # Params
//   - k: the key to be added
func (l *LRU[T]) addToFront(k T) {
	if l.head.isEmpty() {
		l.head.key = k
	} else {
		n := &lruNode[T]{key: k}
		l.head.prev = n
		n.next = l.head
		n.prev = &lruNode[T]{}
		l.head = n
	}
}

// Moves a node to the MRU position if it is not already in it.
//
// # Params
//   - node: the node to be moved
func (l *LRU[T]) moveToFront(node *lruNode[T]) {
	if l.size <= 1 || node == l.head {
		return
	}

	l.remove(node)

	l.head.prev = node
	node.next = l.head
	node.prev = &lruNode[T]{}
	l.head = node
}

// Removes a node.
//
// This method does not update the size variable.
//
// # Params
//   - node: the node to be removed
func (l *LRU[T]) remove(node *lruNode[T]) {
	var empty T
	if l.size == 0 {
		return
	}

	if l.size == 1 {
		l.head.key = empty
		return
	}

	if node == l.head {
		l.head = node.next
		l.head.prev = nil
		return
	}

	if node == l.tail {
		l.tail = node.prev
		l.tail.next = nil
		return
	}

	node.prev.next = node.next
	node.next.prev = node.prev
}

// Generates a string representation of LRU
func (l *LRU[T]) String() string {
	sb := strings.Builder{}
	for node := l.head; node != nil && !node.isEmpty(); node = node.next {
		sb.WriteString(fmt.Sprintf("%v -> ", node.key))
	}

	if sb.Len() == 0 {
		return "<EMPTY LRU>"
	}

	return sb.String()
}
