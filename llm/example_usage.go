package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

// ExampleJSONUsage 演示JSON格式输出的使用方法
func ExampleJSONUsage() {
	ctx := context.Background()

	// 注意：实际使用时需要先注册DeepSeek工厂
	// llm.DefaultManager().RegisterFactory(llm.ProviderDeepSeek, &deepseek.Factory{})

	// 1. 创建JSON模式配置
	config := GetJSONConfig(ProviderDeepSeek)
	config.APIKey = "your-deepseek-api-key" // 请替换为实际的API密钥

	fmt.Println("=== JSON格式输出示例 ===")

	// 2. 创建客户端
	_, err := CreateClient(config)
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		return
	}

	// 3. 基础JSON聊天
	fmt.Println("\n1. 基础JSON聊天:")
	jsonResult, err := QuickJSONChat(ctx, ProviderDeepSeek,
		"生成一个用户信息JSON，包含name(string)、age(int)、email(string)字段")
	if err != nil {
		fmt.Printf("JSON聊天失败: %v\n", err)
	} else {
		fmt.Printf("结果: %s\n", jsonResult)

		// 验证是否为有效JSON
		var temp interface{}
		if json.Unmarshal([]byte(jsonResult), &temp) == nil {
			fmt.Println("✓ JSON格式有效")
		} else {
			fmt.Println("✗ JSON格式无效")
		}
	}

	// 4. 代码生成（JSON格式）
	fmt.Println("\n2. 代码生成（JSON格式）:")
	codeResult, err := GenerateCode(ctx, ProviderDeepSeek,
		"创建一个简单的HTTP健康检查接口")
	if err != nil {
		fmt.Printf("代码生成失败: %v\n", err)
	} else {
		// 解析代码生成结果
		type CodeGenResult struct {
			Code         string   `json:"code"`
			Filename     string   `json:"filename"`
			Description  string   `json:"description"`
			Dependencies []string `json:"dependencies"`
		}

		var result CodeGenResult
		if err := json.Unmarshal([]byte(codeResult), &result); err != nil {
			fmt.Printf("解析JSON失败: %v\n", err)
			fmt.Printf("原始结果: %s\n", codeResult)
		} else {
			fmt.Printf("✓ 解析成功:\n")
			fmt.Printf("  文件名: %s\n", result.Filename)
			fmt.Printf("  描述: %s\n", result.Description)
			fmt.Printf("  依赖: %v\n", result.Dependencies)
			fmt.Printf("  代码长度: %d 字符\n", len(result.Code))
		}
	}

	// 5. 结构化数据生成
	fmt.Println("\n3. 结构化数据生成:")
	schema := `{
  "type": "object",
  "properties": {
    "endpoints": {
      "type": "array",
      "items": {
        "type": "object",
        "properties": {
          "name": {"type": "string"},
          "method": {"type": "string"},
          "path": {"type": "string"},
          "description": {"type": "string"}
        },
        "required": ["name", "method", "path", "description"]
      }
    }
  },
  "required": ["endpoints"]
}`

	structuredResult, err := GenerateStructuredData(ctx, ProviderDeepSeek,
		"设计一个简单的用户管理API，包含用户注册、登录、获取信息功能", schema)
	if err != nil {
		fmt.Printf("结构化数据生成失败: %v\n", err)
	} else {
		type APIDesign struct {
			Endpoints []struct {
				Name        string `json:"name"`
				Method      string `json:"method"`
				Path        string `json:"path"`
				Description string `json:"description"`
			} `json:"endpoints"`
		}

		var apiResult APIDesign
		if err := json.Unmarshal([]byte(structuredResult), &apiResult); err != nil {
			fmt.Printf("解析JSON失败: %v\n", err)
			fmt.Printf("原始结果: %s\n", structuredResult)
		} else {
			fmt.Printf("✓ 解析成功，生成了 %d 个API端点:\n", len(apiResult.Endpoints))
			for i, endpoint := range apiResult.Endpoints {
				fmt.Printf("  %d. %s %s - %s\n", i+1, endpoint.Method, endpoint.Path, endpoint.Description)
			}
		}
	}
}

// SimpleJSONExample 最简单的JSON使用示例
func SimpleJSONExample() {
	fmt.Println("=== 最简单的JSON使用示例 ===")

	// 1. 手动创建JSON请求
	req := NewJSONRequest("deepseek-coder",
		NewUserMessage("生成一个包含title和content字段的博客文章JSON"))

	fmt.Printf("创建的JSON请求配置:\n")
	fmt.Printf("- 模型: %s\n", req.Model)
	fmt.Printf("- 响应格式: %s\n", req.ResponseFormat.Type)
	fmt.Printf("- 消息数量: %d\n", len(req.Messages))

	// 2. 代码生成请求
	codeReq := NewCodeGenRequest("deepseek-coder", "创建一个简单的计算器函数")

	fmt.Printf("\n代码生成请求配置:\n")
	fmt.Printf("- 模型: %s\n", codeReq.Model)
	fmt.Printf("- 响应格式: %s\n", codeReq.ResponseFormat.Type)
	fmt.Printf("- 温度: %.1f\n", codeReq.Temperature)
	fmt.Printf("- 系统消息: %s\n", codeReq.Messages[0].Content[:50]+"...")

	// 3. 结构化请求
	schema := `{"type": "object", "properties": {"result": {"type": "string"}}}`
	structReq := NewStructuredRequest("deepseek-coder", "生成一个问候语", schema)

	fmt.Printf("\n结构化请求配置:\n")
	fmt.Printf("- 模型: %s\n", structReq.Model)
	fmt.Printf("- 响应格式: %s\n", structReq.ResponseFormat.Type)
	fmt.Printf("- JSON Schema: %s\n", structReq.ResponseFormat.Schema)
}
