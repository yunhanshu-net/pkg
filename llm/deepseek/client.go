package deepseek

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/yunhanshu-net/pkg/llm"
	"github.com/yunhanshu-net/pkg/logger"
)

// Client DeepSeek客户端实现
type Client struct {
	config     llm.Config
	httpClient *http.Client
}

// NewClient 创建DeepSeek客户端
func NewClient(config llm.Config) (*Client, error) {
	if config.APIKey == "" {
		return nil, fmt.Errorf("DeepSeek API key is required")
	}

	// 设置默认值
	if config.BaseURL == "" {
		config.BaseURL = "https://api.deepseek.com"
	}
	if config.Timeout == 0 {
		config.Timeout = 30 * time.Second
	}

	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}, nil
}

// ChatCompletion 实现聊天完成接口
func (c *Client) ChatCompletion(ctx context.Context, req *llm.ChatCompletionRequest) (*llm.ChatCompletionResponse, error) {
	// 设置默认值
	if req.Model == "" {
		req.Model = c.config.DefaultModel
	}
	if req.Temperature == 0 {
		req.Temperature = c.config.DefaultTemperature
	}
	if req.MaxTokens == 0 {
		req.MaxTokens = c.config.DefaultMaxTokens
	}
	if req.TopP == 0 {
		req.TopP = c.config.DefaultTopP
	}

	// 处理RAG上下文
	if req.RAGContext != nil && c.config.EnableRAG {
		err := c.processRAGContext(ctx, req)
		if err != nil {
			logger.ErrorContextf(ctx, "[DeepSeek] RAG处理失败: %v", err)
		}
	}

	// 处理提示词模板
	if req.PromptTemplate != "" && req.Variables != nil {
		err := c.processPromptTemplate(req)
		if err != nil {
			logger.ErrorContextf(ctx, "[DeepSeek] 提示词模板处理失败: %v", err)
		}
	}

	// 构建请求
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request failed: %w", err)
	}

	url := fmt.Sprintf("%s/chat/completions", c.config.BaseURL)
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("create request failed: %w", err)
	}

	// 设置请求头
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.config.APIKey))

	// 记录请求日志
	logger.Infof(ctx, "[DeepSeek] 发送请求: %s, Model: %s, Messages: %d条, Temperature: %.2f",
		url, req.Model, len(req.Messages), req.Temperature)

	// 发送请求（带重试）
	startTime := time.Now()
	resp, err := c.doRequestWithRetry(ctx, httpReq)
	if err != nil {
		logger.ErrorContextf(ctx, "[DeepSeek] 请求失败: %v", err)
		return nil, err
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response failed: %w", err)
	}

	// 记录响应日志
	duration := time.Since(startTime)
	logger.Infof(ctx, "[DeepSeek] 响应完成: 状态=%d, 耗时=%v, 响应长度=%d",
		resp.StatusCode, duration, len(respBody))

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		var apiErr llm.APIError
		if json.Unmarshal(respBody, &apiErr) == nil {
			apiErr.StatusCode = resp.StatusCode
			return nil, &apiErr
		}
		return nil, fmt.Errorf("API error: status=%d, body=%s", resp.StatusCode, string(respBody))
	}

	// 解析响应
	var response llm.ChatCompletionResponse
	if err := json.Unmarshal(respBody, &response); err != nil {
		return nil, fmt.Errorf("unmarshal response failed: %w", err)
	}

	// 记录token使用情况
	if response.Usage.TotalTokens > 0 {
		logger.Infof(ctx, "[DeepSeek] Token使用: 输入=%d, 输出=%d, 总计=%d",
			response.Usage.PromptTokens, response.Usage.CompletionTokens, response.Usage.TotalTokens)
	}

	return &response, nil
}

// GetClientInfo 获取客户端信息
func (c *Client) GetClientInfo() llm.ClientInfo {
	return llm.ClientInfo{
		Provider: "deepseek",
		Version:  "v1.0.0",
		SupportedModels: []string{
			"deepseek-coder",
			"deepseek-chat",
			"deepseek-reasoner",
		},
		Features: []string{
			"chat_completion",
			"rag_support",
			"template_rendering",
			"retry_mechanism",
		},
		APIEndpoint: c.config.BaseURL,
		RateLimit: llm.RateLimit{
			RequestsPerMinute: c.config.RateLimitRPS * 60,
			TokensPerMinute:   100000,
			TokensPerDay:      1000000,
		},
	}
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
	req := &llm.ChatCompletionRequest{
		Model: c.config.DefaultModel,
		Messages: []llm.Message{
			{Role: "user", Content: "ping"},
		},
		MaxTokens: 10,
	}

	_, err := c.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("health check failed: %w", err)
	}

	return nil
}

// doRequestWithRetry 带重试的请求
func (c *Client) doRequestWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error

	for i := 0; i <= c.config.MaxRetries; i++ {
		if i > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(c.config.RetryInterval):
			}
			logger.Infof(ctx, "[DeepSeek] 重试第%d次", i)
		}

		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		// 如果是临时错误，重试
		if resp.StatusCode >= 500 || resp.StatusCode == 429 {
			resp.Body.Close()
			lastErr = fmt.Errorf("server error: %d", resp.StatusCode)
			continue
		}

		return resp, nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

// processRAGContext 处理RAG上下文
func (c *Client) processRAGContext(ctx context.Context, req *llm.ChatCompletionRequest) error {
	if req.RAGContext == nil || len(req.RAGContext.RetrievedDocs) == 0 {
		return nil
	}

	// 构建RAG上下文消息
	var ragContent string
	ragContent += "参考文档:\n"
	for i, doc := range req.RAGContext.RetrievedDocs {
		if i >= c.config.MaxRAGDocuments {
			break
		}
		ragContent += fmt.Sprintf("文档%d: %s\n", i+1, doc.Content)
	}
	ragContent += "\n基于以上文档回答用户问题:"

	// 在用户消息前插入RAG上下文
	if len(req.Messages) > 0 {
		ragMessage := llm.Message{
			Role:    "system",
			Content: ragContent,
		}
		// 插入到最后一条用户消息之前
		req.Messages = append([]llm.Message{ragMessage}, req.Messages...)
	}

	return nil
}

// processPromptTemplate 处理提示词模板
func (c *Client) processPromptTemplate(req *llm.ChatCompletionRequest) error {
	if req.PromptTemplate == "" {
		return nil
	}

	// 简单的模板变量替换
	template := req.PromptTemplate
	for key, value := range req.Variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		template = fmt.Sprintf("%s", fmt.Sprintf(template, placeholder, value))
	}

	// 替换最后一条用户消息
	if len(req.Messages) > 0 {
		for i := len(req.Messages) - 1; i >= 0; i-- {
			if req.Messages[i].Role == "user" {
				req.Messages[i].Content = template
				break
			}
		}
	}

	return nil
}

// Factory DeepSeek工厂
type Factory struct{}

// CreateClient 创建客户端
func (f *Factory) CreateClient(config llm.Config) (llm.LLMClient, error) {
	return NewClient(config)
}

// SupportedModels 支持的模型
func (f *Factory) SupportedModels() []string {
	return []string{"deepseek-coder", "deepseek-chat", "deepseek-reasoner"}
}

// init 自动注册DeepSeek工厂
func init() {
	llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &Factory{})
}
