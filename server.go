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

	try := 1
	ctx.Request.headerAndBody.resetBuf(GetBytes(poolPage))
	buf := ctx.Request.headerAndBody.buf
	pos := 0
	for {
		n, err := c.Read((*buf)[pos:])
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}

		if n == 0 {
			fmt.Printf("find eof\n")
			return
		}

		fmt.Printf("n = %d\n", n)
		*buf = (*buf)[:n]
		sucess, err := ctx.execute()
		if err != nil {
			fmt.Printf("ctx.execute %v\n", err)
			return
		}
		if sucess {
			*buf = (*buf)[:cap(*buf)]
			ctx.Response.Header.SetServer(defaultServerName)
			ctx.Response.init()

			s.Handler(ctx)

			ctx.write(c)
			ctx.reset()
			ctx.Request.headerAndBody.resetBuf(GetBytes(poolPage))
			pos = 0
			try = 1
		} else {
			try *= 2
			newBuf := GetBytes(poolPage * try)
			oldBuf := ctx.Request.headerAndBody.getBufPtr()
			copy(*newBuf, *oldBuf)
			ctx.Request.headerAndBody.resetBuf(newBuf)
			pos += n
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
