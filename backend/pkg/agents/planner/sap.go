package planner

import (
	"context"
	"strings"

	"pentagi/pkg/agents/router"
)

// SapPlanner SAP ERP系统专用测试计划器
type SapPlanner struct{}

// NewSapPlanner 创建SAP ERP测试计划器
func NewSapPlanner() *SapPlanner {
	return &SapPlanner{}
}

// Name 返回计划器名称
func (p *SapPlanner) Name() string {
	return "SAP ERP系统专用测试计划"
}

// Match 判断是否匹配SAP ERP系统
func (p *SapPlanner) Match(meta router.TargetMeta) bool {
	// 匹配条件：
	// 1. 域名包含 sap, erp
	// 2. AppID 为 erp
	// 3. SAP特有端口 (3200, 3300, 8000等)
	return strings.Contains(meta.Domain, "sap") ||
		strings.Contains(meta.Domain, "erp") ||
		meta.AppID == "erp" ||
		meta.Port == 3200 || meta.Port == 3300 || meta.Port == 8000
}

// Plan 生成SAP ERP系统专用测试计划
func (p *SapPlanner) Plan(ctx context.Context, in interface{}) ([]interface{}, error) {
	var steps []interface{}

	// 步骤1: SAP电力企业系统识别和信息收集
	step1 := map[string]interface{}{
		"type": "sap_power_system_recon",
		"content": `🏢 SAP电力企业系统深度侦察

【系统识别】领域感知规划器识别目标为电力企业SAP ERP系统，激活企业级安全测试剧本。

🔍 **SAP电力企业组件发现**:
- SAP ECC (企业核心组件) - 财务、采购、销售
- SAP IS-U (公用事业解决方案) - 电力行业专用模块
- SAP CRM (客户关系管理) - 客户服务和营销
- SAP BW (商业智能) - 数据仓库和分析

⚡ **电力行业特有模块**:
- FI-CA (合同账户和发票) - 电费计算和账单
- CS (客户服务) - 电力客户管理
- PM (设备维护) - 电力设备维护管理
- PS (项目系统) - 电力工程项目管理

🎯 **攻击面分析**:
- SAP GUI 7.x 客户端漏洞
- RFC (远程函数调用) 未授权访问
- Web Dynpro 应用安全缺陷
- ICM (互联网通信管理器) 配置问题

请使用scan_sap_security工具开始SAP系统深度侦察。`,
	}
	steps = append(steps, step1)

	// 步骤2: SAP默认凭据和权限提升测试
	step2 := map[string]interface{}{
		"type": "sap_privilege_escalation",
		"content": `🔐 SAP默认凭据和权限提升攻击

【攻击策略】针对SAP系统进行系统性权限提升测试：

👤 **SAP经典后门账户测试**:
- SAP* / 06071992 (SAP经典默认密码)
- DDIC / 19920706 (数据字典用户)
- DEVELOPER / ch4ngeme (开发者账户)
- EARLYWATCH / support (支持账户)

🔑 **RFC权限提升测试**:
- RFC_READ_TABLE - 直接读取数据库表
- RFC_PING - 系统连通性测试
- RFC_SYSTEM_INFO - 系统信息泄露
- RFC_ABAP_INSTALL_AND_RUN - 代码执行

⚡ **电力行业敏感事务代码**:
- FPL9 (电力计划) - 电力生产计划访问
- IS-U 事务 - 公用事业计费和客户管理
- FI52 (总账余额) - 财务数据访问
- SE16 (数据浏览器) - 直接数据库访问

🚨 **预期发现**: SAP经典后门账户可用，RFC未授权访问，敏感事务代码权限控制不足

请使用power_pentester工具执行SAP权限提升测试。`,
	}
	steps = append(steps, step2)

	// 步骤3: 电力企业财务和业务数据安全测试
	step3 := map[string]interface{}{
		"type": "sap_financial_data_security",
		"content": `💼 电力企业财务和业务数据安全测试

【业务逻辑攻击】针对电力企业SAP系统核心业务数据进行安全测试：

💰 **财务数据访问控制测试**:
- 电力销售收入数据 (表: VBRK, VBRP)
- 客户应收账款信息 (表: KNA1, BSID)
- 电力采购成本数据 (表: EKKO, EKPO)
- 固定资产和设备投资 (表: ANLA, ANLC)

⚡ **电力业务数据越权测试**:
- 客户用电数据和计费信息
- 电力生产和调度数据
- 设备维护和故障记录
- 电力交易和合同信息

🔍 **采购流程权限绕过测试**:
- 电力设备采购审批流程
- 供应商主数据篡改测试
- 采购订单金额修改验证
- 发票验证流程绕过

👥 **人力资源数据保护测试**:
- 电力企业员工薪资信息
- 特殊岗位人员安全等级
- 培训记录和资质认证
- 绩效评估和奖金数据

🎯 **关键风险**: 财务敏感数据无充分保护，采购审批可被绕过，员工薪资信息泄露

请使用biz_logic_tester工具执行SAP业务数据安全测试。`,
	}
	steps = append(steps, step3)

	// 简化为单步骤实现

	return steps, nil
}

// extractTarget 从输入中提取目标信息 (简化版本)
func (p *SapPlanner) extractTarget(in interface{}) string {
	// 简化实现，返回默认值
	return "SAP ERP系统目标"
}
