package chat

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"time"

	"github.com/sashabaranov/go-openai"
)

type retryAction int

const (
	retryAbort retryAction = iota
	retryAgain
)

// classifyError decides whether an error is worth retrying.
// classifyError 判断错误是否值得重试。
func classifyError(err error) (retryAction, error) {
	if err == nil {
		return retryAbort, nil
	}

	var apiErr *openai.APIError
	if errors.As(err, &apiErr) {
		switch apiErr.HTTPStatusCode {
		case 401:
			return retryAbort, fmt.Errorf("invalid or expired API key: %w", err)
		case 429:
			log.Println("rate limit hit, will retry")
			return retryAgain, nil
		case 500, 502, 503:
			log.Printf("server error (HTTP %d), will retry", apiErr.HTTPStatusCode)
			return retryAgain, nil
		default:
			return retryAbort, fmt.Errorf("API error (HTTP %d): %w", apiErr.HTTPStatusCode, err)
		}
	}

	if errors.Is(err, context.DeadlineExceeded) {
		log.Println("request timeout, will retry")
		return retryAgain, nil
	}

	return retryAbort, fmt.Errorf("request failed: %w", err)
}

func sleepBackoff(attempt int) {
	backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
	log.Printf("retry attempt=%d backoff=%v", attempt, backoff)
	time.Sleep(backoff)
}

// withRetry executes op with per-attempt timeout and exponential backoff.
// withRetry 在单次超时和指数退避保护下执行操作。
func withRetry[T any](cfg Config, parentCtx context.Context, op func(context.Context) (T, error)) (T, error) {
	var zero T
	var lastErr error

	for attempt := 0; attempt <= cfg.MaxRetries; attempt++ {
		if attempt > 0 {
			sleepBackoff(attempt)
		}

		reqCtx, cancel := context.WithTimeout(parentCtx, cfg.Timeout)
		result, err := op(reqCtx)
		cancel()

		if err == nil {
			return result, nil
		}

		lastErr = err
		action, fatalErr := classifyError(err)
		if action == retryAgain {
			continue
		}
		return zero, fatalErr
	}

	return zero, fmt.Errorf("failed after %d retries: %w", cfg.MaxRetries, lastErr)
}
