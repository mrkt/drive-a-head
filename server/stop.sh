#!/bin/bash

echo "正在停止游戏服务器..."

# 查找并停止所有 server-linux 进程
pkill -f server-linux

if [ $? -eq 0 ]; then
    echo "✅ 服务器已停止"
else
    echo "⚠️  没有找到运行中的服务器进程"
fi
