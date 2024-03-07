package evict

import (
	"testing"
)

func TestRecordAccess(t *testing.T) {
	lru := NewLRU[string](2)
	a, b, c, d, e := "a", "b", "c", "d", "e"

	if evictedKey, eviction := lru.RecordAccess(a); eviction {
		t.Errorf("Expected no eviction, got %s", evictedKey)
	} else if lru.size != 1 {
		t.Errorf("Expected size to be 1, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(b); eviction {
		t.Errorf("Expected no eviction, got %s", evictedKey)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(c); !eviction || evictedKey != "a" {
		t.Errorf("Expected (evictedKey, eviction) to be (a, true), got (%s, %v)", evictedKey, eviction)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(d); !eviction || evictedKey != "b" {
		t.Errorf("Expected (evictedKey, eviction) to be (b, true), got (%s, %v)", evictedKey, eviction)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(e); !eviction || evictedKey != "c" {
		t.Errorf("Expected (evictedKey, eviction) to be (c, true), got (%s, %v)", evictedKey, eviction)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(d); eviction || evictedKey != "" {
		t.Errorf("Expected (evictedKey, eviction) to be (\"\", false), got (%s, %v)", evictedKey, eviction)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}

	if evictedKey, eviction := lru.RecordAccess(e); eviction || evictedKey != "" {
		t.Errorf("Expected (evictedKey, eviction) to be (\"\", false), got (%s, %v)", evictedKey, eviction)
	} else if lru.size != 2 {
		t.Errorf("Expected size to be 2, got %d", lru.size)
	}
}

func TestRemove(t *testing.T) {
	lru := NewLRU[string](2)
	a, b := "a", "b"

	lru.RecordAccess(a)
	lru.RecordAccess(b)

	lru.Remove(b)
	if lru.size != 1 {
		t.Errorf("Expected size to be 1, got %d", lru.size)
	} else if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	}

	lru.Remove("non existent key")
	if lru.size != 1 {
		t.Errorf("Expected size to still be 1 when trying to remove a non existent key, got %d", lru.size)
	} else if lru.head.key != "a" {
		t.Errorf("Expected head to still be a when trying to remove a non existent key, got %s", lru.head.key)
	}

	lru.Remove(a)
	if lru.size != 0 {
		t.Errorf("Expected size to be 0, got %d", lru.size)
	} else if lru.head.key != "" {
		t.Errorf("Expected head to be empty, got %s", lru.head.key)
	}
}

func TestAddToFront(t *testing.T) {
	lru := NewLRU[string](2)
	a, b := "a", "b"

	lru.addToFront(a)
	if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	}

	lru.addToFront(b)
	if lru.head.key != "b" {
		t.Errorf("Expected head to be b, got %s", lru.head.key)
	} else if lru.head.next.key != "a" {
		t.Errorf("Expected head.next to be a, got %s", lru.head.next.key)
	}
}

func TestMoveToFront(t *testing.T) {
	lru := NewLRU[string](2)
	a, b, c := "a", "b", "c"

	lru.RecordAccess(a)
	lru.moveToFront(lru.head)
	if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	} else if lru.head.next != nil {
		t.Errorf("Expected head.next to be nil, got %s", lru.head.next.key)
	}

	lru.RecordAccess(b)
	lru.moveToFront(lru.head.next)
	if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	} else if lru.head.next.key != "b" {
		t.Errorf("Expected head.next to be b, got %s", lru.head.next.key)
	}

	lru.RecordAccess(c)
	lru.moveToFront(lru.head.next)
	if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	} else if lru.head.next.key != "c" {
		t.Errorf("Expected head.next to be c, got %s", lru.head.next.key)
	}

	lru.moveToFront(lru.head)
	if lru.head.key != "a" {
		t.Errorf("Expected head to be a, got %s", lru.head.key)
	} else if lru.head.next.key != "c" {
		t.Errorf("Expected head.next to be c, got %s", lru.head.next.key)
	}
}
