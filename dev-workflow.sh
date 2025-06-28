#!/bin/bash

# PentAGI 开发工作流脚本
echo "🛠️  PentAGI 开发工作流"
echo "===================="

show_menu() {
    echo ""
    echo "选择开发模式："
    echo "1. 🔥 热重载模式 (推荐) - 自动重启"
    echo "2. 🚀 普通开发模式 - 手动重启"
    echo "3. 📦 生产模式 - 编译后运行"
    echo "4. 🛑 停止所有服务"
    echo "5. 📊 查看服务状态"
    echo "6. 🧪 运行API测试"
    echo "0. 退出"
    echo ""
}

start_database() {
    if ! docker ps | grep -q pgvector; then
        echo "🗄️  启动数据库..."
        docker-compose up -d pgvector
        sleep 3
        echo "✅ 数据库已启动"
    else
        echo "✅ 数据库已运行"
    fi
}

hot_reload_mode() {
    echo "🔥 启动热重载模式..."
    start_database
    
    cd backend
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable"
    
    echo "📝 热重载模式说明："
    echo "• 修改 .go 文件后自动重新编译和重启"
    echo "• 无需手动操作，专注于代码开发"
    echo "• 按 Ctrl+C 停止"
    echo ""
    
    if ! command -v air &> /dev/null; then
        echo "安装 Air..."
        go install github.com/air-verse/air@latest
    fi
    
    air
}

dev_mode() {
    echo "🚀 启动普通开发模式..."
    start_database
    
    cd backend
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable"
    
    echo "📝 普通开发模式说明："
    echo "• 使用 go run 启动"
    echo "• 修改代码后需要手动重启 (Ctrl+C 然后重新运行)"
    echo "• 适合偶尔修改的情况"
    echo ""
    
    go run ./cmd/pentagi
}

production_mode() {
    echo "📦 启动生产模式..."
    start_database
    
    cd backend
    export DATABASE_URL="postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable"
    
    echo "🔨 编译中..."
    go build -o pentagi ./cmd/pentagi
    
    echo "📝 生产模式说明："
    echo "• 编译后运行，性能最佳"
    echo "• 修改代码后需要重新编译"
    echo "• 适合稳定版本运行"
    echo ""
    
    ./pentagi
}

stop_services() {
    echo "🛑 停止所有服务..."
    
    # 停止后端进程
    pkill -f "go run.*pentagi" 2>/dev/null || true
    pkill -f "./pentagi" 2>/dev/null || true
    pkill -f "air" 2>/dev/null || true
    
    # 停止前端进程
    pkill -f "vite" 2>/dev/null || true
    pkill -f "npm run dev" 2>/dev/null || true
    
    # 停止隧道
    pkill -f "cloudflared" 2>/dev/null || true
    
    # 停止数据库
    docker-compose stop pgvector 2>/dev/null || true
    
    echo "✅ 所有服务已停止"
}

show_status() {
    echo "📊 服务状态："
    echo ""
    
    # 数据库状态
    if docker ps | grep -q pgvector; then
        echo "🗄️  数据库: ✅ 运行中"
    else
        echo "🗄️  数据库: ❌ 未运行"
    fi
    
    # 后端状态
    if pgrep -f "go run.*pentagi\|./pentagi\|air" >/dev/null; then
        echo "🔧 后端: ✅ 运行中"
        echo "   进程: $(pgrep -f "go run.*pentagi\|./pentagi\|air" | head -1)"
    else
        echo "🔧 后端: ❌ 未运行"
    fi
    
    # 前端状态
    if pgrep -f "vite" >/dev/null; then
        echo "🎨 前端: ✅ 运行中"
    else
        echo "🎨 前端: ❌ 未运行"
    fi
    
    # 隧道状态
    if pgrep -f "cloudflared" >/dev/null; then
        echo "🌐 隧道: ✅ 运行中 ($(pgrep -f "cloudflared" | wc -l) 个)"
    else
        echo "🌐 隧道: ❌ 未运行"
    fi
    
    echo ""
    echo "🔗 访问地址："
    echo "• 本地前端: http://localhost:8000"
    echo "• 本地后端: http://localhost:8080"
    echo "• 互联网前端: https://frost-attempted-midlands-invitations.trycloudflare.com"
    echo "• 互联网后端: https://fitted-platform-compressed-cholesterol.trycloudflare.com"
}

run_tests() {
    echo "🧪 运行API测试..."
    if [ -f "./test-api.sh" ]; then
        ./test-api.sh
    else
        echo "❌ 测试脚本不存在"
    fi
}

# 主循环
while true; do
    show_menu
    read -p "请选择 (0-6): " choice
    
    case $choice in
        1)
            hot_reload_mode
            ;;
        2)
            dev_mode
            ;;
        3)
            production_mode
            ;;
        4)
            stop_services
            ;;
        5)
            show_status
            ;;
        6)
            run_tests
            ;;
        0)
            echo "👋 再见！"
            exit 0
            ;;
        *)
            echo "❌ 无效选择，请重试"
            ;;
    esac
    
    echo ""
    read -p "按 Enter 继续..."
done
