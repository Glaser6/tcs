<div align="center">
	<h1>Go Gateway + Go Wasm Snake</h1>
</div>

<p align="center">
	<span>English</span>
	<span> | </span>
	<a href="#zh-cn">中文</a>
</p>

<p align="center">Use Go to build a reverse proxy gateway and a WebAssembly Snake game</p>

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![WebAssembly](https://img.shields.io/badge/WebAssembly-enabled-654FF0?style=flat&logo=webassembly)](https://webassembly.org/)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat)]([#](https://github.com/Glaser6/tcs.git))

</div>

This project helps you quickly practice two useful Go skills in one repo:

- Build a reverse proxy gateway with middleware
- Compile Go code to WebAssembly and run a browser game

All you need to do is run both services, open the gateway, and play.

## English

### Project Layout

```text
tcs/
├─ gateway/   # reverse proxy, listens on :8080
└─ gowasm/    # snake game + static server, listens on :8088
```

### Features

- Reverse proxy with logging middleware, timeout middleware, and custom headers
- Go + Wasm Snake game with keyboard control and score rendering
- Clear separation: gateway layer and backend static/game service
- Good starter template for full request path learning:
	Browser -> Gateway(:8080) -> Wasm Server(:8088)

### Run Locally

1. Clone the repository

```powershell
git clone https://github.com/Glaser6/tcs.git
cd tcs
```

2. Build Wasm binary

```powershell
cd gowasm
$env:GOOS="js"
$env:GOARCH="wasm"
go build -o snake.wasm .
```

3. Start Wasm static server (port 8088)

```powershell
go run server.go
```

4. Open a new terminal and start gateway (port 8080)

```powershell
cd ..\gateway
go run main.go
```

5. Open in browser

```text
http://localhost:8080
```

### Controls

- Arrow keys or WASD: move snake
- R: restart after game over

### Notes

- The gateway forwards all traffic from :8080 to :8088
- The gateway adds header: X-Gateway: Go-Reverse-Proxy
- The gateway also injects request header: X-Forwarded-By: Go-Gateway

---

<a id="zh-cn"></a>

## 中文

### 项目结构

```text
tcs/
├─ gateway/   # 反向代理网关，监听 :8080
└─ gowasm/    # 贪吃蛇 Wasm + 静态服务，监听 :8088
```

### 特性

- 网关包含日志中间件、超时中间件、自定义响应头
- Go 编写贪吃蛇并编译为 Wasm，在浏览器运行
- 前后层职责清晰，便于学习完整请求链路
- 适合练手 Go Web + Wasm 的入门项目

### 本地运行

1. 克隆仓库

```powershell
git clone https://github.com/Glaser6/tcs.git
cd tcs
```

2. 构建 Wasm 文件

```powershell
cd gowasm
$env:GOOS="js"
$env:GOARCH="wasm"
go build -o snake.wasm .
```

3. 启动 Wasm 静态服务器（8088）

```powershell
go run server.go
```

4. 新开终端启动网关（8080）

```powershell
cd ..\gateway
go run main.go
```

5. 浏览器访问

```text
http://localhost:8080
```

### 操作说明

- 方向键 或 WASD：控制移动
- R：游戏结束后重开

### 关键代码示例

```go
// 网关把所有请求转发到 Wasm 服务
targetURL, _ := url.Parse("http://127.0.0.1:8088")
proxy := httputil.NewSingleHostReverseProxy(targetURL)
```

```go
// Wasm 游戏主循环
for range ticker.C {
		g.update()
		g.render()
}
```

## Future Ideas

- Add level speed-up and best score persistence
- Split gateway routes for static and API paths
- Add CI to auto-build snake.wasm
