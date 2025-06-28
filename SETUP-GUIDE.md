# PentAGI 本地开发环境设置指南

## 🚀 快速开始

### 前提条件

确保您的系统已安装以下软件：

- **Docker** (20.10+) 和 **Docker Compose** (2.0+)
- **Go** (1.24.0+)
- **Node.js** (18+) 和 **npm**
- **Git**

### 一键设置

1. **克隆项目**（如果还没有）：
   ```bash
   git clone <your-repo-url>
   cd pentagi
   ```

2. **运行设置脚本**：
   ```bash
   ./setup-local-dev.sh
   ```

这个脚本会自动：
- 检查系统要求
- 创建环境配置文件
- 启动数据库容器
- 构建后端应用
- 安装前端依赖
- 启动所有服务

## 📱 访问地址

服务启动后，您可以通过以下地址访问：

- **前端应用**: http://localhost:8000
- **后端API**: http://localhost:8080
- **API文档 (Swagger)**: http://localhost:8080/api/v1/swagger/index.html
- **GraphQL Playground**: http://localhost:8080/api/v1/graphql/playground

## 🧪 测试API

运行API测试脚本来验证所有服务是否正常工作：

```bash
./test-api.sh
```

## 🔧 手动设置（如果自动脚本失败）

### 1. 环境配置

创建 `.env` 文件：
```bash
cp .env.example .env
```

编辑 `.env` 文件，确保数据库URL正确：
```
DATABASE_URL=postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable
```

创建前端环境文件 `frontend/.env`：
```
VITE_PORT=8000
VITE_HOST=0.0.0.0
VITE_USE_HTTPS=false
VITE_API_URL=localhost:8080
VITE_APP_NAME=PentAGI
```

### 2. 启动数据库

```bash
docker-compose up -d pgvector
```

### 3. 构建和启动后端

```bash
cd backend
go mod download
go build -o pentagi ./cmd/pentagi
DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable" ./pentagi
```

### 4. 启动前端

在新终端中：
```bash
cd frontend
npm install
npm run dev
```

## 🌐 外部访问配置

### 当前配置状态

✅ **已配置外部访问支持**
- 前端绑定到: `0.0.0.0:8000`
- 后端绑定到: `*:8080`
- CORS已配置支持多个IP地址

### 验证外部访问

运行验证脚本检查配置：
```bash
./verify-external-setup.sh
```

### 外部访问地址

如果服务运行在IP `172.17.0.2`，您可以通过以下地址访问：

- **前端应用**: http://172.17.0.2:8000
- **后端API**: http://172.17.0.2:8080
- **API文档**: http://172.17.0.2:8080/api/v1/swagger/index.html
- **GraphQL**: http://172.17.0.2:8080/api/v1/graphql/playground

### 访问方法

#### 方法1: 直接IP访问
如果在同一网络中，直接使用服务器IP访问

#### 方法2: 端口转发（推荐）
如果使用VS Code Remote、GitHub Codespaces等：
1. 转发端口 8000 (前端)
2. 转发端口 8080 (后端)

#### 方法3: SSH隧道
```bash
ssh -L 8000:localhost:8000 -L 8080:localhost:8080 user@server
```

#### 方法4: 反向代理
使用nginx等配置反向代理

### 故障排除

如果无法外部访问：
1. 检查防火墙设置
2. 确认云服务商安全组配置
3. 验证网络路由
4. 使用 `./test-external-access.sh` 诊断

## 🔄 环境迁移

### 保存当前环境

1. **提交代码更改**：
   ```bash
   git add .
   git commit -m "本地开发环境配置"
   git push
   ```

2. **备份重要文件**：
   - `.env` - 环境配置
   - `frontend/.env` - 前端配置
   - `setup-local-dev.sh` - 设置脚本
   - `test-api.sh` - 测试脚本

### 在新机器上恢复

1. **克隆项目**：
   ```bash
   git clone <your-repo-url>
   cd pentagi
   ```

2. **运行设置脚本**：
   ```bash
   ./setup-local-dev.sh
   ```

3. **配置API密钥**（如需要）：
   编辑 `.env` 文件添加您的API密钥：
   ```
   OPEN_AI_KEY=your_openai_key
   ANTHROPIC_API_KEY=your_anthropic_key
   ```

## 🛠️ 常用命令

### 启动服务
```bash
# 只启动数据库
docker-compose up -d pgvector

# 启动后端
cd backend && DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable" ./pentagi

# 启动前端
cd frontend && npm run dev
```

### 停止服务
```bash
# 停止数据库
docker-compose stop pgvector

# 停止后端和前端
# 在各自的终端中按 Ctrl+C
```

### 查看日志
```bash
# 数据库日志
docker-compose logs pgvector

# 后端日志在运行终端中查看
# 前端日志在运行终端中查看
```

## 🐛 故障排除

### 数据库连接问题
```bash
# 检查数据库是否运行
docker-compose ps pgvector

# 重启数据库
docker-compose restart pgvector
```

### 端口占用问题
```bash
# 检查端口占用
lsof -i :8000  # 前端端口
lsof -i :8080  # 后端端口
lsof -i :5432  # 数据库端口
```

### 依赖问题
```bash
# 重新安装Go依赖
cd backend && go mod download

# 重新安装Node.js依赖
cd frontend && rm -rf node_modules && npm install
```

## 📞 获取帮助

如果遇到问题：

1. 运行测试脚本检查状态：`./test-api.sh`
2. 检查服务日志
3. 确认所有前提条件已满足
4. 查看本指南的故障排除部分
