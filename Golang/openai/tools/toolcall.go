package tools

import (
	"encoding/json"
	"fmt"

	"github.com/sashabaranov/go-openai"
)

// AllTools returns every tool exposed to the model.
// AllTools 返回暴露给模型的全部工具定义。
func AllTools() []openai.Tool {
	return []openai.Tool{GetWeatherTool()}
}

// MergeToolCallDelta merges one streamed tool-call fragment into the accumulator.
// MergeToolCallDelta 将流式 tool call 分片合并到累积结果中。
func MergeToolCallDelta(existing []openai.ToolCall, delta openai.ToolCall) []openai.ToolCall {
	if delta.Index == nil {
		return existing
	}

	idx := *delta.Index
	for len(existing) <= idx {
		existing = append(existing, openai.ToolCall{})
	}

	if delta.ID != "" {
		existing[idx].ID = delta.ID
	}
	if delta.Type != "" {
		existing[idx].Type = delta.Type
	}
	if delta.Function.Name != "" {
		existing[idx].Function.Name = delta.Function.Name
	}
	if delta.Function.Arguments != "" {
		existing[idx].Function.Arguments += delta.Function.Arguments
	}

	return existing
}

// Dispatch executes a tool call and returns the result as a JSON string.
// Dispatch 执行工具调用并返回 JSON 字符串结果。
func Dispatch(tc openai.ToolCall) (string, error) {
	switch tc.Function.Name {
	case "get_weather":
		var args struct {
			City string `json:"city"`
		}
		if err := json.Unmarshal([]byte(tc.Function.Arguments), &args); err != nil {
			return "", fmt.Errorf("parse get_weather args: %w", err)
		}
		return GetWeather(args.City), nil
	default:
		return "", fmt.Errorf("unknown tool: %s", tc.Function.Name)
	}
}
