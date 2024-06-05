package quickhttp

// 减少copy, 如果超过一定大小的数据包，直接引用底层的read buffer
// 如果是很小的buf，直接copy
type buf struct {
	start int
	end   int
	buf   *[]byte
}

// 设置偏移位置
func (b *buf) setPos(start, end int) {
	b.start = start
	b.end = end
}

// 设置数据
func (b *buf) setBytes(buf []byte) {
	b.start, b.end = 0, 0

	if b.buf == nil {
		if len(buf) >= poolPage {
			b.buf = GetAndResetBytes(len(buf))
		} else {
			b2 := make([]byte, len(buf))
			b2 = b2[:0]
			b.buf = &b2
		}
	}
	*b.buf = append(*b.buf, buf...)
}

// 根据limit的大小设置偏移位置和[]byte
func (b *buf) setPosOrBytes(start, end int, buf []byte, limit int) {
	if len(buf) < limit {
		b.setBytes(buf)
		return
	}

	b.setPos(start, end)
}

// reset函数
func (b *buf) reset() {
	b.start, b.end = 0, 0
	if b.buf != nil && cap(*b.buf) > poolPage {
		PutBytes(b.buf)
	}
}
