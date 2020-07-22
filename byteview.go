package vcache
//这里封装了只读的数据结构Byteview，用来表示缓存值
//实现了缓存值的抽象与封装

type ByteView struct {
	b []byte
}

func (v ByteView) Len() int {
	return len(v.b)
}

//该方法用来实现ByteView的只读，防止缓存值呗修改
func (v ByteView) ByteSlice() []byte {
	return cloneBytes(v.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}

func (v ByteView) String() string {
	return string(v.b)
}
