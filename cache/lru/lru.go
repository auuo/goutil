package lru

import (
	"container/list"
	"sync"
)

type LRU interface {
	Add(key string, value interface{})

	Get(key string) (interface{}, bool)
}

type item struct {
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

	// todo 太丑陋
	ele, ok := c.dataMap[key]
	it := item{
		key:   key,
		value: value,
	}
	if ok {
		*ele.Value.(*item) = it
		c.Get(key)
	} else {
		ele = c.list.PushBack(&it)
		c.dataMap[key] = ele
	}

	if len(c.dataMap) > c.cap {
		ele := c.list.Front()
		c.list.Remove(ele)
		delete(c.dataMap, ele.Value.(*item).key)
	}
}

func (c *cache) Get(key string) (interface{}, bool) {
	c.lock.Lock()
	c.lock.Unlock()

	ele, ok := c.dataMap[key]
	if ok {
		c.list.MoveToBack(ele)
		return ele.Value.(*item).value, true
	}
	return nil, false
}
