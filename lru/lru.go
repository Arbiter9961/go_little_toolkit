package lru

import (
	"container/list"
	"sync"
	"time"
)

type item struct {
	Key   string
	Value interface{}
	Time  time.Time
}

type Cache struct {
	sync.Mutex
	size     int
	overtime int64
	List     *list.List
	Keys     map[string]*list.Element
}

func New(Size int, Overtime int64) (Cache, error) {
	return Cache{
		size:     Size,
		overtime: Overtime,
		List:     list.New(),
		Keys:     make(map[string]*list.Element),
	}, nil
}

func (c *Cache) Add(key string, value interface{}) bool {
	c.Lock()
	defer c.Unlock()
	if it, ok := c.Keys[key]; ok {
		c.List.MoveToFront(it)
		return false
	}

	it := &item{
		Key:   key,
		Value: value,
		Time:  time.Now(),
	}
	ite := c.List.PushFront(it)
	c.Keys[key] = ite

	if c.List.Len() > c.size {
		nouse := c.List.Back()
		c.List.Remove(nouse)
		delete(c.Keys, nouse.Value.(*item).Key)
	}
	return true
}

func (c *Cache) Get(key string) interface{} {
	c.Lock()
	defer c.Unlock()
	if it, ok := c.Keys[key]; ok {
		if time.Now().Unix()-it.Value.(*item).Time.Unix() > c.overtime {
			return nil
		}
		return it
	}
	return nil
}

func (c *Cache) Remove(key string) bool {
	c.Lock()
	defer c.Unlock()
	if it, ok := c.Keys[key]; ok {
		c.List.Remove(it)
		delete(c.Keys, key)
		return true
	}
	return false
}
