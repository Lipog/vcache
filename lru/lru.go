package lru
//这里存放的是缓存淘汰策略
//最底层的缓存，进行增删改查等操作

import (
	"container/list"
)

type Cache struct {
	maxBytes int64  //允许最大使用的内存
	nbytes int64  //当前使用的内存
	ll *list.List
	cache map[string]*list.Element
	//该函数可选，当条目被清除时执行
	OnEvicted func(key string, value Value)
}

//因为从链表删除节点需要key来删除对应的映射
type entry struct {
	key string
	value Value
}

type Value interface {
	Len() int
}

//Cache的构造函数，用来实例化一个Cache对象
func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	cache := &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
	return cache
}

//如果查询到对应的节点，那么就将其放到队首
func (c *Cache) Get(key string) (Value, bool) {
	//如果该缓存存在于cache中，ok就为true
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return nil, false
}

func (c *Cache) Add(key string, value Value) {
	//如果节点已经存在了，那么就更新其值
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		//用kv把旧的元素的相关信息储存起来
		kv := ele.Value.(*entry)
		c.nbytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//键值对不存在，那么就将创建并添加到队首
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.nbytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

//这里是删除最近最少访问的节点，即队尾的元素，实际就是缓存淘汰策略
func (c *Cache) RemoveOldest() {
	//找到链表的最后一个元素并返回
	ele := c.ll.Back()
	//如果最后一个元素存在
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.nbytes -= int64(len(kv.key)) + int64(kv.value.Len())
		//如果回调函数不为nil，那么就调用回调函数
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

//统计缓存中条目的数量
func (c *Cache) Len() int {
	return c.ll.Len()
}

