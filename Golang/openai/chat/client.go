package chat

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/sashabaranov/go-openai"
	"test/tools"
)

var ErrMaxToolRounds = errors.New("max tool rounds exceeded")

// Client wraps OpenAI-compatible chat with retry, timeout, streaming, and tools.
// Client 封装带重试、超时、流式输出和工具调用的 OpenAI 兼容客户端。
type Client struct {
	cfg    Config
	client *openai.Client
	tools  []openai.Tool
}

// New creates a chat client from config.
// New 根据配置创建聊天客户端。
func New(cfg Config) *Client {
	openaiCfg := openai.DefaultConfig(cfg.Token)
	openaiCfg.BaseURL = cfg.BaseURL
	openaiCfg.HTTPClient = &http.Client{Timeout: cfg.HTTPTimeout}

	return &Client{
		cfg:    cfg,
		client: openai.NewClientWithConfig(openaiCfg),
		tools:  tools.AllTools(),
	}
}

// Chat performs a single non-streaming completion with retry.
// Chat 执行一次非流式补全，失败时自动重试。
func (c *Client) Chat(ctx context.Context, messages []openai.ChatCompletionMessage) (string, error) {
	resp, err := withRetry(c.cfg, ctx, func(reqCtx context.Context) (openai.ChatCompletionResponse, error) {
		return c.client.CreateChatCompletion(reqCtx, openai.ChatCompletionRequest{
			Model:       c.cfg.Model,
			Temperature: c.cfg.Temperature,
			TopP:        c.cfg.TopP,
			Messages:    messages,
			Tools:       c.tools,
		})
	})
	if err != nil {
		return "", err
	}
	if len(resp.Choices) == 0 {
		return "", errors.New("empty response choices")
	}
	return resp.Choices[0].Message.Content, nil
}

// RunTurn executes one user turn: stream reply, handle tool calls, update history.
// RunTurn 执行一轮用户对话：流式回复、处理工具调用并更新历史。
func (c *Client) RunTurn(ctx context.Context, messages []openai.ChatCompletionMessage, w io.Writer) ([]openai.ChatCompletionMessage, error) {
	for round := 0; round < c.cfg.MaxToolRounds; round++ {
		result, err := c.streamOnce(ctx, messages, w)
		if err != nil {
			return messages, err
		}

		if result.FinishReason == openai.FinishReasonToolCalls && len(result.ToolCalls) > 0 {
			messages = append(messages, openai.ChatCompletionMessage{
				Role:      openai.ChatMessageRoleAssistant,
				ToolCalls: result.ToolCalls,
			})

			for _, tc := range result.ToolCalls {
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
			Content: result.Content,
		})
		return messages, nil
	}

	return messages, ErrMaxToolRounds
}

// streamOnce creates a stream and consumes it with retry.
// streamOnce 创建并消费流式响应，失败时自动重试。
func (c *Client) streamOnce(ctx context.Context, messages []openai.ChatCompletionMessage, w io.Writer) (StreamResult, error) {
	return withRetry(c.cfg, ctx, func(reqCtx context.Context) (StreamResult, error) {
		stream, err := c.client.CreateChatCompletionStream(reqCtx, openai.ChatCompletionRequest{
			Model:       c.cfg.Model,
			Temperature: c.cfg.Temperature,
			TopP:        c.cfg.TopP,
			Messages:    messages,
			Tools:       c.tools,
		})
		if err != nil {
			return StreamResult{}, fmt.Errorf("create stream: %w", err)
		}

		result, err := ConsumeStream(stream, w)
		stream.Close()
		if err != nil {
			return StreamResult{}, fmt.Errorf("consume stream: %w", err)
		}
		return result, nil
	})
}
