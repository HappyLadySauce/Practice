package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"test/tools"

	"github.com/sashabaranov/go-openai"
)

const (
	defaultBaseURL = "http://127.0.0.1:11434/v1"
	model          = "gemma-4-e4b-it-uncensored"
	maxToolRounds  = 5
)

type streamResult struct {
	content      string
	toolCalls    []openai.ToolCall
	finishReason openai.FinishReason
}

// consumeStream reads SSE chunks, prints assistant text, and accumulates tool calls.
// consumeStream 读取 SSE 分片，打印助手文本，并累积工具调用信息。
func consumeStream(stream *openai.ChatCompletionStream) (streamResult, error) {
	var result streamResult

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		if len(response.Choices) == 0 {
			continue
		}

		choice := response.Choices[0]
		if choice.FinishReason != "" {
			result.finishReason = choice.FinishReason
		}

		for _, tc := range choice.Delta.ToolCalls {
			result.toolCalls = tools.MergeToolCallDelta(result.toolCalls, tc)
		}

		chunk := choice.Delta.Content
		fmt.Print(chunk)
		result.content += chunk
	}
}

func main() {
	cfg := openai.DefaultConfig("")
	cfg.BaseURL = defaultBaseURL
	client := openai.NewClientWithConfig(cfg)

	messages := []openai.ChatCompletionMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: "你是一位友好的 Go 语言助手。",
		},
	}

	availableTools := tools.AllTools()
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话 - 流式模式（输入 quit 退出）:")

	for {
		fmt.Print("用户: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			break
		}

		messages = append(messages, openai.ChatCompletionMessage{
			Role:    openai.ChatMessageRoleUser,
			Content: input,
		})

		for round := 0; round < maxToolRounds; round++ {
			stream, err := client.CreateChatCompletionStream(
				context.Background(),
				openai.ChatCompletionRequest{
					Model:       model,
					Temperature: 0.1,
					Messages:    messages,
					Tools:       availableTools,
				},
			)
			if err != nil {
				log.Printf("create stream failed: %v", err)
				break
			}

			fmt.Print("助手: ")
			result, err := consumeStream(stream)
			stream.Close()
			fmt.Println()
			if err != nil {
				log.Printf("stream receive failed: %v", err)
				break
			}

			if result.finishReason == openai.FinishReasonToolCalls && len(result.toolCalls) > 0 {
				messages = append(messages, openai.ChatCompletionMessage{
					Role:      openai.ChatMessageRoleAssistant,
					ToolCalls: result.toolCalls,
				})

				for _, tc := range result.toolCalls {
					toolResult, err := tools.Dispatch(tc)
					if err != nil {
						log.Printf("tool dispatch failed: tool=%s err=%v", tc.Function.Name, err)
						toolResult = fmt.Sprintf(`{"error":%q}`, err.Error())
					}

					messages = append(messages, openai.ChatCompletionMessage{
						Role:       openai.ChatMessageRoleTool,
						ToolCallID: tc.ID,
						Name:       tc.Function.Name,
						Content:    toolResult,
					})
				}
				continue
			}

			messages = append(messages, openai.ChatCompletionMessage{
				Role:    openai.ChatMessageRoleAssistant,
				Content: result.content,
			})
			break
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read stdin failed: %v", err)
	}
}
