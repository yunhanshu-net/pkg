package httpx

import (
	"testing"
	"time"
)

func TestHttpRequest_Basic(t *testing.T) {
	// 测试GET请求
	req := Get("https://httpbin.org/get")
	if req.method != "GET" {
		t.Errorf("期望method为GET，实际为%s", req.method)
	}

	// 测试链式调用
	req = req.Header("User-Agent", "test-agent").
		Timeout(30 * time.Second)

	if req.headers["User-Agent"] != "test-agent" {
		t.Errorf("期望User-Agent为test-agent，实际为%s", req.headers["User-Agent"])
	}

	if req.timeout != 30*time.Second {
		t.Errorf("期望timeout为30s，实际为%s", req.timeout)
	}
}

func TestHttpRequest_DryRun(t *testing.T) {
	req := Post("https://api.example.com/users").
		Header("Authorization", "Bearer token123").
		Body(map[string]interface{}{
			"name": "张三",
			"age":  25,
		})

	dryRun := req.DryRun()

	if dryRun.Map()["method"] != "POST" {
		t.Errorf("期望Method为POST，实际为%s", dryRun.Map()["method"])
	}

	if dryRun.Map()["url"] != "https://api.example.com/users" {
		t.Errorf("期望Url为https://api.example.com/users，实际为%s", dryRun.Map()["url"])
	}

	if dryRun.Map()["headers"].(map[string]string)["Authorization"] != "Bearer token123" {
		t.Errorf("期望Authorization为Bearer token123，实际为%s", dryRun.Map()["headers"].(map[string]string)["Authorization"])
	}

	// 验证Body包含JSON
	if dryRun.Map()["body"] == "" {
		t.Error("期望Body不为空")
	}
}

func TestRequestContext_String(t *testing.T) {
	ctx := &RequestContext{
		Method:        "GET",
		Url:           "https://example.com",
		Cost:          100 * time.Millisecond,
		Code:          200,
		ResBodyString: "{\"status\":\"ok\"}",
	}

	result := ctx.String("测试请求")
	if result == "" {
		t.Error("期望String方法返回非空字符串")
	}
}

func TestRequestContext_OK(t *testing.T) {
	// 测试成功情况
	ctx := &RequestContext{Code: 200}
	if !ctx.OK() {
		t.Error("期望OK()返回true")
	}

	// 测试失败情况
	ctx = &RequestContext{Code: 404}
	if ctx.OK() {
		t.Error("期望OK()返回false")
	}

	// 测试nil情况
	var nilCtx *RequestContext
	if nilCtx.OK() {
		t.Error("期望nil的OK()返回false")
	}
}
