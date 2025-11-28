#!/bin/bash

echo "Starting multiplayer car game server..."

# 进入服务器目录
cd "$(dirname "$0")"

# 检查服务器文件是否存在
if [ ! -f "server-linux" ]; then
    echo "❌ 错误: 找不到 server-linux 文件"
    echo "请确保 server-linux 文件存在于当前目录"
    exit 1
fi

# 给服务器文件添加执行权限
chmod +x server-linux

# 运行服务器（前台运行）
echo ""
echo "🚀 Starting server on http://0.0.0.0:9988"
echo "📡 WebSocket endpoint: ws://0.0.0.0:8899/ws"
echo ""
echo "Press Ctrl+C to stop the server"
echo ""

./server-linux
