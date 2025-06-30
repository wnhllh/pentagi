package router

import (
	"context"
	"fmt"
)

// TargetMeta 目标元数据，用于Planner匹配
type TargetMeta struct {
	Domain string
	Port   int
	AppID  string
	URL    string
}

// Planner 接口定义
type Planner interface {
	// Match 判断是否匹配当前目标
	Match(meta TargetMeta) bool
	// Plan 生成测试计划，返回多个测试步骤
	Plan(ctx context.Context, in interface{}) ([]interface{}, error)
	// Name 返回Planner名称
	Name() string
}

// Router 路由器，在Pentester Agent和执行层之间插入规划层
type Router struct {
	planners []Planner
	fallback interface{} // 现有 Pentester (简化为interface{})
	enabled  bool
}

// NewRouter 创建新的路由器
func NewRouter(enabled bool, fallback interface{}, planners ...Planner) *Router {
	return &Router{
		enabled:  enabled,
		fallback: fallback,
		planners: planners,
	}
}

// Execute 执行路由逻辑 (简化版本，主要用于匹配检查)
func (r *Router) Execute(ctx context.Context, in interface{}) (interface{}, error) {
	// 这个方法主要在tryPowerPlanner中使用，用于检查是否匹配
	// 实际的执行逻辑在performers.go中实现
	return nil, fmt.Errorf("not implemented - use tryPowerPlanner instead")
}

// extractTargetMeta 从输入中提取目标元数据 (简化版本)
func (r *Router) extractTargetMeta(in interface{}) TargetMeta {
	// 这个方法已经在performers.go中实现，这里保留接口兼容性
	return TargetMeta{}
}

// GetName 实现Agent接口
func (r *Router) GetName() string {
	return "PlannerRouter"
}

// GetDescription 实现Agent接口  
func (r *Router) GetDescription() string {
	return "电力行业智能测试路由器，根据目标自动选择专用测试计划"
}
