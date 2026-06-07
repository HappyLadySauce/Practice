package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/adk"
	"github.com/cloudwego/eino/components/tool"
    "github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/callbacks"
	"github.com/cloudwego/eino/compose"
)

func GetModelCallbacksHandler() callbacks.Handler {
	return callbacks.NewHandlerBuilder().
		OnStartFn(func(ctx context.Context, info *callbacks.RunInfo, input callbacks.CallbackInput) context.Context {
			// 只关注 Agent 组件的回调
			if info.Component == adk.ComponentOfAgent{
				agentInput := adk.ConvAgentCallbackInput(input)
				if agentInput != nil {
					fmt.Printf("[%s] 🚀 Agent [%s] 开始执行，输入消息数: %d\n",
                	time.Now().Format("15:04:05"), info.Name, len(agentInput.Input.Messages))
				} else {
                	fmt.Printf("[%s] 🔄 Agent [%s] 从中断恢复执行\n",
                	time.Now().Format("15:04:05"), info.Name)
				}
			}
			return ctx
		}).
		OnEndFn(func(ctx context.Context, info *callbacks.RunInfo, output callbacks.CallbackOutput) context.Context {
			if info.Component == adk.ComponentOfAgent {
        		fmt.Printf("[%s] ✅ Agent [%s] 执行完成\n",
            	time.Now().Format("15:04:05"), info.Name)
			}
			return ctx
		}).Build()
}


type memCheckPointStore struct {
	m map[string][]byte
}

func (m *memCheckPointStore) Get(ctx context.Context, id string) ([]byte, bool, error) {
	v, ok := m.m[id]
	return v, ok, nil
}

func (m *memCheckPointStore) Set(_ context.Context, id string, data []byte) error {
    m.m[id] = data
    return nil
}

// 删除文件工具的参数
type DeleteFileParams struct {
    FilePath string `json:"file_path" jsonschema:"description=要删除的文件路径"`
}

func deleteFile(ctx context.Context, params *DeleteFileParams) (string, error) {
	// 检查是否处于中断恢复流程
	wasInterrupted, _, _ := compose.GetInterruptState[string](ctx)
	if !wasInterrupted {
		// 首次调用，中断等待用户确认
		return "", compose.Interrupt(ctx, fmt.Sprintf("即将删除文件: %s，是否确认？", params.FilePath))
	}

	// 从中断恢复，检查用户的确认结果
	isTarget, hasData, approval := compose.GetResumeContext[string](ctx)
	if !isTarget || !hasData {
		// 不是恢复目标，重新中断
		return "", compose.Interrupt(ctx, fmt.Sprintf("即将删除文件: %s，是否确认？", params.FilePath))
	}

	if approval == "yes" {
		return fmt.Sprintf("文件 %s 已成功删除", params.FilePath), nil
	}
	return fmt.Sprintf("用户取消了删除操作，文件 %s 未被删除", params.FilePath), nil
}

const MODEL = "gemma-4-e4b-it-uncensored"
const APIURL = "http://127.0.0.1:11434/v1"

func main() {
	ctx := context.Background()

	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
		APIKey: "",
		Model: MODEL,
		BaseURL: APIURL,
	})
	if err != nil {
		log.Fatal("创建 ChatModel 失败: %w", err)
	}

    // 创建删除文件工具
    deleteTool, err := utils.InferTool(
       "delete_file",
       "删除指定路径的文件，执行前需要用户确认",
       deleteFile,
    )
    if err != nil {
       log.Fatal(err)
    }

// 	// 技术方案评审系统

// 	// ===== 第一阶段：三个并行分析师 =====

// 	perfAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
// 		Name: "perf_reviewer",
// 		Description: "性能评审专家",
// 		Instruction: `你是性能评审专家。请从性能维度评审用户提出的技术方案，输出格式：
// 【性能评分】1-10分
// 【核心发现】列出2-3个关键点
// 【改进建议】如有性能风险，给出具体建议`,
// 		Model: chatModel,
// 	})

// 	secAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
// 		Name:        "sec_reviewer",
// 		Description: "安全评审专家",
// 		Instruction: `你是安全评审专家。请从安全维度评审用户提出的技术方案，输出格式：
// 【安全评分】1-10分
// 【核心发现】列出2-3个关键点
// 【改进建议】如有安全风险，给出具体建议`,
// 		Model: chatModel,
// 	})

// 	costAgent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
// 		Name:        "cost_reviewer",
// 		Description: "成本评审专家",
// 		Instruction: `你是成本评审专家。请从成本维度评审用户提出的技术方案，输出格式：
// 【成本评分】1-10分
// 【核心发现】列出2-3个关键点
// 【改进建议】如有成本优化空间，给出具体建议`,
// 		Model: chatModel,
// 	})


// 	parallelReview, _ := adk.NewParallelAgent(ctx, &adk.ParallelAgentConfig{
// 		Name: "parallel_review",
// 		Description: "并行技术评审",
// 		SubAgents: []adk.Agent{perfAgent, secAgent, costAgent},
// 	})


// 	// ===== 第二阶段：综合汇总 =====

// 	summarizer, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
// 		Name: "summarizer",
// 		Description: "评审报告汇总专家",
// 		Instruction: `你是技术评审委员会主席。你会收到三位专家（性能、安全、成本）的独立评审意见。
// 请综合所有评审意见，生成一份结构化的最终评审报告，格式如下：

// # 技术方案评审报告

// ## 综合评分
// （三个维度的加权平均，权重：性能40%、安全35%、成本25%）

// ## 各维度概要
// （简要汇总每个维度的核心发现）

// ## 关键风险
// （列出最需要关注的风险项）

// ## 最终结论
// （通过/有条件通过/不通过，并说明理由）

// ## 行动项
// （按优先级列出需要落实的改进项）`,
// 		Model: chatModel,
// 	})

// 	// ===== 组装完整流水线 =====

// 	fullPipeline, _ := adk.NewSequentialAgent(ctx, &adk.SequentialAgentConfig{
// 		Name: "tech_review_system",
// 		Description: "完整技术方案评审系统",
// 		SubAgents: []adk.Agent{parallelReview, summarizer},
// 	})

// 	// ===== 运行 =====

// 	runner := adk.NewRunner(ctx, adk.RunnerConfig{Agent: fullPipeline})

// 	proposal := `技术方案：将公司核心交易系统从单体架构迁移到微服务架构

// 关键设计决策：
// 1. 使用 Kubernetes 作为容器编排平台
// 2. 服务间通信采用 gRPC + Protobuf
// 3. 数据库从单一 MySQL 拆分为每个服务独立的数据库（Database per Service）
// 4. 引入 Apache Kafka 作为异步消息队列
// 5. 使用 Istio 服务网格处理流量治理
// 6. API Gateway 使用 Kong

// 预计影响范围：核心交易链路、用户中心、库存管理、支付系统`

// 	iter := runner.Query(ctx, proposal, adk.WithCallbacks(GetModelCallbacksHandler()))
	
// 	for {
// 		event, ok := iter.Next()
// 		if !ok {
// 			break
// 		}
// 		if event.Err != nil {
// 			log.Fatal(event.Err)
// 		}
// 		if event.Output != nil && event.Output.MessageOutput != nil {
// 			fmt.Printf("\n========== [%s] ==========\n%s\n",
// 				event.AgentName,
// 				event.Output.MessageOutput.Message.Content)
// 		}
// 	}

	agent, _ := adk.NewChatModelAgent(ctx, &adk.ChatModelAgentConfig{
		Name: "file_manager",
		Description: "文件管理助手",
		Instruction: "你是一个文件管理助手。当用户要求删除文件时，使用 delete_file 工具执行操作。",
		Model: chatModel,
		ToolsConfig: adk.ToolsConfig{
			ToolsNodeConfig: compose.ToolsNodeConfig{
				Tools: []tool.BaseTool{deleteTool},
			},
		},
	})

	// 创建 Runner，必须配置 CheckPointStore 才能使用中断恢复
	store := &memCheckPointStore{m: make(map[string][]byte)}
	runner := adk.NewRunner(ctx, adk.RunnerConfig{
		Agent: agent,
		CheckPointStore: store,
	})

	// 第一次执行：Agent 会调用 delete_file 工具，触发中断
	checkPointID := "session-001"
	iter := runner.Query(ctx, "请帮我删除 /tmp/old_logs.txt 这个文件",
       adk.WithCheckPointID(checkPointID))

	var interruptID string
	for {
		event, ok := iter.Next()
		if !ok {
			break
		}
		if event.Err != nil {
			log.Fatal(err)
		}

		// 检测中断事件
		if event.Action != nil && event.Action.Interrupted != nil {
			interruptID = event.Action.Interrupted.InterruptContexts[0].ID
			fmt.Println("⏸️  Agent 中断，等待用户确认...")
			fmt.Printf("   中断信息: %v\n", event.Action.Interrupted.InterruptContexts[0].Info)
			fmt.Printf("   中断ID: %s\n", interruptID)
		}
		if event.Output != nil && event.Output.MessageOutput != nil {
			fmt.Printf("[%s] %s\n", event.AgentName, event.Output.MessageOutput.Message.Content)
		}

		// 模拟用户确认
		fmt.Println("\n--- 用户确认: yes ---\n")

		// 恢复执行，传入用户的确认结果
		resumeIter, err := runner.ResumeWithParams(ctx, checkPointID, &adk.ResumeParams{
		Targets: map[string]any{
			interruptID: "yes",
		}})
		if err != nil {
			log.Fatal(err)
		}

		for {
			event, ok := resumeIter.Next()
			if !ok {
				break
			}
			if event.Err != nil {
				log.Fatal(event.Err)
			}
			if event.Output != nil && event.Output.MessageOutput != nil {
				fmt.Printf("[%s] %s\n", event.AgentName, event.Output.MessageOutput.Message.Content)
			}
		}
	}
}