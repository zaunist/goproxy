package main

import (
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/elazarl/goproxy"
)

func main() {
	proxy := goproxy.NewProxyHttpServer()
	proxy.Verbose = true

	// 添加 CONNECT 处理器，只允许 cursor.sh 域名的 HTTPS 连接
	proxy.OnRequest().HandleConnectFunc(func(host string, ctx *goproxy.ProxyCtx) (*goproxy.ConnectAction, string) {
		log.Printf("[%03d] 收到 CONNECT 请求: %s", ctx.Session, host)

		if strings.Contains(host, ".cursor.sh") {
			log.Printf("[%03d] 允许连接到: %s", ctx.Session, host)
			return goproxy.OkConnect, host
		}

		log.Printf("[%03d] 拒绝连接到: %s", ctx.Session, host)
		return goproxy.RejectConnect, host
	})

	// 处理 HTTP 请求
	proxy.OnRequest().DoFunc(func(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
		log.Printf("[%03d] 收到 HTTP 请求: %s %s", ctx.Session, req.Method, req.URL)

		// 允许 cursor.sh 域名的请求
		if strings.Contains(req.Host, ".cursor.sh") {
			log.Printf("[%03d] 转发请求: %s %s", ctx.Session, req.Method, req.URL)
			return req, nil
		}

		log.Printf("[%03d] 重定向请求: %s %s 到 baidu.com", ctx.Session, req.Method, req.URL)
		// 其他域名重定向到百度
		return req, goproxy.NewResponse(req,
			goproxy.ContentTypeText,
			http.StatusMovedPermanently,
			"Redirecting to Baidu")
	})

	port := ":8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = ":" + envPort
	}
	log.Printf("HTTP Proxy Server 启动在: 0.0.0.0%s", port)
	log.Fatal(http.ListenAndServe("0.0.0.0"+port, proxy))
}
