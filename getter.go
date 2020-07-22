package vcache

//回调要求实现Getter接口，只要实现了该接口的Get方法即可
type Getter interface {
	Get(key string) ([]byte, error)
}
//GetterFunc用来实现Getter接口，实现Get方法,
//该函数实现了通过Get来调用自己，从而获得数据
type GetterFunc func(key string) ([]byte, error)
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}