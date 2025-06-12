package llm

import (
	"context"
	"fmt"
	"sync"
)

// Manager 大模型客户端管理器
type Manager struct {
	clients   map[ProviderType]LLMClient
	factories map[ProviderType]LLMClientFactory
	mutex     sync.RWMutex
}

// NewManager 创建管理器
func NewManager() *Manager {
	m := &Manager{
		clients:   make(map[ProviderType]LLMClient),
		factories: make(map[ProviderType]LLMClientFactory),
	}

	return m
}

// RegisterFactory 注册工厂
func (m *Manager) RegisterFactory(provider ProviderType, factory LLMClientFactory) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.factories[provider] = factory
}

// CreateClient 创建客户端
func (m *Manager) CreateClient(config Config) (LLMClient, error) {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 验证配置
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	factory, exists := m.factories[config.Provider]
	if !exists {
		return nil, fmt.Errorf("unsupported provider: %s", config.Provider)
	}

	client, err := factory.CreateClient(config)
	if err != nil {
		return nil, fmt.Errorf("create client failed: %w", err)
	}

	m.clients[config.Provider] = client
	return client, nil
}

// GetClient 获取客户端
func (m *Manager) GetClient(provider ProviderType) (LLMClient, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	client, exists := m.clients[provider]
	if !exists {
		return nil, fmt.Errorf("client not found for provider: %s", provider)
	}

	return client, nil
}

// GetOrCreateClient 获取或创建客户端
func (m *Manager) GetOrCreateClient(config Config) (LLMClient, error) {
	// 先尝试获取已存在的客户端
	if client, err := m.GetClient(config.Provider); err == nil {
		return client, nil
	}

	// 不存在则创建新客户端
	return m.CreateClient(config)
}

// QuickChat 快速聊天（便利方法）
func (m *Manager) QuickChat(ctx context.Context, provider ProviderType, userMessage string) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := &ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: userMessage},
		},
	}

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// QuickChatWithTemplate 使用模板的快速聊天
func (m *Manager) QuickChatWithTemplate(ctx context.Context, provider ProviderType, template string, variables map[string]string) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := &ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: "请根据模板生成内容"},
		},
		PromptTemplate: template,
		Variables:      variables,
	}

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// QuickChatWithRAG 使用RAG的快速聊天
func (m *Manager) QuickChatWithRAG(ctx context.Context, provider ProviderType, query string, docs []RetrievedDocument) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := &ChatCompletionRequest{
		Messages: []Message{
			{Role: "user", Content: query},
		},
		RAGContext: &RAGContext{
			Query:         query,
			RetrievedDocs: docs,
		},
	}

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// QuickJSONChat JSON格式的快速聊天
func (m *Manager) QuickJSONChat(ctx context.Context, provider ProviderType, userMessage string) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := NewJSONRequest("", NewUserMessage(userMessage))

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateCode 生成代码（JSON格式）
func (m *Manager) GenerateCode(ctx context.Context, provider ProviderType, userRequest string) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := NewCodeGenRequest("", userRequest)

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// GenerateStructuredData 生成结构化数据
func (m *Manager) GenerateStructuredData(ctx context.Context, provider ProviderType, query string, schema string) (string, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return "", err
	}

	req := NewStructuredRequest("", query, schema)

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response choices")
	}

	return resp.Choices[0].Message.Content, nil
}

// HealthCheckAll 检查所有客户端健康状态
func (m *Manager) HealthCheckAll(ctx context.Context) map[ProviderType]error {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	results := make(map[ProviderType]error)
	for provider, client := range m.clients {
		results[provider] = client.HealthCheck(ctx)
	}

	return results
}

// ListProviders 列出支持的提供商
func (m *Manager) ListProviders() []ProviderType {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	providers := make([]ProviderType, 0, len(m.factories))
	for provider := range m.factories {
		providers = append(providers, provider)
	}

	return providers
}

// GetClientInfo 获取客户端信息
func (m *Manager) GetClientInfo(provider ProviderType) (ClientInfo, error) {
	client, err := m.GetClient(provider)
	if err != nil {
		return ClientInfo{}, err
	}

	return client.GetClientInfo(), nil
}

// 全局管理器实例
var defaultManager = NewManager()

// DefaultManager 获取默认管理器
func DefaultManager() *Manager {
	return defaultManager
}

// 便利函数
func CreateClient(config Config) (LLMClient, error) {
	return defaultManager.CreateClient(config)
}

func GetClient(provider ProviderType) (LLMClient, error) {
	return defaultManager.GetClient(provider)
}

func GetOrCreateClient(config Config) (LLMClient, error) {
	return defaultManager.GetOrCreateClient(config)
}

func QuickChat(ctx context.Context, provider ProviderType, message string) (string, error) {
	return defaultManager.QuickChat(ctx, provider, message)
}

func QuickChatWithTemplate(ctx context.Context, provider ProviderType, template string, variables map[string]string) (string, error) {
	return defaultManager.QuickChatWithTemplate(ctx, provider, template, variables)
}

func QuickChatWithRAG(ctx context.Context, provider ProviderType, query string, docs []RetrievedDocument) (string, error) {
	return defaultManager.QuickChatWithRAG(ctx, provider, query, docs)
}

func QuickJSONChat(ctx context.Context, provider ProviderType, userMessage string) (string, error) {
	return defaultManager.QuickJSONChat(ctx, provider, userMessage)
}

func GenerateCode(ctx context.Context, provider ProviderType, userRequest string) (string, error) {
	return defaultManager.GenerateCode(ctx, provider, userRequest)
}

func GenerateStructuredData(ctx context.Context, provider ProviderType, query string, schema string) (string, error) {
	return defaultManager.GenerateStructuredData(ctx, provider, query, schema)
}
