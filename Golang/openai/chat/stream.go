package chat

import (
	"io"

	"github.com/sashabaranov/go-openai"
	"test/tools"
)

// StreamResult holds accumulated output from one streamed completion.
// StreamResult 保存一次流式补全的累积结果。
type StreamResult struct {
	Content      string
	ToolCalls    []openai.ToolCall
	FinishReason openai.FinishReason
}

// ConsumeStream reads SSE chunks, writes assistant text to w, and accumulates tool calls.
// ConsumeStream 读取 SSE 分片，将助手文本写入 w，并累积工具调用信息。
func ConsumeStream(stream *openai.ChatCompletionStream, w io.Writer) (StreamResult, error) {
	var result StreamResult

	for {
		response, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				return result, nil
			}
			return result, err
		}
		if len(response.Choices) == 0 {
			continue
		}

		choice := response.Choices[0]
		if choice.FinishReason != "" {
			result.FinishReason = choice.FinishReason
		}

		for _, tc := range choice.Delta.ToolCalls {
			result.ToolCalls = tools.MergeToolCallDelta(result.ToolCalls, tc)
		}

		chunk := choice.Delta.Content
		if chunk != "" && w != nil {
			_, _ = io.WriteString(w, chunk)
		}
		result.Content += chunk
	}
}
