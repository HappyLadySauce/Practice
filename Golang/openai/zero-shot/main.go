package main

import (
	"context"
	"sync"
	"fmt"
	"log"

	"github.com/sashabaranov/go-openai"
)

// 最简单的Prompt 形式叫 Zero-shot Prompting不给任何示例，直接让模型完成任务。
// “Zero-shot”的意思就是"零样本"，模型完全依靠自身的预训练知识来理解和执行你的指令。

const defaultBaseURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

func ask(ctx context.Context, client *openai.Client, model, label, prompt string) {
	resp, err := client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:       model,
			Temperature: 0.1,
			TopP:        0.95,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("=== %s ===\n%s\n\n", label, resp.Choices[0].Message.Content)
}

func main() {
	var wg sync.WaitGroup
	wg.Add(2)
	cfg := openai.DefaultConfig("")
	cfg.BaseURL = defaultBaseURL
	client := openai.NewClientWithConfig(cfg)

	ctx := context.Background()
	go func() {
		defer wg.Done()
		ask(ctx, client, model, "模糊Prompt", "帮我写个Go函数处理错误。")
	}()

	go func() {
		defer wg.Done()
		ask(ctx, client, model, "清晰Prompt",
		`用Go语言编写一个函数 WrapError，功能如下：
- 接收一个 error 和一条描述字符串
- 如果 error 为 nil，直接返回 nil
- 如果 error 不为 nil，用 fmt.Errorf 包装错误并添加描述信息
- 请给出函数签名、实现和一个简单的使用示例`,
	)}()

	wg.Wait()
}
