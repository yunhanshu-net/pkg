package httpx

import (
	"fmt"
	"testing"
)

func TestSSLVerification(t *testing.T) {
	// 测试有效的SSL证书
	req := Get("https://httpbin.org/get").ConnectivityCheck()
	connectivity := req.GetConnectivity()

	if connectivity == nil {
		t.Error("ConnectivityCheck should return a result")
		return
	}

	fmt.Printf("SSL验证测试 - httpbin.org:\n")
	fmt.Printf("  可达性: %t\n", connectivity.Reachable)
	fmt.Printf("  SSL有效: %t\n", connectivity.SSLValid)
	fmt.Printf("  状态码: %d\n", connectivity.StatusCode)

	// httpbin.org应该有有效的SSL证书
	if !connectivity.SSLValid {
		t.Error("httpbin.org should have valid SSL certificate")
	}
}

func TestHTTPWithoutSSL(t *testing.T) {
	// 测试HTTP请求（不应该有SSL验证）
	req := Get("http://httpbin.org/get").ConnectivityCheck()
	connectivity := req.GetConnectivity()

	if connectivity == nil {
		t.Error("ConnectivityCheck should return a result")
		return
	}

	fmt.Printf("HTTP测试 - httpbin.org:\n")
	fmt.Printf("  可达性: %t\n", connectivity.Reachable)
	fmt.Printf("  SSL有效: %t\n", connectivity.SSLValid)
	fmt.Printf("  状态码: %d\n", connectivity.StatusCode)

	// HTTP请求的SSLValid应该为false
	if connectivity.SSLValid {
		t.Error("HTTP request should not have SSL validation")
	}
}

func TestInvalidSSL(t *testing.T) {
	// 测试无效的SSL证书（使用自签名证书的网站）
	req := Get("https://self-signed.badssl.com/").ConnectivityCheck()
	connectivity := req.GetConnectivity()

	if connectivity == nil {
		t.Error("ConnectivityCheck should return a result")
		return
	}

	fmt.Printf("无效SSL测试 - self-signed.badssl.com:\n")
	fmt.Printf("  可达性: %t\n", connectivity.Reachable)
	fmt.Printf("  SSL有效: %t\n", connectivity.SSLValid)
	fmt.Printf("  错误: %s\n", connectivity.Error)

	// 自签名证书应该被检测为无效
	if connectivity.SSLValid {
		t.Error("Self-signed certificate should be detected as invalid")
	}
}
