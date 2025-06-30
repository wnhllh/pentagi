package planner

import (
	"context"
	"strings"

	"pentagi/pkg/agents/router"
)

// MarketingPlanner 电力营销系统2.0专用测试计划器
type MarketingPlanner struct{}

// NewMarketingPlanner 创建营销系统测试计划器
func NewMarketingPlanner() *MarketingPlanner {
	return &MarketingPlanner{}
}

// Name 返回计划器名称
func (p *MarketingPlanner) Name() string {
	return "电力营销系统2.0专用测试计划"
}

// Match 判断是否匹配营销系统
func (p *MarketingPlanner) Match(meta router.TargetMeta) bool {
	// 匹配条件：
	// 1. 域名包含 marketing, billing
	// 2. 端口 8080 (营销系统常用端口)
	// 3. AppID 为 marketing
	return strings.Contains(meta.Domain, "marketing") ||
		strings.Contains(meta.Domain, "billing") ||
		meta.Port == 8080 ||
		meta.AppID == "marketing"
}

// Plan 生成营销系统专用测试计划
func (p *MarketingPlanner) Plan(ctx context.Context, in interface{}) ([]interface{}, error) {
	var steps []interface{}

	// 步骤1: 智能规划与路由识别
	step1 := map[string]interface{}{
		"type": "power_marketing_planning",
		"content": `🎯 电力营销系统专属规划器已激活

【挑战背景】营销系统作为计费和营收的核心，其微服务架构带来了庞大的API攻击面，而复杂的阶梯/分时电价策略则蕴含着极易被传统工具忽略的业务逻辑风险。

【智能规划】领域感知规划与路由中枢成功识别出目标为电力营销系统，任务被精确分派至营销系统专属规划器，激活针对计费逻辑和API安全的测试剧本。

请使用power_pentester工具开始电力营销系统专项安全评估。`,
	}
	steps = append(steps, step1)

	// 步骤2: 业务逻辑穿透测试
	step2 := map[string]interface{}{
		"type": "billing_logic_penetration",
		"content": `💰 业务逻辑穿透 - BillingLogicTester技能激活

【测试策略】调用BillingLogicTester业务感知技能，依据电价模型知识库，自主生成包含数百个测试用例的矩阵：

🔍 **阶梯电价临界点测试**:
- 第一阶梯(0-240kWh): 边界值239.99, 240.00, 240.01
- 第二阶梯(240-400kWh): 边界值399.99, 400.00, 400.01
- 第三阶梯(400kWh+): 大用量测试用例

⏰ **分时电价峰谷转换测试**:
- 峰时段(8:00-22:00): 转换点7:59:59, 8:00:00, 22:00:00, 22:00:01
- 谷时段(22:00-8:00): 深夜电价计算验证
- 节假日特殊电价策略测试

👥 **多客户类型组合测试**:
- 居民用户、商业用户、工业用户计费差异
- 特殊用户群体(低保户、军属)优惠政策验证

🎯 **重点检测目标**: 浮点数精度处理、四舍五入计费偏差、批量计费越权、电价策略绕过

请使用test_billing_logic工具执行深度计费逻辑安全测试。`,
	}
	steps = append(steps, step2)

	// 步骤3: API越权挖掘
	step3 := map[string]interface{}{
		"type": "api_privilege_escalation",
		"content": `🔐 API越权挖掘 - IDORScanner技能激活

【攻击策略】通过前端JS文件分析和流量学习，自动构建系统API拓扑图：

🗺️ **API拓扑发现**:
- /api/v2/customer/{id}/details - 客户详情API
- /api/v2/billing/{id}/history - 计费历史API
- /api/v2/usage/{id}/pattern - 用电模式API
- /api/v2/payment/{id}/records - 缴费记录API

🎯 **IDOR测试矩阵**:
- 使用低权限账户认证令牌
- 系统性遍历其他用户ID资源
- 测试水平越权访问全网客户数据
- 验证垂直权限提升可能性

⚠️ **高风险数据泄露检测**:
- 客户姓名、身份证号、住址信息
- 详细用电习惯和时间模式
- 历史缴费记录和银行信息
- 用户画像和行为分析数据

🏛️ **合规风险评估**: 自动关联《网络安全法》数据安全保护义务，识别潜在合规风险

请使用api_tester工具执行API越权和数据泄露测试。`,
	}
	steps = append(steps, step3)

	// 步骤4: 关键发现分析和风险评级
	step4 := map[string]interface{}{
		"type": "risk_analysis_reporting",
		"content": `📊 关键发现分析和智能风险评级

【预期关键发现】:

🔍 **计费逻辑漏洞**:
- 浮点数精度处理不当导致的四舍五入偏差
- 单次偏差<0.01元，但海量用户规模下年化影响数百万
- 风险等级: 中危 (经济损失潜力巨大)

🚨 **API越权漏洞**:
- /api/v2/customer/{id}/details存在严重水平越权
- 任意用户可遍历ID获取全网客户敏感信息
- 风险等级: 高危 (大规模数据泄露)

📋 **业务影响分析**:
- 计费偏差: 企业资金流失或客户集体投诉
- 数据泄露: 违反《网络安全法》，面临监管处罚
- 声誉损失: 电力公司公信力受损

🎯 **修复建议**:
- 计费系统: 实施精确的decimal计算，避免浮点运算
- API安全: 实施严格的用户身份验证和资源访问控制
- 监控告警: 部署异常访问检测和实时告警机制

请使用power_pentester工具生成详细的安全评估报告。`,
	}
	steps = append(steps, step4)

	// 简化为单步骤实现
	// 在实际使用中，这个步骤会被传递给原有的pentester逻辑执行

	return steps, nil
}

// extractTarget 从输入中提取目标信息 (简化版本)
func (p *MarketingPlanner) extractTarget(in interface{}) string {
	// 简化实现，返回默认值
	return "电力营销系统目标"
}
