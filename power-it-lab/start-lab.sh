#!/bin/bash

# 电力IT系统安全靶场启动脚本
# 用于启动三个电力行业IT系统的安全测试环境

set -e

echo "🏗️  电力IT系统安全靶场启动器"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查Docker和Docker Compose
check_dependencies() {
    echo -e "${BLUE}🔍 检查依赖...${NC}"
    
    if ! command -v docker &> /dev/null; then
        echo -e "${RED}❌ Docker 未安装${NC}"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        echo -e "${RED}❌ Docker Compose 未安装${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}✅ 依赖检查通过${NC}"
}

# 创建必要的目录
create_directories() {
    echo -e "${BLUE}📁 创建必要目录...${NC}"
    
    mkdir -p marketing-system/gateway/uploads
    mkdir -p iguowang-app/uploads
    mkdir -p erp-system/reports
    mkdir -p tools
    
    # 创建示例报表文件
    echo "电力系统运行报表" > erp-system/reports/power_report.txt
    echo "财务月度报表" > erp-system/reports/finance_monthly.txt
    echo "../../etc/passwd" > erp-system/reports/sensitive_file.txt
    
    echo -e "${GREEN}✅ 目录创建完成${NC}"
}

# 显示菜单
show_menu() {
    echo ""
    echo -e "${YELLOW}请选择要启动的靶场:${NC}"
    echo "1) 🏭 电力营销系统 2.0 靶场"
    echo "2) 📱 i国网APP 后端靶场"
    echo "3) 🏢 ERP系统靶场 (SAP风格)"
    echo "4) 🚀 启动所有靶场"
    echo "5) 📊 查看靶场状态"
    echo "6) 🛑 停止所有靶场"
    echo "7) 🧹 清理环境"
    echo "8) 📖 显示使用指南"
    echo "0) 退出"
    echo ""
}

# 启动营销系统靶场
start_marketing() {
    echo -e "${BLUE}🏭 启动电力营销系统 2.0 靶场...${NC}"
    docker-compose up -d marketing-gateway marketing-user-service marketing-billing-service marketing-db marketing-cache
    echo -e "${GREEN}✅ 营销系统靶场启动完成${NC}"
    echo -e "${YELLOW}访问地址: http://localhost:8080${NC}"
}

# 启动i国网靶场
start_iguowang() {
    echo -e "${BLUE}📱 启动i国网APP 后端靶场...${NC}"
    docker-compose up -d iguowang-api iguowang-mobile-sim iguowang-db
    echo -e "${GREEN}✅ i国网靶场启动完成${NC}"
    echo -e "${YELLOW}API地址: http://localhost:9080${NC}"
    echo -e "${YELLOW}移动端模拟器: http://localhost:9090${NC}"
}

# 启动ERP靶场
start_erp() {
    echo -e "${BLUE}🏢 启动ERP系统靶场...${NC}"
    docker-compose up -d erp-app-server erp-web-gui erp-db
    echo -e "${GREEN}✅ ERP靶场启动完成${NC}"
    echo -e "${YELLOW}ERP服务器: http://localhost:8000${NC}"
    echo -e "${YELLOW}Web界面: http://localhost:8001${NC}"
}

# 启动所有靶场
start_all() {
    echo -e "${BLUE}🚀 启动所有靶场...${NC}"
    docker-compose up -d
    echo -e "${GREEN}✅ 所有靶场启动完成${NC}"
    show_access_info
}

# 显示访问信息
show_access_info() {
    echo ""
    echo -e "${GREEN}🌐 靶场访问信息:${NC}"
    echo "=================================="
    echo -e "${YELLOW}电力营销系统 2.0:${NC}"
    echo "  - API网关: http://localhost:8080"
    echo "  - 用户服务: http://localhost:8081"
    echo "  - 计费服务: http://localhost:8082"
    echo "  - 数据库: localhost:3306"
    echo ""
    echo -e "${YELLOW}i国网APP:${NC}"
    echo "  - API服务器: http://localhost:9080"
    echo "  - 移动端模拟器: http://localhost:9090"
    echo "  - 数据库: localhost:5432"
    echo ""
    echo -e "${YELLOW}ERP系统:${NC}"
    echo "  - 应用服务器: http://localhost:8000"
    echo "  - Web界面: http://localhost:8001"
    echo "  - 数据库: localhost:5433"
    echo ""
    echo -e "${YELLOW}攻击工具:${NC}"
    echo "  - Kali Linux: docker exec -it power-kali-attacker bash"
    echo "  - 流量监控: http://localhost:3000"
    echo ""
}

# 查看状态
show_status() {
    echo -e "${BLUE}📊 靶场状态:${NC}"
    echo "=================================="
    docker-compose ps
}

# 停止所有服务
stop_all() {
    echo -e "${BLUE}🛑 停止所有靶场...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ 所有靶场已停止${NC}"
}

# 清理环境
cleanup() {
    echo -e "${BLUE}🧹 清理环境...${NC}"
    docker-compose down -v --remove-orphans
    docker system prune -f
    echo -e "${GREEN}✅ 环境清理完成${NC}"
}

# 显示使用指南
show_guide() {
    echo -e "${GREEN}📖 电力IT系统安全靶场使用指南${NC}"
    echo "============================================"
    echo ""
    echo -e "${YELLOW}🏭 电力营销系统 2.0 靶场${NC}"
    echo "主要漏洞类型:"
    echo "  - SQL注入: /api/auth/login"
    echo "  - 越权访问: /api/users/:userId"
    echo "  - 命令执行: /api/admin/backdoor"
    echo "  - IDOR: /api/billing/:accountId"
    echo "  - 任意文件上传: /api/upload"
    echo "  - 信息泄露: /api/system/info"
    echo ""
    echo "测试账户:"
    echo "  - admin/admin123 (管理员)"
    echo "  - test/test (测试账户)"
    echo "  - guest/ (空密码)"
    echo ""
    echo -e "${YELLOW}📱 i国网APP 靶场${NC}"
    echo "主要漏洞类型:"
    echo "  - 短信验证码绕过: /api/auth/send-sms"
    echo "  - SQL注入: /api/auth/login"
    echo "  - 越权访问: /api/user/profile"
    echo "  - IDOR: /api/billing/query"
    echo "  - 金额篡改: /api/payment/create"
    echo "  - 命令执行: /api/admin/debug"
    echo "  - 任意文件上传: /api/file/upload"
    echo "  - 信息泄露: /api/system/config"
    echo ""
    echo "测试手机号:"
    echo "  - 13800138000 (有历史数据)"
    echo "  - 10000000000 (测试账户)"
    echo ""
    echo -e "${YELLOW}🏢 ERP系统靶场${NC}"
    echo "主要漏洞类型:"
    echo "  - 默认密码: SAP*/06071992"
    echo "  - SQL注入: /api/finance/query"
    echo "  - 越权访问: /api/hr/employee"
    echo "  - 命令执行: /api/admin/execute"
    echo "  - 信息泄露: /api/system/config"
    echo "  - 路径遍历: /api/report/generate"
    echo ""
    echo "测试账户:"
    echo "  - SAP*/06071992 (后门账户)"
    echo "  - ADMIN/ADMIN123 (管理员)"
    echo "  - DDIC/DDIC (开发账户)"
    echo ""
    echo -e "${BLUE}🔧 测试工具建议:${NC}"
    echo "  - Burp Suite (Web应用测试)"
    echo "  - SQLMap (SQL注入测试)"
    echo "  - Postman (API测试)"
    echo "  - Nmap (端口扫描)"
    echo "  - Wireshark (流量分析)"
    echo ""
}

# 主程序
main() {
    check_dependencies
    create_directories
    
    while true; do
        show_menu
        read -p "请输入选择 [0-8]: " choice
        
        case $choice in
            1)
                start_marketing
                ;;
            2)
                start_iguowang
                ;;
            3)
                start_erp
                ;;
            4)
                start_all
                ;;
            5)
                show_status
                ;;
            6)
                stop_all
                ;;
            7)
                cleanup
                ;;
            8)
                show_guide
                ;;
            0)
                echo -e "${GREEN}👋 再见!${NC}"
                exit 0
                ;;
            *)
                echo -e "${RED}❌ 无效选择，请重新输入${NC}"
                ;;
        esac
        
        echo ""
        read -p "按回车键继续..."
    done
}

# 运行主程序
main "$@"
