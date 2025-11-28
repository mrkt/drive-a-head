#!/bin/bash

echo "=========================================="
echo "  多人在线碰撞车游戏服务器"
echo "=========================================="
echo ""

# 获取脚本所在目录
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
cd "$SCRIPT_DIR"

# 检查服务器文件是否存在
if [ ! -f "server-linux" ]; then
    echo "❌ 错误: 找不到 server-linux 文件"
    exit 1
fi

# 给服务器文件添加执行权限
chmod +x server-linux

# 检查端口是否被占用
PORT=8899
if lsof -Pi :$PORT -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo "⚠️  警告: 端口 $PORT 已被占用"
    echo "正在尝试停止旧进程..."
    pkill -f server-linux
    sleep 2
fi

# 启动服务器
echo "🚀 启动游戏服务器..."
echo ""
echo "📡 WebSocket 服务: ws://0.0.0.0:8899/ws"
echo "🌐 HTTP 服务: http://0.0.0.0:9988"
echo ""
echo "玩家可以通过以下地址访问游戏:"
echo "  http://你的服务器IP:9988"
echo ""
echo "按 Ctrl+C 停止服务器"
echo "=========================================="
echo ""

# 运行服务器（后台运行并记录日志）
nohup ./server-linux > server.log 2>&1 &
SERVER_PID=$!

echo "✅ 服务器已启动 (PID: $SERVER_PID)"
echo "📝 日志文件: $SCRIPT_DIR/server.log"
echo ""
echo "查看实时日志: tail -f $SCRIPT_DIR/server.log"
echo "停止服务器: kill $SERVER_PID 或 pkill -f server-linux"
echo ""
