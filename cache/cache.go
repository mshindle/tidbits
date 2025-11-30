package cache

import (
	"math"
	"sync"
	"time"
)

// The maximum value for nanoseconds.
const maxNanos = int64(math.MaxInt64)

// Create the time.Time object representing the maximum possible time & convert to Unix()
var maxTime = time.Unix(0, maxNanos).UnixNano()

type Item[V any] struct {
	Value  V
	Expiry int64
}

type Cache[K comparable, V any] struct {
	mu    sync.RWMutex
	items map[K]*Item[V]
}

func New[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{items: make(map[K]*Item[V])}
}

// Set allows items to be stored with an expiry set to max Unix time.
func (c *Cache[K, V]) Set(k K, v V) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.items[k] = &Item[V]{Value: v, Expiry: maxTime}
}

// SetTTL allows items to be stored with an expiry time.
func (c *Cache[K, V]) SetTTL(k K, v V, ttl time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()

	expiry := time.Now().Add(ttl).UnixNano()
	c.items[k] = &Item[V]{Value: v, Expiry: expiry}
}

func (c *Cache[K, V]) Get(k K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	item, ok := c.items[k]
	if !ok || time.Now().UnixNano() > item.Expiry {
		return *new(V), false
	}

	return item.Value, true
}

// Cleanup should be run in a go function which triggers cleanup based on some constraint: time interval, memory usage
// too high, etc.
func (c *Cache[K, V]) Cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now().UnixNano()
	for k, item := range c.items {
		if now > item.Expiry {
			delete(c.items, k)
		}
	}
}
