package quickhttp

import (
	"github.com/antlabs/httparser"
)

type RequestCtx struct {
	parser    *httparser.Parser
	setting   *httparser.Setting
	buf       *[]byte
	bodyStart int
	bodyEnd   int
}

func newRequestCtx() *RequestCtx {

	buf := GetBytes(1024)
	r := &RequestCtx{
		parser: httparser.New(httparser.REQUEST),
		// buffer:  make([]byte, 1024),
		// request: &http.Request{},
	}
	r.buf = buf

	r.buf = buf
	r.setting = &httparser.Setting{
		MessageBegin: func(p *httparser.Parser) {
		},
		URL: func(p *httparser.Parser, buf []byte) {
			//url数据
			//fmt.Printf("url->%s\n", buf)
			// hConn.request.RequestURI = string(buf)
		},
		Status: func(p *httparser.Parser, buf []byte) {
			// 响应包才需要用到
		},
		HeaderField: func(p *httparser.Parser, buf []byte) {
			// http header field
			// fmt.Printf("header field:%s\n", buf)
			// hConn.lastHeader = string(buf)
		},
		HeaderValue: func(p *httparser.Parser, buf []byte) {
			// http header value
			//fmt.Printf("header value:%s\n", buf)
			// if "Host" == hConn.lastHeader {
			// 	hConn.request.Host = string(buf)
			// } else {
			// 	hConn.request.Header.Add(hConn.lastHeader, string(buf))
			// }
			// hConn.lastHeader = ""
		},
		HeadersComplete: func(p *httparser.Parser, pos int) {
			r.bodyStart = pos
		},
		Body: func(p *httparser.Parser, buf []byte, pos int) {
			r.bodyEnd = pos
		},
		MessageComplete: func(p *httparser.Parser, pos int) {
			r.bodyEnd = pos
		},
	}

	return r
}

func (r *RequestCtx) PostBody() []byte {
	return (*r.buf)[r.bodyStart:r.bodyEnd]
}

func (r *RequestCtx) execute() (bool, error) {
	_, err := r.parser.Execute(r.setting, *r.buf)
	return r.parser.EOF(), err
}
