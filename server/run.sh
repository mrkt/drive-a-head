#!/bin/bash

echo "Starting multiplayer car game server..."

# 进入服务器目录
cd "$(dirname "$0")"

# 检查是否已编译
if [ ! -f "game-server" ]; then
    echo "Server not built yet. Building now..."
    ./build.sh
    if [ $? -ne 0 ]; then
        exit 1
    fi
fi

# 运行服务器
echo ""
echo "🚀 Starting server on http://localhost:9988"
echo "📡 WebSocket endpoint: ws://localhost:9988/ws"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./game-server
