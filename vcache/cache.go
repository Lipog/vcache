package vcache
//这里是并发控制，通过互斥锁来实现
//这是中间层的缓存，通过锁来保证底层的缓存的增删改查的并发安全

import (
	"Vcache/lru"
	"sync"
)


type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

//cache主要实例化一个lru，封装了get和add方法，并添加了互斥锁MUTEX
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//如果lru缓存不存在，那么就新建一个，其大小为cache的cacheBytes大小
	//不为nil才创建，称为延迟初始化
	//主要用于提高性能，减少程序内存需要
	if c.lru == nil {
		c.lru = lru.New(c.cacheBytes, nil)
	}
	c.lru.Add(key, value)
}

func (c *cache) get(key string) (ByteView, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//如果底层缓存为空，那么就直接返回nil
	if c.lru == nil {
		return ByteView{}, false
	}
	//如果存在，就返回对应的值
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	//不存在则返回空和false
	return ByteView{}, false
}