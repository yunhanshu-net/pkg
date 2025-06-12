# LLM 标签使用指南

## 概述

我们已经实现了专门的 `llm` 标签支持，让你可以直接在结构体上定义字段的LLM描述信息。这种方式比传统的 `description` 标签更专业，明确表示是为LLM服务的。

## 标签语法

### 1. 基础用法

```go
type Bd struct {
    ItemId      string `json:"itemId" llm:"desc:值班列表的id"`
    GroupNotice string `json:"groupNotice" llm:"desc:服务组的通知"`
    ID          string `json:"id" llm:"-"`                          // 忽略此字段
    CreateTime  string `json:"createTime,omitempty" llm:"desc:创建时间"`
}
```

### 2. 支持的llm标签格式

- `llm:"desc:字段描述"` - 为字段添加描述信息，会包含在生成的JSON Schema中
- `llm:"-"` - 忽略此字段，不包含在生成的JSON Schema中

### 3. 标签组合使用

```go
type MixedExample struct {
    Name        string `json:"name" llm:"desc:用户姓名"`
    Age         int    `json:"age" description:"用户年龄"`           // 传统标签仍支持
    Email       string `json:"email" llm:"desc:邮箱地址"`
    Phone       string `json:"phone,omitempty" llm:"desc:手机号码"`
    InternalID  int    `json:"-" llm:"-"`                         // 完全忽略
    Password    string `json:"password" llm:"-"`                  // JSON中有但LLM忽略
    Description string `json:"description,omitempty" llm:"desc:个人描述" example:"热爱编程"`
}
```

## 标签优先级

1. `llm:"desc:xxx"` - **优先使用**，LLM专用标签描述
2. `description:"xxx"` - 如果没有LLM标签描述，则使用传统描述标签
3. `example:"xxx"` - 示例值标签，任何情况下都支持，可与llm标签组合

## 使用示例

### 生成JSON Schema

```go
// 自动生成JSON Schema
schema, err := llm.GenerateJSONSchema(Bd{})
if err != nil {
    log.Fatal(err)
}
fmt.Println(schema)

// 输出的Schema会：
// 1. 包含 llm:"desc:xxx" 字段的描述
// 2. 忽略 llm:"-" 的字段
// 3. 保持JSON结构完整
```

### 使用结构体模板生成数据

```go
// 方式1: 简化版本（推荐！只需传递一个参数）
var bd Bd
err := llm.ChatWithStruct(ctx, llm.ProviderDeepSeek,
    "生成一个值班管理系统的数据", &bd)  // 只传递结果指针

if err != nil {
    log.Fatal(err)
}

// 方式2: 完整版本（如果需要不同的模板和结果类型）
var bd2 Bd
err = llm.ChatWithStructTemplate(ctx, llm.ProviderDeepSeek,
    "生成一个值班管理系统的数据",
    Bd{}, // 模板结构体
    &bd2) // 结果存储

// bd.ID 字段会是空值，因为使用了 llm:"-"
// 其他字段会根据 llm:"desc:xxx" 的描述生成合适的值
```

## 实际应用场景

### 1. 忽略系统字段

```go
type User struct {
    ID         int       `json:"id" llm:"-"`                    // 系统生成，不需要AI生成
    UUID       string    `json:"uuid" llm:"-"`                 // 系统生成
    Name       string    `json:"name" llm:"desc:用户真实姓名"`
    Email      string    `json:"email" llm:"desc:邮箱地址"`
    CreatedAt  time.Time `json:"created_at" llm:"-"`           // 系统时间戳
    UpdatedAt  time.Time `json:"updated_at" llm:"-"`           // 系统时间戳
}
```

### 2. 忽略敏感信息

```go
type UserAuth struct {
    Username     string `json:"username" llm:"desc:用户名"`
    Email        string `json:"email" llm:"desc:邮箱地址"`
    Password     string `json:"password" llm:"-"`              // 敏感信息，不让AI生成
    Salt         string `json:"salt" llm:"-"`                  // 安全相关
    LastLoginIP  string `json:"last_login_ip" llm:"-"`         // 隐私信息
    Profile      string `json:"profile" llm:"desc:个人简介"`
}
```

### 3. 业务数据生成

```go
type Product struct {
    SKU         string  `json:"sku" llm:"desc:商品SKU编码"`
    Name        string  `json:"name" llm:"desc:商品名称"`
    Description string  `json:"description" llm:"desc:商品详细描述"`
    Price       float64 `json:"price" llm:"desc:商品价格（元）"`
    Category    string  `json:"category" llm:"desc:商品分类"`
    InStock     bool    `json:"in_stock" llm:"desc:是否有库存"`
    InternalID  string  `json:"-" llm:"-"`                     // 内部字段，完全不暴露
}
```

## 迁移指南

如果你之前使用的是 `description` 标签，可以逐步迁移：

```go
// 旧方式
type OldStyle struct {
    Name string `json:"name" description:"用户姓名"`
    Age  int    `json:"age" description:"用户年龄"`
}

// 新方式（推荐）
type NewStyle struct {
    Name string `json:"name" llm:"desc:用户姓名"`
    Age  int    `json:"age" llm:"desc:用户年龄"`
}

// 混合方式（过渡期）
type MixedStyle struct {
    Name string `json:"name" llm:"desc:用户姓名"`           // 新标签
    Age  int    `json:"age" description:"用户年龄"`          // 旧标签仍然有效
}
```

## 优势

1. **专业性**：`llm` 标签明确表示是为LLM服务的
2. **灵活性**：支持字段级别的忽略控制（`llm:"-"`）
3. **兼容性**：完全向后兼容传统 `description` 标签
4. **清晰性**：代码意图更明确，维护性更好
5. **控制力**：可以精确控制哪些字段参与LLM生成

这种设计完全符合你的需求，既专业又实用！

## 高级用法

### 多种API选择

根据不同场景，我们提供了4种API方法：

#### 1. 简单场景 - `ChatWithStruct` ⭐️ **推荐**
```go
// 最简单的用法，只需4个参数
var bd Bd
err := llm.ChatWithStruct(ctx, llm.ProviderDeepSeek,
    "生成一个值班管理系统的数据", &bd)
```

#### 2. RAG场景 - `ChatWithStructRAG` 🔥 **RAG专用**
```go
// 从数据库获取的多个文档
ragDocs := []string{
    "文档1：API设计规范...",
    "文档2：业务需求说明...",
    "文档3：数据库设计标准...",
}

var apiDesign APISpec
err := llm.ChatWithStructRAG(ctx, llm.ProviderDeepSeek,
    "根据这些文档设计一个用户管理API", ragDocs, &apiDesign)
```

#### 3. 自定义系统提示词 - `ChatWithStructContext`
```go
systemPrompt := `你是数据库设计专家，设计时需要遵循：
1. 所有表都需要主键
2. 支持软删除
3. 包含审计字段`

var dbTable DatabaseTable
err := llm.ChatWithStructContext(ctx, llm.ProviderDeepSeek,
    systemPrompt, "设计订单表结构", &dbTable)
```

#### 4. 完全自定义 - `ChatWithStructMessages`
```go
// 支持完整的对话历史
messages := []llm.Message{
    llm.NewSystemMessage("你是Go语言专家"),
    llm.NewUserMessage("我需要设计博客系统"),
    llm.NewAssistantMessage("好的，你需要哪些具体功能？"),
    llm.NewUserMessage("需要文章的CRUD操作，生成数据结构"),
}

var blogPost BlogPost
err := llm.ChatWithStructMessages(ctx, llm.ProviderDeepSeek,
    messages, &blogPost)
```

### 🎯 API选择指南

| 场景 | 推荐API | 优势 | 示例 |
|------|--------|------|------|
| 简单数据生成 | `ChatWithStruct` | 最简单 | 生成用户信息、产品数据 |
| 基于文档生成 | `ChatWithStructRAG` | RAG优化 | 根据需求文档生成API |
| 需要专业角色 | `ChatWithStructContext` | 系统提示词 | 数据库设计、代码规范 |
| 复杂交互 | `ChatWithStructMessages` | 最灵活 | 多轮对话、完整上下文 |

### 🔥 RAG使用场景

特别适合以下情况：
- 从数据库查询到的多个相关文档
- 需要基于历史数据生成新内容
- 文档内容超出单个消息长度限制
- 需要综合多个信息源的知识

```go
// 实际业务示例：基于历史API文档生成新API
func GenerateAPIFromDocs(ctx context.Context, requirement string) (*APISpec, error) {
    // 从数据库获取相关文档
    docs, err := database.GetRelatedDocs(requirement)
    if err != nil {
        return nil, err
    }
    
    var api APISpec
    err = llm.ChatWithStructRAG(ctx, llm.ProviderDeepSeek,
        requirement, docs, &api)
    
    return &api, err
}
``` 