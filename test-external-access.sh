#!/bin/bash

# 测试外部访问脚本
echo "🌐 测试外部访问配置"
echo "==================="

# 获取所有IP地址
echo "📍 当前服务器IP地址："
ip addr show | grep -E "inet.*scope global" | awk '{print $2}' | cut -d'/' -f1

echo ""
echo "🔍 检查服务绑定状态："
echo "前端服务 (端口8000):"
ss -tlnp | grep :8000

echo ""
echo "后端服务 (端口8080):"
ss -tlnp | grep :8080

echo ""
echo "🧪 测试本地访问："

# 测试本地访问
echo -n "本地前端 (localhost:8000): "
if curl -s http://localhost:8000 >/dev/null; then
    echo "✅ 可访问"
else
    echo "❌ 不可访问"
fi

echo -n "本地后端 (localhost:8080): "
if curl -s http://localhost:8080/api/v1/info >/dev/null; then
    echo "✅ 可访问"
else
    echo "❌ 不可访问"
fi

echo ""
echo "🌍 测试外部IP访问："

# 获取主要IP地址
MAIN_IP=$(ip route get 8.8.8.8 | awk '{print $7; exit}' 2>/dev/null || echo "172.17.0.2")

echo "使用主要IP: $MAIN_IP"

echo -n "外部前端 ($MAIN_IP:8000): "
if curl -s http://$MAIN_IP:8000 >/dev/null; then
    echo "✅ 可访问"
else
    echo "❌ 不可访问"
fi

echo -n "外部后端 ($MAIN_IP:8080): "
if curl -s http://$MAIN_IP:8080/api/v1/info >/dev/null; then
    echo "✅ 可访问"
else
    echo "❌ 不可访问"
fi

echo ""
echo "🔗 外部访问地址："
echo "前端应用: http://$MAIN_IP:8000"
echo "后端API:  http://$MAIN_IP:8080"
echo "Swagger:  http://$MAIN_IP:8080/api/v1/swagger/index.html"
echo "GraphQL:  http://$MAIN_IP:8080/api/v1/graphql/playground"

echo ""
echo "📋 CORS配置检查："
echo "当前CORS设置："
grep "CORS_ORIGINS" .env

echo ""
echo "🔧 如果外部访问失败，可能的原因："
echo "1. 防火墙阻止了端口8000和8080"
echo "2. 网络配置不允许外部访问"
echo "3. 需要在云服务商控制台开放端口"
echo "4. 需要配置反向代理"

echo ""
echo "💡 建议的解决方案："
echo "1. 使用端口转发（推荐）"
echo "2. 配置nginx反向代理"
echo "3. 使用ngrok等隧道工具"
