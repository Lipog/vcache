package vcache
//cache实现并发控制，保证并发读写的安全
import (
	"sync"
	"vcache/lru"
)

//该包封装了lru，可以支持并发读写lru,同时封装了add和get等方法
//并添加了互斥锁来保证并发安全

type cache struct {
	mu sync.Mutex
	lru *lru.Cache
	cacheBytes int64
}

//add添加的值的类型就是Byteview类型，实现了Value接口的Len的方法，
//所以在底层存入的就是结构体
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//进行添加时，如果lru还未创建，说明是第一次访问，那么就创建一个lru
	if c.lru == nil {
		c.lru = lru.NewCahce(c.cacheBytes, nil)
	}
	//在确保lru存在的情况下，调用lru的Add方法添加kv
	c.lru.Add(key, value)
}

//
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}

	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}


