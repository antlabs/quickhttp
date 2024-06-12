package quickhttp

import (
	"fmt"
	"net"
)

type RequestHandler func(ctx *RequestCtx)

type Server struct {
	Handler RequestHandler

	Name string //服务名
}

func (s *Server) serve(c net.Conn) {
	ctx := newRequestCtx(c)
	defer c.Close()

	ctx.Request.headerAndBody.resetBuf(GetBytes(1024))
	buf := ctx.Request.headerAndBody.buf
	for {
		n, err := c.Read(*buf)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		fmt.Printf("n = %d\n", n)
		*buf = (*buf)[:n]
		sucess, err := ctx.execute()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		if sucess {
			*buf = (*buf)[:cap(*buf)]
			ctx.Response.Header.SetServer(defaultServerName)
			ctx.Response.init()

			s.Handler(ctx)

			ctx.write(c)
			ctx.reset()
		}
	}
}

func (s *Server) Serve(ln net.Listener) error {
	for {
		con, err := ln.Accept()
		if err != nil {
			fmt.Printf("accept:%v\n", err)
			continue
		}
		go s.serve(con)
	}
}

func (s *Server) ListenAndServe(addr string) error {
	ln, err := net.Listen("tcp4", addr)
	if err != nil {
		return err
	}
	return s.Serve(ln)
}

func ListenAndServe(addr string, handler RequestHandler) error {
	s := &Server{
		Handler: handler,
	}
	return s.ListenAndServe(addr)
}
