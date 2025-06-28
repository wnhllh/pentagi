#!/bin/bash

# 开发模式后端启动脚本（确保环境变量正确加载）
echo "🔥 启动后端热重载开发模式（加载环境变量）"
echo "============================================"

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

# 检查.env文件是否存在
if [ ! -f "../.env" ]; then
    echo "❌ .env 文件不存在，请先创建 .env 文件"
    exit 1
fi

echo "📋 加载环境变量..."

# 从.env文件中读取环境变量
set -a  # 自动导出所有变量
source ../.env
set +a  # 停止自动导出

# 显示关键环境变量状态
echo "🔍 环境变量检查："
echo "DATABASE_URL: ${DATABASE_URL:-未设置}"
echo "ANTHROPIC_API_KEY: $([ -n "$ANTHROPIC_API_KEY" ] && echo "已设置 (${#ANTHROPIC_API_KEY} 字符)" || echo "未设置")"
echo "OPEN_AI_KEY: $([ -n "$OPEN_AI_KEY" ] && echo "已设置 (${#OPEN_AI_KEY} 字符)" || echo "未设置")"

echo ""
echo "🔥 热重载模式特点："
echo "• 监控文件变化，自动重新编译和重启"
echo "• 修改代码后无需手动重启"
echo "• 环境变量已正确加载"
echo ""
echo "📁 监控的文件类型："
echo "• .go 文件"
echo "• .tpl, .tmpl, .html 模板文件"
echo ""
echo "🔄 使用方法："
echo "• 修改代码后自动重启"
echo "• 修改.env文件后需要手动重启此脚本"
echo "• 按 Ctrl+C 完全停止"
echo ""
echo "🌐 服务地址："
echo "• 本地API: http://localhost:8080"
echo "• 隧道API: https://fitted-platform-compressed-cholesterol.trycloudflare.com"
echo ""
echo "🚀 启动热重载..."

# 使用 Air 启动热重载
air
