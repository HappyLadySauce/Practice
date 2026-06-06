package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/schema"
)

const MODEL = "gemma-4-e4b-it-uncensored"
const APIURL = "http://127.0.0.1:11434/v1"

func trimHistory(history []*schema.Message, maxRounds int) []*schema.Message {
	maxMessages := maxRounds * 2
	if len(history) <= maxMessages {
		return history
	}
	return history[len(history) - maxMessages:]
}

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

	template := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个友好的AI助手，名叫小秀。回答简洁，不超过100字。"),
		schema.MessagesPlaceholder("history", false),
		schema.UserMessage("{input}"),
	)

	history := make([]*schema.Message, 0)

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话(输入quit退出):")

	for {
		fmt.Print("用户: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())

		messages, err := template.Format(ctx, map[string]any{
			"history": history,
			"input": input,
		})
		if err != nil {
			log.Printf("格式化模板失败: %w", err)
			continue
		}

		resp, err := chatModel.Generate(ctx, messages)
		if err != nil {
			log.Printf("模型调用失败: %w", err)
			continue
		}

		fmt.Print("小秀: ")
		fmt.Println(resp.Content)

		history = append(history, schema.UserMessage(input))
		history = append(history, schema.AssistantMessage(resp.Content, nil))
		trimHistory(history, 2)
	}
}