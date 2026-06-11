package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sashabaranov/go-openai"
)

// 最简单的CoT 使用方式只需要在Prompt末尾加一句话-"Let's think step by step”(让我们一步一步地思考)。
// 这就是著名的Zero-shotCoT。别看这句话简单，它的效果已经被大量实验证实了。

const defaultBaseURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

func generate(client *openai.Client, label, prompt string) {
    resp, err := client.CreateChatCompletion(
       context.Background(),
       openai.ChatCompletionRequest{
          Model:       model,
          Temperature: 0,
          Messages: []openai.ChatCompletionMessage{
             {Role: openai.ChatMessageRoleUser, Content: prompt},
          },
       },
    )
    if err != nil {
       log.Fatal(err)
    }
    fmt.Printf("=== %s ===\n%s\n\n", label, resp.Choices[0].Message.Content)
}

func main() {
    cfg := openai.DefaultConfig("")
    cfg.BaseURL = defaultBaseURL
    client := openai.NewClientWithConfig(cfg)

    problem := `小明的书架上有3层，第一层放了12本书，第二层放的书是第一层的2倍，
第三层放的书比第一层和第二层加起来少5本。书架上一共有多少本书？`

    // 不用CoT，直接问
    generate(client, "直接回答", problem+"\n请直接给出答案。")

    // 使用CoT
    generate(client, "CoT推理", problem+"\n请一步一步地思考，展示你的推理过程。")
}