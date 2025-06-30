# PentAGI 电力行业安全测试系统 - 重启指南和故障排除

## 🚀 正确的重启顺序

### 1. 停止所有服务
```bash
# 停止前端
cd frontend
# 按 Ctrl+C 停止 npm run dev

# 停止后端
cd ../backend
# 按 Ctrl+C 停止 go run cmd/pentagi/main.go

# 停止cloudflare隧道
# 按 Ctrl+C 停止 cloudflared tunnel

# 停止电力IT实验室（可选）
cd ../power-it-lab
docker-compose down
```

### 2. 启动服务（正确顺序）

#### 步骤1: 启动数据库和电力IT实验室
```bash
# 启动PostgreSQL数据库
docker-compose up -d pgvector

# 启动电力IT实验室
cd power-it-lab
docker-compose up -d
cd ..
```

#### 步骤2: 启动后端
```bash
# 进入backend目录
cd backend

# 加载环境变量
source .env

# 启动后端
go run cmd/pentagi/main.go
```

#### 步骤3: 启动前端
```bash
# 新终端，进入frontend目录
cd frontend

# 启动前端
npm run dev
```

#### 步骤4: 创建cloudflare隧道
```bash
# 新终端，在根目录
cloudflared tunnel --url http://localhost:8000
```

#### 步骤5: 连接Docker网络（如果需要）
```bash
# 连接PentAGI容器到电力网络
docker network connect power-it-lab_power-network pentagi-terminal-1
docker network connect power-it-lab_power-network pentagi-terminal-2
```

## 🌐 WebSocket和实时通信配置

### 智能连接模式

系统实现了智能WebSocket配置，根据部署环境自动选择最佳连接方式：

#### 本地开发环境
- **检测条件**: `hostname` 不包含 `trycloudflare.com`
- **连接方式**: WebSocket (ws:// 或 wss://)
- **实时性**: 真正的实时推送，延迟 < 100ms
- **适用场景**: 开发调试、本地测试

#### Cloudflare隧道环境
- **检测条件**: `hostname` 包含 `trycloudflare.com`
- **连接方式**: HTTP轮询 (每3秒)
- **实时性**: 准实时推送，延迟 ≤ 3秒
- **适用场景**: 演示部署、远程访问

### 配置文件位置
- **Apollo配置**: `frontend/src/lib/apollo.ts`
- **环境检测**: 自动检测，无需手动配置
- **轮询间隔**: 3000ms (可在代码中调整)

### 故障排除
如果消息推送不工作：
1. 检查浏览器控制台是否有WebSocket错误
2. 确认当前环境模式（控制台会显示日志）
3. 在Cloudflare环境下，等待最多3秒查看新消息
4. 如果仍有问题，手动刷新页面

## 🔧 常见故障排除

### 问题1: 看不到Provider选择

**症状**: 登录后Provider下拉框为空或显示错误

**原因**: 
- 环境变量未正确加载
- Anthropic API密钥配置错误

**解决方案**:
```bash
# 1. 检查环境变量
cd backend
source .env
echo $ANTHROPIC_API_KEY

# 2. 确保.env文件存在且包含正确的API密钥
cat .env | grep ANTHROPIC_API_KEY

# 3. 重启后端
# 按 Ctrl+C 停止后端
source .env
go run cmd/pentagi/main.go
```

### 问题2: WebSocket连接失败

**症状**: 控制台显示WebSocket连接错误，消息推送不工作

**原因**:
- Cloudflare免费隧道不支持WebSocket连接
- 本地开发环境可以使用WebSocket

**解决方案**:
系统已实现智能WebSocket配置：
- **本地环境** (localhost): 自动使用WebSocket实时连接
- **Cloudflare环境** (*.trycloudflare.com): 自动使用HTTP轮询（3秒间隔）
- 系统会自动检测环境并选择最佳连接方式

**重要注意事项**:
⚠️ **消息推送机制**:
- 本地开发: 真正的WebSocket实时推送
- Cloudflare部署: HTTP轮询模拟实时推送（3秒延迟）
- 用户无需手动刷新，新消息会自动显示

### 问题3: 工具执行权限错误

**症状**: nmap等工具显示"Operation not permitted"

**解决方案**:
```bash
# 自动修复nmap权限问题
docker exec pentagi-terminal-1 bash -c "
cat > /usr/bin/nmap << 'EOF'
#!/usr/bin/env sh
set -e
if [ ! -f /tmp/nmap ] || [ /usr/lib/nmap/nmap -nt /tmp/nmap ]; then
  cp /usr/lib/nmap/nmap /tmp/nmap
  chmod +x /tmp/nmap
fi
if [ \"\$(id -u)\" -eq 0 ] || [ \"\$1\" = \"--resume\" ]; then
  exec /tmp/nmap \"\$@\"
else
  exec /tmp/nmap --privileged \"\$@\"
fi
EOF
chmod +x /usr/bin/nmap"

# 对第二个容器执行相同操作
docker exec pentagi-terminal-2 bash -c "
cat > /usr/bin/nmap << 'EOF'
#!/usr/bin/env sh
set -e
if [ ! -f /tmp/nmap ] || [ /usr/lib/nmap/nmap -nt /tmp/nmap ]; then
  cp /usr/lib/nmap/nmap /tmp/nmap
  chmod +x /tmp/nmap
fi
if [ \"\$(id -u)\" -eq 0 ] || [ \"\$1\" = \"--resume\" ]; then
  exec /tmp/nmap \"\$@\"
else
  exec /tmp/nmap --privileged \"\$@\"
fi
EOF
chmod +x /usr/bin/nmap"
```

### 问题4: 端口未开放/无法连接目标系统

**症状**: curl连接失败，nmap显示filtered

**解决方案**:
```bash
# 1. 检查电力IT实验室是否运行
cd power-it-lab
docker-compose ps

# 2. 如果容器停止，重新启动
docker-compose up -d

# 3. 检查网络连接
docker network ls
docker network inspect power-it-lab_power-network

# 4. 连接PentAGI容器到电力网络
docker network connect power-it-lab_power-network pentagi-terminal-1
docker network connect power-it-lab_power-network pentagi-terminal-2

# 5. 获取最新IP地址
docker network inspect power-it-lab_power-network | jq '.[0].Containers | to_entries | map({name: .value.Name, ip: .value.IPv4Address})'
```

### 问题5: 524超时错误

**症状**: GraphQL请求返回524错误

**解决方案**: 
- 已增加超时时间到30秒
- 减少轮询频率到5秒
- 如果仍有问题，检查后端日志

## 📋 当前系统配置

### 网络地址
- **电力营销系统2.0**: `http://172.25.0.10:8080`
- **i国网APP**: `http://172.25.0.9:8080`
- **ERP系统**: `http://172.25.0.8:8080`

### 端口配置
- **后端**: localhost:8080
- **前端**: localhost:8001 (如果8000被占用会自动切换)
- **数据库**: localhost:5432
- **Cloudflare隧道**: 动态生成

### WebSocket配置
- **本地环境**: WebSocket实时连接 (ws://localhost:8001/api/v1/graphql)
- **Cloudflare环境**: HTTP轮询 (3秒间隔)
- **自动检测**: 基于hostname自动选择连接方式
- **日志标识**:
  - 🔌 本地环境: "运行在本地环境，使用WebSocket实时模式"
  - 🌐 Cloudflare环境: "运行在Cloudflare隧道环境，使用HTTP轮询模式"

### 重要文件
- **环境变量**: `.env` (根目录) 和 `backend/.env`
- **前端配置**: `frontend/vite.config.ts`
- **Apollo配置**: `frontend/src/lib/apollo.ts` (包含智能WebSocket配置)
- **Docker配置**: `docker-compose.yml`, `power-it-lab/docker-compose.yml`

## 🔍 健康检查命令

```bash
# 检查所有服务状态
echo "=== 数据库 ==="
docker ps | grep pgvector

echo "=== 电力IT实验室 ==="
cd power-it-lab && docker-compose ps && cd ..

echo "=== 后端API ==="
curl -s http://localhost:8080/api/v1/info | jq .status

echo "=== 前端代理 ==="
curl -s http://localhost:8000/api/v1/info | jq .status

echo "=== 网络连接 ==="
docker exec pentagi-terminal-2 curl -I http://172.25.0.10:8080

echo "=== 工具权限 ==="
docker exec pentagi-terminal-2 nmap --version
```

## 🚨 紧急恢复

如果系统完全无法启动，执行完整重置：

```bash
# 1. 停止所有容器
docker stop $(docker ps -q)

# 2. 重新启动必要服务
docker-compose up -d pgvector
cd power-it-lab && docker-compose up -d && cd ..

# 3. 连接网络
docker network connect power-it-lab_power-network pentagi-terminal-1
docker network connect power-it-lab_power-network pentagi-terminal-2

# 4. 修复工具权限
# (执行上面的nmap修复脚本)

# 5. 按正确顺序启动服务
# (按照重启顺序执行)
```

## 🛠️ 高级故障排除

### 问题6: 聊天记录不显示或不更新

**症状**: 登录后看不到历史聊天记录，或新消息不显示

**原因**:
- 数据库连接问题
- GraphQL查询失败
- Apollo缓存问题

**解决方案**:
```bash
# 1. 检查数据库连接
docker exec pgvector psql -U postgres -d pentagidb -c "SELECT COUNT(*) FROM flows;"

# 2. 清除浏览器缓存和localStorage
# 在浏览器开发者工具中执行:
# localStorage.clear(); sessionStorage.clear(); location.reload();

# 3. 检查后端日志中的GraphQL错误
# 查看后端终端输出

# 4. 重启后端和前端
```

### 问题7: 前端编译错误

**症状**: npm run dev失败或显示编译错误

**解决方案**:
```bash
cd frontend

# 1. 清除node_modules和重新安装
rm -rf node_modules package-lock.json
npm install

# 2. 清除Vite缓存
rm -rf .vite

# 3. 检查TypeScript错误
npm run type-check

# 4. 重新启动
npm run dev
```

### 问题8: Docker容器网络问题

**症状**: 容器之间无法通信

**解决方案**:
```bash
# 1. 检查所有网络
docker network ls

# 2. 检查容器网络配置
docker inspect pentagi-terminal-1 | jq '.[0].NetworkSettings.Networks'

# 3. 重新创建网络连接
docker network disconnect power-it-lab_power-network pentagi-terminal-1 2>/dev/null || true
docker network disconnect power-it-lab_power-network pentagi-terminal-2 2>/dev/null || true
docker network connect power-it-lab_power-network pentagi-terminal-1
docker network connect power-it-lab_power-network pentagi-terminal-2

# 4. 验证连接
docker exec pentagi-terminal-1 ping -c 1 172.25.0.10
```

### 问题9: Cloudflare隧道连接问题

**症状**: 无法通过cloudflare地址访问，或连接不稳定

**解决方案**:
```bash
# 1. 重新创建隧道
# 按 Ctrl+C 停止当前隧道
cloudflared tunnel --url http://localhost:8000

# 2. 检查本地服务是否正常
curl http://localhost:8000

# 3. 更新前端allowedHosts配置
# 编辑 frontend/vite.config.ts，添加新的隧道地址

# 4. 如果仍有问题，尝试不同的端口
cloudflared tunnel --url http://localhost:8001
```

### 问题10: 权限和安全问题

**症状**: 各种权限拒绝错误

**解决方案**:
```bash
# 1. 检查文件权限
ls -la .env backend/.env

# 2. 修复权限
chmod 600 .env backend/.env
chmod +x dev-backend.sh

# 3. 检查Docker权限
docker info | grep -i error

# 4. 重启Docker服务（如果需要）
sudo systemctl restart docker
```

## 📊 监控和日志

### 查看日志
```bash
# 后端日志
# 直接在后端终端查看

# 前端日志
# 在浏览器开发者工具的Console中查看

# Docker容器日志
docker logs pentagi-terminal-1
docker logs power-marketing-gateway

# 数据库日志
docker logs pgvector
```

### 性能监控
```bash
# 检查系统资源
docker stats

# 检查端口占用
netstat -tulpn | grep :8080
netstat -tulpn | grep :8000

# 检查磁盘空间
df -h
docker system df
```

## 🔄 定期维护

### 每日检查
```bash
# 运行健康检查
./health-check.sh

# 清理Docker
docker system prune -f

# 检查日志大小
du -sh /var/lib/docker/containers/*/
```

### 每周维护
```bash
# 更新Docker镜像
docker-compose pull
cd power-it-lab && docker-compose pull && cd ..

# 备份数据库
docker exec pgvector pg_dump -U postgres pentagidb > backup_$(date +%Y%m%d).sql

# 清理未使用的资源
docker system prune -a -f
```
