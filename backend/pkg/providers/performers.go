package providers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"pentagi/pkg/cast"
	"pentagi/pkg/database"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/tools"
	"pentagi/pkg/agents/router"
	"pentagi/pkg/agents/planner"

	"strconv"
	"github.com/vxcontrol/langchaingo/llms"
)

func (fp *flowProvider) performTaskResultReporter(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemReporterTmpl, userReporterTmpl, input string,
) (*tools.TaskResult, error) {
	var (
		taskResult   tools.TaskResult
		optAgentType = provider.OptionsTypeSimple
		msgChainType = database.MsgchainTypeReporter
	)

	chain := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemReporterTmpl),
		llms.TextParts(llms.ChatMessageTypeHuman, userReporterTmpl),
	}
	cfg := tools.ReporterExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		ReportResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &taskResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal task result: %w", err)
			}
			return "report result successfully processed", nil
		},
	}
	executor, err := fp.executor.GetReporterExecutor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get reporter executor: %w", err)
	}

	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal msg chain: %w", err)
	}

	msgChain, err := fp.db.CreateMsgChain(ctx, database.CreateMsgChainParams{
		Type:          msgChainType,
		Model:         fp.Model(optAgentType),
		ModelProvider: string(fp.Type()),
		Chain:         chainBlob,
		FlowID:        fp.flowID,
		TaskID:        database.Int64ToNullInt64(taskID),
		SubtaskID:     database.Int64ToNullInt64(subtaskID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create msg chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChain.ID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return nil, fmt.Errorf("failed to get task reporter result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			input,
			taskResult.Result,
			taskID,
			subtaskID,
		)
	}

	return &taskResult, nil
}

func (fp *flowProvider) performSubtasksGenerator(
	ctx context.Context,
	taskID int64,
	systemGeneratorTmpl, userGeneratorTmpl, input string,
) ([]tools.SubtaskInfo, error) {
	var (
		subtaskList  tools.SubtaskList
		optAgentType = provider.OptionsTypeGenerator
		msgChainType = database.MsgchainTypeGenerator
	)

	chain := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemGeneratorTmpl),
		llms.TextParts(llms.ChatMessageTypeHuman, userGeneratorTmpl),
	}

	memorist, err := fp.GetMemoristHandler(ctx, &taskID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetTaskSearcherHandler(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.GeneratorExecutorConfig{
		TaskID:   taskID,
		Memorist: memorist,
		Searcher: searcher,
		SubtaskList: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &subtaskList)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal subtask list: %w", err)
			}
			return "subtask list successfully processed", nil
		},
	}
	executor, err := fp.executor.GetGeneratorExecutor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get generator executor: %w", err)
	}

	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal msg chain: %w", err)
	}

	msgChain, err := fp.db.CreateMsgChain(ctx, database.CreateMsgChainParams{
		Type:          msgChainType,
		Model:         fp.Model(optAgentType),
		ModelProvider: string(fp.Type()),
		Chain:         chainBlob,
		FlowID:        fp.flowID,
		TaskID:        database.Int64ToNullInt64(&taskID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create msg chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChain.ID, &taskID, nil, chain, executor, fp.summarizer)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtasks generator result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			input,
			fp.subtasksToMarkdown(subtaskList.Subtasks),
			&taskID,
			nil,
		)
	}

	return subtaskList.Subtasks, nil
}

func (fp *flowProvider) performSubtasksRefiner(
	ctx context.Context,
	taskID int64,
	systemRefinerTmpl, userRefinerTmpl, input string,
) ([]tools.SubtaskInfo, error) {
	var (
		subtaskList  tools.SubtaskList
		chain        []llms.MessageContent
		optAgentType = provider.OptionsTypeRefiner
		msgChainType = database.MsgchainTypeRefiner
	)

	restoreChain := func(msgChain json.RawMessage) ([]llms.MessageContent, error) {
		var msgList []llms.MessageContent
		err := json.Unmarshal(msgChain, &msgList)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal msg chain: %w", err)
		}

		ast, err := cast.NewChainAST(msgList, true)
		if err != nil {
			return nil, fmt.Errorf("failed to create refiner chain ast: %w", err)
		}

		if len(ast.Sections) == 0 {
			return nil, fmt.Errorf("failed to get sections from refiner chain ast")
		}

		systemSection := ast.Sections[0] // there may be multiple sections due to reflector agent
		systemMessage := llms.TextParts(llms.ChatMessageTypeSystem, systemRefinerTmpl)
		systemSection.Header.SystemMessage = &systemMessage
		humanMessage := llms.TextParts(llms.ChatMessageTypeHuman, userRefinerTmpl)
		systemSection.Header.HumanMessage = &humanMessage
		// remove the last report with subtasks list
		for idx := len(systemSection.Body) - 1; idx >= 0; idx-- {
			if systemSection.Body[idx].Type == cast.RequestResponse {
				systemSection.Body = systemSection.Body[:idx]
				break
			}
		}
		// remove all past completions
		for idx := len(systemSection.Body) - 1; idx >= 0; idx-- {
			if systemSection.Body[idx].Type != cast.Completion {
				systemSection.Body = systemSection.Body[:idx+1]
				break
			}
		}

		// restore the chain
		return systemSection.Messages(), nil
	}

	msgChain, err := fp.db.GetFlowTaskTypeLastMsgChain(ctx, database.GetFlowTaskTypeLastMsgChainParams{
		FlowID: fp.flowID,
		TaskID: database.Int64ToNullInt64(&taskID),
		Type:   msgChainType,
	})
	if err != nil || isEmptyChain(msgChain.Chain) {
		// fallback to generator chain if refiner chain is not found or empty
		msgChain, err = fp.db.GetFlowTaskTypeLastMsgChain(ctx, database.GetFlowTaskTypeLastMsgChainParams{
			FlowID: fp.flowID,
			TaskID: database.Int64ToNullInt64(&taskID),
			Type:   database.MsgchainTypeGenerator,
		})
		if err != nil || isEmptyChain(msgChain.Chain) {
			// is unexpected, but we should fallback to empty chain
			chain = []llms.MessageContent{
				llms.TextParts(llms.ChatMessageTypeSystem, systemRefinerTmpl),
				llms.TextParts(llms.ChatMessageTypeHuman, userRefinerTmpl),
			}
		} else {
			if chain, err = restoreChain(msgChain.Chain); err != nil {
				return nil, fmt.Errorf("failed to restore chain from generator state: %w", err)
			}
		}
	} else {
		if chain, err = restoreChain(msgChain.Chain); err != nil {
			return nil, fmt.Errorf("failed to restore chain from refiner state: %w", err)
		}
	}

	memorist, err := fp.GetMemoristHandler(ctx, &taskID, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetTaskSearcherHandler(ctx, taskID)
	if err != nil {
		return nil, fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.GeneratorExecutorConfig{
		TaskID:   taskID,
		Memorist: memorist,
		Searcher: searcher,
		SubtaskList: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &subtaskList)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal subtask list: %w", err)
			}
			return "subtask list successfully processed", nil
		},
	}
	executor, err := fp.executor.GetGeneratorExecutor(cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get generator executor: %w", err)
	}

	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal msg chain: %w", err)
	}

	msgChain, err = fp.db.CreateMsgChain(ctx, database.CreateMsgChainParams{
		Type:          msgChainType,
		Model:         fp.Model(optAgentType),
		ModelProvider: string(fp.Type()),
		Chain:         chainBlob,
		FlowID:        fp.flowID,
		TaskID:        database.Int64ToNullInt64(&taskID),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create msg chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChain.ID, &taskID, nil, chain, executor, fp.summarizer)
	if err != nil {
		return nil, fmt.Errorf("failed to get subtasks refiner result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			input,
			fp.subtasksToMarkdown(subtaskList.Subtasks),
			&taskID,
			nil,
		)
	}

	return subtaskList.Subtasks, nil
}

func (fp *flowProvider) performCoder(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemCoderTmpl, userCoderTmpl, question string,
) (string, error) {
	var (
		codeResult   tools.CodeResult
		optAgentType = provider.OptionsTypeCoder
		msgChainType = database.MsgchainTypeCoder
	)

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	installer, err := fp.GetInstallerHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get installer handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.CoderExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Installer: installer,
		Memorist:  memorist,
		Searcher:  searcher,
		CodeResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &codeResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "code result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetCoderExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get coder executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemCoderTmpl, userCoderTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task coder result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			codeResult.Result,
			taskID,
			subtaskID,
		)
	}

	return codeResult.Result, nil
}

func (fp *flowProvider) performInstaller(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemInstallerTmpl, userInstallerTmpl, question string,
) (string, error) {
	var (
		maintenanceResult tools.MaintenanceResult
		optAgentType      = provider.OptionsTypeInstaller
		msgChainType      = database.MsgchainTypeInstaller
	)

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.InstallerExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Memorist:  memorist,
		Searcher:  searcher,
		MaintenanceResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &maintenanceResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "maintenance result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetInstallerExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get installer executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemInstallerTmpl, userInstallerTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task installer result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			maintenanceResult.Result,
			taskID,
			subtaskID,
		)
	}

	return maintenanceResult.Result, nil
}

func (fp *flowProvider) performMemorist(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemMemoristTmpl, userMemoristTmpl, question string,
) (string, error) {
	var (
		memoristResult tools.MemoristResult
		optAgentType   = provider.OptionsTypeSearcher
		msgChainType   = database.MsgchainTypeMemorist
	)

	cfg := tools.MemoristExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		SearchResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &memoristResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "memorist result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetMemoristExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemMemoristTmpl, userMemoristTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task memorist result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			memoristResult.Result,
			taskID,
			subtaskID,
		)
	}

	return memoristResult.Result, nil
}

func (fp *flowProvider) performPentester(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemPentesterTmpl, userPentesterTmpl, question string,
) (string, error) {
	var (
		hackResult   tools.HackResult
		optAgentType = provider.OptionsTypePentester
		msgChainType = database.MsgchainTypePentester
	)

	// 检查是否启用电力行业Planner
	if fp.isPowerPlannerEnabled() {
		// 尝试使用电力行业专用Planner
		if result, handled := fp.tryPowerPlanner(ctx, taskID, subtaskID, question); handled {
			return result, nil
		}
	}

	adviser, err := fp.GetAskAdviceHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get adviser handler: %w", err)
	}

	coder, err := fp.GetCoderHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get coder handler: %w", err)
	}

	installer, err := fp.GetInstallerHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get installer handler: %w", err)
	}

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.PentesterExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Adviser:   adviser,
		Coder:     coder,
		Installer: installer,
		Memorist:  memorist,
		Searcher:  searcher,
		HackResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &hackResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "hack result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetPentesterExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get pentester executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemPentesterTmpl, userPentesterTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task pentester result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			hackResult.Result,
			taskID,
			subtaskID,
		)
	}

	return hackResult.Result, nil
}

// isPowerPlannerEnabled 检查是否启用电力行业Planner
func (fp *flowProvider) isPowerPlannerEnabled() bool {
	// 暂时返回true作为默认值，实际应该从配置中读取
	return true
}

// tryPowerPlanner 尝试使用电力行业专用Planner
func (fp *flowProvider) tryPowerPlanner(ctx context.Context, taskID, subtaskID *int64, question string) (string, bool) {
	// 创建电力行业Planner
	marketingPlanner := planner.NewMarketingPlanner()
	iguowangPlanner := planner.NewIGuowangPlanner()
	sapPlanner := planner.NewSapPlanner()

	// 不需要创建Router，直接检查Planner匹配

	// 提取目标元数据
	meta := fp.extractTargetMetaFromQuestion(question)

	// 检查是否有Planner匹配
	for _, p := range []router.Planner{marketingPlanner, iguowangPlanner, sapPlanner} {
		if p.Match(meta) {
			// 匹配成功，使用专用Planner
			result, err := fp.executePowerPlan(ctx, taskID, subtaskID, p, question)
			if err != nil {
				// 规划执行失败，返回false让系统使用默认逻辑
				return "", false
			}
			return result, true
		}
	}

	// 没有匹配的Planner
	return "", false
}

// extractTargetMetaFromQuestion 从问题中提取目标元数据
func (fp *flowProvider) extractTargetMetaFromQuestion(question string) router.TargetMeta {
	meta := router.TargetMeta{}

	questionLower := strings.ToLower(question)

	// 提取域名和端口
	if strings.Contains(questionLower, "http://") || strings.Contains(questionLower, "https://") {
		// 简单的URL解析
		if idx := strings.Index(questionLower, "://"); idx != -1 {
			urlPart := questionLower[idx+3:]
			if colonIdx := strings.Index(urlPart, ":"); colonIdx != -1 {
				meta.Domain = urlPart[:colonIdx]
				// 提取端口
				if slashIdx := strings.Index(urlPart[colonIdx+1:], "/"); slashIdx != -1 {
					portStr := urlPart[colonIdx+1 : colonIdx+1+slashIdx]
					if port, err := strconv.Atoi(portStr); err == nil {
						meta.Port = port
					}
				}
			} else if slashIdx := strings.Index(urlPart, "/"); slashIdx != -1 {
				meta.Domain = urlPart[:slashIdx]
			} else {
				meta.Domain = urlPart
			}
		}
	}

	// 提取端口信息
	if strings.Contains(questionLower, ":8080") {
		meta.Port = 8080
	} else if strings.Contains(questionLower, ":8443") {
		meta.Port = 8443
	}

	// 识别应用类型
	if strings.Contains(questionLower, "营销") || strings.Contains(questionLower, "marketing") || strings.Contains(questionLower, "billing") {
		meta.AppID = "marketing"
	} else if strings.Contains(questionLower, "国网") || strings.Contains(questionLower, "iguowang") {
		meta.AppID = "iguowang"
	} else if strings.Contains(questionLower, "erp") || strings.Contains(questionLower, "sap") {
		meta.AppID = "erp"
	}

	return meta
}

// executePowerPlan 执行电力行业专用测试计划
func (fp *flowProvider) executePowerPlan(ctx context.Context, taskID, subtaskID *int64, planner router.Planner, question string) (string, error) {
	// 简化的输入结构
	agentInput := map[string]interface{}{
		"question": question,
		"target":   question, // 简化处理
	}

	// 生成测试计划
	steps, err := planner.Plan(ctx, agentInput)
	if err != nil {
		return "", fmt.Errorf("failed to generate power plan: %w", err)
	}

	var finalResult strings.Builder
	finalResult.WriteString(fmt.Sprintf("🎯 使用电力行业专用测试计划: %s\n\n", planner.Name()))

	// 执行每个步骤
	for i, step := range steps {
		stepResult, err := fp.executeStep(ctx, taskID, subtaskID, step, i+1)
		if err != nil {
			finalResult.WriteString(fmt.Sprintf("步骤 %d 执行失败: %v\n", i+1, err))
			continue
		}
		finalResult.WriteString(stepResult)
		finalResult.WriteString("\n\n")
	}

	return finalResult.String(), nil
}

// executeStep 执行单个测试步骤
func (fp *flowProvider) executeStep(ctx context.Context, taskID, subtaskID *int64, step interface{}, stepNum int) (string, error) {
	// 将step转换为map以提取内容
	stepMap, ok := step.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("invalid step format")
	}

	content, ok := stepMap["content"].(string)
	if !ok {
		return "", fmt.Errorf("step content not found")
	}

	stepType, _ := stepMap["type"].(string)

	// 构建结果
	result := fmt.Sprintf("📋 步骤 %d: %s\n\n%s\n\n✅ 电力行业专用测试计划已生成，请使用相应的安全工具执行测试。",
		stepNum, stepType, content)

	return result, nil
}

func (fp *flowProvider) performSearcher(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemSearcherTmpl, userSearcherTmpl, question string,
) (string, error) {
	var (
		searchResult tools.SearchResult
		optAgentType = provider.OptionsTypeSearcher
		msgChainType = database.MsgchainTypeSearcher
	)

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	cfg := tools.SearcherExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Memorist:  memorist,
		SearchResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &searchResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "search result successfully processed", nil
		},
		Summarizer: fp.GetSummarizeResultHandler(taskID, subtaskID),
	}
	executor, err := fp.executor.GetSearcherExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemSearcherTmpl, userSearcherTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task searcher result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			searchResult.Result,
			taskID,
			subtaskID,
		)
	}

	return searchResult.Result, nil
}

func (fp *flowProvider) performEnricher(
	ctx context.Context,
	taskID, subtaskID *int64,
	systemEnricherTmpl, userEnricherTmpl, question string,
) (string, error) {
	var (
		enricherResult tools.EnricherResult
		optAgentType   = provider.OptionsTypeEnricher
		msgChainType   = database.MsgchainTypeEnricher
	)

	memorist, err := fp.GetMemoristHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get memorist handler: %w", err)
	}

	searcher, err := fp.GetSubtaskSearcherHandler(ctx, taskID, subtaskID)
	if err != nil {
		return "", fmt.Errorf("failed to get searcher handler: %w", err)
	}

	cfg := tools.EnricherExecutorConfig{
		TaskID:    taskID,
		SubtaskID: subtaskID,
		Memorist:  memorist,
		Searcher:  searcher,
		EnricherResult: func(ctx context.Context, name string, args json.RawMessage) (string, error) {
			err := json.Unmarshal(args, &enricherResult)
			if err != nil {
				return "", fmt.Errorf("failed to unmarshal result: %w", err)
			}
			return "enrich result successfully processed", nil
		},
	}
	executor, err := fp.executor.GetEnricherExecutor(cfg)
	if err != nil {
		return "", fmt.Errorf("failed to get enricher executor: %w", err)
	}

	msgChainID, chain, err := fp.restoreChain(
		ctx, taskID, subtaskID, optAgentType, msgChainType, systemEnricherTmpl, userEnricherTmpl,
	)
	if err != nil {
		return "", fmt.Errorf("failed to restore chain: %w", err)
	}

	ctx = tools.PutAgentContext(ctx, msgChainType)
	err = fp.performAgentChain(ctx, optAgentType, msgChainID, taskID, subtaskID, chain, executor, fp.summarizer)
	if err != nil {
		return "", fmt.Errorf("failed to get task enricher result: %w", err)
	}

	if agentCtx, ok := tools.GetAgentContext(ctx); ok {
		fp.agentLog.PutLog(
			ctx,
			agentCtx.ParentAgentType,
			agentCtx.CurrentAgentType,
			question,
			enricherResult.Result,
			taskID,
			subtaskID,
		)
	}

	return enricherResult.Result, nil
}

func (fp *flowProvider) performSimpleChain(
	ctx context.Context,
	taskID, subtaskID *int64,
	opt provider.ProviderOptionsType,
	msgChainType database.MsgchainType,
	systemTmpl, userTmpl string,
) (string, error) {
	var (
		resp *llms.ContentResponse
		err  error
	)

	chain := []llms.MessageContent{
		llms.TextParts(llms.ChatMessageTypeSystem, systemTmpl),
		llms.TextParts(llms.ChatMessageTypeHuman, userTmpl),
	}

	for idx := 0; idx <= maxRetriesToCallSimpleChain; idx++ {
		if idx == maxRetriesToCallSimpleChain {
			return "", fmt.Errorf("failed to call simple chain: %w", err)
		}

		resp, err = fp.CallEx(ctx, opt, chain, nil)
		if err == nil {
			break
		} else {
			if errors.Is(err, context.Canceled) {
				return "", err
			}

			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(time.Second * 5):
			default:
			}
		}
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	var parts []string
	var inputTokens, outputTokens int64
	for _, choice := range resp.Choices {
		parts = append(parts, choice.Content)
		inputTokens, outputTokens = fp.GetUsage(choice.GenerationInfo)
	}
	chain = append(chain, llms.TextParts(llms.ChatMessageTypeAI, parts...))

	chainBlob, err := json.Marshal(chain)
	if err != nil {
		return "", fmt.Errorf("failed to marshal summarizer msg chain: %w", err)
	}

	_, err = fp.db.CreateMsgChain(ctx, database.CreateMsgChainParams{
		Type:          msgChainType,
		Model:         fp.Model(opt),
		ModelProvider: string(fp.Type()),
		UsageIn:       inputTokens,
		UsageOut:      outputTokens,
		Chain:         chainBlob,
		FlowID:        fp.flowID,
		TaskID:        database.Int64ToNullInt64(taskID),
		SubtaskID:     database.Int64ToNullInt64(subtaskID),
	})

	return strings.Join(parts, "\n\n"), nil
}
