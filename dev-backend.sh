#!/bin/bash

# 开发模式后端启动脚本
echo "🚀 启动后端开发模式"
echo "=================="

# 加载环境变量
if [ -f "../.env" ]; then
    echo "📋 加载环境变量..."
    export $(grep -v '^#' ../.env | xargs)
    echo "✅ 环境变量已加载"
fi

# 设置工作目录
cd backend

# 检查数据库是否运行
if ! docker ps | grep -q pgvector; then
    echo "⚠️  数据库未运行，正在启动..."
    cd ..
    docker-compose up -d pgvector
    cd backend
    echo "✅ 数据库已启动"
fi

# 确保数据库URL正确
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable"

echo "📝 开发模式特点："
echo "• 修改代码后自动重新编译"
echo "• 无需手动构建二进制文件"
echo "• 适合频繁修改和测试"
echo ""
echo "🔄 重启方法：按 Ctrl+C 停止，然后重新运行此脚本"
echo ""
echo "🌐 服务地址："
echo "• 本地API: http://localhost:8080"
echo "• 隧道API: https://fitted-platform-compressed-cholesterol.trycloudflare.com"
echo ""
echo "🚀 启动中..."

# 使用 go run 启动
go run ./cmd/pentagi
