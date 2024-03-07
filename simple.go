package fooche

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/fooche/evict"
	"github.com/joaovictorsl/fooche/storage"
)

// A cache which does no eviction and has no expiration.
type SimpleCache struct {
	lock    *sync.RWMutex
	storage storage.Storage
	policy  evict.EvictionPolicy[string]
}

/*
Creates a bounded SimpleCache, see BoundedStorage in [dcache.core.cache.storage] for more info on allocated bytes.
*/
func NewSimpleBounded(sizeAndCapMap map[int]int, createPolicy func(capacity int) evict.EvictionPolicy[string]) *SimpleCache {
	s := storage.NewBoundedStorage(sizeAndCapMap)
	return &SimpleCache{
		lock:    &sync.RWMutex{},
		storage: s,
		policy:  createPolicy(s.Capacity()),
	}
}

// Creates a unbounded SimpleCache.
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

func (c *SimpleCache) Set(k string, v []byte, ttl time.Duration) error {
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
	c.lock.RLock()
	defer c.lock.RUnlock()

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

func (c *SimpleCache) Delete(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.storage.Remove(k)
	c.policy.Remove(k)
}
