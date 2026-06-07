package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"test/tools"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
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
		schema.MessagesPlaceholder("history", false),
		schema.UserMessage("{input}"),
	)

	history := make([]*schema.Message, 0)

	weatherTool := utils.NewTool(tools.GetWeatherTool(), tools.GetWeather)
	weatherToolInfo, _ := weatherTool.Info(ctx)
    fmt.Printf("工具名: %s\n", weatherToolInfo.Name)
    fmt.Printf("工具描述: %s\n", weatherToolInfo.Desc)

	searchTool, err := utils.InferTool("web_search", "搜索互联网上的信息，返回相关网页的标题、链接和摘要", tools.SearchWeb)
	if err != nil {
		log.Print("工具构建错误: %w", err)
	}
	searchToolInfo, _ := searchTool.Info(ctx)
    fmt.Printf("工具名: %s\n", searchToolInfo.Name)
    fmt.Printf("工具描述: %s\n", searchToolInfo.Desc)

	agent, err := react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: chatModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{weatherTool, searchTool},
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message{
			messages := make([]*schema.Message, 0, len(input)+1)
			messages = append(messages, schema.SystemMessage(
				"你是一个AI助手, 名叫小秀。回答简洁, 不超过100字。",
			))
			messages = append(messages, input...)
			return messages
		},
	})
	if err != nil {
		log.Fatalf("创建 Agent 失败: %v", err)
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("开始对话(输入quit退出):")

	for {
		fmt.Print("用户: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		if input == "quit" {
			println("再见")
			break
		}

		messages, err := template.Format(ctx, map[string]any{
			"history": history,
			"input": input,
		})
		if err != nil {
			log.Printf("格式化模板失败: %w", err)
			continue
		}

		resp, err := agent.Generate(ctx, messages)
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