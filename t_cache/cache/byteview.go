package cache

type ByteView struct {
	// 存储真实的缓存值，byte 类型能够支持任意的数据类型的存储，例如字符串、图片等
	b []byte
}

// Len 实现Len方法，达到LRU的Value接口要求
func (bv ByteView) Len() int {
	return len(bv.b)
}

// NewByteView 将任意对象序列化为 []byte
func NewByteView(objBytes []byte) (ByteView, error) {
	return ByteView{b: objBytes}, nil
}

func (bv ByteView) ByteSlice() []byte {
	return cloneBytes(bv.b)
}

func cloneBytes(b []byte) []byte {
	c := make([]byte, len(b))
	copy(c, b)
	return c
}
