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

		// r.Request.Header.method = append(r.Request.Header.method[:0], p.Method.String()...)
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

func (ctx *RequestCtx) Method() []byte {
	if len(ctx.Request.Header.method) == 0 {
		ctx.Request.Header.method = append(ctx.Request.Header.method[:0], ctx.parser.Method.String()...)
	}
	return ctx.Request.Header.method
}

func (ctx *RequestCtx) PostBody() []byte {
	return (*ctx.buf)[ctx.Request.bodyStart:ctx.Request.bodyEnd]
}

func (ctx *RequestCtx) RequestURI() []byte {

	return nil
}

func (ctx *RequestCtx) Path() []byte {
	return nil
}

func (ctx *RequestCtx) Host() []byte {
	return nil
}

func (ctx *RequestCtx) QueryArgs() []byte {
	return nil
}

func (ctx *RequestCtx) UserAgent() []byte {
	return nil
}

func (ctx *RequestCtx) RemoteIP() net.IP {
	return nil
}

func (ctx *RequestCtx) execute() (bool, error) {
	_, err := ctx.parser.Execute(ctx.reqSetting, *ctx.buf)
	return ctx.parser.EOF(), err
}

func (ctx *RequestCtx) Write(p []byte) (int, error) {
	ctx.Response.AppendBody(p)
	return len(p), nil
}

func (ctx *RequestCtx) write(w io.Writer) error {
	return ctx.Response.Write(w)
}

func (ctx *RequestCtx) reset() {
	ctx.Request.Reset()
	ctx.Response.Reset()
	ctx.parser.Reset()
}
