package main

import (
	"embed"
	"log"
	"net/http"
)

//go:embed snake.wasm wasm_exec.js index.html
var fs embed.FS

func main() {
	// 静态文件服务，当前目录作为根目录
	http.Handle("/", http.FileServer(http.Dir(".")))
	log.Println("服务器启动: http://localhost:8088")
	log.Fatal(http.ListenAndServe(":8088", nil))
}
