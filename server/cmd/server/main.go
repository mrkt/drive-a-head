package main

import (
	"log"
	"net/http"
	"server/internal/server"
	"strings"
)

// addCompressionHeaders 添加压缩头
func addCompressionHeaders(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		// Gzip 压缩文件
		if strings.HasSuffix(path, ".gz") {
			w.Header().Set("Content-Encoding", "gzip")
			if strings.Contains(path, ".js") {
				w.Header().Set("Content-Type", "application/javascript")
			} else if strings.Contains(path, ".wasm") {
				w.Header().Set("Content-Type", "application/wasm")
			} else if strings.Contains(path, ".data") {
				w.Header().Set("Content-Type", "application/octet-stream")
			}
		}
		h.ServeHTTP(w, r)
	})
}

func main() {
	// 创建游戏服务器
	gameServer := server.NewServer()

	// 启动服务器协程
	go gameServer.Run()

	// 设置HTTP路由
	http.HandleFunc("/ws", gameServer.HandleWebSocket)

	// 提供静态文件服务（用于托管WebGL构建）
	fs := http.FileServer(http.Dir("../Builds/WebGL"))
	http.Handle("/", addCompressionHeaders(fs))

	// 启动HTTP服务器
	addr := ":8899"
	log.Printf("Server starting on %s", addr)
	log.Printf("WebSocket endpoint: ws://localhost%s/ws", addr)
	log.Printf("Game URL: http://localhost:9988")

	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe error: ", err)
	}
}
