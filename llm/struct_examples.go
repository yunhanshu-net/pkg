package llm

import (
	"context"
	"fmt"
	"log"
)

// 新的llm标签使用示例
type Bd struct {
	ItemId      string `json:"itemId" llm:"desc:值班列表的id"`
	GroupNotice string `json:"groupNotice" llm:"desc:服务组的通知"`
	ID          string `json:"id" llm:"-"` // llm:"-" 表示忽略此字段
	CreateTime  string `json:"createTime,omitempty" llm:"desc:创建时间"`
}

// 混合使用llm标签和传统标签的示例
type MixedExample struct {
	Name        string `json:"name" llm:"desc:用户姓名"`
	Age         int    `json:"age" description:"用户年龄"` // 传统description标签仍然支持
	Email       string `json:"email" llm:"desc:邮箱地址"`
	Phone       string `json:"phone,omitempty" llm:"desc:手机号码"`
	InternalID  int    `json:"-" llm:"-"`        // 完全忽略
	Password    string `json:"password" llm:"-"` // JSON中有但LLM忽略
	Description string `json:"description,omitempty" llm:"desc:个人描述" example:"热爱编程"`
}

// 用户信息结构体 - 使用标签定义JSON结构
type UserProfile struct {
	Name    string   `json:"name" description:"用户姓名" example:"张三"`
	Age     int      `json:"age" description:"用户年龄" example:"25"`
	Email   string   `json:"email" description:"邮箱地址" example:"zhangsan@example.com"`
	Phone   string   `json:"phone,omitempty" description:"手机号码"`
	Address string   `json:"address,omitempty" description:"家庭住址"`
	Hobbies []string `json:"hobbies,omitempty" description:"兴趣爱好列表"`
}

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

type APIParameter struct {
	Name        string `json:"name" description:"参数名称"`
	Type        string `json:"type" description:"参数类型"`
	Required    bool   `json:"required" description:"是否必填"`
	Description string `json:"description" description:"参数描述"`
}

// 代码生成结构体
type CodeGenSpec struct {
	Language     string     `json:"language" description:"编程语言"`
	Framework    string     `json:"framework,omitempty" description:"使用的框架"`
	Files        []CodeFile `json:"files" description:"生成的文件列表"`
	Dependencies []string   `json:"dependencies,omitempty" description:"依赖包列表"`
	Instructions string     `json:"instructions,omitempty" description:"使用说明"`
}

type CodeFile struct {
	Filename    string `json:"filename" description:"文件名"`
	Content     string `json:"content" description:"文件内容"`
	Description string `json:"description" description:"文件用途说明"`
}

// 数据库设计结构体
type DatabaseSchema struct {
	Name        string  `json:"name" description:"数据库名称"`
	Description string  `json:"description" description:"数据库描述"`
	Tables      []Table `json:"tables" description:"数据表列表"`
}

type Table struct {
	Name        string   `json:"name" description:"表名"`
	Description string   `json:"description" description:"表描述"`
	Fields      []Field  `json:"fields" description:"字段列表"`
	Indexes     []string `json:"indexes,omitempty" description:"索引列表"`
}

type Field struct {
	Name     string `json:"name" description:"字段名"`
	Type     string `json:"type" description:"字段类型"`
	Length   int    `json:"length,omitempty" description:"字段长度"`
	Nullable bool   `json:"nullable" description:"是否允许为空"`
	Default  string `json:"default,omitempty" description:"默认值"`
	Comment  string `json:"comment,omitempty" description:"字段注释"`
}

// FunctionGoCodeGen function-go框架代码生成结构体
type FunctionGoCodeGen struct {
	ModuleName   string               `json:"module_name" llm:"desc:模块名称，如user、order等"`
	Description  string               `json:"description" llm:"desc:模块功能描述"`
	ApiCode      string               `json:"api_code" llm:"desc:API层代码，包含gin路由处理"`
	ServiceCode  string               `json:"service_code" llm:"desc:Service层代码，包含业务逻辑"`
	RepoCode     string               `json:"repo_code" llm:"desc:Repo层代码，包含数据库操作"`
	ModelCode    string               `json:"model_code" llm:"desc:Model定义代码，包含gorm标签"`
	DTOCode      string               `json:"dto_code" llm:"desc:DTO结构体代码，请求响应参数"`
	FunctionInfo string               `json:"function_info" llm:"desc:FunctionInfo配置代码"`
	TestCode     string               `json:"test_code,omitempty" llm:"desc:单元测试代码"`
	Dependencies []string             `json:"dependencies,omitempty" llm:"desc:需要的依赖包"`
	Instructions string               `json:"instructions,omitempty" llm:"desc:部署和使用说明"`
	Files        []FunctionGoCodeFile `json:"files" llm:"desc:生成的文件列表"`
}

type FunctionGoCodeFile struct {
	Path        string `json:"path" llm:"desc:文件路径，如api/v1/user.go"`
	Content     string `json:"content" llm:"desc:完整的文件内容"`
	Description string `json:"description" llm:"desc:文件功能说明"`
	Type        string `json:"type" llm:"desc:文件类型：api/service/repo/model/dto"`
}

// ExampleStructTemplateUsage 演示使用结构体模板的示例
func ExampleStructTemplateUsage() {
	ctx := context.Background()
	provider := ProviderDeepSeek

	fmt.Println("=== 结构体模板使用示例 ===")

	// 1. 用户信息生成
	fmt.Println("\n1. 用户信息生成:")
	var userProfile UserProfile
	err := ChatWithStructTemplate(ctx, provider,
		"生成一个软件工程师的用户信息",
		UserProfile{}, // 模板结构体
		&userProfile)  // 结果存储

	if err != nil {
		log.Printf("用户信息生成失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的用户信息:\n")
		fmt.Printf("  姓名: %s\n", userProfile.Name)
		fmt.Printf("  年龄: %d\n", userProfile.Age)
		fmt.Printf("  邮箱: %s\n", userProfile.Email)
		fmt.Printf("  爱好: %v\n", userProfile.Hobbies)
	}

	// 2. API设计
	fmt.Println("\n2. API设计:")
	var apiSpec APISpec
	err = ChatWithStructTemplate(ctx, provider,
		"设计一个博客管理系统的RESTful API，包含文章的增删改查功能",
		APISpec{}, // 模板结构体
		&apiSpec)

	if err != nil {
		log.Printf("API设计失败: %v", err)
	} else {
		fmt.Printf("✓ API设计结果:\n")
		fmt.Printf("  标题: %s\n", apiSpec.Title)
		fmt.Printf("  版本: %s\n", apiSpec.Version)
		fmt.Printf("  端点数量: %d\n", len(apiSpec.Endpoints))
		for i, endpoint := range apiSpec.Endpoints {
			fmt.Printf("  %d. %s %s - %s\n",
				i+1, endpoint.Method, endpoint.Path, endpoint.Description)
		}
	}

	// 3. 代码生成
	fmt.Println("\n3. 代码生成:")
	var codeSpec CodeGenSpec
	err = ChatWithStructTemplate(ctx, provider,
		"生成一个Go语言的HTTP健康检查服务",
		CodeGenSpec{}, // 模板结构体
		&codeSpec)

	if err != nil {
		log.Printf("代码生成失败: %v", err)
	} else {
		fmt.Printf("✓ 代码生成结果:\n")
		fmt.Printf("  语言: %s\n", codeSpec.Language)
		fmt.Printf("  框架: %s\n", codeSpec.Framework)
		fmt.Printf("  文件数量: %d\n", len(codeSpec.Files))
		fmt.Printf("  依赖数量: %d\n", len(codeSpec.Dependencies))
		for i, file := range codeSpec.Files {
			fmt.Printf("  文件%d: %s - %s\n", i+1, file.Filename, file.Description)
		}
	}

	// 4. 数据库设计
	fmt.Println("\n4. 数据库设计:")
	var dbSchema DatabaseSchema
	err = ChatWithStructTemplate(ctx, provider,
		"设计一个在线商城的数据库结构，包含用户、商品、订单相关表",
		DatabaseSchema{}, // 模板结构体
		&dbSchema)

	if err != nil {
		log.Printf("数据库设计失败: %v", err)
	} else {
		fmt.Printf("✓ 数据库设计结果:\n")
		fmt.Printf("  数据库名: %s\n", dbSchema.Name)
		fmt.Printf("  描述: %s\n", dbSchema.Description)
		fmt.Printf("  表数量: %d\n", len(dbSchema.Tables))
		for i, table := range dbSchema.Tables {
			fmt.Printf("  表%d: %s - %d个字段\n", i+1, table.Name, len(table.Fields))
		}
	}
}

// ExampleSchemaGeneration 演示Schema生成
func ExampleSchemaGeneration() {
	fmt.Println("=== JSON Schema生成示例 ===")

	// 1. 用户信息Schema
	userSchema, err := GenerateJSONSchema(UserProfile{})
	if err != nil {
		log.Printf("生成用户Schema失败: %v", err)
	} else {
		fmt.Printf("\n1. 用户信息Schema:\n%s\n", userSchema)
	}

	// 2. API设计Schema
	apiSchema, err := GenerateJSONSchema(APISpec{})
	if err != nil {
		log.Printf("生成API Schema失败: %v", err)
	} else {
		fmt.Printf("\n2. API设计Schema长度: %d 字符\n", len(apiSchema))
	}

	// 3. 简单结构体Schema
	type SimpleStruct struct {
		ID       int    `json:"id" description:"唯一标识符"`
		Name     string `json:"name" description:"名称"`
		Optional string `json:"optional,omitempty" description:"可选字段"`
	}

	simpleSchema, err := GenerateJSONSchema(SimpleStruct{})
	if err != nil {
		log.Printf("生成简单Schema失败: %v", err)
	} else {
		fmt.Printf("\n3. 简单结构体Schema:\n%s\n", simpleSchema)
	}
}

// ExampleStructCompareApproaches 比较不同方法
func ExampleStructCompareApproaches() {
	ctx := context.Background()
	provider := ProviderDeepSeek

	fmt.Println("=== 结构体方法对比示例 ===")

	// 方法1: 手写JSON Schema字符串
	fmt.Println("\n方法1: 手写JSON Schema")
	manualSchema := `{
  "type": "object",
  "properties": {
    "name": {"type": "string"},
    "age": {"type": "integer"},
    "email": {"type": "string"}
  },
  "required": ["name", "age", "email"]
}`

	var result1 UserProfile
	err := ChatWithStructuredSchema(ctx, provider, "生成用户信息", manualSchema, &result1)
	if err != nil {
		log.Printf("手写Schema失败: %v", err)
	} else {
		fmt.Printf("✓ 手写Schema成功: %s\n", result1.Name)
	}

	// 方法2: 结构体自动生成Schema
	fmt.Println("\n方法2: 结构体自动生成Schema")
	var result2 UserProfile
	err = ChatWithStructTemplate(ctx, provider, "生成用户信息", UserProfile{}, &result2)
	if err != nil {
		log.Printf("结构体模板失败: %v", err)
	} else {
		fmt.Printf("✓ 结构体模板成功: %s\n", result2.Name)
	}

	// 方法3: 普通JSON模式
	fmt.Println("\n方法3: 普通JSON模式")
	var result3 UserProfile
	err = ChatWithJSONResult(ctx, provider, "生成用户信息，包含name、age、email字段", &result3)
	if err != nil {
		log.Printf("普通JSON失败: %v", err)
	} else {
		fmt.Printf("✓ 普通JSON成功: %s\n", result3.Name)
	}
}

// ExampleCustomTags 自定义标签示例
func ExampleCustomTags() {
	fmt.Println("=== 自定义标签使用示例 ===")

	// 带详细标签的结构体
	type ProductInfo struct {
		ID          int      `json:"id" description:"商品ID，唯一标识符" example:"12345"`
		Name        string   `json:"name" description:"商品名称，不能为空" example:"iPhone 15"`
		Price       float64  `json:"price" description:"商品价格，单位为元" example:"6999.99"`
		Category    string   `json:"category" description:"商品分类" example:"手机"`
		Tags        []string `json:"tags,omitempty" description:"商品标签列表"`
		InStock     bool     `json:"in_stock" description:"是否有库存"`
		Description string   `json:"description,omitempty" description:"商品详细描述"`
		Images      []string `json:"images,omitempty" description:"商品图片URL列表"`
	}

	// 生成Schema并查看
	schema, err := GenerateJSONSchema(ProductInfo{})
	if err != nil {
		log.Printf("生成商品Schema失败: %v", err)
		return
	}

	fmt.Printf("商品信息Schema:\n%s\n", schema)

	// 使用该结构体生成数据
	ctx := context.Background()
	provider := ProviderDeepSeek

	var product ProductInfo
	err = ChatWithStructTemplate(ctx, provider,
		"生成一个笔记本电脑的商品信息",
		ProductInfo{}, &product)

	if err != nil {
		log.Printf("生成商品信息失败: %v", err)
	} else {
		fmt.Printf("\n✓ 生成的商品信息:\n")
		fmt.Printf("  ID: %d\n", product.ID)
		fmt.Printf("  名称: %s\n", product.Name)
		fmt.Printf("  价格: %.2f\n", product.Price)
		fmt.Printf("  分类: %s\n", product.Category)
		fmt.Printf("  标签: %v\n", product.Tags)
		fmt.Printf("  库存: %t\n", product.InStock)
	}
}

// ExampleLLMTags 演示新的llm标签使用方法
func ExampleLLMTags() {
	fmt.Println("=== LLM标签使用示例 ===")

	// 1. 测试Bd结构体的Schema生成
	fmt.Println("\n1. Bd结构体Schema生成:")
	bdSchema, err := GenerateJSONSchema(Bd{})
	if err != nil {
		log.Printf("生成Bd Schema失败: %v", err)
	} else {
		fmt.Printf("Bd结构体Schema:\n%s\n", bdSchema)
	}

	// 2. 测试MixedExample结构体的Schema生成
	fmt.Println("\n2. 混合标签结构体Schema生成:")
	mixedSchema, err := GenerateJSONSchema(MixedExample{})
	if err != nil {
		log.Printf("生成Mixed Schema失败: %v", err)
	} else {
		fmt.Printf("Mixed结构体Schema:\n%s\n", mixedSchema)
	}

	// 3. 实际使用llm标签的结构体生成数据
	ctx := context.Background()
	provider := ProviderDeepSeek

	fmt.Println("\n3. 使用Bd结构体模板生成数据:")
	var bd Bd
	err = ChatWithStructTemplate(ctx, provider,
		"生成一个值班管理系统的数据",
		Bd{}, &bd)

	if err != nil {
		log.Printf("生成Bd数据失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的Bd数据:\n")
		fmt.Printf("  值班列表ID: %s\n", bd.ItemId)
		fmt.Printf("  服务组通知: %s\n", bd.GroupNotice)
		fmt.Printf("  创建时间: %s\n", bd.CreateTime)
		fmt.Printf("  ID字段(应该为空，因为llm:\"-\"): %s\n", bd.ID)
	}

	fmt.Println("\n3.1. 使用简化版本（推荐）:")
	var bd2 Bd
	err = ChatWithStruct(ctx, provider,
		"生成一个值班管理系统的数据", &bd2) // 只需要传递一个参数！

	if err != nil {
		log.Printf("生成Bd2数据失败: %v", err)
	} else {
		fmt.Printf("✓ 简化版本生成的Bd数据:\n")
		fmt.Printf("  值班列表ID: %s\n", bd2.ItemId)
		fmt.Printf("  服务组通知: %s\n", bd2.GroupNotice)
		fmt.Printf("  创建时间: %s\n", bd2.CreateTime)
		fmt.Printf("  ID字段(应该为空，因为llm:\"-\"): %s\n", bd2.ID)
	}

	fmt.Println("\n4. 使用Mixed结构体模板生成数据:")
	var mixed MixedExample
	err = ChatWithStructTemplate(ctx, provider,
		"生成一个用户的信息",
		MixedExample{}, &mixed)

	if err != nil {
		log.Printf("生成Mixed数据失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的Mixed数据:\n")
		fmt.Printf("  姓名: %s\n", mixed.Name)
		fmt.Printf("  年龄: %d\n", mixed.Age)
		fmt.Printf("  邮箱: %s\n", mixed.Email)
		fmt.Printf("  手机: %s\n", mixed.Phone)
		fmt.Printf("  个人描述: %s\n", mixed.Description)
		fmt.Printf("  Password字段(应该为空，因为llm:\"-\"): %s\n", mixed.Password)
		fmt.Printf("  InternalID字段(应该为0，因为llm:\"-\"): %d\n", mixed.InternalID)
	}

	fmt.Println(`
=== LLM标签语法说明 ===

支持的llm标签格式：
- llm:"desc:字段描述"     - 为字段添加描述信息
- llm:"-"               - 忽略此字段，不包含在生成的Schema中

标签优先级：
1. llm:"desc:xxx"       - 优先使用
2. description:"xxx"    - 如果没有llm标签描述，则使用description
3. example:"xxx"        - 示例值标签仍然支持

使用场景：
- ID字段通常使用 llm:"-" 忽略，因为不需要AI生成
- 密码等敏感字段使用 llm:"-" 忽略
- 业务字段使用 llm:"desc:xxx" 提供清晰描述

这种方式比传统的description标签更专业，明确表示是为LLM服务的标签。
`)
}

// ExampleAdvancedStructUsage 演示高级结构体使用方法
func ExampleAdvancedStructUsage() {
	fmt.Println("=== 高级结构体使用示例 ===")

	ctx := context.Background()
	provider := ProviderDeepSeek

	// 1. RAG场景：基于多个文档生成API设计
	fmt.Println("\n1. RAG场景示例:")

	// 模拟从数据库获取的文档
	ragDocs := []string{
		`API设计规范文档：
- 所有API都应该遵循RESTful规范
- 使用标准HTTP状态码
- 请求和响应都应该是JSON格式
- 需要支持分页、排序、过滤功能`,

		`用户管理业务需求：
- 用户注册、登录、注销
- 用户信息的增删改查
- 用户权限管理
- 用户状态管理（激活/禁用）`,

		`数据库设计标准：
- 每个表都需要id主键
- 创建时间和更新时间字段是必须的
- 软删除字段deleted_at
- 用户相关表需要tenant_id进行多租户隔离`,
	}

	var apiDesign APISpec
	err := ChatWithStructRAG(ctx, provider,
		"根据提供的文档，设计一个完整的用户管理API系统",
		ragDocs, &apiDesign)

	if err != nil {
		log.Printf("RAG生成失败: %v", err)
	} else {
		fmt.Printf("✓ 基于RAG生成的API设计:\n")
		fmt.Printf("  标题: %s\n", apiDesign.Title)
		fmt.Printf("  端点数量: %d\n", len(apiDesign.Endpoints))
		for i, endpoint := range apiDesign.Endpoints {
			fmt.Printf("  %d. %s %s\n", i+1, endpoint.Method, endpoint.Path)
		}
	}

	// 2. 多轮对话场景
	fmt.Println("\n2. 多轮对话场景示例:")

	messages := []Message{
		NewSystemMessage("你是专业的Go语言开发专家"),
		NewUserMessage("我需要设计一个博客系统"),
		NewAssistantMessage("好的，博客系统通常需要文章管理、用户管理、评论系统等功能。你具体需要哪些功能？"),
		NewUserMessage("我需要文章的CRUD操作，请生成相应的数据结构"),
	}

	type BlogPost struct {
		ID          int    `json:"id" llm:"-"`
		Title       string `json:"title" llm:"desc:文章标题"`
		Content     string `json:"content" llm:"desc:文章正文内容"`
		AuthorID    int    `json:"author_id" llm:"desc:作者ID"`
		CategoryID  int    `json:"category_id" llm:"desc:分类ID"`
		Status      string `json:"status" llm:"desc:文章状态(draft/published/archived)"`
		PublishedAt string `json:"published_at,omitempty" llm:"desc:发布时间"`
		CreatedAt   string `json:"created_at" llm:"-"`
		UpdatedAt   string `json:"updated_at" llm:"-"`
	}

	var blogPost BlogPost
	err = ChatWithStructMessages(ctx, provider, messages, &blogPost)

	if err != nil {
		log.Printf("多轮对话生成失败: %v", err)
	} else {
		fmt.Printf("✓ 多轮对话生成的文章结构:\n")
		fmt.Printf("  标题: %s\n", blogPost.Title)
		fmt.Printf("  作者ID: %d\n", blogPost.AuthorID)
		fmt.Printf("  状态: %s\n", blogPost.Status)
		fmt.Printf("  发布时间: %s\n", blogPost.PublishedAt)
	}

	// 3. 自定义系统提示词场景
	fmt.Println("\n3. 自定义系统提示词场景:")

	systemPrompt := `你是数据库设计专家，专门负责设计符合以下规范的数据结构：
1. 所有ID字段使用int类型
2. 时间字段使用string类型，格式为ISO8601
3. 状态字段使用枚举值
4. 必须包含audit字段（created_at, updated_at）
5. 支持软删除（deleted_at字段）`

	type DatabaseField struct {
		Name       string `json:"name" llm:"desc:字段名"`
		Type       string `json:"type" llm:"desc:字段类型"`
		Length     int    `json:"length,omitempty" llm:"desc:字段长度"`
		Nullable   bool   `json:"nullable" llm:"desc:是否允许为空"`
		DefaultVal string `json:"default_value,omitempty" llm:"desc:默认值"`
		Comment    string `json:"comment" llm:"desc:字段注释"`
		IsPrimary  bool   `json:"is_primary" llm:"desc:是否为主键"`
		IsIndex    bool   `json:"is_index" llm:"desc:是否需要索引"`
	}

	type DatabaseTable struct {
		TableName   string          `json:"table_name" llm:"desc:数据表名称"`
		Description string          `json:"description" llm:"desc:表的用途说明"`
		Fields      []DatabaseField `json:"fields" llm:"desc:字段列表"`
	}

	var dbTable DatabaseTable
	err = ChatWithStructContext(ctx, provider, systemPrompt,
		"设计一个电商系统的订单表结构", &dbTable)

	if err != nil {
		log.Printf("自定义系统提示词生成失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的数据库表设计:\n")
		fmt.Printf("  表名: %s\n", dbTable.TableName)
		fmt.Printf("  说明: %s\n", dbTable.Description)
		fmt.Printf("  字段数量: %d\n", len(dbTable.Fields))
		for i, field := range dbTable.Fields {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, field.Name, field.Type, field.Comment)
		}
	}
}

// ExampleAPIComparison 比较不同API的使用场景
func ExampleAPIComparison() {
	fmt.Println("=== API方法对比 ===")

	fmt.Println(`
API方法选择指南：

1. ChatWithStruct(ctx, provider, message, &result)
   🎯 适用场景：简单的单轮对话生成
   ✅ 优点：最简单，4个参数
   ❌ 限制：只能单轮对话，无法利用上下文
   📝 示例：生成用户信息、产品数据等

2. ChatWithStructMessages(ctx, provider, messages, &result)  
   🎯 适用场景：复杂的多轮对话、完全自定义
   ✅ 优点：最灵活，支持完整对话历史
   ❌ 限制：需要手动构建Message数组
   📝 示例：多轮交互、复杂业务逻辑

3. ChatWithStructRAG(ctx, provider, question, docs, &result)
   🎯 适用场景：基于文档的知识生成
   ✅ 优点：专门优化RAG场景，自动处理文档格式
   ❌ 限制：只适用于RAG场景
   📝 示例：基于数据库文档生成API、根据需求文档生成代码

4. ChatWithStructContext(ctx, provider, systemPrompt, message, &result)
   🎯 适用场景：需要自定义系统提示词的单轮对话
   ✅ 优点：方便设置专业角色和规范
   ❌ 限制：仍然是单轮对话
   📝 示例：专业领域生成、特定格式要求

推荐使用顺序：
- 简单场景 → ChatWithStruct
- RAG场景 → ChatWithStructRAG  
- 需要系统提示词 → ChatWithStructContext
- 复杂交互 → ChatWithStructMessages
`)
}

// ExampleFunctionGoRAG 演示function-go框架的RAG代码生成
func ExampleFunctionGoRAG() {
	fmt.Println("=== Function-Go框架RAG代码生成示例 ===")

	ctx := context.Background()
	provider := ProviderDeepSeek

	// 模拟从function-go项目中获取的现有代码作为参考
	functionGoRagDocs := []string{
		`// 现有API层示例代码
package runner

import (
	"github.com/gin-gonic/gin"
	"github.com/yunhanshu-net/pkg/response"
)

func (r *Runner) GetUserList(c *gin.Context) {
	var req base.PageInfoReq
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.ParamError(c, err.Error())
		return
	}
	
	user := c.GetString("user")
	list, total, err := r.userService.GetUserList(c, user, &req)
	if err != nil {
		response.Error(c, err.Error())
		return
	}
	
	response.PageSuccess(c, list, total)
}`,

		`// 现有Service层示例代码
package service

import (
	"context"
	"github.com/pkg/errors"
)

type UserService struct {
	userRepo *repo.UserRepo
}

func NewUserService(userRepo *repo.UserRepo) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserReq) (*dto.CreateUserResp, error) {
	user := &model.User{
		Name:  req.Name,
		Email: req.Email,
		User:  req.User,
	}
	
	err := s.userRepo.Create(ctx, user)
	if err != nil {
		return nil, errors.Wrap(err, "创建用户失败")
	}
	
	return &dto.CreateUserResp{ID: user.ID}, nil
}`,

		`// 现有Model层示例代码
package model

import "github.com/yunhanshu-net/pkg/model"

type User struct {
	model.Base
	Name     string ` + "`json:\"name\" gorm:\"column:name;comment:用户名\"`" + `
	Email    string ` + "`json:\"email\" gorm:\"column:email;comment:邮箱\"`" + `
	Phone    string ` + "`json:\"phone\" gorm:\"column:phone;comment:手机号\"`" + `
	Status   int    ` + "`json:\"status\" gorm:\"column:status;comment:状态 1正常 2禁用\"`" + `
	User     string ` + "`json:\"user\" gorm:\"column:user;comment:归属用户\"`" + `
}

func (User) TableName() string {
	return "users"
}`,

		`// 现有FunctionInfo配置示例
func GetUserListInfo() *runner.FunctionInfo {
	return &runner.FunctionInfo{
		Router:      "/api/v1/users",
		Method:      "GET", 
		ApiDesc:     "获取用户列表",
		ChineseName: "用户列表",
		EnglishName: "user_list",
		Classify:    "user",
		Tags:        []string{"user", "list"},
		RenderType:  "table",
		UseTables:   []interface{}{&model.User{}},
		OperateTables: map[interface{}][]runner.OperateTableType{
			&model.User{}: {runner.OperateTableTypeGet},
		},
		Request:  base.PageInfoReq{},
		Response: []model.User{},
	}
}`,

		`// 项目编码规范
规范要求：
1. API层只做参数解析和响应，业务逻辑放Service层
2. 所有数据库操作必须在Repo层
3. 错误处理使用github.com/pkg/errors包装
4. Model必须继承model.Base，包含gorm标签
5. 多租户隔离使用user字段
6. FunctionInfo必须配置完整的路由和操作表信息`,
	}

	// 用户需求
	userRequest := "帮我生成一个商品管理(Product)模块，包含商品的增删改查功能，商品字段包括：名称、价格、分类、库存、描述"

	var codeGen FunctionGoCodeGen
	err := ChatWithStructRAG(ctx, provider, userRequest, functionGoRagDocs, &codeGen)

	if err != nil {
		log.Printf("Function-Go代码生成失败: %v", err)
	} else {
		fmt.Printf("✓ 基于Function-Go框架生成的代码:\n")
		fmt.Printf("  模块名称: %s\n", codeGen.ModuleName)
		fmt.Printf("  功能描述: %s\n", codeGen.Description)
		fmt.Printf("  生成文件数量: %d\n", len(codeGen.Files))
		fmt.Printf("  依赖包数量: %d\n", len(codeGen.Dependencies))

		fmt.Println("\n生成的文件列表:")
		for i, file := range codeGen.Files {
			fmt.Printf("  %d. %s (%s) - %s\n", i+1, file.Path, file.Type, file.Description)
		}

		if len(codeGen.Files) > 0 {
			fmt.Printf("\n首个文件内容预览:\n%s\n",
				truncateString(codeGen.Files[0].Content, 300))
		}
	}

	fmt.Println(`
🎯 这种方式的优势：
1. ✅ AI完全理解你的项目架构和编码规范  
2. ✅ 生成的代码风格与现有代码完全一致
3. ✅ 自动遵循三层架构模式
4. ✅ 正确配置FunctionInfo和路由
5. ✅ 包含完整的错误处理和多租户支持
6. ✅ 生成的代码可以直接使用，无需大幅修改

这比让AI从零开始猜测你的架构要准确得多！
`)
}

// truncateString 截断字符串用于预览
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
