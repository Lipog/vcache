package vcache
//http服务短对要传输的数据进行proto.Marshal编码

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"log"
	"net/http"
	"strings"
	"sync"
	"vcache/onehash"
	"vcache/vcachepb"
)

//该包作为http的服务端，提供被其他节点访问的能力

const defaultBasePath = "/_vcache/"
const defaultReplicas = 50

//HttpPool是作为节点间进行通信的核心数据
//包括了客户端和服务端
type HttpPool struct {
	self_addr string //self_addr记录当前节点的ip地址和端口号
	basePath string //basePath是当前节点的api入口，即保证传来的URL是用作缓存的
	mu sync.Mutex //保证节点选择时的并发安全
	hashNodes *onehash.OneHash //nodes映射一致性hash，用来进行节点选择
	httpGetters map[string]*httpGetter //保存节点对应的本身的httpGetter，用来实现远程节点Get数据的方法
}



func NewHttpPool(self_addr string) *HttpPool {
	return &HttpPool{
		self_addr: self_addr,
		basePath:  defaultBasePath,
		hashNodes: onehash.NewOneHash(defaultReplicas, nil),
		httpGetters: make(map[string]*httpGetter),
	}
}

func (p *HttpPool) Set(nodes ...string) {
	//设置新的节点时，要保证数据的安全
	p.mu.Lock()
	defer p.mu.Unlock()
	//如果服务端的一致hash还为nil，说明还没有初始化，要对其进行初始化
	if p.hashNodes == nil {
		p.hashNodes = onehash.NewOneHash(defaultReplicas, nil)
		p.httpGetters = make(map[string]*httpGetter)
	}
	p.hashNodes.AddNode(nodes...)
	for _, node := range nodes {
		p.httpGetters[node] = &httpGetter{baseURL:node + p.basePath}
	}
}

func (p *HttpPool) PickNode(key string) (NodeGetter, bool) {
	p.mu.Lock()
	defer p.mu.Unlock()
	//从hash环中找到key对应的节点的地址
	//如果节点的地址不为空并且不是本身的地址，那么就返回该地址对应的httpGetter
	//httpGetter实现了Get方法，可以根据其地址从对应的值获得数据
	if node := p.hashNodes.Get(key); node != "" && node != p.self_addr {
		p.Log("Pick node %s", node)
		return p.httpGetters[node], true
	}
	return nil, false
}

func (p *HttpPool) Log(format string, v ...interface{}) {
	log.Printf("[Serve %s] %s", p.self_addr, fmt.Sprintf(format, v...))
}

//实现serveHTTP接口，从而接管w和r
func (p *HttpPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HttpPool serving unexpected path: " + r.URL.Path)
	}

	//如果requet的请求是缓存请求，即含有basePath，那么就对第其进行处理
	p.Log("%s %s", r.Method, r.URL.Path)
	//然后对请求的Group和key进行分离
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	//如果请求的路径不完整，那么就400
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	//从全局的groups中寻找该group，如果没有找到该组，那么就404
	group := GetGruop(groupName)
	if group == nil {
		http.Error(w, "no such group:" + groupName, http.StatusNotFound)
		return
	}

	//如果查找数据过程中出现错误，那么就500
	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//这里通过protobuf将响应的数据写入到body中
	body, err := proto.Marshal(&vcachepb.Response{
		Value: view.ByteSlice(),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Write(body)
}

