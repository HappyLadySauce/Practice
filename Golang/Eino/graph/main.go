package main

import (
	"context"
	"fmt"
	"log"
	"strings"

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

	// 分类器：用 Lambda 做简单的关键词分类
	classifier := compose.InvokableLambda(func(ctx context.Context, input map[string]any) (map[string]any, error) {
		question := input["question"].(string)

		if strings.Contains(question, "代码") || strings.Contains(question, "编程") || strings.Contains(question, "bug") {
			input["category"] = "code"
		} else if strings.Contains(question, "部署") || strings.Contains(question, "运维") || strings.Contains(question, "服务器") {
			input["category"] = "ops"
		} else {
			input["category"] = "general"
		}
		return input, nil	
	})

    // 三个不同角色的 Prompt 模板
    codeTpl := prompt.FromMessages(schema.FString,
       schema.SystemMessage("你是一个资深Go语言开发者，擅长代码审查和问题排查，请回答用户的编程问题。"),
       schema.UserMessage("{question}"),
    )

    opsTpl := prompt.FromMessages(schema.FString,
       schema.SystemMessage("你是一个运维专家，精通Linux、Docker和K8s，请回答用户的运维问题。"),
       schema.UserMessage("{question}"),
    )

    generalTpl := prompt.FromMessages(schema.FString,
       schema.SystemMessage("你是一个友好的技术助手，请简洁地回答用户的问题。"),
       schema.UserMessage("{question}"),
    )

	graph := compose.NewGraph[map[string]any, *schema.Message]()

	// 添加 Graph 节点
	_ = graph.AddLambdaNode("classifier", classifier)
	_ = graph.AddChatTemplateNode("code", codeTpl)
	_ = graph.AddChatTemplateNode("ops", opsTpl)
	_ = graph.AddChatTemplateNode("general", generalTpl)
	_ = graph.AddChatModelNode("model", chatModel)

	// 连接边: 定义数据流向
	_ = graph.AddEdge(compose.START, "classifier")

	_ = graph.AddBranch("classifier", compose.NewGraphBranch(
		// 条件函数：根据分类结果决定走哪个分支
		func(ctx context.Context, input map[string]any) (string, error) {
			category := input["category"].(string)
			return category, nil
		},
		// 分支映射：声明所有可能的下游节点
		map[string]bool{
			"code":		true,
			"ops":		true,
			"general":	true,
		},
	))

	_ = graph.AddEdge("code", "model")
	_ = graph.AddEdge("ops", "model")
	_ = graph.AddEdge("general", "model")
	_ = graph.AddEdge("model", compose.END)
	
	// 编译运行
	runner, err := graph.Compile(ctx)
	if err != nil {
			log.Fatal("编译失败:", err)
	}

    // 测试不同类型的问题
    questions := []string{
       "Go代码里怎么避免goroutine泄漏？",
       "Docker容器部署时端口映射不生效怎么办？",
       "推荐几本学习分布式系统的书？",
    }

    for _, q := range questions {
       result, err := runner.Invoke(ctx, map[string]any{"question": q})
       if err != nil {
          log.Printf("问题: %s, 错误: %v\n", q, err)
          continue
       }
       fmt.Printf("问题: %s\n回答: %s\n\n", q, result.Content)
    }
}
