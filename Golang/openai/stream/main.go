package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
	"test/chat"
)

func main() {
	cfg := chat.DefaultConfig()
	client := chat.New(cfg)

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "你是一位友好的 Go 语言助手。当用户询问天气时，请调用 get_weather 工具获取数据后再回答。",
		},
	}

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

		fmt.Print("助手: ")
		var err error
		messages, err = client.RunTurn(context.Background(), messages, os.Stdout)
		fmt.Println()
		if err != nil {
			log.Printf("chat turn failed: %v", err)
			messages = messages[:len(messages)-1]
			continue
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("read stdin failed: %v", err)
	}
}
