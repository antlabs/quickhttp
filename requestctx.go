package quickhttp

import (
	"io"
	"net"

	"github.com/antlabs/httparser"
)

type RequestCtx struct {
	parser     *httparser.Parser  // http 解析器
	reqSetting *httparser.Setting // 状态回调函数
	buf        *[]byte            // header+小body, 或者header
	Request    Request            // 请求
	Response   Response           // 响应
}

var defaultRequestSetting = &httparser.Setting{
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
		r := p.GetUserData().(*RequestCtx)

		r.Request.bodyStart = pos
	},
	Body: func(p *httparser.Parser, buf []byte, pos int) {
		r := p.GetUserData().(*RequestCtx)
		r.Request.bodyEnd = pos
	},
	MessageComplete: func(p *httparser.Parser, pos int) {
		r := p.GetUserData().(*RequestCtx)
		r.Request.bodyEnd = pos
	},
}

func newRequestCtx() *RequestCtx {

	buf := GetBytes(1024)
	r := &RequestCtx{
		parser: httparser.New(httparser.REQUEST),
		// buffer:  make([]byte, 1024),
		// request: &http.Request{},
	}

	r.parser.SetUserData(r)
	r.reqSetting = defaultRequestSetting
	r.buf = buf

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
	_, err := r.parser.Execute(r.reqSetting, *r.buf)
	return r.parser.EOF(), err
}

func (r *RequestCtx) Write(p []byte) (int, error) {
	r.Response.AppendBody(p)
	return len(p), nil
}

func (r *RequestCtx) write(w io.Writer) error {
	return r.Response.Write(w)
}
