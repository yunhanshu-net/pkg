package workflow

import (
	"context"
	"testing"
)

func TestExecutor_ParameterMapping(t *testing.T) {
	// 1. 创建工作流代码
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
        return
    }
    fmt.Printf("工号: %s, 用户名: %s\n", 工号, 用户名)
}`

	// 2. 解析工作流
	parser := NewSimpleParser()
	parseResult := parser.ParseWorkflow(code)
	if !parseResult.Success {
		t.Fatalf("解析失败: %s", parseResult.Error)
	}

	// 3. 设置FlowID
	parseResult.FlowID = "test-flow-001"

	// 4. 创建执行器
	executor := NewExecutor()

	// 5. 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s", step.Name)
		t.Logf("输入参数: %+v", in.RealInput)

		// 验证输入参数映射
		expectedInput := map[string]interface{}{
			"username": "张三",        // 形参名 -> 实际值
			"phone":    13800138000, // 形参名 -> 实际值
		}

		for key, expectedValue := range expectedInput {
			if actualValue, exists := in.RealInput[key]; !exists {
				t.Errorf("缺少输入参数: %s", key)
			} else if actualValue != expectedValue {
				t.Errorf("输入参数 %s 期望值 %v，实际值 %v", key, expectedValue, actualValue)
			}
		}

		// 返回输出参数（使用形参名）
		wantOutput := map[string]interface{}{
			"workId":   "EMP001", // 形参名
			"username": "张三",     // 形参名
			"err":      nil,      // 形参名
		}

		return &ExecutorOut{
			Success:    true,
			WantOutput: wantOutput,
			Error:      "",
			Logs:       []string{"用户创建成功"},
		}, nil
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

	// 7. 验证输出参数映射
	expectedVariables := map[string]interface{}{
		"工号":       "EMP001", // 实例名 -> 实际值
		"用户名":      "张三",     // 实例名 -> 实际值
		"step1Err": nil,      // 实例名 -> 实际值
	}

	for varName, expectedValue := range expectedVariables {
		if varInfo, exists := parseResult.Variables[varName]; !exists {
			t.Errorf("缺少变量: %s", varName)
		} else if varInfo.Value != expectedValue {
			t.Errorf("变量 %s 期望值 %v，实际值 %v", varName, expectedValue, varInfo.Value)
		} else {
			t.Logf("变量 %s: 类型=%s, 值=%v, 来源=%s", varName, varInfo.Type, varInfo.Value, varInfo.Source)
		}
	}
}
