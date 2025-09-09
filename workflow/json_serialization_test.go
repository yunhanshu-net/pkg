package workflow

import (
	"encoding/json"
	"testing"
)

// 测试 JSON 序列化和反序列化
func TestJSONSerialization(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);

func main() {
    //desc: 开始用户注册流程
    fmt.Println("开始用户注册流程...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    
    //desc: 检查用户创建是否成功
    if step1Err != nil {
        //desc: 用户创建失败，记录错误并退出
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    
    //desc: 用户创建成功，记录成功日志
    fmt.Printf("用户创建成功，工号: %s\n", 工号)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 测试 JSON 序列化
	jsonData, err := json.Marshal(result)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	// 验证 JSON 格式
	var jsonResult map[string]interface{}
	err = json.Unmarshal(jsonData, &jsonResult)
	if err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	// 验证关键字段
	if jsonResult["success"] != true {
		t.Error("success 字段不正确")
	}

	if jsonResult["error"] != "" {
		t.Error("error 字段应该为空")
	}

	// 验证输入变量
	inputVars, ok := jsonResult["input_vars"].(map[string]interface{})
	if !ok {
		t.Error("input_vars 字段类型不正确")
	}

	if inputVars["用户名"] != "张三" {
		t.Error("input_vars 中的用户名不正确")
	}

	// 验证步骤
	steps, ok := jsonResult["steps"].([]interface{})
	if !ok {
		t.Error("steps 字段类型不正确")
	}

	if len(steps) != 1 {
		t.Errorf("steps 长度不正确，期望 1，实际 %d", len(steps))
	}

	// 验证主函数
	mainFunc, ok := jsonResult["main_func"].(map[string]interface{})
	if !ok {
		t.Error("main_func 字段类型不正确")
	}

	statements, ok := mainFunc["statements"].([]interface{})
	if !ok {
		t.Error("main_func.statements 字段类型不正确")
	}

	if len(statements) < 3 {
		t.Errorf("statements 长度不正确，期望至少 3，实际 %d", len(statements))
	}

	// 验证变量映射
	variables, ok := jsonResult["variables"].(map[string]interface{})
	if !ok {
		t.Error("variables 字段类型不正确")
	}

	if len(variables) < 2 {
		t.Errorf("variables 长度不正确，期望至少 2，实际 %d", len(variables))
	}

	t.Logf("JSON 序列化成功，大小: %d 字节", len(jsonData))
}

// 测试 JSON 序列化的完整性
func TestJSONCompleteness(t *testing.T) {
	code := `var input = map[string]interface{}{
    "项目名称": "test-project",
    "版本号": "v1.0.0",
}

step1 = beiluo.test1.devops.git_push[用例001] -> (err 是否失败);
step2 = beiluo.test1.devops.deploy_test(string 项目名称, string 版本号) -> (string 部署地址, err 是否失败);

func main() {
    step1Err := step1()
    if step1Err != nil {
        step1.Printf("推送失败: %v", step1Err)
        return
    }
    
    部署地址, step2Err := step2(input["项目名称"], input["版本号"]){retry:2, timeout:3000}
    if step2Err != nil {
        step2.Printf("部署失败: %v", step2Err)
        return
    }
    
    fmt.Printf("部署成功，地址: %s\n", 部署地址)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 序列化为 JSON
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	// 验证 JSON 包含所有必要字段
	jsonStr := string(jsonData)

	requiredFields := []string{
		`"success": true`,
		`"input_vars"`,
		`"steps"`,
		`"main_func"`,
		`"variables"`,
		`"error": ""`,
		`"statements"`,
		`"type"`,
		`"content"`,
		`"line_number"`,
		`"function"`,
		`"args"`,
		`"returns"`,
		`"metadata"`,
	}

	for _, field := range requiredFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON 缺少必要字段: %s", field)
		}
	}

	t.Logf("JSON 序列化完整性验证通过")
	t.Logf("JSON 大小: %d 字节", len(jsonData))
}

// 测试 JSON 反序列化
func TestJSONDeserialization(t *testing.T) {
	// 创建一个解析结果
	code := `var input = map[string]interface{}{
    "用户名": "测试用户",
}

step1 = beiluo.test1.test.test_func(string 用户名) -> (string 结果, err 是否失败);

func main() {
    结果, step1Err := step1(input["用户名"])
    if step1Err != nil {
        return
    }
    fmt.Printf("结果: %s\n", 结果)
}`

	parser := NewSimpleParser()
	originalResult := parser.ParseWorkflow(code)

	if !originalResult.Success {
		t.Fatalf("解析失败: %s", originalResult.Error)
	}

	// 序列化为 JSON
	jsonData, err := json.Marshal(originalResult)
	if err != nil {
		t.Fatalf("JSON 序列化失败: %v", err)
	}

	// 反序列化
	var newResult SimpleParseResult
	err = json.Unmarshal(jsonData, &newResult)
	if err != nil {
		t.Fatalf("JSON 反序列化失败: %v", err)
	}

	// 验证反序列化结果
	if newResult.Success != originalResult.Success {
		t.Error("Success 字段不匹配")
	}

	if newResult.Error != originalResult.Error {
		t.Error("Error 字段不匹配")
	}

	if len(newResult.Steps) != len(originalResult.Steps) {
		t.Error("Steps 长度不匹配")
	}

	if len(newResult.MainFunc.Statements) != len(originalResult.MainFunc.Statements) {
		t.Error("MainFunc.Statements 长度不匹配")
	}

	if len(newResult.Variables) != len(originalResult.Variables) {
		t.Error("Variables 长度不匹配")
	}

	// 验证步骤信息
	if len(newResult.Steps) > 0 {
		if newResult.Steps[0].Name != originalResult.Steps[0].Name {
			t.Error("Step Name 不匹配")
		}
		if newResult.Steps[0].Function != originalResult.Steps[0].Function {
			t.Error("Step Function 不匹配")
		}
	}

	t.Logf("JSON 反序列化验证通过")
}

// 辅助函数：检查字符串是否包含子字符串
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
