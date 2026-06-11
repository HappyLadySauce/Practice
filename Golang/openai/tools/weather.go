package tools

import (
	"fmt"

	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
)

func GetWeather(city string) string {
	// 这里用硬编码数据模拟，实际开发中替换为真实API调用
	weatherData := map[string]string{
			"北京": `{"city":"北京","temp":"22°C","weather":"晴","wind":"北风3级"}`,
			"上海": `{"city":"上海","temp":"26°C","weather":"多云","wind":"东南风2级"}`,
			"深圳": `{"city":"深圳","temp":"30°C","weather":"雷阵雨","wind":"南风4级"}`,
	}
	if data, ok := weatherData[city]; ok {
			return data
	}
	return fmt.Sprintf(`{"city":"%s","error":"暂无该城市的天气数据"}`, city)
}

func GetWeatherTool() openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name: "get_weather",
			Description: "Get the current weather for a specified city.",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"city": {
						Type: jsonschema.String,
						Description: "城市名称，例如：北京、上海、深圳",
					},
				},
				Required: []string{"city"},
			},
		},
	}
}