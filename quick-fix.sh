#!/bin/bash

# PentAGI 快速修复脚本

echo "🔧 PentAGI 快速修复脚本"
echo "======================"

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 修复函数
fix_nmap_permissions() {
    echo -e "${BLUE}修复 nmap 权限问题...${NC}"
    
    for container in pentagi-terminal-1 pentagi-terminal-2; do
        echo "修复容器: $container"
        docker exec $container bash -c "
cat > /usr/bin/nmap << 'EOF'
#!/usr/bin/env sh
set -e
if [ ! -f /tmp/nmap ] || [ /usr/lib/nmap/nmap -nt /tmp/nmap ]; then
  cp /usr/lib/nmap/nmap /tmp/nmap
  chmod +x /tmp/nmap
fi
if [ \"\$(id -u)\" -eq 0 ] || [ \"\$1\" = \"--resume\" ]; then
  exec /tmp/nmap \"\$@\"
else
  exec /tmp/nmap --privileged \"\$@\"
fi
EOF
chmod +x /usr/bin/nmap"
    done
    echo -e "${GREEN}✓ nmap 权限修复完成${NC}"
}

fix_docker_network() {
    echo -e "${BLUE}修复 Docker 网络连接...${NC}"
    
    # 重新连接网络
    for container in pentagi-terminal-1 pentagi-terminal-2; do
        echo "连接容器 $container 到电力网络..."
        docker network disconnect power-it-lab_power-network $container 2>/dev/null || true
        docker network connect power-it-lab_power-network $container
    done
    echo -e "${GREEN}✓ Docker 网络修复完成${NC}"
}

fix_power_lab() {
    echo -e "${BLUE}重启电力IT实验室...${NC}"
    
    cd power-it-lab
    docker-compose down
    docker-compose up -d
    cd ..
    
    # 等待服务启动
    echo "等待服务启动..."
    sleep 10
    
    echo -e "${GREEN}✓ 电力IT实验室重启完成${NC}"
}

fix_env_variables() {
    echo -e "${BLUE}修复环境变量配置...${NC}"
    
    # 确保backend目录有.env文件
    if [ ! -f backend/.env ]; then
        echo "复制 .env 文件到 backend 目录..."
        cp .env backend/.env
    fi
    
    # 检查关键环境变量
    if ! grep -q "ANTHROPIC_API_KEY=sk-" backend/.env; then
        echo -e "${YELLOW}⚠️  警告: Anthropic API密钥可能未正确配置${NC}"
    fi
    
    echo -e "${GREEN}✓ 环境变量检查完成${NC}"
}

restart_services() {
    echo -e "${BLUE}重启 PentAGI 服务...${NC}"
    
    echo "请手动执行以下命令重启服务:"
    echo ""
    echo "1. 停止当前运行的前端和后端服务 (Ctrl+C)"
    echo ""
    echo "2. 重启后端:"
    echo "   cd backend"
    echo "   source .env"
    echo "   go run cmd/pentagi/main.go"
    echo ""
    echo "3. 重启前端 (新终端):"
    echo "   cd frontend"
    echo "   npm run dev"
    echo ""
    echo "4. 重启 cloudflare 隧道 (新终端):"
    echo "   cloudflared tunnel --url http://localhost:8000"
    echo ""
}

# 主菜单
show_menu() {
    echo ""
    echo "请选择要执行的修复操作:"
    echo "1) 修复 nmap 工具权限问题"
    echo "2) 修复 Docker 网络连接"
    echo "3) 重启电力IT实验室"
    echo "4) 修复环境变量配置"
    echo "5) 显示服务重启指南"
    echo "6) 执行所有修复 (推荐)"
    echo "7) 运行健康检查"
    echo "0) 退出"
    echo ""
    read -p "请输入选项 (0-7): " choice
}

# 执行所有修复
fix_all() {
    echo -e "${BLUE}执行所有修复操作...${NC}"
    echo ""
    
    fix_env_variables
    echo ""
    
    fix_power_lab
    echo ""
    
    fix_docker_network
    echo ""
    
    fix_nmap_permissions
    echo ""
    
    echo -e "${GREEN}🎉 所有自动修复完成！${NC}"
    echo ""
    echo "接下来请手动重启 PentAGI 服务:"
    restart_services
}

# 主循环
while true; do
    show_menu
    
    case $choice in
        1)
            fix_nmap_permissions
            ;;
        2)
            fix_docker_network
            ;;
        3)
            fix_power_lab
            ;;
        4)
            fix_env_variables
            ;;
        5)
            restart_services
            ;;
        6)
            fix_all
            break
            ;;
        7)
            echo -e "${BLUE}运行健康检查...${NC}"
            ./health-check.sh
            ;;
        0)
            echo "退出修复脚本"
            break
            ;;
        *)
            echo -e "${RED}无效选项，请重新选择${NC}"
            ;;
    esac
    
    echo ""
    read -p "按 Enter 键继续..."
done
