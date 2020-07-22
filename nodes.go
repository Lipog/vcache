package vcache

import "vcache/vcachepb"

//该包的作用是抽象两个接口
//节点选择接口和节点缓存查找接口

type NodePicker interface {
	PickNode(key string) (node NodeGetter, ok bool)
}

type NodeGetter interface {
	Get(in *vcachepb.Request, out *vcachepb.Response) error
}

//type NodeGetter interface {
//	Get(group string, key string) ([]byte, error)
//}
