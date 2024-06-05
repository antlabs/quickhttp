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

// var bytesBody = []byte("HTTP/1.1 200 OK \r\nContent-Length: 0\r\n\r\n")

func (s *Server) serve(c net.Conn) {
	ctx := newRequestCtx()
	defer c.Close()
	for {
		n, err := c.Read(*ctx.buf)
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		*ctx.buf = (*ctx.buf)[:n]
		sucess, err := ctx.execute()
		if err != nil {
			fmt.Printf("%v\n", err)
			return
		}
		if sucess {
			*ctx.buf = (*ctx.buf)[:cap(*ctx.buf)]
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
