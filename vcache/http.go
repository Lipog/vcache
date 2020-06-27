package vcache

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

const defaultBasePath = "/_vcache/"

type HTTPPool struct {
	self string  //用来记录自己的地址，包括主机名/IP和端口
	basePath string //作为节点间通讯的前缀
}

func NewHTTPPool(self string) *HTTPPool {
	return &HTTPPool{
		self:     self,
		basePath: defaultBasePath,
	}
}

//记录对应服务名的信息
func (p *HTTPPool) Log(format string, v ...interface{}) {
	log.Printf("[Serve %s] %s", p.self, fmt.Sprintf(format, v...))
}

func (p *HTTPPool) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//如果请求不含有基础前缀的话，就报错
	//约定访问路径格式为 /<basepath>/<groupname>/<key>，
	// 通过 groupname 得到 group 实例，
	// 再使用 group.Get(key) 获取缓存数据
	if !strings.HasPrefix(r.URL.Path, p.basePath) {
		panic("HTTPPool serving unexpected path:" + r.URL.Path)
	}

	p.Log("%s %s", r.Method, r.URL.Path)
	//将除去前缀的部分，分割为两个部分
	parts := strings.SplitN(r.URL.Path[len(p.basePath):], "/", 2)
	if len(parts) != 2 {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	groupName := parts[0]
	key := parts[1]

	group := GetGroup(groupName)
	if group == nil {
		http.Error(w, "no such group:" + groupName, http.StatusNotFound)
		return
	}

	view, err := group.Get(key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(view.ByteSlice())
}





