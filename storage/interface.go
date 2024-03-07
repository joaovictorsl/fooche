package storage

type Storage interface {
	// Add a key and its value to the storage
	Put(k string, v []byte) (err error)
	// Get a value associated with the given key.
	//
	// Returns a boolean indicating if the key was found or not.
	Get(k string) (v []byte, ok bool)
	// Removes a key from the storage
	Remove(k string)
	// Returns the amount of entries a storage can hold. This may change
	// over time depending on the storage implementation.
	Capacity() int
}
