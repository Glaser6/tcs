package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

// 日志中间件：记录请求方法、路径和耗时
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("--> %s %s", r.Method, r.URL.Path)

		// 包装 ResponseWriter 以捕获状态码（可选）
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wrapped, r)

		log.Printf("<-- %s %s %d %v", r.Method, r.URL.Path, wrapped.statusCode, time.Since(start))
	})
}

// 用于捕获 HTTP 状态码的 ResponseWriter 包装器
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// 超时中间件：为每个请求设置最大处理时间
func timeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// 使用 http.TimeoutHandler 实现超时控制
			http.TimeoutHandler(next, timeout, "Gateway Timeout").ServeHTTP(w, r)
		})
	}
}

// 自定义中间件示例：添加一个响应头标识网关
func addHeaderMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Gateway", "Go-Reverse-Proxy")
		next.ServeHTTP(w, r)
	})
}

func main() {
	// 后端服务地址
	targetURL, err := url.Parse("http://127.0.0.1:8088")
	if err != nil {
		log.Fatal("Invalid target URL:", err)
	}

	// 创建标准反向代理
	proxy := httputil.NewSingleHostReverseProxy(targetURL)

	// 可选：修改 Director 以添加自定义请求头或路径改写
	// 默认 Director 已经设置了 Host 头和 URL 重写
	// 这里演示如何在转发前添加一个请求头
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-By", "Go-Gateway")
	}

	// 可选：自定义错误处理（如后端不可达时返回友好错误）
	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("Proxy error: %v", err)
		http.Error(w, "Bad Gateway", http.StatusBadGateway)
	}

	// 创建路由器并注册中间件链
	mux := http.NewServeMux()
	mux.Handle("/", proxy) // 所有请求都走代理

	// 组装中间件：注意顺序 —— 请求先经过外层，响应后经过内层
	// 执行顺序：addHeader -> logging -> timeout -> proxy
	var handler http.Handler = mux
	handler = addHeaderMiddleware(handler)
	handler = loggingMiddleware(handler)
	handler = timeoutMiddleware(30 * time.Second)(handler)

	addr := ":8080"
	log.Printf("Gateway listening on %s, forwarding to %s", addr, targetURL)

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal("Server failed:", err)
	}
}
