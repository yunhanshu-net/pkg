package httpx

import (
	"fmt"
	"log"
)

// ExampleUsage 展示httpx包的使用方法
func ExampleUsage() {
	// 1. 基本GET请求
	fmt.Println("=== 基本GET请求 ===")
	req := Get("https://httpbin.org/get")
	resp, err := req.DoString()
	if err != nil {
		log.Printf("GET请求失败: %v", err)
	} else {
		fmt.Printf("状态码: %d, 耗时: %v\n", resp.Code, resp.Cost)
	}

	// 2. POST请求带JSON数据
	fmt.Println("\n=== POST请求带JSON数据 ===")
	postReq := Post("https://httpbin.org/post").
		Header("Content-Type", "application/json").
		Body(map[string]interface{}{
			"name": "张三",
			"age":  25,
		})

	// 先进行DryRun预览
	dryRun := postReq.DryRun()
	fmt.Printf("DryRun预览:\n")
	fmt.Printf("  方法: %s\n", dryRun.Map()["method"])
	fmt.Printf("  URL: %s\n", dryRun.Map()["url"])
	fmt.Printf("  请求体: %s\n", dryRun.Map()["body"])

	// 实际执行请求
	var result map[string]interface{}
	resp, err = postReq.Do(&result)
	if err != nil {
		log.Printf("POST请求失败: %v", err)
	} else {
		fmt.Printf("请求成功，状态码: %d\n", resp.Code)
	}

	// 3. DELETE请求
	fmt.Println("\n=== DELETE请求 ===")
	deleteReq := Delete("https://httpbin.org/delete").
		Header("Authorization", "Bearer token123")

	dryRun = deleteReq.DryRun()
	fmt.Printf("DryRun预览:\n")
	fmt.Printf("  方法: %s\n", dryRun.Map()["method"])
	fmt.Printf("  URL: %s\n", dryRun.Map()["url"])
	fmt.Printf("  认证头: %s\n", dryRun.Map()["headers"].(map[string]string)["Authorization"])

	// 4. 带超时的请求
	fmt.Println("\n=== 带超时的请求 ===")
	timeoutReq := Get("https://httpbin.org/delay/1").
		Timeout(500) // 500ms超时，应该会失败

	resp, err = timeoutReq.DoString()
	if err != nil {
		fmt.Printf("请求超时: %v\n", err)
	} else {
		fmt.Printf("请求成功: %d\n", resp.Code)
	}
}

// ExampleResponseBinding 展示响应绑定
func ExampleResponseBinding() {
	fmt.Println("=== 响应绑定示例 ===")

	// 定义响应结构体
	type UserResponse struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	// 发送请求并绑定响应
	req := Post("https://httpbin.org/post").
		Body(map[string]interface{}{
			"id":   1,
			"name": "李四",
			"age":  30,
		})

	var userResp UserResponse
	resp, err := req.Do(&userResp)
	if err != nil {
		log.Printf("请求失败: %v", err)
		return
	}

	fmt.Printf("请求成功，状态码: %d\n", resp.Code)
	fmt.Printf("响应体: %+v\n", userResp)
}
