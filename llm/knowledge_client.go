package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// KnowledgeClient 知识库客户端
type KnowledgeClient struct {
	BaseURL string
	Client  *http.Client
}

// NewKnowledgeClient 创建知识库客户端
func NewKnowledgeClient(baseURL string) *KnowledgeClient {
	return &KnowledgeClient{
		BaseURL: baseURL,
		Client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// KnowledgeRequest 知识库查询请求
type KnowledgeRequest struct {
	Category string `json:"category"`
	Keyword  string `json:"keyword"`
	Limit    int    `json:"limit"`
	Role     string `json:"role"`
	SortBy   string `json:"sort_by"`
}

// KnowledgeResponse 知识库响应
type KnowledgeResponse struct {
	Code int `json:"code"`
	Data struct {
		FormattedContent string `json:"formatted_content"`
		TotalCount       int    `json:"total_count"`
		Categories       string `json:"categories"`
	} `json:"data"`
	Message string `json:"message"`
}

// GetKnowledge 获取知识库内容
func (kc *KnowledgeClient) GetKnowledge(ctx context.Context, req *KnowledgeRequest) ([]string, error) {
	// 构建请求URL
	url := fmt.Sprintf("%s/function/run/beiluo/rag_knowledge/knowledge/get/",
		strings.TrimSuffix(kc.BaseURL, "/"))

	// 序列化请求体
	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 创建HTTP请求
	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := kc.Client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 解析响应
	var knowledgeResp KnowledgeResponse
	if err := json.Unmarshal(respBody, &knowledgeResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if knowledgeResp.Code != 0 {
		return nil, fmt.Errorf("接口返回错误: %s", knowledgeResp.Message)
	}

	// 如果没有内容，返回空数组
	if knowledgeResp.Data.FormattedContent == "" ||
		knowledgeResp.Data.FormattedContent == "未找到匹配的知识内容" {
		return []string{}, nil
	}

	// 按 </split> 分割知识条目
	knowledge := strings.Split(knowledgeResp.Data.FormattedContent, "</split>")

	// 清理空白字符
	var result []string
	for _, item := range knowledge {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}

	return result, nil
}

// GetKnowledgeForRAG 获取用于RAG的知识内容（包含用户消息）
func (kc *KnowledgeClient) GetKnowledgeForRAG(ctx context.Context, userMessage string, req *KnowledgeRequest) ([]string, error) {
	// 获取知识库内容
	knowledge, err := kc.GetKnowledge(ctx, req)
	if err != nil {
		return nil, err
	}

	// 构建RAG文档数组
	ragDocs := make([]string, 0, len(knowledge)+1)

	// 首先添加用户需求（用message标签包裹）
	if userMessage != "" {
		ragDocs = append(ragDocs, fmt.Sprintf("<message>%s</message>", userMessage))
	}

	// 添加知识库内容
	ragDocs = append(ragDocs, knowledge...)

	return ragDocs, nil
}

// DefaultKnowledgeClient 默认知识库客户端（使用本地地址）
var DefaultKnowledgeClient = NewKnowledgeClient("http://localhost:8080")

// GetKnowledge 便捷函数：使用默认客户端获取知识
func GetKnowledge(ctx context.Context, req *KnowledgeRequest) ([]string, error) {
	return DefaultKnowledgeClient.GetKnowledge(ctx, req)
}

// GetKnowledgeForRAG 便捷函数：使用默认客户端获取RAG知识
func GetKnowledgeForRAG(ctx context.Context, userMessage string, req *KnowledgeRequest) ([]string, error) {
	return DefaultKnowledgeClient.GetKnowledgeForRAG(ctx, userMessage, req)
}

// 预定义的常用查询
var (
	// AllKnowledge 获取所有知识
	AllKnowledge = &KnowledgeRequest{
		Category: "",
		Keyword:  "",
		Limit:    100,
		Role:     "all",
		SortBy:   "",
	}

	// APIKnowledge 获取API相关知识
	APIKnowledge = &KnowledgeRequest{
		Category: "API",
		Keyword:  "",
		Limit:    50,
		Role:     "user",
		SortBy:   "sort_order",
	}

	// SystemKnowledge 获取系统级知识
	SystemKnowledge = &KnowledgeRequest{
		Category: "",
		Keyword:  "",
		Limit:    20,
		Role:     "system",
		SortBy:   "sort_order",
	}
)
