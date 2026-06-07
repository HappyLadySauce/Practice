package tools

import (
	"context"
	"strings"
)

type SearchRequest struct {
	Query	string	`json:"query" jsonschema:"required" jsonschema_description:"搜索关键词"`
    MaxCount int    `json:"max_count" jsonschema_description:"最多返回的结果数量, 默认5"`
    Language string `json:"language" jsonschema:"enum=zh,enum=en" jsonschema_description:"结果语言, zh为中文, en为英文"`
}

type SearchResponse struct {
	Items []SearchItem	`json:"items"`
	Total int			`json:"total"`
}

type SearchItem struct {
	Title	string	`json:"title"`
	URL		string	`json:"url"`
	Summary	string	`json:"summary"`
}

func SearchWeb(ctx context.Context, req *SearchRequest) (*SearchResponse, error) {
    // 模拟搜索逻辑
    maxCount := req.MaxCount
    if maxCount <= 0 {
       maxCount = 5
    }

    items := []SearchItem{
       {Title: "Go语言官方文档", URL: "https://go.dev/doc/", Summary: "Go编程语言官方文档和教程"},
       {Title: "Eino框架指南", URL: "https://cloudwego.io/docs/eino", Summary: "字节跳动开源的Go语言LLM应用开发框架"},
       {Title: "Go并发编程实战", URL: "https://example.com/go-concurrency", Summary: "深入讲解goroutine和channel的用法"},
    }

    // 简单过滤
    filtered := make([]SearchItem, 0)
    for _, item := range items {
       if strings.Contains(strings.ToLower(item.Title+item.Summary), strings.ToLower(req.Query)) {
          filtered = append(filtered, item)
       }
       if len(filtered) >= maxCount {
          break
       }
    }

    return &SearchResponse{Items: filtered, Total: len(filtered)}, nil
}
