package vcache
//这里是缓存值的抽象于封装

//ByteView用来表示只读的数据结构，表示缓存
type ByteView struct {
	b []byte  //b将会存储真实的缓存值
}

func (v ByteView) Len() int {
	return len(v.b)
}

//返回一个v.b的深拷贝
//因为b是只读的，返回一个拷贝，可以防止缓存值被外部程序修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func (v ByteView) String() string {
	return string(v.b)
}

//深拷贝函数
func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
