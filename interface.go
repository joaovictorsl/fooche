package fooche

import "time"

type ICache interface {
	// Adds an element to cache.
	//
	// # Params
	// 	- k: object key.
	// 	- v: object value.
	//
	// # Returns
	//
	// An error if something goes wrong.
	Set(k string, v []byte, ttl time.Duration) error

	// Checks if an element is in cache.
	//
	// # Params
	// 	- k: object key.
	//
	// # Returns
	//
	// Boolean indicating wether or not the object is present.
	Has(k string) bool

	// Gets an element from cache.
	//
	// # Params
	// 	- k: object key.
	//
	// # Returns
	//
	// The object value and an error if anything goes wrong.
	Get(k string) (v []byte, err error)

	// Gets an element from cache.
	//
	// # Params
	// 	- k: object key.
	//  - computeValue: function to compute the value if it's not present in cache.
	//
	// # Returns
	//
	// The object value and an error if anything goes wrong.
	ComputeIfAbsent(k string, computeValue func() []byte) (v []byte, err error)

	// Removes an element from cache.
	//
	// # Params
	// 	- k: object key.
	Delete(k string)

	// Returns a string representation of the cache.
	String() string
}
