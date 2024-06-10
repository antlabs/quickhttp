package quickhttp

type Request struct {
	Header        RequestHeader
	headerAndBody buf[int] //存在header和body

	uri       URI
	parsedURI bool
	isTLS     bool
}

func (r *Request) Reset() {
	r.Header.Reset()
}

func (req *Request) parseURI() error {
	if req.parsedURI {
		return nil
	}
	req.parsedURI = true

	return req.uri.parse(req.Header.Host(), req.Header.RequestURI(), req.isTLS)
}

func (req *Request) URI() *URI {
	req.parseURI() //nolint:errcheck
	return &req.uri
}
