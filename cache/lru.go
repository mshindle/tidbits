package cache

type inode struct {
	key   string
	value int
	prev  *inode
	next  *inode
}

func (i *inode) removeBindings() {
	if i.prev != nil {
		i.prev.next = i.next
	}
	if i.next != nil {
		i.next.prev = i.prev
	}
	i.prev, i.next = nil, nil
}

type cacheIndex struct {
	head *inode
	tail *inode
}

func (ci *cacheIndex) setHeadTo(node *inode) {
	if ci.head == node {
		return
	}
	if ci.head == nil {
		ci.head, ci.tail = node, node
		return
	}
	if ci.head == ci.tail {
		ci.tail.prev = node
		ci.head = node
		ci.head.next = ci.tail
		return
	}
	if ci.tail == node {
		ci.removeTail()
	}
	ci.head.prev = node
	node.next = ci.head
	ci.head = node
}

func (ci *cacheIndex) removeTail() {
	if ci.tail == nil {
		return
	}
	if ci.tail == ci.head {
		ci.head, ci.tail = nil, nil
		return
	}
	ci.tail = ci.tail.prev
	ci.tail.next.removeBindings()
}

type LRUCache struct {
	maxSize int
	size    int
	data    map[string]*inode
	recent  *cacheIndex
}

func NewLRUCache(size int) *LRUCache {
	data := make(map[string]*inode, size)
	recent := &cacheIndex{}
	return &LRUCache{maxSize: size, data: data, recent: recent}
}

func (cache *LRUCache) InsertKeyValuePair(key string, value int) {
	defer cache.updateMostRecent(cache.data[key])

	if _, found := cache.data[key]; found {
		cache.replaceKey(key, value)
		return
	}
	if cache.size == cache.maxSize {
		cache.evictLeastRecent()
	} else {
		cache.size += 1
	}
	cache.data[key] = &inode{key: key, value: value}
}

// The second return value indicates whether or not the key was found
// in the cache.
func (cache *LRUCache) GetValueFromKey(key string) (int, bool) {
	node, found := cache.data[key]
	if !found {
		return 0, false
	}
	cache.updateMostRecent(node)
	return node.value, true
}

// The second return value is false if the cache is empty.
func (cache *LRUCache) GetMostRecentKey() (string, bool) {
	if cache.recent.head == nil {
		return "", false
	}
	return cache.recent.head.key, true
}

func (cache *LRUCache) replaceKey(key string, value int) {
	node, found := cache.data[key]
	if !found {
		// better error handling here as we should have already confirmed this. just fail silently
		return
	}
	node.value = value
}

func (cache *LRUCache) updateMostRecent(node *inode) {
	cache.recent.setHeadTo(node)
}

func (cache *LRUCache) evictLeastRecent() {
	key := cache.recent.tail.key
	cache.recent.removeTail()
	delete(cache.data, key)
}
