package llm

import (
	"context"
	"fmt"
	"log"
)

// 用户信息结构体
type UserInfo struct {
	Name  string `json:"name"`
	Age   int    `json:"age"`
	Email string `json:"email"`
}

// API设计结构体
type APIEndpoint struct {
	Name        string `json:"name"`
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type APIDesign struct {
	Title     string        `json:"title"`
	Endpoints []APIEndpoint `json:"endpoints"`
}

// 代码生成结果结构体
type CodeResult struct {
	Code         string   `json:"code"`
	Filename     string   `json:"filename"`
	Description  string   `json:"description"`
	Dependencies []string `json:"dependencies"`
}

// ExampleUniversalUsage 演示通用调用方法的使用
func ExampleUniversalUsage() {
	ctx := context.Background()
	provider := ProviderDeepSeek

	fmt.Println("=== 通用LLM调用方法演示 ===")

	// 注意：实际使用时需要先注册工厂和设置API密钥
	// llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})
	// config := llm.GetDefaultConfig(llm.ProviderDeepSeek)
	// config.APIKey = "your-api-key"
	// llm.CreateClient(config)

	// 1. 字符串结果（普通聊天）
	fmt.Println("\n1. 普通文本聊天:")
	var textResult string
	err := ChatWithStringResult(ctx, provider, "介绍一下Go语言的特点", &textResult)
	if err != nil {
		log.Printf("普通聊天失败: %v", err)
	} else {
		fmt.Printf("回答长度: %d 字符\n", len(textResult))
		fmt.Printf("回答前100字符: %s...\n", textResult[:min(100, len(textResult))])
	}

	// 2. JSON结果（自动反序列化到结构体）
	fmt.Println("\n2. JSON格式用户信息生成:")
	var userInfo UserInfo
	err = ChatWithJSONResult(ctx, provider,
		"生成一个示例用户信息，包含name、age、email字段", &userInfo)
	if err != nil {
		log.Printf("JSON聊天失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的用户信息:\n")
		fmt.Printf("  姓名: %s\n", userInfo.Name)
		fmt.Printf("  年龄: %d\n", userInfo.Age)
		fmt.Printf("  邮箱: %s\n", userInfo.Email)
	}

	// 3. 智能选择模式（根据目标类型自动选择）
	fmt.Println("\n3. 智能模式API设计:")
	var apiDesign APIDesign
	err = ChatWithResult(ctx, provider,
		"设计一个用户管理系统的RESTful API，包含用户注册、登录、获取用户信息、更新用户信息功能",
		&apiDesign)
	if err != nil {
		log.Printf("API设计失败: %v", err)
	} else {
		fmt.Printf("✓ 设计了 '%s' 包含 %d 个端点:\n", apiDesign.Title, len(apiDesign.Endpoints))
		for i, endpoint := range apiDesign.Endpoints {
			fmt.Printf("  %d. %s %s - %s\n",
				i+1, endpoint.Method, endpoint.Path, endpoint.Description)
		}
	}

	// 4. 自定义系统提示词
	fmt.Println("\n4. 自定义提示词代码生成:")
	var codeResult CodeResult
	systemPrompt := `你是一个专业的Go语言开发专家。请根据用户需求生成高质量的Go代码。

返回JSON格式，包含以下字段：
- code: 完整的Go代码
- filename: 建议的文件名  
- description: 功能描述
- dependencies: 需要的依赖包列表`

	err = ChatWithCustomPrompt(ctx, provider, systemPrompt,
		"创建一个HTTP中间件用于JWT认证", &codeResult, true)
	if err != nil {
		log.Printf("代码生成失败: %v", err)
	} else {
		fmt.Printf("✓ 生成的代码信息:\n")
		fmt.Printf("  文件名: %s\n", codeResult.Filename)
		fmt.Printf("  描述: %s\n", codeResult.Description)
		fmt.Printf("  依赖数量: %d\n", len(codeResult.Dependencies))
		fmt.Printf("  代码行数: 约%d行\n", estimateLines(codeResult.Code))
	}

	// 5. 同样的请求，不同的处理方式
	fmt.Println("\n5. 同样请求的不同处理:")

	// 5a. 作为字符串处理
	var stringVersion string
	err = ChatWithResult(ctx, provider, "创建一个Hello World函数", &stringVersion)
	if err != nil {
		log.Printf("字符串版本失败: %v", err)
	} else {
		fmt.Printf("字符串版本长度: %d\n", len(stringVersion))
	}

	// 5b. 作为结构化JSON处理
	var structVersion CodeResult
	err = ChatWithResult(ctx, provider, "创建一个Hello World函数", &structVersion)
	if err != nil {
		log.Printf("结构化版本失败: %v", err)
	} else {
		fmt.Printf("结构化版本文件名: %s\n", structVersion.Filename)
	}
}

// ExampleCompareApproaches 比较不同调用方式的例子
func ExampleCompareApproaches() {
	ctx := context.Background()
	provider := ProviderDeepSeek
	userQuery := "设计一个简单的博客系统数据库表结构"

	fmt.Println("=== 不同调用方式比较 ===")

	// 方式1: 原始的QuickChat（返回字符串）
	fmt.Println("\n方式1: 传统QuickChat")
	result1, err := QuickChat(ctx, provider, userQuery)
	if err != nil {
		log.Printf("QuickChat失败: %v", err)
	} else {
		fmt.Printf("结果类型: 字符串, 长度: %d\n", len(result1))
	}

	// 方式2: 新的通用方法（字符串）
	fmt.Println("\n方式2: 通用方法(字符串)")
	var result2 string
	err = ChatWithResult(ctx, provider, userQuery, &result2)
	if err != nil {
		log.Printf("通用方法失败: %v", err)
	} else {
		fmt.Printf("结果类型: 字符串, 长度: %d\n", len(result2))
	}

	// 方式3: 新的通用方法（结构体）
	fmt.Println("\n方式3: 通用方法(结构体)")
	type DatabaseDesign struct {
		Tables []struct {
			Name        string   `json:"name"`
			Description string   `json:"description"`
			Fields      []string `json:"fields"`
		} `json:"tables"`
	}

	var result3 DatabaseDesign
	err = ChatWithResult(ctx, provider, userQuery, &result3)
	if err != nil {
		log.Printf("结构化方法失败: %v", err)
	} else {
		fmt.Printf("结果类型: 结构体, 表数量: %d\n", len(result3.Tables))
	}

	// 方式4: 强制JSON模式
	fmt.Println("\n方式4: 强制JSON模式")
	var result4 DatabaseDesign
	err = ChatWithJSONResult(ctx, provider, userQuery, &result4)
	if err != nil {
		log.Printf("强制JSON失败: %v", err)
	} else {
		fmt.Printf("结果类型: 强制JSON结构体, 表数量: %d\n", len(result4.Tables))
	}
}

// ExampleErrorHandling 错误处理示例
func ExampleErrorHandling() {
	ctx := context.Background()
	provider := ProviderDeepSeek

	fmt.Println("=== 错误处理示例 ===")

	// 1. 无效的结构体（测试JSON解析错误）
	fmt.Println("\n1. 测试JSON解析错误处理:")
	type InvalidStruct struct {
		RequiredField int `json:"required_field"`
	}

	var invalid InvalidStruct
	err := ChatWithJSONResult(ctx, provider, "返回一个简单的问候语", &invalid)
	if err != nil {
		fmt.Printf("✓ 正确捕获错误: %v\n", err)
	}

	// 2. 非指针参数（测试参数验证）
	fmt.Println("\n2. 测试参数验证:")
	var notPointer string
	err = ChatWithResult(ctx, provider, "hello", notPointer) // 故意传非指针
	if err != nil {
		fmt.Printf("✓ 正确捕获参数错误: %v\n", err)
	}
}

// 工具函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func estimateLines(code string) int {
	lines := 1
	for _, char := range code {
		if char == '\n' {
			lines++
		}
	}
	return lines
}

// ExampleQuickStart 快速开始示例
func ExampleQuickStart() {
	fmt.Println("=== 快速开始示例 ===")
	fmt.Println(`
// 1. 注册DeepSeek工厂
llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})

// 2. 创建配置
config := llm.GetDefaultConfig(llm.ProviderDeepSeek)
config.APIKey = "your-deepseek-api-key"
llm.CreateClient(config)

// 3. 使用方式示例:

// 方式A: 字符串结果
var answer string
err := llm.ChatWithStringResult(ctx, llm.ProviderDeepSeek, "解释Go语言的特点", &answer)

// 方式B: JSON结构体结果  
type UserInfo struct {
    Name string ` + "`json:\"name\"`" + `
    Age  int    ` + "`json:\"age\"`" + `
}
var user UserInfo
err = llm.ChatWithJSONResult(ctx, llm.ProviderDeepSeek, "生成用户信息", &user)

// 方式C: 智能选择（根据目标类型自动选择JSON或文本模式）
err = llm.ChatWithResult(ctx, llm.ProviderDeepSeek, "生成用户信息", &user) // 自动使用JSON
err = llm.ChatWithResult(ctx, llm.ProviderDeepSeek, "介绍Go语言", &answer) // 自动使用文本

// 方式D: 自定义系统提示词
err = llm.ChatWithCustomPrompt(ctx, llm.ProviderDeepSeek, 
    "你是代码生成专家", "创建HTTP服务器", &codeResult, true)
`)
}
