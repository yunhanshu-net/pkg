package llm

import (
	"context"
	"fmt"
)

// ExampleKnowledgeUsage 演示知识库客户端的使用
func ExampleKnowledgeUsage() {
	ctx := context.Background()

	// 示例1：获取所有知识
	fmt.Println("=== 获取所有知识 ===")
	allKnowledge, err := GetKnowledge(ctx, AllKnowledge)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("获取到 %d 条知识\n", len(allKnowledge))

	// 示例2：获取API相关知识
	fmt.Println("\n=== 获取API知识 ===")
	apiKnowledge, err := GetKnowledge(ctx, APIKnowledge)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("获取到 %d 条API知识\n", len(apiKnowledge))

	// 示例3：自定义查询
	fmt.Println("\n=== 自定义查询 ===")
	customReq := &KnowledgeRequest{
		Category: "API",
		Keyword:  "table",
		Limit:    10,
		Role:     "user",
		SortBy:   "sort_order",
	}
	customKnowledge, err := GetKnowledge(ctx, customReq)
	if err != nil {
		fmt.Printf("错误: %v\n", err)
		return
	}
	fmt.Printf("获取到 %d 条自定义知识\n", len(customKnowledge))
}

// ExampleRAGWithKnowledge 演示如何结合知识库使用RAG
func ExampleRAGWithKnowledge() {
	ctx := context.Background()
	userRequest := "帮我创建一个用户管理的table接口"

	// 1. 获取相关知识作为RAG文档
	ragDocs, err := GetKnowledgeForRAG(ctx, userRequest, APIKnowledge)
	if err != nil {
		fmt.Printf("获取知识失败: %v\n", err)
		return
	}

	fmt.Printf("RAG文档数量: %d\n", len(ragDocs))
	fmt.Printf("第一个文档: %s\n", ragDocs[0]) // 应该是用户需求

	// 2. 使用RAG生成代码
	var result FunctionGoCodeGen
	err = ChatWithStructRAG(ctx, ProviderDeepSeek, userRequest, ragDocs, &result)
	if err != nil {
		fmt.Printf("生成代码失败: %v\n", err)
		return
	}

	fmt.Printf("生成的API代码: %s\n", result.ApiCode)
	fmt.Printf("生成的Service代码: %s\n", result.ServiceCode)
}

// GenerateCodeWithKnowledge 完整的代码生成流程
func GenerateCodeWithKnowledge(ctx context.Context, userRequest string, category string) (*FunctionGoCodeGen, error) {
	// 1. 构建知识查询请求
	knowledgeReq := &KnowledgeRequest{
		Category: category,
		Keyword:  "",
		Limit:    20,
		Role:     "all", // 包含系统知识和用户示例
		SortBy:   "sort_order",
	}

	// 2. 获取RAG文档
	ragDocs, err := GetKnowledgeForRAG(ctx, userRequest, knowledgeReq)
	if err != nil {
		return nil, fmt.Errorf("获取知识库内容失败: %w", err)
	}

	// 3. 使用RAG生成代码
	var result FunctionGoCodeGen
	err = ChatWithStructRAG(ctx, ProviderDeepSeek, userRequest, ragDocs, &result)
	if err != nil {
		return nil, fmt.Errorf("生成代码失败: %w", err)
	}

	return &result, nil
}

// QuickGenerateAPI 快速生成API代码（便捷函数）
func QuickGenerateAPI(ctx context.Context, userRequest string) (*FunctionGoCodeGen, error) {
	return GenerateCodeWithKnowledge(ctx, userRequest, "API")
}

// QuickGenerateModel 快速生成Model代码（便捷函数）
func QuickGenerateModel(ctx context.Context, userRequest string) (*FunctionGoCodeGen, error) {
	return GenerateCodeWithKnowledge(ctx, userRequest, "Model")
}
