#!/bin/bash

# Cloudflare Tunnel 设置脚本（免费方案）
echo "☁️  Cloudflare Tunnel 设置"
echo "========================"

# 安装 cloudflared
install_cloudflared() {
    echo "📦 安装 cloudflared..."
    
    if command -v cloudflared >/dev/null; then
        echo "✅ cloudflared 已安装"
        return 0
    fi
    
    # 下载并安装 cloudflared
    wget -q https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-amd64.deb
    sudo dpkg -i cloudflared-linux-amd64.deb
    rm cloudflared-linux-amd64.deb
    
    if command -v cloudflared >/dev/null; then
        echo "✅ cloudflared 安装成功"
        return 0
    else
        echo "❌ cloudflared 安装失败"
        return 1
    fi
}

# 检查服务状态
check_services() {
    echo "🔍 检查本地服务状态..."
    
    if curl -s http://localhost:8000 >/dev/null; then
        echo "✅ 前端服务运行正常"
    else
        echo "❌ 前端服务未运行"
        return 1
    fi
    
    if curl -s http://localhost:8080/api/v1/info >/dev/null; then
        echo "✅ 后端服务运行正常"
    else
        echo "❌ 后端服务未运行"
        return 1
    fi
    
    return 0
}

# 启动隧道
start_tunnel() {
    echo ""
    echo "🚇 启动 Cloudflare 隧道..."
    echo ""
    echo "选择要暴露的服务："
    echo "1. 前端应用 (端口 8000)"
    echo "2. 后端API (端口 8080)"
    echo "3. 两个都暴露"
    echo ""
    read -p "请选择 (1/2/3): " choice
    
    case $choice in
        1)
            echo "🌐 启动前端隧道..."
            echo "⚠️  注意：前端需要连接到后端API，可能需要配置CORS"
            cloudflared tunnel --url http://localhost:8000
            ;;
        2)
            echo "🔧 启动后端API隧道..."
            cloudflared tunnel --url http://localhost:8080
            ;;
        3)
            echo "🔄 启动两个隧道..."
            echo "前端隧道将在后台运行，后端隧道在前台运行"
            echo ""
            
            # 启动前端隧道（后台）
            cloudflared tunnel --url http://localhost:8000 > frontend-tunnel.log 2>&1 &
            FRONTEND_PID=$!
            sleep 3
            
            # 从日志中提取前端URL
            FRONTEND_URL=$(grep -o 'https://[^[:space:]]*\.trycloudflare\.com' frontend-tunnel.log | head -1)
            echo "前端隧道: $FRONTEND_URL"
            echo "前端日志: frontend-tunnel.log"
            echo ""
            
            # 启动后端隧道（前台）
            echo "启动后端隧道..."
            cloudflared tunnel --url http://localhost:8080
            ;;
        *)
            echo "无效选择"
            exit 1
            ;;
    esac
}

# 主函数
main() {
    if ! install_cloudflared; then
        exit 1
    fi
    
    if ! check_services; then
        echo ""
        echo "请先启动服务："
        echo "前端: cd frontend && npm run dev"
        echo "后端: cd backend && DATABASE_URL=\"postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable\" ./pentagi"
        exit 1
    fi
    
    start_tunnel
}

# 运行主函数
main "$@"
