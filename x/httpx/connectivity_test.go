package httpx

import (
	"fmt"
	"testing"
	"time"
)

func TestConnectivityCheck(t *testing.T) {
	// 测试连通性检查
	req := Get("https://httpbin.org/get").
		Header("User-Agent", "test-agent").
		ConnectivityCheck()

	connectivity := req.GetConnectivity()
	if connectivity == nil {
		t.Error("ConnectivityCheck should return a result")
		return
	}

	fmt.Printf("连通性检查结果:\n")
	fmt.Printf("  可达性: %t\n", connectivity.Reachable)
	fmt.Printf("  状态码: %d\n", connectivity.StatusCode)
	fmt.Printf("  响应时间: %v\n", connectivity.ResponseTime)
	fmt.Printf("  服务器: %s\n", connectivity.Server)
	fmt.Printf("  内容类型: %s\n", connectivity.ContentType)
	fmt.Printf("  SSL有效: %t\n", connectivity.SSLValid)

	if !connectivity.Reachable {
		t.Error("httpbin.org should be reachable")
	}

	if connectivity.StatusCode != 200 {
		t.Errorf("Expected status code 200, got %d", connectivity.StatusCode)
	}
}

func TestConnectivityCheckWithError(t *testing.T) {
	// 测试不可达的URL
	req := Get("https://nonexistent-domain-12345.com").
		ConnectivityCheck()

	connectivity := req.GetConnectivity()
	if connectivity == nil {
		t.Error("ConnectivityCheck should return a result even for errors")
		return
	}

	fmt.Printf("错误连通性检查结果:\n")
	fmt.Printf("  可达性: %t\n", connectivity.Reachable)
	fmt.Printf("  错误: %s\n", connectivity.Error)
	fmt.Printf("  DNS解析: %t\n", connectivity.DNSResolved)

	if connectivity.Reachable {
		t.Error("Non-existent domain should not be reachable")
	}

	if connectivity.Error == "" {
		t.Error("Should have an error message")
	}
}

func TestDryRunWithConnectivity(t *testing.T) {
	// 测试DryRun包含连通性检查
	dryRunCase := Delete("https://httpbin.org/delete").
		Header("Authorization", "Bearer test-token").
		Body(map[string]interface{}{
			"user_id": "123",
			"reason":  "测试删除",
		}).
		ConnectivityCheck().
		DryRun()

	fmt.Printf("DryRun案例:\n")
	fmt.Printf("  方法: %s\n", dryRunCase.Method)
	fmt.Printf("  URL: %s\n", dryRunCase.Url)
	fmt.Printf("  描述: %s\n", dryRunCase.Description)

	if dryRunCase.Connectivity == nil {
		t.Error("DryRun should include connectivity check result")
		return
	}

	fmt.Printf("  连通性检查:\n")
	fmt.Printf("    可达性: %t\n", dryRunCase.Connectivity.Reachable)
	fmt.Printf("    状态码: %d\n", dryRunCase.Connectivity.StatusCode)
	fmt.Printf("    响应时间: %v\n", dryRunCase.Connectivity.ResponseTime)
}

func TestChainCalls(t *testing.T) {
	// 测试链式调用
	req := Post("https://httpbin.org/post").
		Header("Content-Type", "application/json").
		Body(map[string]interface{}{
			"name": "张三",
			"age":  25,
		}).
		Timeout(10 * time.Second).
		ConnectivityCheck()

	// 验证链式调用后的状态
	if req.GetConnectivity() == nil {
		t.Error("Chain calls should preserve connectivity check result")
	}

	// 继续链式调用
	dryRunCase := req.DryRun()
	if dryRunCase.Connectivity == nil {
		t.Error("DryRun should include connectivity from chain calls")
	}
}
