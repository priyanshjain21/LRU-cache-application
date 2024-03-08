package main

import (
	"container/list"
	"sync"
	"time"
)

type entry struct {
	key        string
	value      interface{}
	expiration *time.Time
}

type LRUCache struct {
	capacity int
	cache    map[string]*list.Element
	ll       *list.List
	lock     sync.Mutex
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		cache:    make(map[string]*list.Element),
		ll:       list.New(),
	}
}

func (c *LRUCache) Get(key string) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if ele, hit := c.cache[key]; hit {
		c.ll.MoveToFront(ele)
		if ele.Value.(*entry).expiration != nil && ele.Value.(*entry).expiration.Before(time.Now()) {
			c.ll.Remove(ele)
			delete(c.cache, key)
			return nil, false
		}
		return ele.Value.(*entry).value, true
	}
	return
}

func (c *LRUCache) Set(key string, value interface{}, duration time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	var expiration *time.Time
	if duration > 0 {
		exp := time.Now().Add(duration)
		expiration = &exp
	}

	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		ele.Value.(*entry).value = value
		ele.Value.(*entry).expiration = expiration
		return
	}
	ele := c.ll.PushFront(&entry{key, value, expiration})
	c.cache[key] = ele
	if c.ll.Len() > c.capacity {
		oldest := c.ll.Back()
		if oldest != nil {
			c.ll.Remove(oldest)
			delete(c.cache, oldest.Value.(*entry).key)
		}
	}
}
