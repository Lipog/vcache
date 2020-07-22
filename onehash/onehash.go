package onehash

import (
	"hash/crc32"
	"log"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

type OneHash struct {
	hash Hash //hash算法
	replicas int //节点的复制倍数
	circle []int //circle用来存储hash环上的所有的节点，并且排序
	hashMap map[int]string //用来存储虚拟节点和真实节点地址的映射
}

//采用依赖注入的方式，允许替换自定义的hash函数，默认为crc32.ChecksumIEEE
func NewOneHash(replicas int, fn Hash) *OneHash {
	m := &OneHash{
		replicas: replicas,
		hash: fn,
		hashMap: make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

//AddNode方法用来添加真实节点，并且添加到hash环中去,
//允许传入0个或者多个真实节点的名称
//对于每一个真实节点key，创建h.replicas个虚拟节点
//然后用h.hash计算虚拟节点的hash值，添加到hash环上，并建立与真实节点的映射
func (h *OneHash) AddNode(keys ...string) {
	for _, key := range keys {
		for i := 0; i < h.replicas; i++ {
			hash := int(h.hash([]byte(strconv.Itoa(i) + key)))
			h.circle = append(h.circle, hash)
			h.hashMap[hash] = key
		}
	}
	//添加完节点后，对节点进行排序
	sort.Ints(h.circle)
}

func (h *OneHash) Get(key string) string {
	//如果节点的hash环为空，那么也就得不到任何节点信息，返回空
	if len(h.circle) == 0 {
		return ""
	}

	//然后计算key的对应的hash转换后的值
	hash := int(h.hash([]byte(key)))
	index := sort.Search(len(h.circle), func(i int) bool {
		return h.circle[i] >= hash
	})
	log.Println("get hash value:", hash)
	//因为或许为最后一个元素，所以最后一个hash对应的是hash环的第一个节点
	//因此要进行取模的运算
	return h.hashMap[h.circle[index % len(h.circle)]]
}