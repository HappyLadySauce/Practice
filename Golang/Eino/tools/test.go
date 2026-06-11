package tools

// const MODEL = "gemma-4-e4b-it-uncensored"
// const APIURL = "http://127.0.0.1:11434/v1"

// func trimHistory(history []*schema.Message, maxRounds int) []*schema.Message {
// 	maxMessages := maxRounds * 2
// 	if len(history) <= maxMessages {
// 		return history
// 	}
// 	return history[len(history) - maxMessages:]
// }

// func main() {
// 	ctx := context.Background()

// 	chatModel, err := openai.NewChatModel(ctx, &openai.ChatModelConfig{
// 		APIKey: "",
// 		Model: MODEL,
// 		BaseURL: APIURL,
// 	})
// 	if err != nil {
// 		log.Fatal("创建 ChatModel 失败: %w", err)
// 	}

// 	template := prompt.FromMessages(
// 		schema.FString,
// 		schema.SystemMessage("你是一个友好的AI助手，名叫小秀。回答简洁，不超过100字。"),
// 		schema.MessagesPlaceholder("history", false),
// 		schema.UserMessage("{input}"),
// 	)

// 	history := make([]*schema.Message, 0)

// 	weatherTool := utils.NewTool(tools.GetWeatherTool(), tools.GetWeather)
// 	weatherToolInfo, _ := weatherTool.Info(ctx)
//     fmt.Printf("工具名: %s\n", weatherToolInfo.Name)
//     fmt.Printf("工具描述: %s\n", weatherToolInfo.Desc)

// 	searchTool, err := utils.InferTool("web_search", "搜索互联网上的信息，返回相关网页的标题、链接和摘要", tools.SearchWeb)
// 	if err != nil {
// 		log.Print("工具构建错误: %w", err)
// 	}
// 	searchToolInfo, _ := searchTool.Info(ctx)
//     fmt.Printf("工具名: %s\n", searchToolInfo.Name)
//     fmt.Printf("工具描述: %s\n", searchToolInfo.Desc)

// 	toolsNode, err := compose.NewToolNode(ctx, &compose.ToolsNodeConfig{
// 		Tools: []tool.BaseTool{weatherTool, searchTool},
// 	})
// 	if err != nil {
// 		log.Print("工具执行节点构建错误: %w", err)
// 	}

// 	chatModel.BindTools([]*schema.ToolInfo{searchToolInfo, weatherToolInfo})

// 	scanner := bufio.NewScanner(os.Stdin)
// 	fmt.Println("开始对话(输入quit退出):")

// 	for {
// 		fmt.Print("用户: ")
// 		if !scanner.Scan() {
// 			break
// 		}
// 		input := strings.TrimSpace(scanner.Text())
//		if input == "quit" {
//			println("再见")
//			break
//		}
// 		messages, err := template.Format(ctx, map[string]any{
// 			"history": history,
// 			"input": input,
// 		})
// 		if err != nil {
// 			log.Printf("格式化模板失败: %w", err)
// 			continue
// 		}

// 		resp, err := chatModel.Generate(ctx, messages)
// 		if err != nil {
// 			log.Printf("模型调用失败: %w", err)
// 			continue
// 		}

// 		if len(resp.ToolCalls) > 0 {
// 			fmt.Printf("模型决定调用工具: %s\n", resp.ToolCalls[0].Function.Name)
// 			fmt.Printf("调用参数: %s\n\n", resp.ToolCalls[0].Function.Arguments)

// 			toolResults, err := toolsNode.Invoke(ctx, resp)
// 			if err != nil {
// 				log.Printf("工具调用失败: %w", err)
// 			}

// 			fmt.Printf("工具返回: %s\n\n", toolResults[0].Content)

// 			messages = append(messages, resp)
// 			messages = append(messages, toolResults...)

// 			finalResp, err := chatModel.Generate(ctx, messages)
// 			if err != nil {
// 				log.Printf("模型调用失败: %w", err)
// 				continue
// 			}

// 			fmt.Print("小秀: ")
// 			fmt.Println(finalResp.Content)
// 		} else {
// 			fmt.Print("小秀: ")
// 			fmt.Println(resp.Content)
// 		}

// 		history = append(history, schema.UserMessage(input))
// 		history = append(history, schema.AssistantMessage(resp.Content, nil))
// 		trimHistory(history, 2)
// 	}
// }