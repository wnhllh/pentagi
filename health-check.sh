#!/bin/bash

# PentAGI 电力行业安全测试系统 - 健康检查脚本

echo "🔍 PentAGI 系统健康检查"
echo "======================="
echo "时间: $(date)"
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查函数
check_service() {
    local service_name=$1
    local check_command=$2
    local expected_result=$3
    
    echo -n "检查 $service_name... "
    
    if eval "$check_command" > /dev/null 2>&1; then
        if [ -n "$expected_result" ]; then
            result=$(eval "$check_command" 2>/dev/null)
            if [[ "$result" == *"$expected_result"* ]]; then
                echo -e "${GREEN}✓ 正常${NC}"
                return 0
            else
                echo -e "${RED}✗ 异常${NC}"
                return 1
            fi
        else
            echo -e "${GREEN}✓ 正常${NC}"
            return 0
        fi
    else
        echo -e "${RED}✗ 失败${NC}"
        return 1
    fi
}

# 检查计数器
total_checks=0
passed_checks=0

# 1. 检查数据库
echo "📊 数据库服务"
echo "-------------"
((total_checks++))
if check_service "PostgreSQL容器" "docker ps | grep pgvector" "pgvector"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "数据库连接" "docker exec pgvector psql -U postgres -d pentagidb -c 'SELECT 1;'" "1"; then
    ((passed_checks++))
fi

echo ""

# 2. 检查电力IT实验室
echo "⚡ 电力IT实验室"
echo "---------------"
((total_checks++))
if check_service "电力营销系统" "curl -I http://172.25.0.10:8080" "HTTP"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "i国网APP" "curl -I http://172.25.0.9:8080" "HTTP"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "ERP系统" "curl -I http://172.25.0.8:8080" "HTTP"; then
    ((passed_checks++))
fi

echo ""

# 3. 检查PentAGI服务
echo "🤖 PentAGI服务"
echo "---------------"
((total_checks++))
if check_service "后端API" "curl -s http://localhost:8080/api/v1/info | jq -r .status" "success"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "前端服务" "curl -I http://localhost:8000" "HTTP"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "前端代理" "curl -s http://localhost:8000/api/v1/info | jq -r .status" "success"; then
    ((passed_checks++))
fi

echo ""

# 4. 检查Docker网络
echo "🌐 网络连接"
echo "-----------"
((total_checks++))
if check_service "PentAGI网络连接" "docker exec pentagi-terminal-2 curl -I http://172.25.0.10:8080" "HTTP"; then
    ((passed_checks++))
fi

echo ""

# 5. 检查安全工具
echo "🔧 安全工具"
echo "-----------"
((total_checks++))
if check_service "nmap工具" "docker exec pentagi-terminal-2 nmap --version" "Nmap version"; then
    ((passed_checks++))
fi

((total_checks++))
if check_service "masscan工具" "docker exec pentagi-terminal-2 masscan --version" "Masscan version"; then
    ((passed_checks++))
fi

echo ""

# 6. 检查环境变量
echo "🔑 环境配置"
echo "-----------"
((total_checks++))
if check_service "Anthropic API密钥" "grep -q 'ANTHROPIC_API_KEY=sk-' backend/.env" ""; then
    ((passed_checks++))
fi

echo ""

# 7. 检查系统资源
echo "💻 系统资源"
echo "-----------"
echo -n "磁盘空间... "
disk_usage=$(df / | awk 'NR==2 {print $5}' | sed 's/%//')
if [ "$disk_usage" -lt 90 ]; then
    echo -e "${GREEN}✓ 正常 (${disk_usage}% 已使用)${NC}"
    ((passed_checks++))
else
    echo -e "${YELLOW}⚠ 警告 (${disk_usage}% 已使用)${NC}"
fi
((total_checks++))

echo -n "内存使用... "
memory_usage=$(free | awk 'NR==2{printf "%.0f", $3*100/$2}')
if [ "$memory_usage" -lt 90 ]; then
    echo -e "${GREEN}✓ 正常 (${memory_usage}% 已使用)${NC}"
    ((passed_checks++))
else
    echo -e "${YELLOW}⚠ 警告 (${memory_usage}% 已使用)${NC}"
fi
((total_checks++))

echo ""

# 8. 检查端口占用
echo "🔌 端口状态"
echo "-----------"
ports=("8080:后端API" "8000:前端服务" "5432:数据库")
for port_info in "${ports[@]}"; do
    port=$(echo $port_info | cut -d: -f1)
    name=$(echo $port_info | cut -d: -f2)
    echo -n "端口 $port ($name)... "
    if netstat -tulpn 2>/dev/null | grep -q ":$port "; then
        echo -e "${GREEN}✓ 已监听${NC}"
        ((passed_checks++))
    else
        echo -e "${RED}✗ 未监听${NC}"
    fi
    ((total_checks++))
done

echo ""

# 总结
echo "📋 检查总结"
echo "==========="
echo "总检查项: $total_checks"
echo "通过检查: $passed_checks"
echo "失败检查: $((total_checks - passed_checks))"

if [ $passed_checks -eq $total_checks ]; then
    echo -e "${GREEN}🎉 所有检查都通过！系统运行正常。${NC}"
    exit 0
elif [ $passed_checks -gt $((total_checks * 3 / 4)) ]; then
    echo -e "${YELLOW}⚠️  大部分检查通过，但有一些问题需要注意。${NC}"
    exit 1
else
    echo -e "${RED}❌ 多个检查失败，系统可能存在严重问题。${NC}"
    echo ""
    echo "建议操作:"
    echo "1. 查看 RESTART_GUIDE.md 中的故障排除部分"
    echo "2. 检查服务日志"
    echo "3. 按照正确顺序重启服务"
    exit 2
fi
