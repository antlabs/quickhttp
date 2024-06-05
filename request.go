package quickhttp

type Request struct {
	Header    RequestHeader
	bodyStart int
	bodyEnd   int
}

func (r *Request) Reset() {
	r.Header.Reset()
}
