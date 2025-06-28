#!/bin/bash

# 热重载开发模式后端启动脚本
echo "🔥 启动后端热重载开发模式"
echo "========================="

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

# 检查Air是否安装
if ! command -v air &> /dev/null; then
    echo "⚠️  Air未安装，正在安装..."
    go install github.com/air-verse/air@latest
    echo "✅ Air已安装"
fi

# 设置环境变量
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable"

echo "🔥 热重载模式特点："
echo "• 监控文件变化，自动重新编译和重启"
echo "• 修改代码后无需手动重启"
echo "• 最适合频繁开发和调试"
echo ""
echo "📁 监控的文件类型："
echo "• .go 文件"
echo "• .tpl, .tmpl, .html 模板文件"
echo ""
echo "🔄 使用方法："
echo "• 修改代码后自动重启"
echo "• 按 Ctrl+C 完全停止"
echo ""
echo "🌐 服务地址："
echo "• 本地API: http://localhost:8080"
echo "• 隧道API: https://fitted-platform-compressed-cholesterol.trycloudflare.com"
echo ""
echo "🚀 启动热重载..."

# 使用 Air 启动热重载
air
