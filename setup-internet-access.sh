#!/bin/bash

# PentAGI 互联网访问设置脚本
echo "🌍 PentAGI 互联网访问设置"
echo "========================"

# 检查服务状态
check_services() {
    echo "🔍 检查本地服务状态..."
    
    # 检查前端
    if curl -s http://localhost:8000 >/dev/null; then
        echo "✅ 前端服务运行正常 (localhost:8000)"
    else
        echo "❌ 前端服务未运行，请先启动前端服务"
        echo "   运行: cd frontend && npm run dev"
        return 1
    fi
    
    # 检查后端
    if curl -s http://localhost:8080/api/v1/info >/dev/null; then
        echo "✅ 后端服务运行正常 (localhost:8080)"
    else
        echo "❌ 后端服务未运行，请先启动后端服务"
        echo "   运行: cd backend && DATABASE_URL=\"postgres://postgres:postgres@localhost:5432/pentagidb?sslmode=disable\" ./pentagi"
        return 1
    fi
    
    return 0
}

# 设置ngrok隧道
setup_ngrok() {
    echo ""
    echo "🚇 设置 ngrok 隧道..."
    
    # 检查ngrok是否已认证
    if ! ngrok config check >/dev/null 2>&1; then
        echo "⚠️  ngrok 需要认证令牌"
        echo "请访问 https://dashboard.ngrok.com/get-started/your-authtoken"
        echo "获取您的认证令牌，然后运行："
        echo "   ngrok config add-authtoken YOUR_TOKEN"
        echo ""
        echo "或者选择其他方案（见下方）"
        return 1
    fi
    
    echo "🔧 启动 ngrok 隧道..."
    
    # 创建ngrok配置文件
    cat > ngrok.yml << EOF
version: "2"
tunnels:
  frontend:
    addr: 8000
    proto: http
    bind_tls: true
  backend:
    addr: 8080
    proto: http
    bind_tls: true
EOF
    
    echo "启动隧道中..."
    ngrok start --all --config ngrok.yml &
    NGROK_PID=$!
    
    # 等待ngrok启动
    sleep 5
    
    # 获取ngrok URL
    if command -v jq >/dev/null; then
        FRONTEND_URL=$(curl -s http://localhost:4040/api/tunnels | jq -r '.tunnels[] | select(.name=="frontend") | .public_url')
        BACKEND_URL=$(curl -s http://localhost:4040/api/tunnels | jq -r '.tunnels[] | select(.name=="backend") | .public_url')
    else
        echo "⚠️  需要安装 jq 来获取隧道URL"
        echo "请访问 http://localhost:4040 查看隧道状态"
        FRONTEND_URL="请查看 http://localhost:4040"
        BACKEND_URL="请查看 http://localhost:4040"
    fi
    
    echo ""
    echo "🎉 ngrok 隧道已启动！"
    echo "前端访问地址: $FRONTEND_URL"
    echo "后端API地址:  $BACKEND_URL"
    echo "ngrok 控制台: http://localhost:4040"
    echo ""
    echo "⚠️  注意: 免费版ngrok有连接限制，生产环境建议升级"
    
    return 0
}

# 其他方案说明
show_alternatives() {
    echo ""
    echo "🔄 其他互联网访问方案："
    echo ""
    echo "方案2: Cloudflare Tunnel (免费，推荐)"
    echo "1. 安装 cloudflared"
    echo "2. 运行: cloudflared tunnel --url http://localhost:8000"
    echo ""
    echo "方案3: 端口转发 + 公网IP"
    echo "1. 如果您有公网IP，配置路由器端口转发"
    echo "2. 转发 8000 → 内网IP:8000"
    echo "3. 转发 8080 → 内网IP:8080"
    echo ""
    echo "方案4: VPS反向代理"
    echo "1. 在VPS上安装nginx"
    echo "2. 配置反向代理到您的本地服务"
    echo "3. 使用SSH隧道连接"
    echo ""
    echo "方案5: 云服务部署"
    echo "1. 部署到 Vercel/Netlify (前端)"
    echo "2. 部署到 Railway/Render (后端)"
    echo "3. 使用云数据库"
}

# 主函数
main() {
    if ! check_services; then
        exit 1
    fi
    
    echo ""
    echo "选择互联网访问方案："
    echo "1. 使用 ngrok (简单快速)"
    echo "2. 查看其他方案"
    echo ""
    read -p "请选择 (1/2): " choice
    
    case $choice in
        1)
            if setup_ngrok; then
                echo ""
                echo "🎯 现在您可以从互联网访问 PentAGI 了！"
                echo "按 Ctrl+C 停止隧道"
                wait $NGROK_PID
            fi
            ;;
        2)
            show_alternatives
            ;;
        *)
            echo "无效选择"
            exit 1
            ;;
    esac
}

# 如果直接运行脚本
if [[ "${BASH_SOURCE[0]}" == "${0}" ]]; then
    main "$@"
fi
