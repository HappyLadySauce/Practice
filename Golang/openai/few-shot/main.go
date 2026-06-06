package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sashabaranov/go-openai"
)

// 有时候，光靠文字描述很难让模型准确理解你要的是什么
// 一尤其是当你需要特定的输出格式、特殊的分类标准，或者模型对你的任务理解有偏差时。
// 这时候就需要Few-shot Prompting了。

// Few-shot的核心思路非常基础:
// 给模型几个"输入一输出”的例子，让它从例子中总结出规律，然后按照同样的规律处理新的输入。
// 这就像你教一个新同事做数据标注，与其给他写一页A4纸的标注规范，不如直接丢几个标注好的样本让他照着做
// --后者往往更有效。

const defaultBaseURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

func main() {
    cfg := openai.DefaultConfig("")
    cfg.BaseURL = defaultBaseURL
    client := openai.NewClientWithConfig(cfg)

    // Few-shot Prompt：用3个例子教会模型转换规则
    prompt := `你的任务是将用户的自然语言指令转换为JSON命令格式。

示例1：
用户：帮我查一下北京明天的天气
输出：{"action": "query_weather", "params": {"city": "北京", "date": "tomorrow"}}

示例2：
用户：把这个文件发给张三
输出：{"action": "send_file", "params": {"recipient": "张三", "file": "current"}}

示例3：
用户：创建一个名为weekly-report的定时任务，每周一早上9点执行
输出：{"action": "create_scheduled_task", "params": {"name": "weekly-report", "cron": "0 9 * * 1"}}

现在请转换以下指令：
用户：帮我订一张下周五从上海到深圳的机票`

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

    fmt.Println(resp.Choices[0].Message.Content)
}
