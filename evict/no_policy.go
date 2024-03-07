package evict

type NoPolicy[T comparable] struct{}

func (n *NoPolicy[T]) RecordAccess(key T) (evictedKey T, eviction bool) {
	return evictedKey, false
}

func (n *NoPolicy[T]) Remove(key T) {
}
