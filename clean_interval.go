package fooche

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/fooche/evict"
	"github.com/joaovictorsl/fooche/storage"
)

// A cache that cleans expired keys every cleanInterval. Keys are evicted based on its policy when bounded.
//
// The cache is thread-safe and can be used concurrently.
type CleanIntervalCache struct {
	lock    sync.RWMutex
	storage storage.Storage
	// Maps a key to its expiration time
	keyExpMap map[string]time.Time
	// Interval in which expired keys will be evicted
	cleanInterval time.Duration
	policy        evict.EvictionPolicy[string]
}

// Creates a bounded CleanIntervalCache, see [storage.BoundedStorage] for more info on storage implementation.
//
// # Params
//   - cleanInterval: the interval in which the cache will be cleaned.
//   - sizeAndCapMap: a map of size in bytes and capacity (how many objects with this size or less can be stored).
//   - createPolicy: a function that creates a new eviction policy based on the storage's capacity.
//
// # Returns
//
// A new CleanIntervalCache.
func NewCleanIntervalBounded(cleanInterval time.Duration, sizeAndCapMap map[int]int, createPolicy func(capacity int) evict.EvictionPolicy[string]) *CleanIntervalCache {
	s := storage.NewBoundedStorage(sizeAndCapMap)
	c := &CleanIntervalCache{
		storage:       s,
		cleanInterval: cleanInterval,
		keyExpMap:     make(map[string]time.Time),
		policy:        createPolicy(s.Capacity()),
	}

	go c.startCleaner()

	return c
}

// Creates a unbounded CleanIntervalCache, see [storage.UnboundedStorage] for more info on storage implementation.
//
// No policy is used ince storage is unbounded.
//
// # Params
//   - cleanInterval: the interval in which the cache will be cleaned.
//
// # Returns
//
// A new CleanIntervalCache.
func NewCleanInterval(cleanInterval time.Duration) *CleanIntervalCache {
	c := &CleanIntervalCache{
		storage:       storage.NewUnboundedStorage(),
		cleanInterval: cleanInterval,
		keyExpMap:     make(map[string]time.Time),
		policy:        &evict.NoPolicy[string]{},
	}

	go c.startCleaner()

	return c
}

// Cleans expired objects from cache every cleanInterval.
func (c *CleanIntervalCache) startCleaner() {
	for {
		<-time.After(c.cleanInterval)

		c.lock.Lock()

		for k, exp := range c.keyExpMap {
			if time.Now().After(exp) {
				delete(c.keyExpMap, k)
				c.storage.Remove(k)
				c.policy.Remove(k)
			}
		}

		c.lock.Unlock()
	}
}

func (c *CleanIntervalCache) String() string {
	return fmt.Sprintf("%v", c.storage)
}

func (c *CleanIntervalCache) Set(k string, v []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	evictedKey, eviction := c.policy.RecordAccess(k)
	if eviction {
		c.storage.Remove(evictedKey)
	}

	exp := time.Now().Add(ttl)
	c.keyExpMap[k] = exp
	if err := c.storage.Put(k, v); err != nil {
		return err
	}

	return nil
}

func (c *CleanIntervalCache) Has(k string) bool {
	c.lock.RLock()
	defer c.lock.RUnlock()

	_, ok := c.storage.Get(k)
	notExpired := !time.Now().After(c.keyExpMap[k])

	return ok && notExpired
}

func (c *CleanIntervalCache) Get(k string) ([]byte, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()

	data, ok := c.storage.Get(k)
	expired := time.Now().After(c.keyExpMap[k])

	if !ok || expired {
		return nil, fmt.Errorf("key (%s) not found", k)
	}

	evictedKey, eviction := c.policy.RecordAccess(k)
	if eviction {
		c.storage.Remove(evictedKey)
	}

	return data, nil
}

func (c *CleanIntervalCache) ComputeIfAbsent(k string, computeValue func() []byte) (v []byte, err error) {
	v, err = c.Get(k)
	if err != nil {
		v = computeValue()
	}

	if err := c.Set(k, v, time.Minute); err != nil {
		return nil, err
	}

	return v, nil
}

func (c *CleanIntervalCache) Delete(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.keyExpMap, k)
	c.storage.Remove(k)
	c.policy.Remove(k)
}
