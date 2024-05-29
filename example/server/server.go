package main

import (
	"fmt"
	"net/http"

	"github.com/antlabs/quickhttp"
)

// handler 函数处理所有传入的 HTTP 请求
func handler(w http.ResponseWriter, r *http.Request) {
	// 写入响应体
	fmt.Fprintf(w, "Hello, World!")
}

func stdMain() {
	// http.HandleFunc 用于将 URL 路径和 handler 函数关联起来
	http.HandleFunc("/", handler)

	// 启动 HTTP 服务器并监听端口 8080
	fmt.Println("Starting server at :8080")
	if err := http.ListenAndServe(":8084", nil); err != nil {
		// 如果启动失败，打印错误信息
		fmt.Printf("Error starting server: %s\n", err)
	}
}

// request handler in quickhttp style, i.e. just plain function.
func quickHTTPHandler(ctx *quickhttp.RequestCtx) {
	// fmt.Printf("(%s)\n", ctx.PostBody())
	// fmt.Fprintf(ctx, "Hi there! RequestURI is %q", ctx.RequestURI())
}

func main() {

	go stdMain()
	// pass plain function to quickhttp
	fmt.Println(quickhttp.ListenAndServe(":8083", quickHTTPHandler))
}
