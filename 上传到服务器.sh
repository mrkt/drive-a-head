#!/bin/bash

# 上传更新到服务器的脚本
# 使用方法：./上传到服务器.sh 你的服务器IP

if [ -z "$1" ]; then
    echo "❌ 错误: 请提供服务器IP地址"
    echo "使用方法: ./上传到服务器.sh 你的服务器IP"
    echo "例如: ./上传到服务器.sh 123.45.67.89"
    exit 1
fi

SERVER_IP=$1
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "=========================================="
echo "  上传游戏更新到服务器"
echo "=========================================="
echo ""
echo "📡 目标服务器: $SERVER_IP"
echo "📦 部署目录: $SCRIPT_DIR"
echo ""

# 确认上传
read -p "确认上传到服务器 $SERVER_IP? (y/n) " -n 1 -r
echo ""
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "❌ 取消上传"
    exit 1
fi

echo ""
echo "🚀 开始上传..."
echo ""

# 上传 WebGL 目录
echo "📤 上传 WebGL 文件..."
scp -r "$SCRIPT_DIR/Builds/WebGL" root@$SERVER_IP:/opt/car-game/deployment/Builds/

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "✅ 上传完成！"
    echo "=========================================="
    echo ""
    echo "🌐 访问游戏："
    echo "   http://$SERVER_IP:8899"
    echo ""
    echo "📝 注意事项："
    echo "   - 服务器无需重启"
    echo "   - 刷新浏览器即可看到更新"
    echo "   - 建议清除浏览器缓存（Ctrl+F5）"
    echo ""
    echo "🧪 测试清单："
    echo "   ✓ PC端：测试转向是否更平滑"
    echo "   ✓ PC端：测试车辆是否不再总飞天"
    echo "   ✓ 手机端：测试虚拟摇杆是否为模拟控制"
    echo "   ✓ 手机端：测试轻推/重推摇杆的速度变化"
    echo ""
else
    echo ""
    echo "❌ 上传失败！"
    echo "请检查："
    echo "   1. 服务器IP是否正确"
    echo "   2. SSH连接是否正常"
    echo "   3. 服务器目录是否存在"
    exit 1
fi
