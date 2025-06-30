package agents

import "context"

// Agent 基础Agent接口
type Agent interface {
	Execute(ctx context.Context, input AgentInput) (*AgentOutput, error)
	GetName() string
	GetDescription() string
}

// BaseAgent 基础Agent实现
type BaseAgent struct {
	name        string
	description string
}

// GetName 返回Agent名称
func (a *BaseAgent) GetName() string {
	return a.name
}

// GetDescription 返回Agent描述
func (a *BaseAgent) GetDescription() string {
	return a.description
}

// Message 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// AgentInput Agent输入
type AgentInput struct {
	Messages []Message   `json:"messages"`
	Tools    interface{} `json:"tools,omitempty"`
	Target   string      `json:"target,omitempty"`
}

// AgentOutput Agent输出
type AgentOutput struct {
	Messages  []Message              `json:"messages"`
	ToolCalls []interface{}          `json:"tool_calls,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
