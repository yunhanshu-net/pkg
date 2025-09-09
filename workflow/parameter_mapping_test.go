package workflow

import (
	"context"
	"testing"
)

func TestParameterMapping_CompleteWorkflow(t *testing.T) {
	// 1. 创建工作流代码 - 包含多个步骤和复杂的参数传递
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "部门": "技术部",
}

// 步骤1：创建用户
step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号"
) -> (
    workId: string "工号",
    username: string "用户名",
    err: error "是否失败"
);

// 步骤2：分配部门
step2 = beiluo.test1.devops.devops_script_assign(
    workId: string "工号",
    department: string "部门"
) -> (
    success: bool "是否成功",
    message: string "消息",
    err: error "是否失败"
);

// 步骤3：发送通知
step3 = beiluo.test1.devops.devops_script_notify(
    workId: string "工号",
    username: string "用户名",
    department: string "部门"
) -> (
    success: bool "是否成功",
    err: error "是否失败"
);

func main() {
    // 步骤1：创建用户
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    if step1Err != nil {
        fmt.Printf("用户创建失败: %v\n", step1Err)
        return
    }
    
    // 步骤2：分配部门
    分配成功, 消息, step2Err := step2(工号, input["部门"]){retry:2, timeout:3000}
    if step2Err != nil {
        fmt.Printf("部门分配失败: %v\n", step2Err)
        return
    }
    
    // 步骤3：发送通知
    通知成功, step3Err := step3(工号, 用户名, input["部门"]){retry:1, timeout:2000}
    if step3Err != nil {
        fmt.Printf("通知发送失败: %v\n", step3Err)
        return
    }
    
    fmt.Printf("用户 %s 创建成功，工号: %s，部门: %s\n", 用户名, 工号, input["部门"])
}`

	// 2. 解析工作流
	parser := NewSimpleParser()
	parseResult := parser.ParseWorkflow(code)
	if !parseResult.Success {
		t.Fatalf("解析失败: %s", parseResult.Error)
	}

	// 3. 设置FlowID
	parseResult.FlowID = "test-complete-workflow-001"

	// 4. 创建执行器
	executor := NewExecutor()

	// 5. 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s - %s", step.Name, in.StepDesc)
		t.Logf("输入参数: %+v", in.RealInput)

		switch step.Name {
		case "step1":
			// 验证step1的输入参数
			expectedInput := map[string]interface{}{
				"username": "张三",        // 形参名 -> 实际值
				"phone":    13800138000, // 形参名 -> 实际值
			}
			validateInput(t, in.RealInput, expectedInput, "step1")

			// 返回step1的输出
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":   "EMP001", // 形参名
					"username": "张三",     // 形参名
					"err":      nil,      // 形参名
				},
				Error: "",
				Logs:  []string{"用户创建成功"},
			}, nil

		case "step2":
			// 验证step2的输入参数
			expectedInput := map[string]interface{}{
				"workId":     "EMP001", // 来自step1的输出
				"department": "技术部",    // 来自input
			}
			validateInput(t, in.RealInput, expectedInput, "step2")

			// 返回step2的输出
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": true,   // 形参名
					"message": "分配成功", // 形参名
					"err":     nil,    // 形参名
				},
				Error: "",
				Logs:  []string{"部门分配成功"},
			}, nil

		case "step3":
			// 验证step3的输入参数
			expectedInput := map[string]interface{}{
				"workId":     "EMP001", // 来自step1的输出
				"username":   "张三",     // 来自step1的输出
				"department": "技术部",    // 来自input
			}
			validateInput(t, in.RealInput, expectedInput, "step3")

			// 返回step3的输出
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"success": true, // 形参名
					"err":     nil,  // 形参名
				},
				Error: "",
				Logs:  []string{"通知发送成功"},
			}, nil

		default:
			t.Errorf("未知步骤: %s", step.Name)
			return &ExecutorOut{Success: false, Error: "未知步骤"}, nil
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s, 变量数量=%d", current.FlowID, len(current.Variables))
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *SimpleParseResult) error {
		t.Log("工作流正常结束")
		return nil
	}

	// 6. 执行工作流
	ctx := context.Background()
	if err := executor.Start(ctx, parseResult); err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 7. 验证最终变量状态
	expectedVariables := map[string]struct {
		Value interface{}
		Type  string
	}{
		"工号":       {"EMP001", "string"},
		"用户名":      {"张三", "string"},
		"step1Err": {nil, "error"},
		"分配成功":     {true, "bool"},
		"消息":       {"分配成功", "string"},
		"step2Err": {nil, "error"},
		"通知成功":     {true, "bool"},
		"step3Err": {nil, "error"},
	}

	for varName, expected := range expectedVariables {
		if varInfo, exists := parseResult.Variables[varName]; !exists {
			t.Errorf("缺少变量: %s", varName)
		} else {
			if varInfo.Value != expected.Value {
				t.Errorf("变量 %s 期望值 %v，实际值 %v", varName, expected.Value, varInfo.Value)
			}
			if varInfo.Type != expected.Type {
				t.Errorf("变量 %s 期望类型 %s，实际类型 %s", varName, expected.Type, varInfo.Type)
			}
			t.Logf("✅ 变量 %s: 类型=%s, 值=%v, 来源=%s", varName, varInfo.Type, varInfo.Value, varInfo.Source)
		}
	}

	// 8. 验证步骤定义中的参数描述
	validateStepParameterDescriptions(t, parseResult)
}

// validateInput 验证输入参数
func validateInput(t *testing.T, actual, expected map[string]interface{}, stepName string) {
	for key, expectedValue := range expected {
		if actualValue, exists := actual[key]; !exists {
			t.Errorf("[%s] 缺少输入参数: %s", stepName, key)
		} else if actualValue != expectedValue {
			t.Errorf("[%s] 输入参数 %s 期望值 %v，实际值 %v", stepName, key, expectedValue, actualValue)
		} else {
			t.Logf("[%s] ✅ 输入参数 %s: %v", stepName, key, actualValue)
		}
	}
}

// validateStepParameterDescriptions 验证步骤参数描述
func validateStepParameterDescriptions(t *testing.T, parseResult *SimpleParseResult) {
	t.Log("验证步骤参数描述:")

	for _, step := range parseResult.Steps {
		t.Logf("步骤 %s:", step.Name)

		// 验证输入参数描述
		for _, param := range step.InputParams {
			if param.Desc == "" {
				t.Errorf("步骤 %s 输入参数 %s 缺少描述", step.Name, param.Name)
			} else {
				t.Logf("  输入参数 %s (%s): %s", param.Name, param.Type, param.Desc)
			}
		}

		// 验证输出参数描述
		for _, param := range step.OutputParams {
			if param.Desc == "" {
				t.Errorf("步骤 %s 输出参数 %s 缺少描述", step.Name, param.Name)
			} else {
				t.Logf("  输出参数 %s (%s): %s", param.Name, param.Type, param.Desc)
			}
		}
	}
}

func TestParameterMapping_ErrorHandling(t *testing.T) {
	// 测试错误处理场景
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(
    username: string "用户名",
    phone: int "手机号"
) -> (
    workId: string "工号",
    username: string "用户名",
    err: error "是否失败"
);

func main() {
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    if step1Err != nil {
        fmt.Printf("用户创建失败: %v\n", step1Err)
        return
    }
    fmt.Printf("用户创建成功: %s\n", 用户名)
}`

	parser := NewSimpleParser()
	parseResult := parser.ParseWorkflow(code)
	if !parseResult.Success {
		t.Fatalf("解析失败: %s", parseResult.Error)
	}

	parseResult.FlowID = "test-error-handling-001"
	executor := NewExecutor()

	// 模拟步骤执行失败
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s", step.Name)

		// 模拟失败场景
		return &ExecutorOut{
			Success:    false,
			WantOutput: map[string]interface{}{},
			Error:      "用户已存在",
			Logs:       []string{"用户创建失败: 用户已存在"},
		}, nil
	}

	executor.OnWorkFlowReturn = func(ctx context.Context, current *SimpleParseResult) error {
		t.Log("工作流因错误中断")
		return nil
	}

	ctx := context.Background()
	err := executor.Start(ctx, parseResult)
	if err == nil {
		t.Error("期望执行失败，但实际成功了")
	}

	// 验证错误信息
	if err.Error() != "步骤执行失败: 用户已存在" {
		t.Errorf("期望错误信息包含 '用户已存在'，实际: %s", err.Error())
	}
}
