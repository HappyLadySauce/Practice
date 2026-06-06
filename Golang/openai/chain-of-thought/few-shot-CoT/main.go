package main

import (
    "context"
    "fmt"
    "log"

    "github.com/sashabaranov/go-openai"
)

const defaultBaseURL = "http://127.0.0.1:11434/v1"
const model = "gemma-4-e4b-it-uncensored"

func main() {
    cfg := openai.DefaultConfig("")
    cfg.BaseURL = defaultBaseURL
    client := openai.NewClientWithConfig(cfg)

    // Few-shot CoT：教Agent做工具调用决策
    prompt := `你是一个AI助手，可以使用以下工具：
- search_web: 搜索互联网获取实时信息
- query_database: 查询内部数据库
- send_email: 发送邮件
- calculate: 执行数学计算

请根据用户的输入，判断是否需要调用工具，以及调用哪个工具。

示例1：
用户：今天上海的股票大盘表现怎么样？
思考：用户问的是"今天"的股票信息，这是实时数据，我的训练数据中不包含今天的信息。我需要使用search_web来获取实时的股票市场数据。
决策：{"need_tool": true, "tool": "search_web", "query": "今天上海股票大盘行情"}

示例2：
用户：Go语言的goroutine和线程有什么区别？
思考：这是一个通用的技术知识问题，我的训练数据中已经包含了这方面的知识，不需要查阅外部资料就能准确回答。
决策：{"need_tool": false, "reason": "已有知识可直接回答"}

示例3：
用户：帮我算一下如果投资10万元，年化收益5.5%，复利计算3年后是多少？
思考：这是一个数学计算问题。虽然我能估算，但复利计算需要精确结果，使用calculate工具更可靠。计算公式是：100000 × (1 + 0.055)^3。
决策：{"need_tool": true, "tool": "calculate", "expression": "100000 * (1 + 0.055) ** 3"}

现在请处理：
用户：帮我查一下我们公司上个月的销售额是多少？`

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