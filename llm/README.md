# LLM 大模型调用包

这是一个通用的大模型调用包，支持多种大模型提供商，目前优先支持 DeepSeek。

## 🚀 **特性**

- ✅ **通用接口设计**：统一的API调用方式，支持多种大模型
- ✅ **DeepSeek 优先支持**：专门优化的DeepSeek客户端实现
- ✅ **RAG 支持**：内置RAG（检索增强生成）功能
- ✅ **提示词模板**：支持模板变量替换
- ✅ **错误重试**：自动重试机制，提高稳定性
- ✅ **类型安全**：完整的类型定义，避免运行时错误
- ✅ **并发安全**：支持并发调用
- ✅ **可扩展**：易于添加新的大模型提供商

## 📦 **包结构**

```
pkg/llm/
├── types.go          # 通用数据类型定义
├── config.go         # 配置管理
├── client.go         # 客户端接口定义
├── manager.go        # 客户端管理器
└── deepseek/
    └── client.go     # DeepSeek具体实现
```

## 🎯 **快速开始**

### 1. 基础使用

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/yunhanshu-net/pkg/llm"
    "github.com/yunhanshu-net/pkg/llm/deepseek"
)

func main() {
    ctx := context.Background()
    
    // 1. 注册DeepSeek工厂
    llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})
    
    // 2. 创建配置
    config := llm.GetDefaultConfig(llm.ProviderDeepSeek)
    config.APIKey = "your-deepseek-api-key"
    
    // 3. 创建客户端
    client, err := llm.CreateClient(config)
    if err != nil {
        log.Fatal(err)
    }
    
    // 4. 发送请求
    req := llm.NewChatRequest("deepseek-coder",
        llm.NewUserMessage("请生成一个Go语言的Hello World程序"),
    )
    
    resp, err := client.ChatCompletion(ctx, req)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("AI回复:", resp.Choices[0].Message.Content)
}
```

### 2. 便利方法

```go
// 快速聊天
answer, err := llm.QuickChat(ctx, llm.ProviderDeepSeek, "什么是Go语言？")

// 使用模板
template := "请为{{language}}语言编写一个{{type}}函数"
variables := map[string]string{
    "language": "Go",
    "type": "HTTP处理",
}
result, err := llm.QuickChatWithTemplate(ctx, llm.ProviderDeepSeek, template, variables)

// 使用RAG
docs := []llm.RetrievedDocument{
    {ID: "doc1", Content: "Go语言文档内容...", Score: 0.95},
}
ragResult, err := llm.QuickChatWithRAG(ctx, llm.ProviderDeepSeek, "Go语言特点？", docs)
```

### 3. 代码生成示例

```go
func generateCode(userRequest string) (string, error) {
    ctx := context.Background()
    
    // 代码生成模板
    template := `基于runner框架生成Go函数文件。

用户需求：{{user_request}}

请生成完整的.go文件代码，包含：
1. Model定义（带gorm标签）
2. FunctionInfo配置
3. Handler处理函数

直接输出代码：`

    variables := map[string]string{
        "user_request": userRequest,
    }
    
    return llm.QuickChatWithTemplate(ctx, llm.ProviderDeepSeek, template, variables)
}
```

### 4. JSON格式输出

```go
// 1. 基础JSON聊天
jsonResult, err := llm.QuickJSONChat(ctx, llm.ProviderDeepSeek, 
    "生成一个用户信息的JSON结构，包含姓名、年龄、邮箱")

// 2. 结构化代码生成（JSON格式）
codeJSON, err := llm.GenerateCode(ctx, llm.ProviderDeepSeek, 
    "创建一个用户管理的CRUD功能")

// 返回格式：
// {
//   "code": "package main\n\nfunc CreateUser()...",
//   "filename": "user_manager.go", 
//   "description": "用户管理CRUD功能",
//   "dependencies": ["gorm.io/gorm", "github.com/gin-gonic/gin"]
// }

// 3. 自定义JSON结构
schema := `{
  "type": "object",
  "properties": {
    "api_list": {
      "type": "array",
      "items": {
        "type": "object", 
        "properties": {
          "name": {"type": "string"},
          "method": {"type": "string"},
          "path": {"type": "string"},
          "description": {"type": "string"}
        }
      }
    }
  }
}`

structuredData, err := llm.GenerateStructuredData(ctx, llm.ProviderDeepSeek,
    "设计一个博客管理系统的RESTful API", schema)
```

### 5. 手动控制JSON格式

```go
// 创建JSON格式请求
req := llm.NewJSONRequest("deepseek-coder",
    llm.NewSystemMessage("返回JSON格式的API设计"),
    llm.NewUserMessage("设计用户注册API"),
)

// 代码生成专用JSON请求
codeReq := llm.NewCodeGenRequest("deepseek-coder", "创建JWT认证中间件")

// 自定义JSON结构请求
customReq := llm.NewStructuredRequest("deepseek-coder", 
    "分析这段代码的复杂度", customSchema)

resp, err := client.ChatCompletion(ctx, req)
```

## ⚙️ **配置说明**

### DeepSeek 配置

```go
config := llm.Config{
    Provider:           llm.ProviderDeepSeek,
    APIKey:            "your-api-key",
    BaseURL:           "https://api.deepseek.com",
    Timeout:           30 * time.Second,
    DefaultModel:      "deepseek-coder",
    DefaultTemperature: 0.1,
    DefaultMaxTokens:  2000,
    DefaultTopP:       0.9,
    MaxRetries:        3,
    RetryInterval:     1 * time.Second,
    EnableRAG:        true,
    MaxRAGDocuments:  5,
    RAGSimilarityMin: 0.7,
    RateLimitRPS:     60,
}
```

### 支持的模型

- `deepseek-coder` - 代码生成专用模型（推荐）
- `deepseek-chat` - 通用对话模型
- `deepseek-reasoner` - 推理专用模型

## 🔧 **高级功能**

### RAG（检索增强生成）

```go
req := &llm.ChatCompletionRequest{
    Model: "deepseek-coder",
    Messages: []llm.Message{
        llm.NewUserMessage("基于文档回答问题"),
    },
    RAGContext: &llm.RAGContext{
        Query: "用户问题",
        RetrievedDocs: []llm.RetrievedDocument{
            {
                ID: "doc1",
                Content: "相关文档内容",
                Score: 0.95,
            },
        },
    },
}
```

### 提示词模板

```go
req := &llm.ChatCompletionRequest{
    Messages: []llm.Message{
        llm.NewUserMessage("请生成代码"),
    },
    PromptTemplate: "为{{language}}生成{{type}}代码",
    Variables: map[string]string{
        "language": "Go",
        "type": "API",
    },
}
```

## 🛡️ **错误处理**

```go
resp, err := client.ChatCompletion(ctx, req)
if err != nil {
    if apiErr, ok := err.(*llm.APIError); ok {
        fmt.Printf("API错误: %d - %s\n", apiErr.Code, apiErr.Message)
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}
```

## 📊 **性能优化建议**

1. **复用客户端**：避免频繁创建客户端实例
2. **合理设置超时**：根据请求复杂度调整超时时间
3. **控制并发**：使用 `MaxConcurrency` 限制并发请求数
4. **优化提示词**：精简提示词，减少token消耗
5. **启用重试**：设置合理的重试次数和间隔

## 🔮 **扩展新提供商**

要添加新的大模型提供商，只需：

1. 实现 `LLMClient` 接口
2. 实现 `LLMClientFactory` 接口  
3. 注册到管理器中

```go
// 实现新提供商
type NewProviderClient struct {
    // ...
}

func (c *NewProviderClient) ChatCompletion(ctx context.Context, req *llm.ChatCompletionRequest) (*llm.ChatCompletionResponse, error) {
    // 具体实现
}

// 注册
llm.DefaultManager().RegisterFactory("new_provider", &NewProviderFactory{})
```

## 📝 **注意事项**

1. **API密钥安全**：不要在代码中硬编码API密钥
2. **速率限制**：注意各提供商的API调用限制
3. **成本控制**：合理设置 `MaxTokens` 控制成本
4. **错误处理**：始终检查错误并进行适当处理
5. **日志记录**：启用日志记录便于调试和监控

## 🚀 **新增：通用调用方法**

### 自动类型选择

根据[DeepSeek官方JSON模式文档](https://api-docs.deepseek.com/zh-cn/guides/json_mode)，我们实现了智能的通用调用方法：

```go
// 智能选择：根据目标类型自动选择JSON或文本模式
var textResult string
err := llm.ChatWithResult(ctx, llm.ProviderDeepSeek, "介绍Go语言", &textResult) // 自动使用文本模式

type UserInfo struct {
    Name  string `json:"name"`
    Age   int    `json:"age"`
    Email string `json:"email"`
}
var userInfo UserInfo
err = llm.ChatWithResult(ctx, llm.ProviderDeepSeek, "生成用户信息", &userInfo) // 自动使用JSON模式
```

### 结构体标签方法（推荐！）

最优雅的方式是使用结构体标签定义JSON结构，现在支持专门的 `llm` 标签：

```go
// 使用专门的llm标签（推荐方式）
type Bd struct {
    ItemId      string `json:"itemId" llm:"desc:值班列表的id"`
    GroupNotice string `json:"groupNotice" llm:"desc:服务组的通知"`
    ID          string `json:"id" llm:"-"`                      // llm:"-" 忽略字段
    CreateTime  string `json:"createTime,omitempty" llm:"desc:创建时间"`
}

// 混合使用（向后兼容）
type UserProfile struct {
    Name     string   `json:"name" llm:"desc:用户姓名"`
    Age      int      `json:"age" description:"用户年龄"`        // 传统标签仍支持
    Email    string   `json:"email" llm:"desc:邮箱地址"`
    Phone    string   `json:"phone,omitempty" llm:"desc:手机号码"`
    Password string   `json:"password" llm:"-"`                // 忽略敏感字段
}

// 直接使用结构体作为模板（推荐：简化版本）
var bd Bd
err := llm.ChatWithStruct(ctx, llm.ProviderDeepSeek,
    "生成一个值班管理系统的数据", &bd)  // 只需传递一个参数！

// 完整版本（如果需要自定义模板）
var bd2 Bd
err = llm.ChatWithStructTemplate(ctx, llm.ProviderDeepSeek,
    "生成值班数据", Bd{}, &bd2)
```

### 支持的标签语法

```go
// LLM专用标签（推荐）
type Example struct {
    Field1 string `json:"field1" llm:"desc:字段描述"`    // 添加字段描述
    Field2 string `json:"field2" llm:"-"`              // 忽略此字段
    Field3 string `json:"field3" llm:"desc:描述" example:"示例值"` // 可与example标签组合
}

// 传统标签（向后兼容）
type Legacy struct {
    Field1 string `json:"field1" description:"字段描述" example:"示例值"`
}

// 标签优先级：
// 1. llm:"desc:xxx"  - 优先使用LLM标签描述
// 2. description:"xxx" - 如果没有LLM标签，使用传统描述标签
// 3. example:"xxx"   - 示例值标签（任何情况下都支持）
```

### 复杂结构体示例

```go
// API设计结构体
type APISpec struct {
    Title       string            `json:"title" description:"API标题"`
    Version     string            `json:"version" description:"API版本"`
    Description string            `json:"description" description:"API描述"`
    Endpoints   []APIEndpointSpec `json:"endpoints" description:"API端点列表"`
}

type APIEndpointSpec struct {
    Name        string            `json:"name" description:"端点名称"`
    Method      string            `json:"method" description:"HTTP方法"`
    Path        string            `json:"path" description:"请求路径"`
    Description string            `json:"description" description:"端点描述"`
    Parameters  []APIParameter    `json:"parameters,omitempty" description:"请求参数"`
    Responses   map[string]string `json:"responses,omitempty" description:"响应说明"`
}

// 使用复杂结构体
var apiSpec APISpec
err := llm.ChatWithStructTemplate(ctx, llm.ProviderDeepSeek,
    "设计一个博客管理系统的RESTful API，包含文章的增删改查功能",
    APISpec{}, &apiSpec)
```

### 专用方法

```go
// 1. 强制JSON模式（自动反序列化）
var apiDesign APIDesign
err := llm.ChatWithJSONResult(ctx, llm.ProviderDeepSeek, 
    "设计用户管理API", &apiDesign)

// 2. 纯文本模式
var description string
err = llm.ChatWithStringResult(ctx, llm.ProviderDeepSeek, 
    "解释RESTful API", &description)

// 3. 自定义系统提示词
var codeResult CodeResult
err = llm.ChatWithCustomPrompt(ctx, llm.ProviderDeepSeek,
    "你是代码生成专家，返回JSON格式", 
    "创建HTTP服务器", &codeResult, true)

// 4. 结构体模板方法（推荐）
var userProfile UserProfile
err = llm.ChatWithStructTemplate(ctx, llm.ProviderDeepSeek,
    "生成用户信息", UserProfile{}, &userProfile)

// 5. 自定义Schema方法
schema := `{"type": "object", "properties": {...}}`
err = llm.ChatWithStructuredSchema(ctx, llm.ProviderDeepSeek,
    "生成数据", schema, &result)
```

### 错误处理增强

```go
// JSON解析失败时提供详细错误信息
var result MyStruct
err := llm.ChatWithJSONResult(ctx, llm.ProviderDeepSeek, "生成数据", &result)
if err != nil {
    // 错误信息包含：原始内容、解析错误、DeepSeek已知问题提示
    log.Printf("详细错误: %v", err)
}
```

### 特性对比

| 方法 | 自动JSON检测 | 类型安全 | 错误处理 | Schema生成 | 使用场景 |
|------|-------------|----------|----------|------------|----------|
| `ChatWithResult` | ✅ | ✅ | ✅ | ❌ | 智能选择，推荐使用 |
| `ChatWithStructTemplate` | N/A | ✅ | ✅ | ✅ | **结构体模板，最推荐** |
| `ChatWithJSONResult` | N/A | ✅ | ✅ | ❌ | 强制JSON，结构化数据 |
| `ChatWithStringResult` | N/A | ✅ | ✅ | ❌ | 纯文本，说明性内容 |
| `ChatWithCustomPrompt` | 可选 | ✅ | ✅ | ❌ | 高度自定义场景 |
| `ChatWithStructuredSchema` | N/A | ✅ | ✅ | 手动 | 手写Schema场景 |
| `QuickChat` | ❌ | ❌ | 基础 | ❌ | 简单快速测试 |

这些新方法完全符合DeepSeek的JSON模式要求，自动处理prompt中的"json"关键词和格式要求！ 