package fooche

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/fooche/evict"
	"github.com/joaovictorsl/fooche/storage"
)

// A simple cache in which keys are evicted based on its policy when bounded.
//
// The cache is thread-safe and can be used concurrently.
type SimpleCache struct {
	lock    *sync.RWMutex
	storage storage.Storage
	policy  evict.EvictionPolicy[string]
}

// Creates a bounded SimpleCache, see [storage.BoundedStorage] for more info on storage implementation.
//
// # Params
//   - sizeAndCapMap: a map of size in bytes and capacity (how many objects with this size or less can be stored).
//   - createPolicy: a function that creates a new eviction policy based on the storage's capacity.
//
// # Returns
//
// A new SimpleCache.
func NewSimpleBounded(sizeAndCapMap map[int]int, createPolicy func(capacity int) evict.EvictionPolicy[string]) *SimpleCache {
	s := storage.NewBoundedStorage(sizeAndCapMap)
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: s,
		policy:  createPolicy(s.Capacity()),
	}
}

// Creates a unbounded SimpleCache, see [storage.UnboundedStorage] for more info on storage implementation.
//
// No policy is used since storage is unbounded.
//
// # Returns
//
// A new SimpleCache.
func NewSimple() *SimpleCache {
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: storage.NewUnboundedStorage(),
		policy:  &evict.NoPolicy[string]{},
	}
}

func (c *SimpleCache) String() string {
	return fmt.Sprintf("%v", c.storage)
}

func (c *SimpleCache) Set(k string, v []byte, _ time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	evictedKey, eviction := c.policy.RecordAccess(k)
	if eviction {
		c.storage.Remove(evictedKey)
	}

	if err := c.storage.Put(k, v); err != nil {
		return err
	}

	return nil
}

func (c *SimpleCache) Has(k string) bool {
	_, err := c.Get(k)
	return err == nil
}

func (c *SimpleCache) Get(k string) (v []byte, err error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	v, ok := c.storage.Get(k)
	if !ok {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	evictedKey, eviction := c.policy.RecordAccess(k)
	if eviction {
		c.storage.Remove(evictedKey)
	}

	return v, nil
}

func (c *SimpleCache) ComputeIfAbsent(k string, computeValue func() []byte) (v []byte, err error) {
	v, err = c.Get(k)
	if err != nil {
		v = computeValue()
	}

	if err := c.Set(k, v, 0); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *SimpleCache) Delete(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.storage.Remove(k)
	c.policy.Remove(k)
}
