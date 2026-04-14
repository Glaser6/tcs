package main

import (
	"embed"
	"log"
	"net/http"
	// "fmt" // 优化方案: 若使用Go生成HTML，需导入此包
)

//go:embed snake.wasm wasm_exec.js index.html
var fs embed.FS

func main() {
	// 静态文件服务，当前目录作为根目录
	http.Handle("/", http.FileServer(http.FS(fs)))
	log.Println("服务器启动: http://localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}

//
// 可选方案：
// 3. 添加导入: "fmt"
// 4. 修改 main() 为如下方式：
//
//    func main() {
//		  ...
//        http.HandleFunc("/", serveIndex)
//        ...
//    }
//
//    func serveIndex(w http.ResponseWriter, r *http.Request) {
//        w.Header().Set("Content-Type", "text/html; charset=utf-8")
//        fmt.Fprint(w, `<!DOCTYPE html>
//    <html lang="en">
//    <head>
//        <meta charset="UTF-8">
//        <meta name="viewport" content="width=device-width, initial-scale=1.0">
//        <title>Go + Wasm 贪吃蛇</title>
//        <style>
//            body { display: flex; flex-direction: column; align-items: center;
//                   background-color: #f5f5f5; margin-top: 50px; }
//            canvas { border: 2px solid #333; background-color: #ffffff; }
//            .tips { margin-top: 20px; font-family: Arial, sans-serif; color: #333; }
//        </style>
//    </head>
//    <body>
//    <canvas id="snakeCanvas" width="600" height="600"></canvas>
//    <div class="tips">
//        方向键/WASD 控制移动 | 游戏结束后按 R 重新开始 | Go + Wasm 实现
//    </div>
//    <script src="wasm_exec.js"><\/script>
//    <script>
//        async function runWasmGame() {
//            const go = new Go();
//            const response = await fetch("snake.wasm");
//            const bytes = await response.arrayBuffer();
//            const result = await WebAssembly.instantiate(bytes, go.importObject);
//            go.run(result.instance);
//        }
//        runWasmGame();
//    </script>
//    </body>
//    </html>`)
//    }
//
