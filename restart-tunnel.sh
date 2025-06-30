#!/bin/bash

# PentAGI Cloudflare隧道重启脚本

echo "🌐 PentAGI Cloudflare隧道重启脚本"
echo "================================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 检查前端服务端口
echo -n "检查前端服务端口... "
if curl -s http://localhost:8000 > /dev/null; then
    FRONTEND_PORT=8000
    echo -e "${GREEN}8000${NC}"
elif curl -s http://localhost:8001 > /dev/null; then
    FRONTEND_PORT=8001
    echo -e "${GREEN}8001${NC}"
else
    echo -e "${RED}前端服务未运行${NC}"
    echo "请先启动前端服务: cd frontend && npm run dev"
    exit 1
fi

# 停止现有的cloudflared进程
echo -n "停止现有隧道... "
pkill -f cloudflared 2>/dev/null
sleep 2
echo -e "${GREEN}完成${NC}"

# 创建新隧道
echo -e "${BLUE}创建新的Cloudflare隧道...${NC}"
echo "目标端口: localhost:$FRONTEND_PORT"
echo ""

# 启动新隧道并捕获输出
cloudflared tunnel --url http://localhost:$FRONTEND_PORT &
TUNNEL_PID=$!

# 等待隧道启动并获取URL
echo "等待隧道启动..."
sleep 5

# 检查隧道是否成功启动
if ps -p $TUNNEL_PID > /dev/null; then
    echo -e "${GREEN}✓ 隧道启动成功${NC}"
    echo -e "${YELLOW}注意: 隧道URL会在cloudflared的输出中显示${NC}"
    echo ""
    echo "请查看上面的输出获取新的隧道URL，格式类似："
    echo "https://xxxxx-xxxxx-xxxxx-xxxxx.trycloudflare.com"
    echo ""
    echo -e "${BLUE}下一步操作:${NC}"
    echo "1. 复制新的隧道URL"
    echo "2. 更新 frontend/vite.config.ts 中的 allowedHosts"
    echo "3. 重启前端服务 (可选，通常不需要)"
    echo ""
    echo -e "${GREEN}隧道进程ID: $TUNNEL_PID${NC}"
    echo "要停止隧道，运行: kill $TUNNEL_PID"
else
    echo -e "${RED}✗ 隧道启动失败${NC}"
    exit 1
fi
