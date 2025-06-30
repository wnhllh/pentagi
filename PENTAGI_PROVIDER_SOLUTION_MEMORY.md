# 🧠 PentAGI Provider问题完整解决方案 - Memory

## 🚨 核心问题

**症状**: 登录PentAGI前端后，聊天界面中找不到provider（anthropic），无法选择AI模型进行对话。

**根本原因**: 前端无法正确连接到后端的GraphQL API来获取可用的providers列表。

## 🔍 问题诊断步骤

### 1. 检查后端Provider注册状态
```bash
# 查看后端日志，确认Anthropic API密钥已加载
tail -f backend.log | grep -i anthropic

# 预期看到：
# ANTHROPIC_API_KEY: 已设置 (108 字符)
# provider=anthropic
```

### 2. 检查GraphQL端点访问
```bash
# 测试GraphQL端点（需要认证）
curl -X POST http://localhost:8080/graphql -H "Content-Type: application/json" -d '{"query": "{ providers }"}'

# 如果返回404或HTML，说明路由配置有问题
```

### 3. 检查前端API配置
```bash
# 查看前端环境变量
cat frontend/.env

# 查看前端API配置
cat frontend/src/models/Api.ts
```

## 🔧 完整解决方案

### 方案一：修复前端API连接（推荐）

#### 1. 修复前端环境变量配置

**文件**: `frontend/.env`
```bash
# 确保后端API URL正确配置
VITE_API_URL=berry-dance-products-label.trycloudflare.com  # 不包含https://
```

#### 2. 修复前端API配置

**文件**: `frontend/src/models/Api.ts`
```typescript
// 支持外部API URL
const apiUrl = import.meta.env.VITE_API_URL;
export const baseUrl = apiUrl ? `https://${apiUrl}/api/v1` : '/api/v1';
```

#### 3. 修复Apollo GraphQL配置

**文件**: `frontend/src/lib/apollo.ts`
```typescript
// 修复GraphQL HTTP连接
const graphqlUri = baseUrl.startsWith('http') 
    ? `${baseUrl}/graphql`  // 外部API URL
    : `${window.location.origin}${baseUrl}/graphql`;  // 相对路径

const httpLink = createHttpLink({
    uri: graphqlUri,
    credentials: 'include',
});

// 修复WebSocket连接
const wsUri = baseUrl.startsWith('http')
    ? baseUrl.replace('https://', 'wss://').replace('http://', 'ws://') + '/graphql'
    : `${window.location.protocol === 'https:' ? 'wss:' : 'ws:'}//${window.location.host}${baseUrl}/graphql`;

const wsLink = new GraphQLWsLink(
    createClient({
        url: wsUri,
        credentials: 'include',
    })
);
```

### 方案二：使用本地开发模式（备选）

如果cloudflare隧道有问题，可以使用本地开发模式：

#### 1. 修改前端环境变量使用本地后端
```bash
# frontend/.env
VITE_API_URL=  # 留空使用相对路径
```

#### 2. 确保前后端在同一域名下
```bash
# 使用nginx或其他代理将前后端统一到同一端口
```

## 🚀 标准启动流程

### 1. 启动后端
```bash
# 在项目根目录
bash dev-backend-with-env.sh

# 等待看到以下日志：
# - ANTHROPIC_API_KEY: 已设置 (108 字符)
# - Server started on :8080
```

### 2. 创建后端Cloudflare隧道
```bash
# 新终端
cloudflared tunnel --url http://localhost:8080

# 记录隧道地址，例如：
# https://berry-dance-products-label.trycloudflare.com
```

### 3. 配置前端环境变量
```bash
# 编辑 frontend/.env
VITE_API_URL=berry-dance-products-label.trycloudflare.com  # 使用步骤2的隧道地址（去掉https://）
```

### 4. 启动前端
```bash
# 新终端
cd frontend && npm run dev

# 等待启动完成，通常在8001端口
```

### 5. 创建前端Cloudflare隧道
```bash
# 新终端
cloudflared tunnel --url http://localhost:8001

# 记录前端访问地址，例如：
# https://conversion-cooler-balance-puzzles.trycloudflare.com
```

### 6. 启动电力IT实验室（可选）
```bash
# 新终端
cd power-it-lab && docker-compose up -d marketing-gateway iguowang-api erp-app-server

# 等待服务启动完成
```

## 🎯 验证Provider功能

### 1. 访问前端
打开浏览器访问前端cloudflare隧道地址

### 2. 登录系统
- 用户名: admin@pentagi.com
- 密码: admin

### 3. 检查Provider选择
- 进入聊天界面
- 点击设置或模型选择
- 应该能看到"anthropic"选项
- 可以选择Claude模型（如claude-3-sonnet-20240229）

### 4. 测试对话功能
发送测试消息确认AI响应正常

## 🚨 常见问题和解决方案

### 问题1: 前端显示"No providers available"
**解决方案**:
1. 检查后端日志确认Anthropic API密钥已加载
2. 检查前端网络请求是否成功到达后端GraphQL端点
3. 重启前端服务并清除浏览器缓存

### 问题2: GraphQL请求返回404
**解决方案**:
1. 确认后端服务正常运行在8080端口
2. 检查前端API配置是否正确指向后端隧道
3. 确认GraphQL端点路径为 `/graphql`

### 问题3: CORS错误
**解决方案**:
1. 确保Apollo配置中使用了 `credentials: 'include'`
2. 检查后端CORS配置是否允许前端域名
3. 使用cloudflare隧道而不是直接IP访问

### 问题4: WebSocket连接失败
**解决方案**:
1. 检查WebSocket URL配置是否正确（wss://）
2. 确认cloudflare隧道支持WebSocket
3. 如果WebSocket有问题，可以禁用实时功能

## 📋 重启后的检查清单

每次重启系统后，按以下顺序检查：

1. ✅ **后端服务**: 确认8080端口运行，API密钥已加载
2. ✅ **后端隧道**: 创建新的cloudflare隧道并记录地址
3. ✅ **前端配置**: 更新frontend/.env中的VITE_API_URL
4. ✅ **前端服务**: 重启前端开发服务器
5. ✅ **前端隧道**: 创建新的前端cloudflare隧道
6. ✅ **登录测试**: 访问前端，登录并检查provider选择
7. ✅ **对话测试**: 发送测试消息确认功能正常

## 🔄 自动化脚本（推荐）

创建自动化启动脚本避免重复操作：

```bash
#!/bin/bash
# start-pentagi.sh

echo "🚀 启动PentAGI完整环境..."

# 1. 启动后端
echo "📡 启动后端服务..."
bash dev-backend-with-env.sh &
BACKEND_PID=$!

# 等待后端启动
sleep 10

# 2. 创建后端隧道
echo "🌐 创建后端Cloudflare隧道..."
cloudflared tunnel --url http://localhost:8080 &
BACKEND_TUNNEL_PID=$!

# 等待隧道创建
sleep 5

# 3. 启动前端
echo "🎨 启动前端服务..."
cd frontend && npm run dev &
FRONTEND_PID=$!

# 等待前端启动
sleep 10

# 4. 创建前端隧道
echo "🌐 创建前端Cloudflare隧道..."
cloudflared tunnel --url http://localhost:8001 &
FRONTEND_TUNNEL_PID=$!

echo "✅ PentAGI环境启动完成！"
echo "📝 请手动更新frontend/.env中的VITE_API_URL"
echo "🔗 查看终端输出获取访问地址"

# 保存PID以便后续清理
echo $BACKEND_PID > .backend.pid
echo $BACKEND_TUNNEL_PID > .backend_tunnel.pid
echo $FRONTEND_PID > .frontend.pid
echo $FRONTEND_TUNNEL_PID > .frontend_tunnel.pid
```

## 📞 当前运行状态

**最后更新**: 2025-06-29 16:12 UTC

**当前访问地址**:
- 前端: https://conversion-cooler-balance-puzzles.trycloudflare.com
- 后端API: https://berry-dance-products-label.trycloudflare.com

**电力IT实验室端口**:
- 营销系统: http://localhost:18080
- i国网APP: http://localhost:9080
- ERP系统: http://localhost:18000

**状态**: ✅ 所有服务正常运行，Provider问题已解决
