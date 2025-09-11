package workflow

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestExecutor_BasicExecution 测试基本执行功能
func TestExecutor_BasicExecution(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 开始执行用户创建流程
    sys.Println("开始执行用户创建流程...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    sys.Println("🎉 用户创建流程完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s - %s", in.StepName, in.StepDesc)
		t.Logf("输入参数: %v", in.RealInput)
		t.Logf("预期返回参数: %d个", len(in.WantParams))

		// 模拟成功执行
		return &ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"workId":   "EMP001",
				"username": in.RealInput["username"],
				"err":      nil,
			},
		}, nil
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s, 变量数量=%d", current.FlowID, len(current.Variables))
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流退出: FlowID=%s", current.FlowID)
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 验证执行结果
	if result.MainFunc.Statements[1].Status != "completed" {
		t.Errorf("步骤执行状态不正确: 期望 completed, 实际 %s", result.MainFunc.Statements[1].Status)
	}

	// 验证变量映射
	if workId, exists := result.Variables["工号"]; !exists {
		t.Error("缺少变量: 工号")
	} else if workId.Value != "EMP001" {
		t.Errorf("工号值不正确: 期望 EMP001, 实际 %v", workId.Value)
	}
}

// TestExecutor_ErrorHandling 测试错误处理
func TestExecutor_ErrorHandling(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行用户创建流程...")
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置返回错误的回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s", in.StepName)
		// 模拟执行失败，返回错误
		return nil, fmt.Errorf("模拟业务错误")
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s", current.FlowID)
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Logf("执行失败（预期）: %v", err)
	}

	// 验证错误处理
	if result.MainFunc.Statements[1].Status != "failed" {
		t.Errorf("步骤执行状态不正确: 期望 failed, 实际 %s", result.MainFunc.Statements[1].Status)
	}
}

// TestExecutor_ContextCancellation 测试上下文取消
func TestExecutor_ContextCancellation(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行用户创建流程...")
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置长时间执行的回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s", in.StepName)
		// 模拟长时间执行，等待上下文取消
		select {
		case <-ctx.Done():
			t.Logf("上下文被取消: %v", ctx.Err())
			return nil, ctx.Err()
		case <-time.After(2 * time.Second):
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":   "EMP001",
					"username": in.RealInput["username"],
					"err":      nil,
				},
			}, nil
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s", current.FlowID)
		return nil
	}

	// 创建可取消的上下文
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	// 执行工作流
	err := executor.Start(ctx, result)
	if err == nil {
		t.Error("期望执行被取消，但实际成功")
	}

	// 验证取消状态
	if result.MainFunc.Statements[1].Status != "cancelled" {
		t.Errorf("步骤执行状态不正确: 期望 cancelled, 实际 %s", result.MainFunc.Statements[1].Status)
	}
}

// TestExecutor_MetadataHandling 测试元数据处理
func TestExecutor_MetadataHandling(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行用户创建流程...")
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high", debug:true}
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置回调函数验证元数据
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s", in.StepName)
		t.Logf("执行选项: %+v", in.Options)

		// 验证元数据传递
		if in.Options == nil {
			t.Error("执行选项不应该为nil")
		}

		// 验证重试次数
		if in.Options.RetryCount != 3 {
			t.Errorf("重试次数不匹配: 期望 3, 实际 %d", in.Options.RetryCount)
		}

		// 验证超时时间
		if in.Options.Timeout == nil || *in.Options.Timeout != 5*time.Second {
			t.Errorf("超时时间不匹配: 期望 5s, 实际 %v", in.Options.Timeout)
		}

		// 验证优先级
		if in.Options.Priority != 1 { // high = 1
			t.Errorf("优先级不匹配: 期望 1, 实际 %d", in.Options.Priority)
		}

		// 验证调试模式
		if !in.Options.Debug {
			t.Error("调试模式应该为true")
		}

		return &ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"workId":   "EMP001",
				"username": in.RealInput["username"],
				"err":      nil,
			},
		}, nil
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s", current.FlowID)
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}
}

// TestExecutor_ComplexWorkflow 测试复杂工作流执行
func TestExecutor_ComplexWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", err: error "是否失败");
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");

func main() {
    sys.Println("开始执行复杂工作流...")
    
    // 创建用户
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名, 部门)
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
    
    // 发送通知
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}"
    step3Err := step3(input["邮箱"], 邮件主题, 邮件内容)
    if step3Err != nil {
        step3.Printf("发送邮件失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 邮件发送成功")
    
    sys.Println("🎉 复杂工作流执行完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		t.Logf("执行步骤: %s - %s", in.StepName, in.StepDesc)
		t.Logf("输入参数: %v", in.RealInput)

		// 根据步骤名称返回不同的结果
		switch in.StepName {
		case "step1":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":     "EMP001",
					"username":   in.RealInput["username"],
					"department": in.RealInput["department"],
					"err":        nil,
				},
			}, nil
		case "step2":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"interviewTime": "2024-01-15 14:00",
					"interviewer":   "李面试官",
					"err":           nil,
				},
			}, nil
		case "step3":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"err": nil,
				},
			}, nil
		default:
			return &ExecutorOut{
				Success: false,
				Error:   "未知步骤",
			}, nil
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s, 变量数量=%d", current.FlowID, len(current.Variables))
		return nil
	}

	executor.OnWorkFlowExit = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流退出: FlowID=%s", current.FlowID)
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 验证所有步骤都执行完成
	expectedSteps := []string{"step1", "step2", "step3"}
	for i, stepName := range expectedSteps {
		stmt := result.MainFunc.Statements[i*2+1] // 每个步骤在function-call语句中
		if stmt.Status != "completed" {
			t.Errorf("步骤 %s 执行状态不正确: 期望 completed, 实际 %s", stepName, stmt.Status)
		}
	}

	// 验证变量映射
	expectedVars := []string{"工号", "用户名", "部门", "面试时间", "面试官名称", "step1Err", "step2Err", "step3Err"}
	for _, varName := range expectedVars {
		if _, exists := result.Variables[varName]; !exists {
			t.Errorf("缺少变量: %s", varName)
		}
	}
}

// TestExecutor_GetMethod 测试Get方法
func TestExecutor_GetMethod(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行用户创建流程...")
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置回调函数
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		return &ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"workId":   "EMP001",
				"username": in.RealInput["username"],
				"err":      nil,
			},
		}, nil
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		// 不执行任何操作
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Fatalf("执行失败: %v", err)
	}

	// 测试Get方法
	retrievedResult, err := executor.Get(result.FlowID)
	if err != nil {
		t.Fatalf("Get方法失败: %v", err)
	}
	if retrievedResult == nil {
		t.Fatal("Get方法返回nil")
	}

	// 验证获取的结果
	if retrievedResult.FlowID != result.FlowID {
		t.Errorf("FlowID不匹配: 期望 %s, 实际 %s", result.FlowID, retrievedResult.FlowID)
	}

	if len(retrievedResult.Steps) != len(result.Steps) {
		t.Errorf("步骤数量不匹配: 期望 %d, 实际 %d", len(result.Steps), len(retrievedResult.Steps))
	}

	if len(retrievedResult.Variables) != len(result.Variables) {
		t.Errorf("变量数量不匹配: 期望 %d, 实际 %d", len(result.Variables), len(retrievedResult.Variables))
	}
}
