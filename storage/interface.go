package storage

type Storage interface {
	// Adds a key and its value to the storage.
	//
	// # Params
	//  - k: the key
	//  - v: the value
	//
	// # Returns
	//
	// An error if the key could not be added.
	Put(k string, v []byte) (err error)
	// Get a value associated with the given key.
	//
	// # Params
	//  - k: the key
	//
	// # Returns
	//
	// The key value, if it was found, and a boolean indicating if the key was found or not.
	Get(k string) (v []byte, ok bool)
	// Removes a key from the storage.
	//
	// # Params
	//  - k: the key
	Remove(k string)
	// Returns the amount of entries a storage can hold. If the storage is unbounded, it returns -1.
	Capacity() int
}
