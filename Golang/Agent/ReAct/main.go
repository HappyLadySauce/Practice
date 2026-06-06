package main

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

const ModelURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

// ReActAgent implements a ReAct-style agent using prompt-based tool invocation.
// ReActAgent 使用基于 Prompt 的工具调用实现 ReAct 架构 Agent
type ReActAgent struct {
	client *openai.Client
}

// ReActStep records Thought-Action-Observation for each step.
// ReActStep 记录每一步的 Thought-Action-Observation
type ReActStep struct {
	Step        int
	Thought     string
	Action      string
	ActionInput string
	Observation string
}

func NewReActAgent(client *openai.Client) *ReActAgent {
	return &ReActAgent{client: client}
}

// buildSystemPrompt embeds tool specs in the system message for models without native tool support.
// buildSystemPrompt 将工具说明写入 system 消息，兼容不支持原生 tool calling 的本地模型
func buildSystemPrompt() string {
	return `你是一个使用 ReAct 模式的 AI 助手。请严格按以下格式回复：

【思考】
分析当前情况，说明下一步计划

如果需要调用工具，继续输出：
Action: 工具名称
Action Input: {"参数名":"参数值"}

如果已有足够信息回答用户，输出：
Final Answer: 你的最终回答

可用工具：
1. get_weather(city) - 查询城市天气，参数: {"city":"城市名"}
2. get_time() - 获取当前时间，参数: {}
3. clothing_advice(temperature, weather) - 穿衣建议，参数: {"temperature":25,"weather":"晴"}

规则：
- 每次只调用一个工具
- Action Input 必须是合法 JSON
- 收到 Observation 后继续思考，直到可以给出 Final Answer`
}

// Run executes the ReAct loop up to maxSteps iterations.
// Run 执行 ReAct 循环，最多 maxSteps 步
func (agent *ReActAgent) Run(goal string, maxSteps int) (string, []ReActStep) {
	var steps []ReActStep

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: buildSystemPrompt(),
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: goal,
		},
	}

	for step := 1; step <= maxSteps; step++ {
		resp, err := agent.client.CreateChatCompletion(
			context.Background(),
			openai.ChatCompletionRequest{
				Model:       model,
				Messages:    messages,
				Temperature: 0.2,
			},
		)
		if err != nil {
			fmt.Printf("Request failed: %v\n", err)
			break
		}

		if len(resp.Choices) == 0 {
			fmt.Println("Request failed: empty response choices")
			break
		}

		content := strings.TrimSpace(resp.Choices[0].Message.Content)
		if content == "" {
			fmt.Println("Request failed: empty model content")
			break
		}

		record := ReActStep{Step: step, Thought: content}
		fmt.Printf("\n--- Step %d ---\n", step)
		fmt.Printf("Model output:\n%s\n", content)

		if finalAnswer := extractFinalAnswer(content); finalAnswer != "" {
			record.Thought = content
			steps = append(steps, record)
			return finalAnswer, steps
		}

		action, actionInput := parseAction(content)
		if action == "" {
			// Treat plain text as final answer when model skips the Final Answer prefix.
			// 模型未使用 Final Answer 前缀时，将纯文本视为最终回答
			steps = append(steps, record)
			return content, steps
		}

		record.Action = action
		record.ActionInput = actionInput
		fmt.Printf("Action: %s(%s)\n", action, actionInput)

		observation := agent.executeTool(action, actionInput)
		record.Observation = observation
		fmt.Printf("Observation: %s\n", observation)
		steps = append(steps, record)

		messages = append(messages,
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: content,
			},
			openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Observation: %s\n请继续思考并决定下一步。", observation),
			},
		)
	}

	return "达到最大步数限制", steps
}

var (
	actionLinePattern = regexp.MustCompile(`(?mi)^Action:\s*(.+?)\s*$`)
	actionInputPattern = regexp.MustCompile(`(?mi)^Action Input:\s*(.+?)\s*$`)
	finalAnswerPattern = regexp.MustCompile(`(?mis)^Final Answer:\s*(.+)$`)
)

// parseAction extracts tool name and JSON input from model output.
// parseAction 从模型输出中解析工具名与 JSON 参数
func parseAction(content string) (string, string) {
	actionMatch := actionLinePattern.FindStringSubmatch(content)
	if len(actionMatch) < 2 {
		return "", ""
	}

	action := strings.TrimSpace(actionMatch[1])
	actionInput := "{}"
	if inputMatch := actionInputPattern.FindStringSubmatch(content); len(inputMatch) >= 2 {
		actionInput = strings.TrimSpace(inputMatch[1])
	}

	return action, actionInput
}

// extractFinalAnswer returns text after the Final Answer marker when present.
// extractFinalAnswer 提取 Final Answer 标记后的文本
func extractFinalAnswer(content string) string {
	match := finalAnswerPattern.FindStringSubmatch(content)
	if len(match) < 2 {
		return ""
	}
	return strings.TrimSpace(match[1])
}

// executeTool dispatches to the registered mock tools.
// executeTool 分发到已注册的模拟工具
func (a *ReActAgent) executeTool(name, args string) string {
	switch strings.TrimSpace(name) {
	case "get_weather":
		var p struct {
			City string `json:"city"`
		}
		if err := json.Unmarshal([]byte(args), &p); err != nil {
			return fmt.Sprintf("invalid Action Input JSON: %v", err)
		}
		weatherData := map[string]string{
			"北京": "晴，25°C，湿度40%，东北风3级",
			"上海": "多云，28°C，湿度65%，东南风2级",
			"广州": "小雨，30°C，湿度80%，南风4级",
		}
		if w, ok := weatherData[p.City]; ok {
			return fmt.Sprintf("%s天气: %s", p.City, w)
		}
		return fmt.Sprintf("未找到%s的天气数据", p.City)

	case "get_time":
		return fmt.Sprintf("当前时间: %s", time.Now().Format("2006-01-02 15:04:05"))

	case "clothing_advice":
		var p struct {
			Temperature int    `json:"temperature"`
			Weather     string `json:"weather"`
		}
		if err := json.Unmarshal([]byte(args), &p); err != nil {
			return fmt.Sprintf("invalid Action Input JSON: %v", err)
		}
		if p.Temperature > 28 {
			return "建议穿短袖、短裤，注意防晒"
		}
		if p.Temperature > 20 {
			return "建议穿薄长袖或T恤，可带一件薄外套"
		}
		return "建议穿厚外套或毛衣，注意保暖"

	default:
		return fmt.Sprintf("未知工具: %s", name)
	}
}

func main() {
	config := openai.DefaultConfig("")
	config.BaseURL = ModelURL
	client := openai.NewClientWithConfig(config)

	agent := NewReActAgent(client)
	answer, steps := agent.Run("今天北京天气怎么样？需要带伞吗？穿什么合适？", 10)

	fmt.Println("\n========== 执行摘要 ==========")
	fmt.Printf("共执行 %d 步\n", len(steps))
	for _, s := range steps {
		fmt.Printf("Step %d: ", s.Step)
		if s.Action != "" {
			fmt.Printf("调用了 %s\n", s.Action)
		} else {
			fmt.Printf("生成了最终回答\n")
		}
	}
	fmt.Printf("\n最终回答:\n%s\n", cleanAnswer(answer))
}

// cleanAnswer strips optional reasoning prefixes from the final response.
// cleanAnswer 去除最终回答中可能的思考前缀
func cleanAnswer(s string) string {
	if idx := strings.LastIndex(s, "】"); idx != -1 && idx < len(s)-3 {
		return strings.TrimSpace(s[idx+len("】"):])
	}
	return s
}
