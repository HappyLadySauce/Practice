package tools

import (
	"context"

	"github.com/cloudwego/eino/schema"
)

type WeatherRequest struct {
	City string `json:"city"`
}

type WeatherResponse struct {
	City 	string	`json:"city"`
	Temp 	string	`json:"temp"`
	Weather string	`json:"weather"`
}

// utils.NewTool接收两个参数:一个*schema.ToolInfo描述工具的元信息，
// 一个函数作为工具的执行逻辑。这个函数的签名有严格要求一
// 第一个参数必须是context.context，
// 第二个参数是一个指向结构体的指针(工具的输入参数)，
// 返回值是一个指向结构体的指针(工具的输出)和error。
func GetWeather(ctx context.Context, req *WeatherRequest) (*WeatherResponse, error) {
	mockData := map[string]WeatherResponse{
       "北京": {City: "北京", Temp: "22°C", Weather: "晴"},
       "上海": {City: "上海", Temp: "26°C", Weather: "多云"},
       "深圳": {City: "深圳", Temp: "30°C", Weather: "阵雨"},
	}

	if data, ok := mockData[req.City]; ok {
		return &data, nil
	}
	return &WeatherResponse{City: req.City, Temp: "未知", Weather: "未知"}, nil
}



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