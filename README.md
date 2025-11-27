# 多人在线碰撞车游戏 - 部署文档

## 📦 部署包内容

```
deployment/
├── server/
│   ├── game-server      # 游戏服务器可执行文件
│   ├── start.sh         # 启动脚本
│   └── stop.sh          # 停止脚本
└── Builds/
    └── WebGL/           # WebGL 客户端文件
        ├── index.html
        ├── Build/
        └── TemplateData/
```

## 🚀 部署步骤

### 1. 上传文件到云服务器

将整个 `deployment` 目录上传到云服务器，例如：

```bash
# 在本地执行（假设服务器IP是 123.45.67.89）
scp -r deployment root@123.45.67.89:/opt/car-game
```

或者使用 FTP/SFTP 工具上传。

### 2. 登录云服务器

```bash
ssh root@你的服务器IP
```

### 3. 进入部署目录

```bash
cd /opt/car-game/deployment/server
```

### 4. 启动服务器

```bash
./start.sh
```

启动成功后会看到：

```
==========================================
  多人在线碰撞车游戏服务器
==========================================

🚀 启动游戏服务器...

📡 WebSocket 服务: ws://0.0.0.0:8899/ws
🌐 HTTP 服务: http://0.0.0.0:9988

玩家可以通过以下地址访问游戏:
  http://你的服务器IP:8899

✅ 服务器已启动 (PID: xxxxx)
📝 日志文件: /opt/car-game/deployment/server/server.log
```

### 5. 配置防火墙

确保云服务器防火墙开放端口 **8899**：

**阿里云/腾讯云控制台：**
- 进入安全组设置
- 添加入站规则：TCP 端口 8899，来源 0.0.0.0/0

**使用 iptables：**
```bash
iptables -A INPUT -p tcp --dport 8899 -j ACCEPT
service iptables save
```

**使用 firewalld：**
```bash
firewall-cmd --permanent --add-port=8899/tcp
firewall-cmd --reload
```

### 6. 访问游戏

玩家在浏览器中访问：

```
http://你的服务器公网IP:8899
```

例如：`http://123.45.67.89:8899`

## 🎮 游戏玩法

- **WASD** 或 **方向键**：控制车辆移动
- **目标**：撞击对手，让对手的头部碰到任何物体（地面、墙壁、其他车辆等）
- **爆炸**：当车辆头部碰到物体时会爆炸并在重生点重生
- **策略**：利用翘头、翻滚等技巧攻击对手

## 🔧 服务器管理

### 查看服务器状态

```bash
ps aux | grep game-server
```

### 查看实时日志

```bash
tail -f /opt/car-game/deployment/server/server.log
```

### 停止服务器

```bash
cd /opt/car-game/deployment/server
./stop.sh
```

或者：

```bash
pkill -f game-server
```

### 重启服务器

```bash
./stop.sh
./start.sh
```

## 📊 端口说明

- **8899**：游戏服务端口（HTTP + WebSocket）
  - HTTP: 提供 WebGL 客户端文件
  - WebSocket: 处理游戏实时通信

## ⚠️ 注意事项

1. **服务器要求**：
   - Linux 系统（推荐 Ubuntu 20.04+、CentOS 7+）
   - 至少 1GB 内存
   - 开放端口 8899

2. **浏览器要求**：
   - 现代浏览器（Chrome、Firefox、Safari、Edge）
   - 支持 WebGL 2.0

3. **网络要求**：
   - 稳定的网络连接
   - 低延迟（建议 < 100ms）

## 🐛 故障排查

### 问题1：无法访问游戏

**检查项：**
1. 服务器是否启动：`ps aux | grep game-server`
2. 端口是否开放：`netstat -tlnp | grep 8899`
3. 防火墙是否配置正确
4. 云服务器安全组是否开放 8899 端口

### 问题2：服务器启动失败

**检查项：**
1. 查看日志：`cat /opt/car-game/deployment/server/server.log`
2. 端口是否被占用：`lsof -i :8899`
3. 文件权限：`chmod +x game-server start.sh stop.sh`

### 问题3：游戏连接不上

**检查项：**
1. 浏览器控制台是否有错误
2. WebSocket 连接是否成功
3. 服务器日志是否有错误信息

## 📝 更新部署

如果需要更新游戏：

1. 停止服务器：`./stop.sh`
2. 备份旧版本：`cp -r deployment deployment.backup`
3. 上传新的部署文件覆盖
4. 启动服务器：`./start.sh`

## 🔒 安全建议

1. 使用非 root 用户运行服务器
2. 配置 HTTPS（使用 Nginx 反向代理）
3. 限制访问来源（如果不是公开游戏）
4. 定期备份日志和数据

## 📞 技术支持

如有问题，请检查：
- 服务器日志：`/opt/car-game/deployment/server/server.log`
- 浏览器控制台（F12）

---

**祝游戏愉快！🎮**
