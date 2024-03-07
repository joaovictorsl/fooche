package fooche

import (
	"fmt"
	"sync"
	"time"

	"github.com/joaovictorsl/fooche/evict"
	"github.com/joaovictorsl/fooche/storage"
)

type CleanIntervalCache struct {
	lock          sync.RWMutex
	storage       storage.Storage
	keyExpMap     map[string]time.Time
	cleanInterval time.Duration
	policy        evict.EvictionPolicy[string]
}

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

func (c *CleanIntervalCache) Delete(k string) {
	c.lock.Lock()
	defer c.lock.Unlock()

	delete(c.keyExpMap, k)
	c.storage.Remove(k)
	c.policy.Remove(k)
}
