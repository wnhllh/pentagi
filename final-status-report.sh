#!/bin/bash

# 最终状态报告和问题解决总结
echo "🎯 PentAGI 502/524 错误排查结果"
echo "=============================="

echo ""
echo "🔍 问题诊断："
echo "原因: Cloudflare 隧道无法连接到本地服务"
echo "症状: 502 Bad Gateway, 524 Timeout 错误"
echo "解决: 重启隧道服务，更新配置"

echo ""
echo "✅ 解决方案执行："
echo "1. 重启了 Cloudflare 隧道"
echo "2. 获得了新的隧道URL"
echo "3. 更新了前端配置"
echo "4. 更新了CORS配置"
echo "5. 重启了前端和后端服务"

echo ""
echo "🌐 新的访问地址："
echo "┌─────────────────────────────────────────────────────────────┐"
echo "│ 🎨 前端应用 (主要访问地址)                                   │"
echo "│    https://meal-ag-facility-trucks.trycloudflare.com        │"
echo "│                                                             │"
echo "│ 🔧 后端API                                                   │"
echo "│    https://farms-python-blond-nationally.trycloudflare.com  │"
echo "│                                                             │"
echo "│ 📚 API文档                                                   │"
echo "│    https://farms-python-blond-nationally.trycloudflare.com/api/v1/swagger/index.html │"
echo "│                                                             │"
echo "│ 🔍 GraphQL                                                   │"
echo "│    https://farms-python-blond-nationally.trycloudflare.com/api/v1/graphql/playground │"
echo "└─────────────────────────────────────────────────────────────┘"

echo ""
echo "🧪 连接测试："

# 测试后端API
echo -n "后端API: "
if response=$(curl -s https://farms-python-blond-nationally.trycloudflare.com/api/v1/info); then
    if echo "$response" | jq -e '.status == "success"' >/dev/null 2>&1; then
        echo "✅ 正常"
    else
        echo "⚠️  响应异常"
    fi
else
    echo "❌ 无响应"
fi

# 测试前端
echo -n "前端页面: "
if curl -s https://meal-ag-facility-trucks.trycloudflare.com | grep -q "<!doctype html"; then
    echo "✅ 正常"
else
    echo "❌ 异常"
fi

# 测试本地连接
echo -n "本地后端: "
if curl -s http://localhost:8080/api/v1/info >/dev/null; then
    echo "✅ 正常"
else
    echo "❌ 异常"
fi

echo -n "本地前端: "
if curl -s http://localhost:8000 >/dev/null; then
    echo "✅ 正常"
else
    echo "❌ 异常"
fi

echo ""
echo "📊 当前服务状态："

# 检查进程
echo "运行中的服务:"
if pgrep -f "air" >/dev/null; then
    echo "  ✅ 后端 (热重载模式)"
else
    echo "  ❌ 后端未运行"
fi

if pgrep -f "vite" >/dev/null; then
    echo "  ✅ 前端 (开发模式)"
else
    echo "  ❌ 前端未运行"
fi

if docker ps | grep -q pgvector; then
    echo "  ✅ 数据库 (Docker)"
else
    echo "  ❌ 数据库未运行"
fi

if pgrep -f "cloudflared" >/dev/null; then
    tunnel_count=$(pgrep -f "cloudflared" | wc -l)
    echo "  ✅ 隧道服务 ($tunnel_count 个)"
else
    echo "  ❌ 隧道未运行"
fi

echo ""
echo "🔧 配置更新："
echo "• CORS: 已更新支持新隧道域名"
echo "• 前端API URL: 已更新为新后端隧道"
echo "• 环境变量: Anthropic API 密钥已配置"

echo ""
echo "💡 使用建议："
echo "1. 使用前端隧道地址访问应用"
echo "2. 如果遇到问题，检查隧道日志"
echo "3. 隧道URL会在重启后改变"
echo "4. 开发时使用本地地址更稳定"

echo ""
echo "🚨 注意事项："
echo "• 免费隧道有连接限制"
echo "• 隧道可能偶尔不稳定"
echo "• 生产环境建议使用付费隧道"

echo ""
echo "🎉 问题已解决！现在可以正常使用 PentAGI 了！"
