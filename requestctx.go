package quickhttp

import (
	"net"

	"github.com/antlabs/httparser"
)

type RequestCtx struct {
	parser  *httparser.Parser
	setting *httparser.Setting
	buf     *[]byte // header+小body, 或者header

	Request Request

	// Outgoing response.
	//
	// Copying Response by value is forbidden. Use pointer to Response instead.
	Response Response
}

func newRequestCtx() *RequestCtx {

	buf := GetBytes(1024)
	r := &RequestCtx{
		parser: httparser.New(httparser.REQUEST),
		// buffer:  make([]byte, 1024),
		// request: &http.Request{},
	}

	r.buf = buf
	r.setting = &httparser.Setting{
		MessageBegin: func(p *httparser.Parser, pos int) {
		},
		URL: func(p *httparser.Parser, buf []byte, pos int) {
		},
		Status: func(p *httparser.Parser, buf []byte, pos int) {
		},
		HeaderField: func(p *httparser.Parser, buf []byte, pos int) {
		},
		HeaderValue: func(p *httparser.Parser, buf []byte, pos int) {
		},
		HeadersComplete: func(p *httparser.Parser, pos int) {
			r.Request.bodyStart = pos
		},
		Body: func(p *httparser.Parser, buf []byte, pos int) {
			r.Request.bodyEnd = pos
		},
		MessageComplete: func(p *httparser.Parser, pos int) {
			r.Request.bodyEnd = pos
		},
	}

	return r
}

func (r *RequestCtx) Method() []byte {
	if len(r.Request.Header.method) == 0 {
		r.Request.Header.method = str2bytes(r.parser.Method.String())
	}
	return r.Request.Header.method
}

func (r *RequestCtx) PostBody() []byte {
	return (*r.buf)[r.Request.bodyStart:r.Request.bodyEnd]
}

func (r *RequestCtx) RequestURI() []byte {

	return nil
}

func (r *RequestCtx) Path() []byte {
	return nil
}

func (r *RequestCtx) Host() []byte {
	return nil
}

func (r *RequestCtx) QueryArgs() []byte {
	return nil
}

func (r *RequestCtx) UserAgent() []byte {
	return nil
}

func (r *RequestCtx) RemoteIP() net.IP {
	return nil
}

func (r *RequestCtx) execute() (bool, error) {
	_, err := r.parser.Execute(r.setting, *r.buf)
	return r.parser.EOF(), err
}
