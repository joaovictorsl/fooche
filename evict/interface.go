package evict

// EvictionPolicy handles which key should be evicted, if any, when recording
// new accesses to keys.
type EvictionPolicy[T comparable] interface {
	// Records an access to a key.
	//
	// # Params
	//   - key: the key that was accessed
	//
	// # Returns
	//
	// The key that should be evicted, if any, and a boolean indicating if an eviction occurred.
	RecordAccess(key T) (evictedKey T, eviction bool)
	// Removes a key
	//
	// # Params
	//   - key: the key to be removed
	Remove(key T)
}
