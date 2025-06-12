package llm

import (
	"fmt"
	"time"
)

// ChatCompletionRequest 聊天完成请求（通用结构，支持多种模型）
type ChatCompletionRequest struct {
	Model            string    `json:"model"`                       // 具体模型名称，如 "deepseek-coder"
	Messages         []Message `json:"messages"`                    // 对话消息列表
	Temperature      float32   `json:"temperature,omitempty"`       // 温度参数 0.0-2.0
	MaxTokens        int       `json:"max_tokens,omitempty"`        // 最大输出token数
	TopP             float32   `json:"top_p,omitempty"`             // 核采样参数
	FrequencyPenalty float32   `json:"frequency_penalty,omitempty"` // 频率惩罚
	PresencePenalty  float32   `json:"presence_penalty,omitempty"`  // 存在惩罚
	Stream           bool      `json:"stream,omitempty"`            // 是否流式输出
	Stop             []string  `json:"stop,omitempty"`              // 停止标记

	// 扩展字段
	ResponseFormat *ResponseFormat `json:"response_format,omitempty"` // 响应格式控制
	Tools          []Tool          `json:"tools,omitempty"`           // 工具调用
	ToolChoice     string          `json:"tool_choice,omitempty"`     // 工具选择策略

	// RAG相关字段
	RAGContext     *RAGContext       `json:"rag_context,omitempty"`     // RAG上下文
	PromptTemplate string            `json:"prompt_template,omitempty"` // 提示词模板
	Variables      map[string]string `json:"variables,omitempty"`       // 模板变量
}

// Message 消息结构
type Message struct {
	Role       string     `json:"role"`                   // system, user, assistant, tool
	Content    string     `json:"content"`                // 消息内容
	Name       string     `json:"name,omitempty"`         // 消息发送者名称
	ToolCalls  []ToolCall `json:"tool_calls,omitempty"`   // 工具调用
	ToolCallID string     `json:"tool_call_id,omitempty"` // 工具调用ID
}

// ChatCompletionResponse 聊天完成响应
type ChatCompletionResponse struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint,omitempty"`
}

// Choice 选择项
type Choice struct {
	Index        int         `json:"index"`
	Message      Message     `json:"message,omitempty"`
	Delta        *Message    `json:"delta,omitempty"` // 流式响应增量
	FinishReason string      `json:"finish_reason"`
	Logprobs     interface{} `json:"logprobs,omitempty"`
}

// Usage 使用统计
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// ResponseFormat 响应格式
type ResponseFormat struct {
	Type   string `json:"type"`             // "text" 或 "json_object"
	Schema string `json:"schema,omitempty"` // JSON Schema
}

// Tool 工具定义
type Tool struct {
	Type     string       `json:"type"` // "function"
	Function ToolFunction `json:"function"`
}

// ToolFunction 工具函数
type ToolFunction struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Parameters  interface{} `json:"parameters"`
}

// ToolCall 工具调用
type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // "function"
	Function FunctionCall `json:"function"`
}

// FunctionCall 函数调用
type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

// RAGContext RAG上下文信息
type RAGContext struct {
	Query           string              `json:"query"`            // 查询内容
	RetrievedDocs   []RetrievedDocument `json:"retrieved_docs"`   // 检索到的文档
	MaxDocuments    int                 `json:"max_documents"`    // 最大文档数量
	SimilarityScore float64             `json:"similarity_score"` // 相似度阈值
}

// RetrievedDocument 检索到的文档
type RetrievedDocument struct {
	ID       string            `json:"id"`
	Content  string            `json:"content"`
	Metadata map[string]string `json:"metadata"`
	Score    float64           `json:"score"`
}

// ClientInfo 客户端信息
type ClientInfo struct {
	Provider        string    `json:"provider"` // deepseek, openai, claude等
	Version         string    `json:"version"`
	SupportedModels []string  `json:"supported_models"`
	Features        []string  `json:"features"` // 支持的功能特性
	APIEndpoint     string    `json:"api_endpoint"`
	RateLimit       RateLimit `json:"rate_limit"`
}

// RateLimit 速率限制信息
type RateLimit struct {
	RequestsPerMinute int `json:"requests_per_minute"`
	TokensPerMinute   int `json:"tokens_per_minute"`
	TokensPerDay      int `json:"tokens_per_day"`
}

// APIError API错误
type APIError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	Type       string `json:"type"`
	Param      string `json:"param,omitempty"`
	Details    string `json:"details,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
}

func (e *APIError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("API错误[%d]: %s - %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("API错误[%d]: %s", e.Code, e.Message)
}

// StreamChoice 流式响应选择项
type StreamChoice struct {
	Index        int     `json:"index"`
	Delta        Message `json:"delta"`
	FinishReason string  `json:"finish_reason"`
}

// StreamResponse 流式响应
type StreamResponse struct {
	ID      string         `json:"id"`
	Object  string         `json:"object"`
	Created int64          `json:"created"`
	Model   string         `json:"model"`
	Choices []StreamChoice `json:"choices"`
}

// PromptTemplate 提示词模板
type PromptTemplate struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Template    string            `json:"template"`
	Variables   []TemplateVar     `json:"variables"`
	Examples    []TemplateExample `json:"examples"`
	Tags        []string          `json:"tags"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// TemplateVar 模板变量定义
type TemplateVar struct {
	Name         string `json:"name"`
	Type         string `json:"type"` // string, int, bool, array
	Required     bool   `json:"required"`
	Description  string `json:"description"`
	DefaultValue string `json:"default_value,omitempty"`
}

// TemplateExample 模板示例
type TemplateExample struct {
	Description string            `json:"description"`
	Variables   map[string]string `json:"variables"`
	Expected    string            `json:"expected"`
}
