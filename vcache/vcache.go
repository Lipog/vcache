package vcache
//这里是顶层的缓存，统筹回调等

import (
	"fmt"
	"log"
	"sync"
)

//负责与外部的交互，控制缓存存储和获取的主流程
//逻辑是现在本地寻找，找不到去远程节点获取，还获取不到的话，就回调，将值添加到缓存


//一个 Group 可以认为是一个缓存的命名空间，每个 Group 拥有一个唯一的名称 name。
// 比如可以创建三个 Group，缓存学生的成绩命名为 scores，
// 缓存学生信息的命名为 info，缓存学生课程的命名为 courses
type Group struct {
	name      string
	getter    Getter //即缓存未命中时获取源数据的回调(callback)
	mainCache cache  //即一开始实现的并发缓存,cache中含有锁，可以保证并发的安全
}

//该接口用来根据key加载数据
type Getter interface {
	Get(key string) ([]byte, error)
}

//GetterFunc用函数实现Getter
type GetterFunc func(key string) ([]byte, error)

//Get实现Getter接口函数
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu sync.RWMutex
	groups = make(map[string]*Group)
)

//NewGroup用来创建一个Group实例,
//同时将创建的group存储在全局变量groups中
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
	}
	//将name和建立的group建立映射关系
	groups[name] = g
	return g
}

//GetGroup返回先前用NewGroup创建的已命名组，如果没有这样的组，则返回nil
func GetGroup(name string) *Group {
	//只读锁 RLock()，因为不涉及任何冲突变量的写操作
	mu.RLock()
	defer mu.RUnlock()
	g := groups[name]
	return g
}

//从cache中获得key对应的值
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("KEY is required")
	}

	//如果能够在缓存中命中key，那么就将其返回
	//读取完数据后，锁就释放掉了，其他的函数可以对其进行操作
	//如果在缓存里找到了，就直接返回了
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[Vcache] hit")
		return v, nil
	}

	//如果在当前的缓存里没找到，就要去全局找了
	//这里进行了深拷贝的操作，以保证返回的数据的安全性
	return g.load(key)
}

func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

//通过回调函数，去对应的地方寻找数据
//找到数据后，再持久化到本地缓存
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	//如果有错，就返回空ByteView，并且返回错误
	if err != nil {
		return ByteView{}, err
	}

	//获取到实例的拷贝，以免其他操作的修改，影响读取的准确性
	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)
	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
