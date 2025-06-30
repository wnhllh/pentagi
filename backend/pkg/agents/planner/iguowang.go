package planner

import (
	"context"
	"strings"

	"pentagi/pkg/agents/router"
)

// IGuowangPlanner i国网APP专用测试计划器
type IGuowangPlanner struct{}

// NewIGuowangPlanner 创建i国网测试计划器
func NewIGuowangPlanner() *IGuowangPlanner {
	return &IGuowangPlanner{}
}

// Name 返回计划器名称
func (p *IGuowangPlanner) Name() string {
	return "i国网APP专用测试计划"
}

// Match 判断是否匹配i国网系统
func (p *IGuowangPlanner) Match(meta router.TargetMeta) bool {
	// 匹配条件：
	// 1. 域名包含 iguowang, sgcc, mobile
	// 2. AppID 为 iguowang
	// 3. 移动端相关特征
	return strings.Contains(meta.Domain, "iguowang") ||
		strings.Contains(meta.Domain, "sgcc") ||
		strings.Contains(meta.Domain, "mobile") ||
		meta.AppID == "iguowang"
}

// Plan 生成i国网APP专用测试计划
func (p *IGuowangPlanner) Plan(ctx context.Context, in interface{}) ([]interface{}, error) {
	var steps []interface{}

	// 步骤1: 移动端API安全拓扑构建
	step1 := map[string]interface{}{
		"type": "mobile_api_topology",
		"content": `📱 i国网APP移动端API安全拓扑构建

【系统识别】领域感知规划器识别目标为i国网移动应用系统，激活移动端专用测试剧本。

🗺️ **移动API拓扑发现**:
- /api/mobile/auth/* - 移动端认证体系
- /api/mobile/sms/* - 短信验证服务
- /api/mobile/payment/* - 移动支付接口
- /api/mobile/query/* - 电费查询服务
- /api/mobile/notice/* - 停电通知推送

📊 **移动端特有攻击面分析**:
- API版本兼容性问题 (v1/v2并存)
- 移动端特有认证机制缺陷
- 推送通知服务安全风险
- 移动网络环境下的中间人攻击

请使用api_tester工具开始移动端API安全拓扑分析。`,
	}
	steps = append(steps, step1)

	// 步骤2: 短信验证码安全测试
	step2 := map[string]interface{}{
		"type": "sms_verification_bypass",
		"content": `🔐 短信验证码安全穿透测试

【攻击策略】针对移动端短信验证机制进行深度安全测试：

🎯 **验证码泄露检测**:
- 响应中debug_code字段暴露
- 错误信息中验证码回显
- 日志文件中验证码记录

⚡ **验证码绕过技术**:
- 暴力破解4-6位数字验证码
- 时间窗口竞争条件利用
- 验证码重放攻击测试
- 多设备并发验证绕过

📱 **移动端特有风险**:
- SIM卡劫持攻击模拟
- 短信拦截恶意软件测试
- 运营商SS7网络攻击向量

🚨 **预期发现**: 验证码在API响应中直接暴露，存在严重的身份认证绕过风险

请使用test_mobile_security工具执行短信验证码安全测试。`,
	}
	steps = append(steps, step2)

	// 步骤3: 电力服务业务逻辑测试
	step3 := map[string]interface{}{
		"type": "power_service_logic_test",
		"content": `⚡ 电力服务业务逻辑安全测试

【业务场景分析】i国网APP核心电力服务业务逻辑安全验证：

💰 **电费查询权限测试**:
- 用户A查询用户B的电费信息
- 批量遍历用户ID获取电费数据
- 历史电费数据越权访问测试

📢 **停电信息推送安全**:
- 恶意停电通知推送测试
- 虚假紧急停电信息发布
- 推送通知权限提升攻击

💳 **在线缴费流程安全**:
- 缴费金额篡改测试
- 缴费状态操纵验证
- 重复缴费和退款逻辑漏洞

🔍 **故障报修流程测试**:
- 虚假故障报告提交
- 报修工单权限越界访问
- 维修人员身份伪造测试

🎯 **关键风险点**: 用户可查询他人电费、恶意推送停电通知、缴费金额篡改

请使用biz_logic_tester工具执行电力服务业务逻辑测试。`,
	}
	steps = append(steps, step3)

	// 简化为单步骤实现

	return steps, nil
}

// extractTarget 从输入中提取目标信息 (简化版本)
func (p *IGuowangPlanner) extractTarget(in interface{}) string {
	// 简化实现，返回默认值
	return "i国网APP目标"
}
