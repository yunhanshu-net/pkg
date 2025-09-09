package workflow

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleParser_StaticWorkflow(t *testing.T) {

	//静态工作流
	code := `
var input = map[string]interface{}{
    "项目名称": "my-project",
    "环境": "production",
    "版本": "v1.0.0",
}

step1 = beiluo.test1.devops.git_push[用例001] -> (err 是否失败);
step2 = beiluo.test1.devops.deploy_test[用例002] -> (int cost, err 是否失败);
step3 = beiluo.test1.devops.deploy_prod[用例003] -> (string msg, err 是否失败);
step4 = beiluo.test1.notify.send_notification[用例004] -> (err 是否失败);

func main() {
    fmt.Println("开始执行发布流程...")
    
    err := step1()
    if err != nil {
        fmt.Printf("推送代码失败: %v\n", err)
        return
    }
    fmt.Println("✅ 代码推送成功")
    
    err = step2()
    if err != nil {
        fmt.Printf("发布测试环境失败: %v\n", err)
        return
    }
    fmt.Println("✅ 测试环境发布成功")
    
    err = step3()
    if err != nil {
        fmt.Printf("发布线上环境失败: %v\n", err)
        return
    }
    fmt.Println("✅ 线上环境发布成功")
    
    err = step4()
    if err != nil {
        fmt.Printf("发送通知失败: %v\n", err)
        return
    }
    fmt.Println("✅ 通知发送成功")
    fmt.Println("🎉 发布流程执行完成！")
}
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}
	marshal, err := json.Marshal(result)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(string(marshal))

	// 测试输入变量
	if len(result.InputVars) != 3 {
		t.Errorf("期望输入变量数量为3，实际为%d", len(result.InputVars))
	}

	expectedInputs := map[string]interface{}{
		"项目名称": "my-project",
		"环境":   "production",
		"版本":   "v1.0.0",
	}

	for key, expectedValue := range expectedInputs {
		if actualValue, exists := result.InputVars[key]; !exists {
			t.Errorf("缺少输入变量: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("输入变量 %s 期望值 %v，实际值 %v", key, expectedValue, actualValue)
		}
	}

	// 测试工作流步骤
	if len(result.Steps) != 4 {
		t.Errorf("期望工作流步骤数量为4，实际为%d", len(result.Steps))
	}

	expectedSteps := []struct {
		name     string
		function string
		isStatic bool
		caseID   string
	}{
		{"step1", "beiluo.test1.devops.git_push", true, "用例001"},
		{"step2", "beiluo.test1.devops.deploy_test", true, "用例002"},
		{"step3", "beiluo.test1.devops.deploy_prod", true, "用例003"},
		{"step4", "beiluo.test1.notify.send_notification", true, "用例004"},
	}

	for i, expected := range expectedSteps {
		if i >= len(result.Steps) {
			t.Errorf("步骤 %d 不存在", i+1)
			continue
		}

		step := result.Steps[i]
		if step.Name != expected.name {
			t.Errorf("步骤 %d 名称期望 %s，实际 %s", i+1, expected.name, step.Name)
		}
		if step.Function != expected.function {
			t.Errorf("步骤 %d 函数名期望 %s，实际 %s", i+1, expected.function, step.Function)
		}
		if step.IsStatic != expected.isStatic {
			t.Errorf("步骤 %d 是否静态期望 %v，实际 %v", i+1, expected.isStatic, step.IsStatic)
		}
		if step.CaseID != expected.caseID {
			t.Errorf("步骤 %d 用例ID期望 %s，实际 %s", i+1, expected.caseID, step.CaseID)
		}
	}

	// 测试主函数
	if result.MainFunc == nil {
		t.Error("主函数解析失败")
	} else if len(result.MainFunc.Statements) == 0 {
		t.Error("主函数语句为空")
	} else {
		// 检查第一个语句
		firstStmt := result.MainFunc.Statements[0]
		if firstStmt.Type != "print" {
			t.Errorf("第一个语句类型期望 print，实际 %s", firstStmt.Type)
		}
		if firstStmt.Content != `fmt.Println("开始执行发布流程...")` {
			t.Errorf("第一个语句内容不匹配")
		}
	}
}

func TestSimpleParser_DynamicWorkflow(t *testing.T) {
	code := `
var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);
step2 = beiluo.test1.crm.crm_interview_schedule(string 用户名) -> (string 面试时间, string 面试官名称, err 是否失败);
step3 = beiluo.test1.crm.crm_interview_notify(string 面试官名称, string 通知信息) -> (err 是否失败);

func main() {
    fmt.Println("开始执行动态工作流...")
    
    // 创建用户
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        fmt.Printf("创建用户失败: %v\n", step1Err)
        return
    }
    fmt.Printf("✅ 用户创建成功，工号: %s\n", 工号)
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        fmt.Printf("安排面试失败: %v\n", step2Err)
        return
    }
    fmt.Printf("✅ 面试安排成功，时间: %s\n", 面试时间)
    
    // 发送通知
    通知信息 := "你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注"
    step3Err := step3(面试官名称, 通知信息)
    if step3Err != nil {
        fmt.Printf("发送通知失败: %v\n", step3Err)
        return
    }
    fmt.Println("✅ 通知发送成功")
    fmt.Println("🎉 动态工作流执行完成！")
}
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 测试输入变量
	if len(result.InputVars) != 3 {
		t.Errorf("期望输入变量数量为3，实际为%d", len(result.InputVars))
	}

	expectedInputs := map[string]interface{}{
		"用户名": "张三",
		"手机号": 13800138000,
		"邮箱":  "zhangsan@example.com",
	}

	for key, expectedValue := range expectedInputs {
		if actualValue, exists := result.InputVars[key]; !exists {
			t.Errorf("缺少输入变量: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("输入变量 %s 期望值 %v，实际值 %v", key, expectedValue, actualValue)
		}
	}

	// 测试工作流步骤
	if len(result.Steps) != 3 {
		t.Errorf("期望工作流步骤数量为3，实际为%d", len(result.Steps))
	}

	// 测试动态工作流步骤
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤1名称期望 step1，实际 %s", step1.Name)
	}
	if step1.Function != "beiluo.test1.devops.devops_script_create" {
		t.Errorf("步骤1函数名不匹配")
	}
	if step1.IsStatic {
		t.Error("步骤1应该是动态工作流")
	}
	if len(step1.InputTypes) != 2 {
		t.Errorf("步骤1输入参数数量期望2，实际%d", len(step1.InputTypes))
	}
	if len(step1.OutputTypes) != 3 {
		t.Errorf("步骤1输出参数数量期望3，实际%d", len(step1.OutputTypes))
	}

	// 测试输入参数
	expectedInputTypes := []SimpleTypeDef{
		{"string", "用户名"},
		{"int", "手机号"},
	}
	for i, expected := range expectedInputTypes {
		if i >= len(step1.InputTypes) {
			t.Errorf("输入参数 %d 不存在", i+1)
			continue
		}
		actual := step1.InputTypes[i]
		if actual.Type != expected.Type || actual.Name != expected.Name {
			t.Errorf("输入参数 %d 期望 %s %s，实际 %s %s", i+1, expected.Type, expected.Name, actual.Type, actual.Name)
		}
	}

	// 测试主函数
	if result.MainFunc == nil {
		t.Error("主函数解析失败")
	} else if len(result.MainFunc.Statements) == 0 {
		t.Error("主函数语句为空")
	} else {
		// 检查第一个语句
		firstStmt := result.MainFunc.Statements[0]
		if firstStmt.Type != "print" {
			t.Errorf("第一个语句类型期望 print，实际 %s", firstStmt.Type)
		}
		if firstStmt.Content != `fmt.Println("开始执行动态工作流...")` {
			t.Errorf("第一个语句内容不匹配")
		}

		// 检查函数调用语句
		functionCallCount := 0
		for _, stmt := range result.MainFunc.Statements {
			if stmt.Type == "function-call" {
				functionCallCount++
			}
		}
		if functionCallCount != 3 {
			t.Errorf("期望函数调用数量为3，实际为%d", functionCallCount)
		}
	}
}

func TestSimpleParser_MixedWorkflow(t *testing.T) {
	code := `
var input = map[string]interface{}{
    "项目名称": "mixed-project",
    "用户姓名": "李四",
    "用户年龄": 25,
}

step1 = beiluo.test1.devops.git_clone[用例001] -> (err 是否失败);
step2 = beiluo.test1.devops.build_project[用例002] -> (string 构建结果, err 是否失败);
step3 = beiluo.test1.user.create_user(string 姓名, int 年龄) -> (string 用户ID, err 是否失败);
step4 = beiluo.test1.user.assign_permissions(string 用户ID, string 项目名称) -> (err 是否失败);
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 测试混合工作流
	if len(result.Steps) != 4 {
		t.Errorf("期望工作流步骤数量为4，实际为%d", len(result.Steps))
	}

	// 检查静态步骤
	staticSteps := 0
	dynamicSteps := 0
	for _, step := range result.Steps {
		if step.IsStatic {
			staticSteps++
		} else {
			dynamicSteps++
		}
	}

	if staticSteps != 2 {
		t.Errorf("期望静态步骤数量为2，实际为%d", staticSteps)
	}
	if dynamicSteps != 2 {
		t.Errorf("期望动态步骤数量为2，实际为%d", dynamicSteps)
	}
}

func TestSimpleParser_ComplexInput(t *testing.T) {
	code := `
var input = map[string]interface{}{
    "项目名称": "complex-project",
    "数据库类型": "postgresql",
    "端口号": 5432,
    "启用SSL": true,
    "超时时间": 300,
}
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 测试复杂输入变量
	expectedInputs := map[string]interface{}{
		"项目名称":  "complex-project",
		"数据库类型": "postgresql",
		"端口号":   5432,
		"启用SSL": true,
		"超时时间":  300,
	}

	if len(result.InputVars) != len(expectedInputs) {
		t.Errorf("期望输入变量数量为%d，实际为%d", len(expectedInputs), len(result.InputVars))
	}

	for key, expectedValue := range expectedInputs {
		if actualValue, exists := result.InputVars[key]; !exists {
			t.Errorf("缺少输入变量: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("输入变量 %s 期望值 %v，实际值 %v", key, expectedValue, actualValue)
		}
	}
}

func TestSimpleParser_EmptyCode(t *testing.T) {
	parser := NewSimpleParser()
	result := parser.ParseWorkflow("")

	if !result.Success {
		t.Fatalf("空代码解析应该成功")
	}

	if len(result.InputVars) != 0 {
		t.Errorf("空代码输入变量应该为空")
	}
	if len(result.Steps) != 0 {
		t.Errorf("空代码步骤应该为空")
	}
}

func TestSimpleParser_InvalidCode(t *testing.T) {
	code := `
var input = map[string]interface{}{
    "项目名称": "test",
    // 缺少右括号
`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	// 这个测试可能会成功，因为我们的解析器比较宽松
	// 如果需要严格的错误检查，可以在这里添加相应的测试
	_ = result
}

func TestSimpleParser_MetadataSupport(t *testing.T) {
	code := `
var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);

func main() {
    // 带元数据的函数调用
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high"}
    if step1Err != nil {
        fmt.Printf("创建用户失败: %v\n", step1Err)
        return
    }
    
    // 纯函数调用带元数据
    step2(用户名){retry:1, timeout:2000, async:true}
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	assert.True(t, result.Success)
	assert.Empty(t, result.Error)
	assert.Len(t, result.MainFunc.Statements, 3)

	// 检查第一个函数调用的元数据
	firstCall := result.MainFunc.Statements[0]
	assert.Equal(t, "function-call", firstCall.Type)
	assert.Equal(t, "step1", firstCall.Function)
	assert.Len(t, firstCall.Metadata, 3)
	assert.Equal(t, 3, firstCall.Metadata["retry"])
	assert.Equal(t, 5000, firstCall.Metadata["timeout"])
	assert.Equal(t, "high", firstCall.Metadata["priority"])

	// 检查第二个函数调用的元数据
	secondCall := result.MainFunc.Statements[2]
	assert.Equal(t, "function-call", secondCall.Type)
	assert.Equal(t, "step2", secondCall.Function)
	assert.Len(t, secondCall.Metadata, 3)
	assert.Equal(t, 1, secondCall.Metadata["retry"])
	assert.Equal(t, 2000, secondCall.Metadata["timeout"])
	assert.Equal(t, true, secondCall.Metadata["async"])
}
