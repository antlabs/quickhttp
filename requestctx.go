package quickhttp

import (
	"bytes"
	"io"
	"net"

	"github.com/antlabs/httparser"
)

type RequestCtx struct {
	parser     *httparser.Parser  // http 解析器
	reqSetting *httparser.Setting // 状态回调函数
	Request    Request            // 请求
	Response   Response           // 响应

	remoteAddr net.Addr

	c net.Conn
}

var zeroTCPAddr = &net.TCPAddr{
	IP: net.IPv4zero,
}

var defaultRequestSetting = &httparser.Setting{
	MessageBegin: func(p *httparser.Parser, pos int) {
	},
	URL: func(p *httparser.Parser, buf []byte, pos int) {
		ctx := p.GetUserData().(*RequestCtx)
		ctx.Request.Header.requestURI.setPosOrBytes(int32(pos-len(buf)), int32(pos), buf, minBufLimit)
	},
	Status: func(p *httparser.Parser, buf []byte, pos int) {
	},
	HeaderField: func(p *httparser.Parser, buf []byte, pos int) {
		// fmt.Printf("###(%s)\n", buf)
		ctx := p.GetUserData().(*RequestCtx)
		switch buf[0] | 0x20 { // A-Z->a-z a-z->a-z
		case 'h':
			if bytes.EqualFold(buf, bytesHost) {
				ctx.Request.Header.hState.set(findHost)
			}

		case 'u':
			if bytes.EqualFold(buf, bytesUserAgent) {
				ctx.Request.Header.hState.set(findUserAgent)
			}
		default:
			appendPair(&ctx.Request.Header.h, buf)
		}
	},
	HeaderValue: func(p *httparser.Parser, buf []byte, pos int) {
		ctx := p.GetUserData().(*RequestCtx)

		if ctx.Request.Header.hState.is(findHost) {
			ctx.Request.Header.host = append(ctx.Request.Header.host[:0], buf...)
			ctx.Request.Header.hState.clear(findHost)
		} else if ctx.Request.Header.hState.is(findUserAgent) {
			ctx.Request.Header.userAgent = append(ctx.Request.Header.userAgent[:0], buf...)
			ctx.Request.Header.hState.clear(findUserAgent)
		} else {
			setLastValue(&ctx.Request.Header.h, int32(pos-len(buf)), int32(pos), buf)
		}
	},
	HeadersComplete: func(p *httparser.Parser, pos int) {

		// r.Request.Header.method = append(r.Request.Header.method[:0], p.Method.String()...)
	},
	Body: func(p *httparser.Parser, buf []byte, pos int) {
		ctx := p.GetUserData().(*RequestCtx)
		if ctx.Request.headerAndBody.start == 0 {
			ctx.Request.headerAndBody.start = pos - len(buf)
		}

		ctx.Request.headerAndBody.end = pos + 1
	},
	MessageComplete: func(p *httparser.Parser, pos int) {
		ctx := p.GetUserData().(*RequestCtx)
		if ctx.Request.headerAndBody.start == 0 {
			ctx.Request.headerAndBody.start = pos + 1
		}
		ctx.Request.headerAndBody.end = pos + 1
	},
}

func newRequestCtx(c net.Conn) *RequestCtx {

	ctx := &RequestCtx{
		parser: httparser.New(httparser.REQUEST),
		c:      c,
		// buffer:  make([]byte, 1024),
		// request: &http.Request{},
	}

	ctx.parser.SetUserData(ctx)
	ctx.Request.Header.setParent(&ctx.Request)
	ctx.reqSetting = defaultRequestSetting

	return ctx
}

func (ctx *RequestCtx) RemoteAddr() net.Addr {
	if ctx.remoteAddr != nil {
		return ctx.remoteAddr
	}
	if ctx.c == nil {
		return zeroTCPAddr
	}
	addr := ctx.c.RemoteAddr()
	if addr == nil {
		return zeroTCPAddr
	}
	return addr
}

func (ctx *RequestCtx) URI() *URI {
	return ctx.Request.URI()
}

func (ctx *RequestCtx) Method() []byte {
	if len(ctx.Request.Header.method) == 0 {
		ctx.Request.Header.method = append(ctx.Request.Header.method[:0], ctx.parser.Method.String()...)
	}
	return ctx.Request.Header.method
}

func (ctx *RequestCtx) PostBody() []byte {
	// fmt.Printf("start = %d, end:%d\n", ctx.Request.headerAndBody.start, ctx.Request.headerAndBody.end)
	return ctx.Request.headerAndBody.getBytes(nil)
}

func (ctx *RequestCtx) RequestURI() []byte {
	return ctx.Request.Header.RequestURI()
}

func (ctx *RequestCtx) Path() []byte {
	return ctx.URI().Path()
}

func (ctx *RequestCtx) Host() []byte {
	return ctx.Request.Header.Host()
}

func (ctx *RequestCtx) QueryArgs() *Args {
	return ctx.URI().QueryArgs()
}

func (ctx *RequestCtx) UserAgent() []byte {
	return ctx.Request.Header.UserAgent()
}

func addrToIP(addr net.Addr) net.IP {
	x, ok := addr.(*net.TCPAddr)
	if !ok {
		return net.IPv4zero
	}
	return x.IP
}

func (ctx *RequestCtx) RemoteIP() net.IP {
	return addrToIP(ctx.RemoteAddr())
}

func (ctx *RequestCtx) execute() (bool, error) {
	_, err := ctx.parser.Execute(ctx.reqSetting, *ctx.Request.headerAndBody.buf)
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
