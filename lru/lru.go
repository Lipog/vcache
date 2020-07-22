package lru
//实现了lru的缓存淘汰策略，是最底层的缓存
import (
	"container/list"
	"fmt"
)

//这是缓存的主结构，用来存放缓存的所有数据
type Cache struct {
	maxBytes int64
	usedBytes int64
	ll *list.List
	cache map[string]*list.Element //用来存放缓存的映射关系，key和Element的映射
	OnEvicted func(key string, value Value) //删除缓存时的回调操作
}

//在list的Element存放缓存时的结构体数据
type entry struct {
	key string
	value Value
}

//entry里的实际的value数据要实现Value接口，即实现Len方法
type Value interface {
	Len() int
}

func NewCahce(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

//然后是Cache应该实现的方法，一共有4种方法，查找，增加，修改，删除（被动删除和主动删除）
//首先是查找，即Get方法,传入key，返回value和是否成功的标志
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.PushFront(ele)
		//取出查询到的元素的Value，然后将Value转换为entry类型
		//因为Ele的Value可以接受任何的类型,但是取出时要进行转换
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	//如果没有找到就返回nil
	return
}

//Add实现了增加和修改的操作，进行操作前要判断对应的key是否存在于缓存种
func (c *Cache) Add(key string, value Value) {
	//首先判断key是否存在于缓存中，存在的话则修改元素，并修改使用的内存大小
	if ele, ok := c.cache[key]; ok {
		//因为进行了操作，首先把元素放置于队首
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.usedBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	} else {
		//创建元素然后放置于队首
		ele := c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.usedBytes += int64(len(key)) + int64(value.Len())
	}
	for c.maxBytes != 0 && c.maxBytes < c.usedBytes {
		c.Remove()
	}
}

//Remove实现了删除队尾元素的操作
func (c *Cache) Remove(key ...string) (bool, error){
	ele := &list.Element{}
	if len(key) == 0 {
		ele = c.ll.Back()//首先获取队尾的元素
	} else {
		ele, _ = c.cache[key[0]]
	}

	//如果队尾的元素存在的话，就对其进行删除，然后根据回调函数再进一步对其进行处理
	if ele != nil {
		//如果存在，则在cache和list中都进行删除，并更新使用内存的大小
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		//然后更新内存的大小
		c.usedBytes -= int64(len(kv.key)) + int64(kv.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
		return true, nil
	}
	return false, fmt.Errorf("key %s does not exists!", key)
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
//Remove用来删除指定的元素
//func (c *Cache) Remove(key ...string) (bool, error) {
//
//
//		c.ll.Remove(ele)
//		delete(c.cache, key)
//		kv := ele.Value.(*entry)
//		c.usedBytes -= int64(len(key)) + int64(kv.value.Len())
//		if c.OnEvicted != nil {
//			c.OnEvicted(kv.key, kv.value)
//		}
//	}
//	return false, fmt.Errorf("key %s does not exists!", key)
//}