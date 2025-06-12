package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// JSONSchemaProperty JSON Schema属性定义
type JSONSchemaProperty struct {
	Type        string                         `json:"type,omitempty"`
	Description string                         `json:"description,omitempty"`
	Items       *JSONSchemaProperty            `json:"items,omitempty"`
	Properties  map[string]*JSONSchemaProperty `json:"properties,omitempty"`
	Required    []string                       `json:"required,omitempty"`
	Example     interface{}                    `json:"example,omitempty"`
}

// JSONSchema JSON Schema定义
type JSONSchema struct {
	Type       string                         `json:"type"`
	Properties map[string]*JSONSchemaProperty `json:"properties,omitempty"`
	Required   []string                       `json:"required,omitempty"`
}

// GenerateJSONSchema 从结构体生成JSON Schema
func GenerateJSONSchema(structType interface{}) (string, error) {
	t := reflect.TypeOf(structType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return "", fmt.Errorf("参数必须是结构体类型")
	}

	schema := &JSONSchema{
		Type:       "object",
		Properties: make(map[string]*JSONSchemaProperty),
		Required:   []string{},
	}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// 跳过非导出字段
		if !field.IsExported() {
			continue
		}

		// 检查llm标签是否为"-"
		llmTag := field.Tag.Get("llm")
		if llmTag == "-" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		// 解析json标签
		jsonName := field.Name
		required := true
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				jsonName = parts[0]
			}
			// 检查omitempty
			for _, part := range parts[1:] {
				if part == "omitempty" {
					required = false
				}
			}
		}

		// 生成属性
		prop := generateProperty(field.Type, field.Tag)
		if prop == nil {
			// generateProperty返回nil表示应该忽略此字段
			continue
		}
		schema.Properties[jsonName] = prop

		// 添加到required列表
		if required {
			schema.Required = append(schema.Required, jsonName)
		}
	}

	// 转换为JSON字符串
	schemaBytes, err := json.Marshal(schema)
	if err != nil {
		return "", fmt.Errorf("生成JSON Schema失败: %w", err)
	}

	return string(schemaBytes), nil
}

// generateProperty 生成单个属性的Schema
func generateProperty(fieldType reflect.Type, tag reflect.StructTag) *JSONSchemaProperty {
	prop := &JSONSchemaProperty{}

	// 优先处理llm标签
	llmTag := tag.Get("llm")
	if llmTag == "-" {
		// llm:"-" 表示忽略此字段
		return nil
	}

	// 解析llm标签中的desc
	if strings.HasPrefix(llmTag, "desc:") {
		prop.Description = strings.TrimPrefix(llmTag, "desc:")
	}

	// 如果没有llm标签的描述，回退到description标签
	if prop.Description == "" {
		if desc := tag.Get("description"); desc != "" {
			prop.Description = desc
		}
	}

	// 获取示例（仍然支持example标签）
	if example := tag.Get("example"); example != "" {
		prop.Example = example
	}

	// 处理指针类型
	if fieldType.Kind() == reflect.Ptr {
		fieldType = fieldType.Elem()
	}

	// 根据Go类型设置JSON Schema类型
	switch fieldType.Kind() {
	case reflect.String:
		prop.Type = "string"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		prop.Type = "integer"
	case reflect.Float32, reflect.Float64:
		prop.Type = "number"
	case reflect.Bool:
		prop.Type = "boolean"
	case reflect.Slice, reflect.Array:
		prop.Type = "array"
		// 递归处理数组元素类型
		elemType := fieldType.Elem()
		if elemType.Kind() == reflect.Struct {
			// 如果是结构体数组，递归生成
			elemProp := generateStructProperty(elemType)
			prop.Items = elemProp
		} else {
			// 基础类型数组
			prop.Items = generateProperty(elemType, "")
		}
	case reflect.Struct:
		prop = generateStructProperty(fieldType)
	case reflect.Map:
		prop.Type = "object"
	default:
		prop.Type = "string" // 默认为string
	}

	return prop
}

// generateStructProperty 生成结构体属性
func generateStructProperty(structType reflect.Type) *JSONSchemaProperty {
	prop := &JSONSchemaProperty{
		Type:       "object",
		Properties: make(map[string]*JSONSchemaProperty),
		Required:   []string{},
	}

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)

		if !field.IsExported() {
			continue
		}

		// 检查llm标签是否为"-"
		llmTag := field.Tag.Get("llm")
		if llmTag == "-" {
			continue
		}

		jsonTag := field.Tag.Get("json")
		if jsonTag == "-" {
			continue
		}

		jsonName := field.Name
		required := true
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			if parts[0] != "" {
				jsonName = parts[0]
			}
			for _, part := range parts[1:] {
				if part == "omitempty" {
					required = false
				}
			}
		}

		fieldProp := generateProperty(field.Type, field.Tag)
		if fieldProp == nil {
			// generateProperty返回nil表示应该忽略此字段
			continue
		}
		prop.Properties[jsonName] = fieldProp

		if required {
			prop.Required = append(prop.Required, jsonName)
		}
	}

	return prop
}

// ChatWithStruct 使用结构体模板进行聊天（简化版本，只需传递一个参数）
func ChatWithStruct(ctx context.Context, provider ProviderType, userMessage string, result interface{}) error {
	// 检查result参数
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	// 获取指针指向的类型作为模板
	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
	}

	// 创建该类型的零值实例作为模板
	templateValue := reflect.New(resultType).Elem().Interface()

	// 调用原有的模板方法
	return ChatWithStructTemplate(ctx, provider, userMessage, templateValue, result)
}

// ChatWithStructTemplate 使用结构体模板进行聊天
func ChatWithStructTemplate(ctx context.Context, provider ProviderType, userMessage string, templateStruct interface{}, result interface{}) error {
	// 生成JSON Schema
	schema, err := GenerateJSONSchema(templateStruct)
	if err != nil {
		return fmt.Errorf("生成JSON Schema失败: %w", err)
	}

	// 调用结构化数据生成
	return ChatWithStructuredSchema(ctx, provider, userMessage, schema, result)
}

// ChatWithStructuredSchema 使用JSON Schema的结构化聊天
func ChatWithStructuredSchema(ctx context.Context, provider ProviderType, userMessage, schema string, result interface{}) error {
	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	// 检查result参数
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	// 构建系统提示词
	systemPrompt := fmt.Sprintf(`请严格按照指定的JSON Schema返回响应。

用户请求：%s

JSON Schema要求：
%s

要求：
1. 必须返回有效的JSON格式
2. 严格遵循Schema定义
3. 确保所有required字段都有值
4. 不要包含Schema之外的字段`, userMessage, schema)

	req := &ChatCompletionRequest{
		Messages: []Message{
			NewSystemMessage(systemPrompt),
			NewUserMessage(userMessage),
		},
		ResponseFormat: &ResponseFormat{
			Type:   "json_object",
			Schema: schema,
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

// ChatWithStructMessages 使用完整消息数组进行结构体生成
func ChatWithStructMessages(ctx context.Context, provider ProviderType, messages []Message, result interface{}) error {
	// 检查result参数
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return fmt.Errorf("result参数必须是指针类型")
	}

	// 获取指针指向的类型作为模板
	resultType := reflect.TypeOf(result)
	if resultType.Kind() == reflect.Ptr {
		resultType = resultType.Elem()
	}

	// 创建该类型的零值实例作为模板
	templateValue := reflect.New(resultType).Elem().Interface()

	// 生成JSON Schema
	schema, err := GenerateJSONSchema(templateValue)
	if err != nil {
		return fmt.Errorf("生成JSON Schema失败: %w", err)
	}

	client, err := GetClient(provider)
	if err != nil {
		return fmt.Errorf("获取客户端失败: %w", err)
	}

	// 构建请求，在现有消息基础上添加Schema说明
	var finalMessages []Message

	// 添加JSON Schema系统提示
	schemaPrompt := fmt.Sprintf(`请严格按照以下JSON Schema返回响应：

JSON Schema:
%s

要求：
1. 必须返回有效的JSON格式
2. 严格遵循Schema定义
3. 确保所有required字段都有值
4. 不要包含Schema之外的字段`, schema)

	// 如果第一个消息不是系统消息，添加系统消息
	if len(messages) == 0 || messages[0].Role != "system" {
		finalMessages = append(finalMessages, NewSystemMessage(schemaPrompt))
		finalMessages = append(finalMessages, messages...)
	} else {
		// 如果已有系统消息，合并Schema说明
		existingSystem := messages[0]
		combinedContent := existingSystem.Content + "\n\n" + schemaPrompt
		finalMessages = append(finalMessages, NewSystemMessage(combinedContent))
		finalMessages = append(finalMessages, messages[1:]...)
	}

	req := &ChatCompletionRequest{
		Messages: finalMessages,
		ResponseFormat: &ResponseFormat{
			Type:   "json_object",
			Schema: schema,
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

// ChatWithStructRAG 使用RAG上下文进行结构体生成
func ChatWithStructRAG(ctx context.Context, provider ProviderType, userQuestion string, ragDocuments []string, result interface{}) error {
	// 构建包含RAG文档的消息
	var messages []Message

	// 系统消息：说明任务和文档
	systemPrompt := fmt.Sprintf(`你是专业的数据分析和结构化生成专家。

用户提供了 %d 个参考文档，请基于这些文档内容回答用户问题。

参考文档：`, len(ragDocuments))

	for i, doc := range ragDocuments {
		systemPrompt += fmt.Sprintf("\n\n--- 文档 %d ---\n%s", i+1, doc)
	}

	systemPrompt += "\n\n请基于以上文档内容，按照指定的JSON格式生成响应。"

	messages = append(messages, NewSystemMessage(systemPrompt))
	messages = append(messages, NewUserMessage(userQuestion))

	return ChatWithStructMessages(ctx, provider, messages, result)
}

// ChatWithStructContext 支持自定义系统提示词的结构体生成
func ChatWithStructContext(ctx context.Context, provider ProviderType, systemPrompt, userMessage string, result interface{}) error {
	messages := []Message{
		NewSystemMessage(systemPrompt),
		NewUserMessage(userMessage),
	}

	return ChatWithStructMessages(ctx, provider, messages, result)
}
