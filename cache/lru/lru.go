package lru

import (
	"container/list"
	"sync"
)

type LRU interface {
	Add(key string, value interface{})

	Get(key string) (interface{}, bool)
}

type entry struct {
	key   string
	value interface{}
}

type cache struct {
	cap int

	list list.List

	dataMap map[string]*list.Element

	lock sync.Mutex
}

func New(cap int) LRU {
	return &cache{cap: cap, dataMap: make(map[string]*list.Element)}
}

func (c *cache) Add(key string, value interface{}) {
	c.lock.Lock()
	c.lock.Unlock()

	if ele, ok := c.dataMap[key]; ok {
		ele.Value.(*entry).value = value
		c.list.MoveToBack(ele)
	} else {
		c.dataMap[key] = c.list.PushBack(&entry{
			key:   key,
			value: value,
		})
	}
	if len(c.dataMap) > c.cap {
		ele := c.list.Front()
		c.list.Remove(ele)
		delete(c.dataMap, ele.Value.(*entry).key)
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	c.lock.Unlock()

	ele, ok := c.dataMap[key]
	if ok {
		c.list.MoveToBack(ele)
		return ele.Value.(*entry).value, true
	}
	return nil, false
}
