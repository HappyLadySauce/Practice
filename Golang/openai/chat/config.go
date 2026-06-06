package chat

import (
	"os"
	"time"
)

// Config holds client-level defaults.
// Config 保存客户端级别的默认配置。
type Config struct {
	BaseURL       string
	Token         string
	Model         string
	MaxRetries    int
	MaxToolRounds int
	Timeout       time.Duration // per-request context deadline
	HTTPTimeout   time.Duration // underlying http.Client timeout
	Temperature   float32
	TopP          float32
}

// DefaultConfig returns defaults, overridable via environment variables.
// DefaultConfig 返回默认配置，可通过环境变量覆盖。
//
// Supported env vars:
//   - OPENAI_BASE_URL / OPENAI_API_KEY / OPENAI_MODEL
//   - DASHSCOPE_API_KEY (fallback token for Aliyun DashScope)
func DefaultConfig() Config {
	baseURL := envOr("OPENAI_BASE_URL", "http://127.0.0.1:11434/v1")

	token := os.Getenv("OPENAI_API_KEY")
	if token == "" {
		token = os.Getenv("DASHSCOPE_API_KEY")
	}

	return Config{
		BaseURL:       baseURL,
		Token:         token,
		Model:         envOr("OPENAI_MODEL", "gemma-4-e4b-it-uncensored"),
		MaxRetries:    3,
		MaxToolRounds: 5,
		Timeout:       60 * time.Second,
		HTTPTimeout:   30 * time.Second,
		Temperature:   0.1,
		TopP:          0.9,
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
