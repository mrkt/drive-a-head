#!/bin/bash

echo "Building multiplayer car game server..."

# 进入服务器目录
cd "$(dirname "$0")"

# 下载依赖
echo "Downloading dependencies..."
go mod download

# 编译服务器
echo "Building server..."
cd cmd/server
go build -o ../../game-server main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "Server binary: server/game-server"
    echo ""
    echo "To run the server:"
    echo "  cd server"
    echo "  ./game-server"
else
    echo "❌ Build failed!"
    exit 1
fi
