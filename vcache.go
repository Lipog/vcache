package Vcache
//负责与外部的交互，控制缓存存储和获取的主流程
//逻辑是现在本地寻找，找不到去远程节点获取，还获取不到的话，就回调，将值添加到缓存

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