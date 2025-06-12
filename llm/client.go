package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// LLMClient 大模型客户端通用接口
type LLMClient interface {
	// 聊天完成接口
	ChatCompletion(ctx context.Context, req *ChatCompletionRequest) (*ChatCompletionResponse, error)

	// 获取客户端信息
	GetClientInfo() ClientInfo

	// 健康检查
	HealthCheck(ctx context.Context) error
}

// LLMClientFactory 客户端工厂接口
type LLMClientFactory interface {
	CreateClient(config Config) (LLMClient, error)
	SupportedModels() []string
}

// RAGProvider RAG提供者接口
type RAGProvider interface {
	// 检索相关文档
	RetrieveDocuments(ctx context.Context, query string, maxDocs int, minScore float64) ([]RetrievedDocument, error)

	// 构建RAG提示词
	BuildRAGPrompt(ctx context.Context, userQuery string, docs []RetrievedDocument, template string) (string, error)
}

// PromptManager 提示词管理器接口
type PromptManager interface {
	// 获取提示词模板
	GetTemplate(name string) (*PromptTemplate, error)

	// 渲染提示词
	RenderTemplate(template string, variables map[string]string) (string, error)

	// 保存提示词模板
	SaveTemplate(template *PromptTemplate) error
}

// ChatWithResult 通用聊天方法，支持多种返回类型
func ChatWithResult(ctx context.Context, provider ProviderType, userMessage string, result interface{}) error {
	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	// 检查result参数类型
	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	// 根据目标类型决定是否使用JSON模式
	var req *ChatCompletionRequest
	targetType := resultValue.Elem().Type()

	if targetType.Kind() == reflect.String {
		// 如果目标是字符串，使用普通模式
		req = NewChatRequest("", NewUserMessage(userMessage))
	} else {
		// 如果目标是结构体，使用JSON模式
		systemPrompt := fmt.Sprintf(`请以JSON格式返回响应。用户输入：%s

要求：返回有效的JSON格式数据。`, userMessage)
		req = &ChatCompletionRequest{
			Messages: []Message{
				NewSystemMessage(systemPrompt),
				NewUserMessage(userMessage),
			},
			ResponseFormat: &ResponseFormat{
				Type: "json_object",
			},
			Temperature: 0.1,
		}
	}

	// 发送请求
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("没有响应结果")
	}

	content := resp.Choices[0].Message.Content

	// 根据目标类型处理结果
	if targetType.Kind() == reflect.String {
		// 直接赋值字符串
		resultValue.Elem().SetString(content)
	} else {
		// JSON反序列化
		if err := json.Unmarshal([]byte(content), result); err != nil {
			return fmt.Errorf("JSON反序列化失败: %w", err)
		}
	}

	return nil
}

// ChatWithJSONResult JSON模式聊天，强制返回JSON并反序列化
func ChatWithJSONResult(ctx context.Context, provider ProviderType, userMessage string, result interface{}) error {
	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	// 检查result参数
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	// 构建JSON模式请求
	systemPrompt := fmt.Sprintf(`请严格按照JSON格式返回响应。用户请求：%s

要求：
1. 必须返回有效的JSON格式
2. 确保JSON结构完整
3. 不要包含其他文本`, userMessage)

	req := &ChatCompletionRequest{
		Messages: []Message{
			NewSystemMessage(systemPrompt),
			NewUserMessage(userMessage),
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.1,
	}

	// 发送请求
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("没有响应结果")
	}

	content := resp.Choices[0].Message.Content
	if content == "" {
		return fmt.Errorf("响应内容为空，这是DeepSeek JSON模式的已知问题，请尝试调整prompt")
	}

	// JSON反序列化
	if err := json.Unmarshal([]byte(content), result); err != nil {
		return fmt.Errorf("JSON反序列化失败，原始内容: %s, 错误: %w", content, err)
	}

	return nil
}

// ChatWithStringResult 普通模式聊天，返回字符串
func ChatWithStringResult(ctx context.Context, provider ProviderType, userMessage string, result *string) error {
	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	req := NewChatRequest("", NewUserMessage(userMessage))

	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("没有响应结果")
	}

	*result = resp.Choices[0].Message.Content
	return nil
}

// ChatWithCustomPrompt 自定义prompt的聊天方法
func ChatWithCustomPrompt(ctx context.Context, provider ProviderType, systemPrompt, userMessage string, result interface{}, useJSON bool) error {
	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	// 构建请求
	req := &ChatCompletionRequest{
		Messages: []Message{
			NewSystemMessage(systemPrompt),
			NewUserMessage(userMessage),
		},
		Temperature: 0.1,
	}

	// 是否使用JSON模式
	if useJSON {
		req.ResponseFormat = &ResponseFormat{
			Type: "json_object",
		}
	}

	// 发送请求
	resp, err := client.ChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}

	if len(resp.Choices) == 0 {
		return fmt.Errorf("没有响应结果")
	}

	content := resp.Choices[0].Message.Content

	// 处理结果
	if result == nil {
		return nil // 如果不需要返回结果
	}

	resultValue := reflect.ValueOf(result)
	if resultValue.Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	if useJSON {
		// JSON反序列化
		if err := json.Unmarshal([]byte(content), result); err != nil {
			return fmt.Errorf("JSON反序列化失败: %w", err)
		}
	} else {
		// 字符串赋值
		if resultValue.Elem().Type().Kind() == reflect.String {
			resultValue.Elem().SetString(content)
		} else {
			return fmt.Errorf("非JSON模式下result必须是字符串指针")
		}
	}

	return nil
}

// 便利函数
func NewChatRequest(model string, messages ...Message) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}
}

// NewJSONRequest 创建JSON格式输出的请求
func NewJSONRequest(model string, messages ...Message) *ChatCompletionRequest {
	return &ChatCompletionRequest{
		Model:    model,
		Messages: messages,
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
	}
}

// NewCodeGenRequest 创建代码生成请求（JSON格式）
func NewCodeGenRequest(model string, userRequest string) *ChatCompletionRequest {
	systemPrompt := `你是一个专业的Go语言代码生成器。请根据用户需求生成完整的Go代码。

要求：
1. 返回的代码必须是有效的JSON格式
2. JSON结构如下：
{
  "code": "完整的Go代码内容",
  "filename": "建议的文件名",
  "description": "代码功能描述",
  "dependencies": ["依赖包列表"]
}

请确保生成的代码符合Go语言规范，包含必要的import语句和错误处理。`

	return &ChatCompletionRequest{
		Model: model,
		Messages: []Message{
			NewSystemMessage(systemPrompt),
			NewUserMessage(userRequest),
		},
		ResponseFormat: &ResponseFormat{
			Type: "json_object",
		},
		Temperature: 0.1,
	}
}

// NewStructuredRequest 创建结构化数据请求
func NewStructuredRequest(model string, userQuery string, jsonSchema string) *ChatCompletionRequest {
	systemPrompt := "请根据用户要求返回结构化的JSON数据。严格按照指定的JSON格式返回。"
	if jsonSchema != "" {
		systemPrompt += "\n\nJSON格式要求：\n" + jsonSchema
	}

	return &ChatCompletionRequest{
		Model: model,
		Messages: []Message{
			NewSystemMessage(systemPrompt),
			NewUserMessage(userQuery),
		},
		ResponseFormat: &ResponseFormat{
			Type:   "json_object",
			Schema: jsonSchema,
		},
		Temperature: 0.1,
	}
}

func NewMessage(role, content string) Message {
	return Message{
		Role:    role,
		Content: content,
	}
}

func NewUserMessage(content string) Message {
	return NewMessage("user", content)
}

func NewSystemMessage(content string) Message {
	return NewMessage("system", content)
}

func NewAssistantMessage(content string) Message {
	return NewMessage("assistant", content)
}
