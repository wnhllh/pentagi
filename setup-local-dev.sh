#!/bin/bash

# PentAGI Local Development Setup Script
# 这个脚本可以帮助您在新环境中快速设置开发环境

set -e

echo "🚀 PentAGI 本地开发环境设置脚本"
echo "=================================="

# 检查必要的工具
check_requirements() {
    echo "📋 检查系统要求..."
    
    # 检查Docker
    if ! command -v docker &> /dev/null; then
        echo "❌ Docker 未安装，请先安装 Docker"
        exit 1
    fi
    
    # 检查Docker Compose
    if ! command -v docker-compose &> /dev/null; then
        echo "❌ Docker Compose 未安装，请先安装 Docker Compose"
        exit 1
    fi
    
    # 检查Go
    if ! command -v go &> /dev/null; then
        echo "❌ Go 未安装，请先安装 Go 1.24.0+"
        exit 1
    fi
    
    # 检查Node.js
    if ! command -v node &> /dev/null; then
        echo "❌ Node.js 未安装，请先安装 Node.js 18+"
        exit 1
    fi
    
    echo "✅ 系统要求检查完成"
}

# 设置环境变量
setup_env() {
    echo "⚙️  设置环境变量..."
    
    if [ ! -f ".env" ]; then
        echo "📝 创建 .env 文件..."
        cp .env.example .env
        
        # 修改数据库URL为本地连接
        sed -i 's/DATABASE_URL=.*/DATABASE_URL=postgres:\/\/postgres:postgres@localhost:5432\/pentagidb?sslmode=disable/' .env
        
        echo "✅ .env 文件已创建，请根据需要修改API密钥"
    else
        echo "✅ .env 文件已存在"
    fi
    
    if [ ! -f "frontend/.env" ]; then
        echo "📝 创建前端 .env 文件..."
        cat > frontend/.env << EOF
# Frontend Environment Variables for Local Development
VITE_PORT=8000
VITE_HOST=0.0.0.0
VITE_USE_HTTPS=false
VITE_API_URL=localhost:8080
VITE_APP_NAME=PentAGI
EOF
        echo "✅ 前端 .env 文件已创建"
    else
        echo "✅ 前端 .env 文件已存在"
    fi
}

# 启动数据库
start_database() {
    echo "🗄️  启动数据库..."
    docker-compose up -d pgvector
    
    # 等待数据库启动
    echo "⏳ 等待数据库启动..."
    sleep 5
    
    # 检查数据库是否运行
    if docker-compose ps pgvector | grep -q "Up"; then
        echo "✅ 数据库启动成功"
    else
        echo "❌ 数据库启动失败"
        exit 1
    fi
}

# 构建后端
build_backend() {
    echo "🔨 构建后端..."
    cd backend
    
    echo "📦 下载Go依赖..."
    go mod download
    
    echo "🔨 构建后端应用..."
    go build -o pentagi ./cmd/pentagi
    
    cd ..
    echo "✅ 后端构建完成"
}

# 安装前端依赖
install_frontend() {
    echo "📦 安装前端依赖..."
    cd frontend
    npm install
    cd ..
    echo "✅ 前端依赖安装完成"
}

# 启动服务
start_services() {
    echo "🚀 启动服务..."
    
    # 启动后端
    echo "🔧 启动后端服务..."
    cd backend
    DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable" ./pentagi &
    BACKEND_PID=$!
    cd ..
    
    # 等待后端启动
    echo "⏳ 等待后端启动..."
    sleep 5
    
    # 检查后端是否运行
    if curl -s http://localhost:8080/api/v1/info > /dev/null; then
        echo "✅ 后端服务启动成功"
    else
        echo "❌ 后端服务启动失败"
        kill $BACKEND_PID 2>/dev/null || true
        exit 1
    fi
    
    # 启动前端
    echo "🎨 启动前端服务..."
    cd frontend
    npm run dev &
    FRONTEND_PID=$!
    cd ..
    
    # 等待前端启动
    echo "⏳ 等待前端启动..."
    sleep 5
    
    # 检查前端是否运行
    if curl -s http://localhost:8000 > /dev/null; then
        echo "✅ 前端服务启动成功"
    else
        echo "❌ 前端服务启动失败"
        kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true
        exit 1
    fi
    
    echo ""
    echo "🎉 所有服务启动成功！"
    echo ""
    echo "📱 访问地址："
    echo "   前端应用: http://localhost:8000"
    echo "   后端API:  http://localhost:8080"
    echo "   Swagger:  http://localhost:8080/api/v1/swagger/index.html"
    echo "   GraphQL:  http://localhost:8080/api/v1/graphql/playground"
    echo ""
    echo "🛑 停止服务请按 Ctrl+C"
    
    # 等待用户中断
    trap "echo ''; echo '🛑 停止服务...'; kill $BACKEND_PID $FRONTEND_PID 2>/dev/null || true; docker-compose stop pgvector; echo '✅ 所有服务已停止'; exit 0" INT
    
    wait
}

# 主函数
main() {
    check_requirements
    setup_env
    start_database
    build_backend
    install_frontend
    start_services
}

# 如果直接运行脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
