# 🚀 PentAGI + 电力IT实验室 完整服务端口大全

## 📋 当前运行状态

### ✅ PentAGI 主系统

| 服务 | 本地端口 | Cloudflare隧道 | 状态 | 说明 |
|------|----------|----------------|------|------|
| **PentAGI后端** | `localhost:8080` | `https://berry-dance-products-label.trycloudflare.com` | ✅ 运行中 | GraphQL API, WebSocket |
| **PentAGI前端** | `localhost:8001` | `https://certification-jet-suddenly-scoop.trycloudflare.com` | ✅ 运行中 | React + Vite开发服务器 |

### ✅ 电力IT实验室（已修复端口冲突）

| 系统 | 本地端口 | 外部访问 | 状态 | 说明 |
|------|----------|----------|------|------|
| **电力营销系统2.0** | `localhost:18080` | 仅内网 | ✅ 运行中 | 微服务架构，计费系统 |
| **i国网APP后端** | `localhost:9080` | 仅内网 | ✅ 运行中 | 移动应用API，SMS验证 |
| **ERP系统(SAP风格)** | `localhost:18000` | 仅内网 | ✅ 运行中 | 企业资源规划系统 |

## 🔧 正确启动方式

### 1. PentAGI系统启动

#### 方法一：使用开发脚本（推荐）

```bash
# 1. 启动后端（在项目根目录）
bash dev-backend-with-env.sh

# 2. 启动前端（新终端）
cd frontend && npm run dev

# 3. 创建前端Cloudflare隧道（新终端）
cloudflared tunnel --url http://localhost:8001

# 4. 创建后端Cloudflare隧道（新终端）
cloudflared tunnel --url http://localhost:8080
```

#### 方法二：使用Docker Compose

```bash
# 启动完整环境
docker-compose up -d

# 创建Cloudflare隧道
cloudflared tunnel --url http://localhost:8443
```

### 2. 电力IT实验室启动

```bash
# 进入实验室目录
cd power-it-lab

# 启动所有电力系统
docker-compose up -d

# 或者启动特定系统
docker-compose up -d marketing-gateway iguowang-api erp-app-server
```

### 3. 环境变量配置

#### PentAGI后端环境变量（.env文件）
```bash
ANTHROPIC_API_KEY=sk-ant-xxx  # 必需
DATABASE_URL=postgres://pentagi:password@localhost:5432/pentagi
OPENAI_API_KEY=sk-xxx  # 可选
```

#### PentAGI前端环境变量（frontend/.env文件）
```bash
VITE_API_URL=berry-dance-products-label.trycloudflare.com  # 后端隧道地址
VITE_PORT=8000
VITE_HOST=0.0.0.0
```

## 🎯 访问地址

### 🌐 外部访问（通过Cloudflare隧道）

- **PentAGI前端**: https://certification-jet-suddenly-scoop.trycloudflare.com
- **PentAGI后端API**: https://berry-dance-products-label.trycloudflare.com

### 🏠 本地访问

#### PentAGI系统
- **前端开发服务器**: http://localhost:8001
- **后端API**: http://localhost:8080
- **GraphQL Playground**: http://localhost:8080/graphql

#### 电力IT实验室
- **电力营销系统2.0**: http://localhost:18080
  - 系统信息: http://localhost:18080/api/system/info
  - 计费API: http://localhost:18080/api/billing
  - 用户API: http://localhost:18080/api/users
  
- **i国网APP**: http://localhost:9080
  - 健康检查: http://localhost:9080/api/health
  - SMS验证: http://localhost:9080/api/auth/send-sms
  - 用户列表: http://localhost:9080/api/user/list
  
- **ERP系统**: http://localhost:18000
  - 健康检查: http://localhost:18000/health
  - 登录API: http://localhost:18000/api/auth/login
  - 系统配置: http://localhost:18000/api/system/config

## 🔍 故障排除

### 常见问题和解决方案

#### 1. 端口冲突
```bash
# 检查端口占用
netstat -tlnp | grep :8080
lsof -i :8080

# 停止冲突服务
docker-compose down
pkill -f "port 8080"
```

#### 2. 前端无法连接后端
```bash
# 检查后端状态
curl http://localhost:8080/api/v1/info

# 检查前端环境变量
cat frontend/.env

# 重启前端服务
cd frontend && npm run dev
```

#### 3. Cloudflare隧道问题
```bash
# 重新创建隧道
pkill cloudflared
cloudflared tunnel --url http://localhost:8001
```

#### 4. 数据库连接问题
```bash
# 检查数据库状态
docker ps | grep postgres

# 重启数据库
docker-compose restart postgres
```

## 🧪 电力行业测试验证

### 快速验证指令

登录PentAGI前端后，使用以下指令验证电力行业增强功能：

#### 1. 电力营销系统测试
```
请使用power_pentester工具测试电力营销系统2.0，目标是http://localhost:18080，进行快速安全评估。
```

#### 2. i国网APP测试
```
请使用api_tester工具测试i国网APP，基础URL是http://localhost:9080，重点测试SMS验证和移动API安全。
```

#### 3. ERP系统合规性测试
```
请使用compliance_agent工具评估SAP ERP系统的合规性，目标是http://localhost:18000，评估ISO27001合规性。
```

#### 4. 计费逻辑测试
```
请使用test_billing_logic工具测试计费逻辑，端点是http://localhost:18080/api/billing，检测价格操纵漏洞。
```

### 预期测试结果

如果电力行业增强功能正常工作，您应该看到：

✅ **智能系统识别**: 自动识别营销2.0/i国网/SAP系统类型  
✅ **专业化分析**: 使用电力行业特定的安全测试方法  
✅ **业务导向报告**: 重点关注计费准确性、客户数据保护、监管合规  
✅ **行业术语**: 大量使用"分层定价"、"用电模式"、"计费逻辑"等专业术语  
✅ **合规性评估**: 自动关联NERC CIP、FERC等电力行业标准  
✅ **风险量化**: 提供具体的财务损失估算和监管风险评估  

## 📝 登录信息

### PentAGI系统登录
- **用户名**: admin@pentagi.com
- **密码**: admin

### 电力系统测试账户
- **营销系统**: 通过API直接测试，无需登录
- **i国网APP**: 通过SMS验证测试
- **ERP系统**: 
  - 默认账户: SAP* / 06071992 / 客户端000
  - 管理员: ADMIN / ADMIN123 / 客户端100

## 🚨 重要提醒

1. **端口冲突**: 确保PentAGI和电力实验室使用不同端口
2. **环境变量**: 前端必须配置正确的后端API地址
3. **Cloudflare隧道**: 前后端都需要独立的隧道
4. **数据库**: 后端需要PostgreSQL数据库连接
5. **API密钥**: 必须配置有效的Anthropic API密钥

## 📞 当前运行状态检查

```bash
# 检查所有服务状态
curl -s http://localhost:8080/api/v1/info | head -3  # PentAGI后端
curl -s http://localhost:18080/api/system/info | head -3  # 营销系统
curl -s http://localhost:9080/api/health  # i国网APP
curl -s http://localhost:18000/health  # ERP系统
```

---

**最后更新**: 2025-06-29 16:03 UTC  
**状态**: 所有服务正常运行，端口冲突已解决，前后端连接已修复
