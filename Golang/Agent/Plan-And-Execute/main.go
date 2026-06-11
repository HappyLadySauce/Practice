package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

const ModelURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

// TaskStep defines one executable step in the plan.
// TaskStep 定义计划中的一个可执行步骤
type TaskStep struct {
	StepNumber  int    `json:"step_number"`
	Description string `json:"description"`
	ToolNeeded  string `json:"tool_needed"`
	DependsOn   []int  `json:"depends_on"`
	Output      string `json:"expected_output"`
}

// TaskPlan holds the decomposed goal and ordered steps.
// TaskPlan 保存分解后的目标与步骤列表
type TaskPlan struct {
	Goal  string     `json:"goal"`
	Steps []TaskStep `json:"steps"`
}

// extractJSON strips markdown fences and returns the first JSON object/array payload.
// extractJSON 去除 Markdown 代码块并提取首个 JSON 对象/数组
func extractJSON(content string) string {
	content = strings.TrimSpace(content)

	fencePattern := regexp.MustCompile("(?s)```(?:json)?\\s*(.*?)\\s*```")
	if match := fencePattern.FindStringSubmatch(content); len(match) >= 2 {
		return strings.TrimSpace(match[1])
	}

	startObj := strings.Index(content, "{")
	startArr := strings.Index(content, "[")
	start := -1
	switch {
	case startObj >= 0 && startArr >= 0:
		if startObj < startArr {
			start = startObj
		} else {
			start = startArr
		}
	case startObj >= 0:
		start = startObj
	case startArr >= 0:
		start = startArr
	}
	if start >= 0 {
		return strings.TrimSpace(content[start:])
	}

	return content
}

func decomposeTask(task string) (*TaskPlan, error) {
	config := openai.DefaultConfig("")
	config.BaseURL = ModelURL
	client := openai.NewClientWithConfig(config)

	systemPrompt := `你是一个任务规划专家。用户会给你一个复杂任务，你需要将其分解为具体的执行步骤。
只返回 JSON，不要添加 Markdown 代码块或额外解释。结构如下：
{
  "goal": "任务目标",
  "steps": [
    {
      "step_number": 1,
      "description": "步骤描述",
      "tool_needed": "需要的工具（如：web_search, database_query, text_generation, code_execution, none）",
      "depends_on": [],
      "expected_output": "这一步的预期输出"
    }
  ]
}
注意：depends_on 表示该步骤依赖哪些前置步骤编号；无依赖时使用空数组 []。`

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: model,
			Messages: []openai.ChatCompletionMessage{
				{Role: openai.ChatMessageRoleSystem, Content: systemPrompt},
				{Role: openai.ChatMessageRoleUser, Content: "请分解这个任务: " + task},
			},
			Temperature: 0.1,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("调用失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("调用失败: empty response choices")
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)
	jsonPayload := extractJSON(content)

	var plan TaskPlan
	if err := json.Unmarshal([]byte(jsonPayload), &plan); err != nil {
		return nil, fmt.Errorf("解析任务失败: %w\nraw content:\n%s", err, content)
	}

	return &plan, nil
}

func main() {
	plan, err := decomposeTask("帮我调研Go语言主流Web框架，对比它们的性能和生态，最终输出一份推荐报告")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("任务目标：%s\n\n", plan.Goal)
	for _, step := range plan.Steps {
		deps := "无"
		if len(step.DependsOn) > 0 {
			depsJSON, _ := json.Marshal(step.DependsOn)
			deps = string(depsJSON)
		}
		fmt.Printf("步骤%d：%s\n", step.StepNumber, step.Description)
		fmt.Printf("  工具：%s\n", step.ToolNeeded)
		fmt.Printf("  依赖：%s\n", deps)
		fmt.Printf("  预期输出：%s\n\n", step.Output)
	}
}
