package quickhttp

type Request struct {
	Header        RequestHeader
	headerAndBody buf[int] //存在header和body
}

func (r *Request) Reset() {
	r.Header.Reset()
}
