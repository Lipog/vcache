package vcache

import (
	"fmt"
	"log"
	"sync"
	"vcache/singleflight"
	"vcache/vcachepb"
)

//负责与外部进行交互，控制缓存存储以及获取的主流程

//首先定义回调，如果缓存不存在，应该从数据源进行数据的获取
//回调后将获得的数据缓存到当前的缓存中

var (
	mu sync.RWMutex
	//定义一个全局的groups，用来存储该节点的所有的Group
	groups = make(map[string]*Group)
)

//Group是缓存的核心数据结构，负责与外部进行交互，是数据的入口和出口
type Group struct {
	name string
	getter GetterFunc
	cache cache
	http_self NodePicker //用来存贮自身的HttpPool，从而来获取对应的远程地址来获取数据
	loader *singleflight.Group
}

//NewGroup用来创建一个Group的实例，一个节点可以创建多个Group
//并且每个Group都放在全局的groups中保存，通过group的name来寻址该缓存
//一个Group可以被认为是一个缓存的命名空间，每个Gruop拥有一个唯一的名称name，
//比如可以创建三个Group，缓存学生成绩命名为scores,缓存学生信息为info,
//缓存学生课程为 courses
func NewGroup(name string, cacheBytes int64, getter GetterFunc) *Group {
	if getter == nil {
		panic("nil Getter")
	}
	//在创建Group之前要进行加锁，保证只有一个组被创建
	mu.Lock()
	defer mu.Unlock()
	g := &Group{
		name:   name,
		getter: getter,
		cache:  cache{cacheBytes:cacheBytes},
		loader: &singleflight.Group{},
	}
	groups[name] = g
	return g
}

func (g *Group) RegisterNode(node NodePicker) {
	if g.http_self != nil {
		panic("RegisterNodePicker called more than once")
	}
	g.http_self = node
}

//Get方法是Group的最核心的方法，因为一个Group就是一个缓存
//所以从该Group中获得数据非常重要
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	//如果key不为空，那么就从该Group的cache中（该cache封装了lru，是并发安全的）来获得数据
	if v, ok := g.cache.get(key); ok {
		log.Printf("[Group - %s] [LocalCache] [hit %s, value: %s] \n", g.name, key, v)
		return v, nil
	}

	log.Println("localcache not hit, searching DB, please wait!")
	//如果没有从本地的缓存中获得数据，就得去回调来获得数据了
	return g.load(key)
}

func (g *Group) load(key string) (value ByteView, err error) {
	//如果本地缓存没有找到的话，就要从远程节点进行寻找了
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.http_self != nil {
			if node, ok := g.http_self.PickNode(key); ok {
				if value, err = g.getFromNode(node, key); err == nil {
					//如果成功回调并获得了结果，那么就要添加到本地的缓存中去
					value_add := ByteView{b:cloneBytes(value.ByteSlice())}
					g.cache.add(key, value_add)
					return value, nil
				}
				//如果没有找到，那么就打印输出信息
				log.Println("[vcache] failed to get from node", err)
			}
		}
		return g.getLocally(key)
	})
	if err == nil {
		return viewi.(ByteView), nil
	}
	return
}

func (g *Group) getFromNode(node NodeGetter, key string) (ByteView, error) {
	req := &vcachepb.Request{
		Group: g.name,
		Key:   key,
	}
	res := &vcachepb.Response{}
	err := node.Get(req, res)
	if err != nil {
		return ByteView{}, err
	}
	return ByteView{b: res.Value}, nil
}
func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err
	}

	//如果成功回调并获得了结果，那么就要添加到本地的缓存中去
	value := ByteView{b:cloneBytes(bytes)}
	g.cache.add(key, value)
	//添加到本地缓存以后再返回
	return value, nil
}

//GetGroup用来从全局的groups中获得对应的name的Group
func GetGruop(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}
