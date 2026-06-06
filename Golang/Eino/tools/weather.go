package tools

import (
    "github.com/cloudwego/eino/schema"
)


func GetWeatherTool() *schema.ToolInfo {
	return &schema.ToolInfo{
		Name: "get_weather",
		Desc: "查询指定城市得当前天气",
		ParamsOneOf: schema.NewParamsOneOfByParams(map[string]*schema.ParameterInfo{
			"city": {
				Type: "string",
				Desc: "城市名称, 如: 北京、上海",
				Required: true,
			},
		}),
	}
}