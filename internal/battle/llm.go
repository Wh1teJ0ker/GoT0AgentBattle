package battle

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
}

type chatChoice struct {
	Message chatMessage `json:"message"`
}

type chatResponse struct {
	Choices []chatChoice `json:"choices"`
}

// ChatCompletion 按统一 OpenAI 兼容接口向指定模型发起一次对话请求。
func ChatCompletion(cfg ProviderConfig, messages []chatMessage) (string, error) {
	timeout := 300 * time.Second
	if cfg.TimeoutSeconds > 0 {
		timeout = time.Duration(cfg.TimeoutSeconds) * time.Second
	}

	retryAttempts := 3
	if cfg.RetryAttempts > 0 {
		retryAttempts = cfg.RetryAttempts
	}

	retryDelay := 5 * time.Second
	if cfg.RetryDelaySeconds > 0 {
		retryDelay = time.Duration(cfg.RetryDelaySeconds) * time.Second
	}

	var lastErr error
	for attempt := 1; attempt <= retryAttempts; attempt++ {
		content, err := doChatCompletion(cfg, messages, timeout)
		if err == nil {
			return content, nil
		}

		lastErr = err
		if attempt == retryAttempts {
			break
		}
		time.Sleep(retryDelay * time.Duration(attempt))
	}

	return "", fmt.Errorf("model request failed after %d attempts: %w", retryAttempts, lastErr)
}

func doChatCompletion(cfg ProviderConfig, messages []chatMessage, timeout time.Duration) (string, error) {
	reqBody := chatRequest{
		Model:    cfg.Model,
		Messages: messages,
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	endpoint := strings.TrimRight(cfg.BaseURL, "/") + "/chat/completions"
	req, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cfg.APIKey)

	client := &http.Client{Timeout: timeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("api error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}

	var parsed chatResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "", err
	}

	if len(parsed.Choices) == 0 {
		return "", fmt.Errorf("no choices returned from model: %s", cfg.Name)
	}

	return strings.TrimSpace(parsed.Choices[0].Message.Content), nil
}
