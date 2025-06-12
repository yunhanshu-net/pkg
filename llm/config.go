package llm

import (
	"fmt"
	"time"
)

// ProviderType 提供商类型
type ProviderType string

const (
	ProviderDeepSeek ProviderType = "deepseek"
	ProviderOpenAI   ProviderType = "openai"
	ProviderClaude   ProviderType = "claude"
	ProviderQwen     ProviderType = "qwen"    // 阿里通义千问
	ProviderChatGLM  ProviderType = "chatglm" // 智谱ChatGLM
)

// Config 通用配置
type Config struct {
	Provider ProviderType  `json:"provider" yaml:"provider"`
	APIKey   string        `json:"api_key" yaml:"api_key"`
	BaseURL  string        `json:"base_url" yaml:"base_url"`
	Timeout  time.Duration `json:"timeout" yaml:"timeout"`

	// 默认请求参数
	DefaultModel       string  `json:"default_model" yaml:"default_model"`
	DefaultTemperature float32 `json:"default_temperature" yaml:"default_temperature"`
	DefaultMaxTokens   int     `json:"default_max_tokens" yaml:"default_max_tokens"`
	DefaultTopP        float32 `json:"default_top_p" yaml:"default_top_p"`

	// 重试配置
	MaxRetries    int           `json:"max_retries" yaml:"max_retries"`
	RetryInterval time.Duration `json:"retry_interval" yaml:"retry_interval"`

	// 日志配置
	EnableLogging bool   `json:"enable_logging" yaml:"enable_logging"`
	LogLevel      string `json:"log_level" yaml:"log_level"` // debug, info, warn, error

	// RAG配置
	EnableRAG        bool    `json:"enable_rag" yaml:"enable_rag"`
	MaxRAGDocuments  int     `json:"max_rag_documents" yaml:"max_rag_documents"`
	RAGSimilarityMin float64 `json:"rag_similarity_min" yaml:"rag_similarity_min"`

	// 性能配置
	MaxConcurrency int `json:"max_concurrency" yaml:"max_concurrency"` // 最大并发数
	RateLimitRPS   int `json:"rate_limit_rps" yaml:"rate_limit_rps"`   // 每秒请求数限制

	// JSON输出配置
	EnableJSONMode bool // 是否默认启用JSON格式输出
}

// DeepSeekConfig DeepSeek特定配置
type DeepSeekConfig struct {
	Config
	// DeepSeek特定配置可以在这里扩展
	EnableReasoner bool `json:"enable_reasoner" yaml:"enable_reasoner"` // 是否启用推理模型
}

// OpenAIConfig OpenAI特定配置
type OpenAIConfig struct {
	Config
	Organization string `json:"organization" yaml:"organization"` // 组织ID
	ProjectID    string `json:"project_id" yaml:"project_id"`     // 项目ID
}

// ClaudeConfig Claude特定配置
type ClaudeConfig struct {
	Config
	AnthropicVersion string `json:"anthropic_version" yaml:"anthropic_version"` // API版本
}

// DefaultConfigs 默认配置
var DefaultConfigs = map[ProviderType]Config{
	ProviderDeepSeek: {
		Provider:           ProviderDeepSeek,
		BaseURL:            "https://api.deepseek.com",
		Timeout:            30 * time.Second,
		DefaultModel:       "deepseek-coder",
		DefaultTemperature: 0.1,
		DefaultMaxTokens:   2000,
		DefaultTopP:        0.9,
		MaxRetries:         3,
		RetryInterval:      1 * time.Second,
		EnableLogging:      true,
		LogLevel:           "info",
		EnableRAG:          true,
		MaxRAGDocuments:    5,
		RAGSimilarityMin:   0.7,
		MaxConcurrency:     10,
		RateLimitRPS:       60,
		EnableJSONMode:     false, // 默认不启用JSON模式
	},
	ProviderOpenAI: {
		Provider:           ProviderOpenAI,
		BaseURL:            "https://api.openai.com/v1",
		Timeout:            60 * time.Second,
		DefaultModel:       "gpt-4",
		DefaultTemperature: 0.7,
		DefaultMaxTokens:   4000,
		DefaultTopP:        1.0,
		MaxRetries:         3,
		RetryInterval:      2 * time.Second,
		EnableLogging:      true,
		LogLevel:           "info",
		EnableRAG:          true,
		MaxRAGDocuments:    10,
		RAGSimilarityMin:   0.8,
		MaxConcurrency:     5,
		RateLimitRPS:       60,
		EnableJSONMode:     false, // 默认不启用JSON模式
	},
	ProviderClaude: {
		Provider:           ProviderClaude,
		BaseURL:            "https://api.anthropic.com",
		Timeout:            60 * time.Second,
		DefaultModel:       "claude-3-5-sonnet-20241022",
		DefaultTemperature: 0.3,
		DefaultMaxTokens:   4000,
		DefaultTopP:        1.0,
		MaxRetries:         3,
		RetryInterval:      2 * time.Second,
		EnableLogging:      true,
		LogLevel:           "info",
		EnableRAG:          true,
		MaxRAGDocuments:    8,
		RAGSimilarityMin:   0.75,
		MaxConcurrency:     3,
		RateLimitRPS:       30,
		EnableJSONMode:     false, // 默认不启用JSON模式
	},
}

// GetDefaultConfig 获取默认配置
func GetDefaultConfig(provider ProviderType) Config {
	if config, exists := DefaultConfigs[provider]; exists {
		return config
	}
	// 返回DeepSeek作为默认配置
	return DefaultConfigs[ProviderDeepSeek]
}

// GetJSONConfig 获取启用JSON模式的配置
func GetJSONConfig(provider ProviderType) Config {
	config := GetDefaultConfig(provider)
	config.EnableJSONMode = true
	config.DefaultTemperature = 0.1 // JSON模式建议使用较低的温度
	return config
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.APIKey == "" {
		return fmt.Errorf("API key is required")
	}

	if c.BaseURL == "" {
		return fmt.Errorf("base URL is required")
	}

	if c.Timeout <= 0 {
		return fmt.Errorf("timeout must be positive")
	}

	if c.DefaultMaxTokens <= 0 {
		return fmt.Errorf("default max tokens must be positive")
	}

	if c.DefaultTemperature < 0 || c.DefaultTemperature > 2 {
		return fmt.Errorf("default temperature must be between 0 and 2")
	}

	if c.DefaultTopP < 0 || c.DefaultTopP > 1 {
		return fmt.Errorf("default top_p must be between 0 and 1")
	}

	return nil
}

// Clone 克隆配置
func (c *Config) Clone() Config {
	return *c
}

// WithAPIKey 设置API密钥
func (c *Config) WithAPIKey(apiKey string) Config {
	c.APIKey = apiKey
	return *c
}

// WithJSONMode 启用JSON模式
func (c *Config) WithJSONMode() Config {
	c.EnableJSONMode = true
	c.DefaultTemperature = 0.1 // JSON模式建议低温度
	return *c
}

// WithModel 设置默认模型
func (c *Config) WithModel(model string) Config {
	c.DefaultModel = model
	return *c
}

// WithTemperature 设置温度
func (c *Config) WithTemperature(temperature float32) Config {
	c.DefaultTemperature = temperature
	return *c
}

// WithMaxTokens 设置默认最大token数
func (c *Config) WithMaxTokens(maxTokens int) *Config {
	c.DefaultMaxTokens = maxTokens
	return c
}

// WithTimeout 设置超时时间
func (c *Config) WithTimeout(timeout time.Duration) *Config {
	c.Timeout = timeout
	return c
}

// WithRAG 启用RAG功能
func (c *Config) WithRAG(maxDocuments int, similarityMin float64) *Config {
	c.EnableRAG = true
	c.MaxRAGDocuments = maxDocuments
	c.RAGSimilarityMin = similarityMin
	return c
}
