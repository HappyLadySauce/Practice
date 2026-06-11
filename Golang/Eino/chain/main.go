package main

import (
	"context"
	"fmt"
	"log"


	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/schema"
)

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

	tpl := prompt.FromMessages(
		schema.FString,
		schema.SystemMessage("你是一个{role}, 请用简洁专业的语言回答问题"),
		schema.UserMessage("{question}"),
	)


	chain := compose.NewChain[map[string]any, string]()
	chain.
			AppendChatTemplate(tpl).
			AppendChatModel(chatModel).
			AppendLambda(compose.InvokableLambda(func(ctx context.Context, msg *schema.Message) (string, error) {
				return fmt.Sprintf("转为string: %s", msg.Content), nil
			}))

	runner, err := chain.Compile(ctx)
	if err != nil {
		log.Fatal("编译失败: %w", err)
	}

	result, err := runner.Invoke(ctx, map[string]any{
		"role": "Go语言专家",
		"question": "Go的channel和mutex各自适合什么场景?",
	})
	if err != nil {
		log.Fatal("运行失败:", err)
	}

	fmt.Println(result)
}