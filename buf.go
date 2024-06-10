package quickhttp

// 类型约束，支持各种整型
type Integer interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | uintptr
}

const minBufLimit = 64

// 减少copy, 如果超过一定大小的数据包，直接引用底层的read buffer
// 如果是很小的buf，直接copy
type buf[T Integer] struct {
	start T
	end   T
	buf   *[]byte
}

// 设置偏移位置
func (b *buf[T]) setPos(start, end T) {
	b.start = start
	b.end = end
}

// 设置数据
func (b *buf[T]) setBytes(data []byte) {
	if b.buf == nil {
		if len(data) >= poolPage {
			b.buf = GetAndResetBytes(len(data))
		} else {
			b2 := make([]byte, len(data))
			b2 = b2[:0]
			b.buf = &b2
		}
	}
	b.start, b.end = 0, T(len(data))
	*b.buf = append(*b.buf, data...)
}

// 根据limit的大小设置偏移位置和[]byte
func (b *buf[T]) setPosOrBytes(start, end T, data []byte, limit int) {
	if len(data) < limit {
		b.setBytes(data)
		return
	}

	b.setPos(start, end)
}

func (b *buf[T]) getBytes(bigBytes []byte) []byte {
	if b.buf != nil {
		return (*b.buf)[b.start:b.end]
	}

	return bigBytes[b.start:b.end]
}

func (b *buf[T]) resetBuf(b2 *[]byte) {
	b.buf = b2
}

// reset函数
func (b *buf[T]) reset() {
	b.start, b.end = 0, 0
	if b.buf != nil {
		if cap(*b.buf) > poolPage {
			PutBytes(b.buf)
		} else {
			*b.buf = (*b.buf)[:0]
		}
	}
}
