package singleflight

import "sync"

//该包的作用是防止缓存穿透，，针对多次请求只向远端节点发送一次请求

//call代表正在进行中，或已经结束的请求，使用sync.WaitGroup锁避免重入
type call struct {
	wg sync.WaitGroup
	val interface{}
	err error
}

//用于管理不通key的请求，call
type Group struct {
	mu sync.Mutex
	m map[string]*call
}

func (g *Group) Do(key string, fn func() (interface{}, error)) (interface{}, error) {
	g.mu.Lock()
	//处理请求时进行上锁
	if g.m == nil {
		g.m = make(map[string]*call)
	}
	//如果请求在g.m中，说明请求正在被处理
	//然后进行解锁，并等待处理请求完毕
	if c, ok := g.m[key]; ok {
		g.mu.Unlock()
		//如果请求还在进行中，则等待
		c.wg.Wait()
		//请求结束后，返回结果
		return c.val, c.err
	}

	//如果key不存在于g.m，即还未被处理请求
	c := new(call)
	//发起请求前进行加锁，并添加到g.m，表明已经有对应的请求在处理
	c.wg.Add(1)
	g.m[key] = c
	g.mu.Unlock()

	//调用fn，发起请求
	c.val, c.err = fn()
	//请求结束
	c.wg.Done()

	g.mu.Lock()
	delete(g.m, key) //更新g.m
	g.mu.Unlock()

	return c.val, c.err
}