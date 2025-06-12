package llm

import (
	"context"
	"fmt"
	"log"
)

// QuickStartExample 快速开始示例
func QuickStartExample() {
	fmt.Println(`
=== LLM包快速开始指南 ===

1. 首先注册DeepSeek工厂并设置配置：

	import (
		"github.com/yunhanshu-net/pkg/llm"
		"github.com/yunhanshu-net/pkg/llm/deepseek"
	)

	// 注册工厂
	llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})
	
	// 创建配置
	config := llm.GetDefaultConfig(llm.ProviderDeepSeek)
	config.APIKey = "your-deepseek-api-key"
	llm.CreateClient(config)

2. 定义结构体模板（推荐方式）：

	type UserInfo struct {
		Name     string   ` + "`json:\"name\" description:\"用户姓名\"`" + `
		Age      int      ` + "`json:\"age\" description:\"用户年龄\"`" + `
		Email    string   ` + "`json:\"email\" description:\"邮箱地址\"`" + `
		Hobbies  []string ` + "`json:\"hobbies,omitempty\" description:\"兴趣爱好\"`" + `
	}

3. 使用结构体模板生成数据（推荐方式）：

	var user UserInfo
	err := llm.ChatWithStruct(ctx, llm.ProviderDeepSeek,
		"生成一个程序员的用户信息", &user)  // 只需传递一个参数

	// 或者使用完整版本
	err = llm.ChatWithStructTemplate(ctx, llm.ProviderDeepSeek,
		"生成用户信息", UserInfo{}, &user)

4. 其他使用方式：

	// 智能选择模式
	err = llm.ChatWithResult(ctx, provider, "生成用户", &user)    // 自动JSON
	err = llm.ChatWithResult(ctx, provider, "介绍Go", &text)     // 自动文本
	
	// 强制JSON模式
	err = llm.ChatWithJSONResult(ctx, provider, "生成数据", &user)
	
	// 纯文本模式
	err = llm.ChatWithStringResult(ctx, provider, "解释概念", &text)
`)
}

// RunQuickDemo 运行快速演示（需要API密钥）
func RunQuickDemo(apiKey string) {
	if apiKey == "" {
		fmt.Println("请提供DeepSeek API密钥来运行演示")
		return
	}

	fmt.Println("=== 快速演示开始 ===")

	// 注意：在实际使用中需要注册工厂
	// llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})

	// 创建配置
	config := GetDefaultConfig(ProviderDeepSeek)
	config.APIKey = apiKey

	// 创建客户端
	_, err := CreateClient(config)
	if err != nil {
		log.Printf("创建客户端失败: %v", err)
		return
	}

	ctx := context.Background()

	// 演示1: 结构体模板方法（简化版本）
	fmt.Println("\n1. 结构体模板方法演示（简化版本）:")

	type SimpleUser struct {
		Name string `json:"name" llm:"desc:用户姓名"`
		Age  int    `json:"age" llm:"desc:用户年龄"`
		Job  string `json:"job" llm:"desc:职业"`
	}

	var user SimpleUser
	err = ChatWithStruct(ctx, ProviderDeepSeek,
		"生成一个25岁的软件工程师用户信息", &user) // 只传递一个参数！

	if err != nil {
		log.Printf("生成用户失败: %v", err)
	} else {
		fmt.Printf("✓ 生成用户: %s, %d岁, 职业: %s\n", user.Name, user.Age, user.Job)
	}

	// 演示2: 智能选择方法
	fmt.Println("\n2. 智能选择方法演示:")

	// 字符串结果
	var explanation string
	err = ChatWithResult(ctx, ProviderDeepSeek, "用一句话解释什么是RESTful API", &explanation)
	if err != nil {
		log.Printf("获取解释失败: %v", err)
	} else {
		fmt.Printf("✓ RESTful API解释: %s\n", explanation)
	}

	// 结构体结果
	var user2 SimpleUser
	err = ChatWithResult(ctx, ProviderDeepSeek, "生成一个医生的用户信息", &user2)
	if err != nil {
		log.Printf("生成医生用户失败: %v", err)
	} else {
		fmt.Printf("✓ 生成医生: %s, %d岁, 职业: %s\n", user2.Name, user2.Age, user2.Job)
	}

	fmt.Println("\n=== 演示结束 ===")
}

// ShowSchemaExample 展示Schema生成示例
func ShowSchemaExample() {
	fmt.Println("=== JSON Schema生成演示 ===")

	// 定义一个复杂的结构体
	type BlogPost struct {
		ID        int      `json:"id" description:"文章ID"`
		Title     string   `json:"title" description:"文章标题"`
		Content   string   `json:"content" description:"文章内容"`
		Author    string   `json:"author" description:"作者名称"`
		Tags      []string `json:"tags,omitempty" description:"文章标签"`
		Published bool     `json:"published" description:"是否已发布"`
		Views     int      `json:"views,omitempty" description:"阅读次数"`
	}

	// 生成Schema
	schema, err := GenerateJSONSchema(BlogPost{})
	if err != nil {
		log.Printf("生成Schema失败: %v", err)
		return
	}

	fmt.Printf("BlogPost结构体生成的JSON Schema:\n%s\n", schema)

	// 显示支持的标签
	fmt.Println(`
支持的结构体标签：
- json:"field_name"          - 指定JSON字段名
- json:"field_name,omitempty" - 字段可选（不是required）
- description:"字段描述"      - 字段说明，会加入Schema
- example:"示例值"           - 示例值，帮助AI理解

最佳实践：
1. 使用有意义的字段名和描述
2. 合理使用omitempty标记可选字段
3. 提供example帮助AI生成更准确的数据
4. 结构体嵌套会自动处理，支持数组和Map
`)
}
