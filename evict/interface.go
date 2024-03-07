package evict

// EvictionPolicy handles which key should be evicted, if any, when recording
// new accesses to keys.
type EvictionPolicy[T comparable] interface {
	// Records an access to a key.
	//
	// Returns the key that should be evicted, if any, and a boolean indicating if an eviction occurred.
	RecordAccess(key T) (evictedKey T, eviction bool)
	// Removes a key
	Remove(key T)
}
