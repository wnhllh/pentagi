#!/bin/bash

# 验证外部访问设置脚本
echo "🔍 验证外部访问设置"
echo "=================="

# 获取主要IP
MAIN_IP=$(ip route get 8.8.8.8 | awk '{print $7; exit}' 2>/dev/null || echo "172.17.0.2")
echo "🌐 主要IP地址: $MAIN_IP"

echo ""
echo "📋 当前配置检查："

echo "1. 后端CORS配置:"
grep "CORS_ORIGINS" .env | head -1

echo ""
echo "2. 前端API配置:"
grep "VITE_API_URL" frontend/.env

echo ""
echo "3. 服务绑定状态:"
echo "   前端: $(ss -tlnp | grep :8000 | awk '{print $4}')"
echo "   后端: $(ss -tlnp | grep :8080 | awk '{print $4}')"

echo ""
echo "🧪 功能测试："

# 测试后端API
echo -n "1. 后端API响应: "
if response=$(curl -s http://$MAIN_IP:8080/api/v1/info); then
    if echo "$response" | jq -e '.status == "success"' >/dev/null 2>&1; then
        echo "✅ 正常"
    else
        echo "⚠️  响应异常"
    fi
else
    echo "❌ 无响应"
fi

# 测试CORS
echo -n "2. CORS配置: "
cors_header=$(curl -s -H "Origin: http://$MAIN_IP:8000" -I http://$MAIN_IP:8080/api/v1/info | grep -i "access-control-allow-origin")
if [ -n "$cors_header" ]; then
    echo "✅ 正常 ($cors_header)"
else
    echo "❌ 未配置"
fi

# 测试前端
echo -n "3. 前端页面: "
if curl -s http://$MAIN_IP:8000 | grep -q "<!doctype html"; then
    echo "✅ 正常"
else
    echo "❌ 异常"
fi

# 测试前端到后端的代理
echo -n "4. 前端API代理: "
if curl -s http://$MAIN_IP:8000/api/v1/info >/dev/null 2>&1; then
    echo "✅ 正常"
else
    echo "⚠️  可能需要配置"
fi

echo ""
echo "🔗 外部访问地址："
echo "┌─────────────────────────────────────────────────────────────┐"
echo "│ 🌐 前端应用                                                  │"
echo "│    http://$MAIN_IP:8000                                    │"
echo "│                                                             │"
echo "│ 🔧 后端API                                                   │"
echo "│    http://$MAIN_IP:8080                                    │"
echo "│                                                             │"
echo "│ 📚 API文档                                                   │"
echo "│    http://$MAIN_IP:8080/api/v1/swagger/index.html         │"
echo "│                                                             │"
echo "│ 🔍 GraphQL                                                   │"
echo "│    http://$MAIN_IP:8080/api/v1/graphql/playground         │"
echo "└─────────────────────────────────────────────────────────────┘"

echo ""
echo "📱 移动端/其他设备访问："
echo "如果您在同一网络中的其他设备上，可以直接使用上述地址访问"

echo ""
echo "🔧 如果无法从外部访问，请检查："
echo "1. 防火墙设置 (端口8000, 8080)"
echo "2. 云服务商安全组配置"
echo "3. 网络路由配置"
echo "4. 使用端口转发或反向代理"

echo ""
echo "💡 推荐的外部访问方案："
echo "1. 端口转发 (VS Code Remote, SSH隧道)"
echo "2. 反向代理 (nginx, Apache)"
echo "3. 隧道服务 (ngrok, cloudflare tunnel)"

# 如果是在容器或云环境中，提供额外建议
if [ -f /.dockerenv ] || [ -n "$KUBERNETES_SERVICE_HOST" ]; then
    echo ""
    echo "🐳 检测到容器环境，建议："
    echo "1. 配置容器端口映射"
    echo "2. 使用LoadBalancer或NodePort服务"
    echo "3. 配置Ingress控制器"
fi
