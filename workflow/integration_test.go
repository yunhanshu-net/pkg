package workflow

import (
	"context"
	"fmt"
	"testing"
	"time"
)

// TestEndToEndWorkflow 测试端到端工作流
func TestEndToEndWorkflow(t *testing.T) {
	// 完整的端到端工作流代码
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
    "职位": "高级工程师",
    "项目名称": "AI工作流引擎",
    "版本": "v1.0.0",
}

// 静态工作流步骤
step1 = beiluo.test1.devops.git_push[用例001] -> (err: error "是否失败");
step2 = beiluo.test1.devops.deploy_test[用例002] -> (cost: int "成本", err: error "是否失败");

// 动态工作流步骤
step3 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
step4 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门", position: string "职位") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", interviewLocation: string "面试地点", err: error "是否失败");
step5 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");
step6 = beiluo.test1.notification.send_sms(phone: int "手机号", content: string "内容") -> (err: error "是否失败");
step7 = beiluo.test1.crm.crm_create_candidate(workId: string "工号", username: string "用户名", department: string "部门", position: string "职位") -> (candidateId: string "候选人ID", err: error "是否失败");
step8 = beiluo.test1.devops.build_project(projectName: string "项目名称", version: string "版本") -> (buildId: string "构建ID", buildStatus: string "构建状态", err: error "是否失败");

func main() {
    //desc: 开始执行完整的工作流
    sys.Println("🚀 开始执行完整的工作流...")
    
    //desc: 推送代码到远程仓库
    err := step1()
    if err != nil {
        step1.Printf("推送代码失败: %v", err)
        return
    }
    step1.Printf("✅ 代码推送成功")
    
    //desc: 部署到测试环境
    cost, err := step2()
    if err != nil {
        step2.Printf("测试环境部署失败: %v", err)
        return
    }
    step2.Printf("✅ 测试环境部署成功，成本: %d", cost)
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, 部门, step3Err := step3(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step3Err != nil {
        step3.Printf("创建用户失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    //desc: 安排面试时间，联系面试官
    面试时间, 面试官名称, 面试地点, step4Err := step4(用户名, 部门, input["职位"])
    if step4Err != nil {
        step4.Printf("安排面试失败: %v", step4Err)
        return
    }
    step4.Printf("✅ 面试安排成功，时间: %s, 地点: %s", 面试时间, 面试地点)
    
    //desc: 发送邮件通知
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}, 地点: {{面试地点}}"
    step5Err := step5(input["邮箱"], 邮件主题, 邮件内容)
    if step5Err != nil {
        step5.Printf("发送邮件失败: %v", step5Err)
        return
    }
    step5.Printf("✅ 邮件发送成功")
    
    //desc: 发送短信通知
    短信内容 := "面试安排: {{面试时间}} {{面试地点}}"
    step6Err := step6(input["手机号"], 短信内容)
    if step6Err != nil {
        step6.Printf("发送短信失败: %v", step6Err)
        return
    }
    step6.Printf("✅ 短信发送成功")
    
    //desc: 创建候选人记录
    候选人ID, step7Err := step7(工号, 用户名, 部门, input["职位"])
    if step7Err != nil {
        step7.Printf("创建候选人记录失败: %v", step7Err)
        return
    }
    step7.Printf("✅ 候选人记录创建成功，ID: %s", 候选人ID)
    
    //desc: 构建项目
    构建ID, 构建状态, step8Err := step8(input["项目名称"], input["版本"])
    if step8Err != nil {
        step8.Printf("构建项目失败: %v", step8Err)
        return
    }
    step8.Printf("✅ 项目构建成功，ID: %s, 状态: %s", 构建ID, 构建状态)
    
    sys.Println("🎉 完整工作流执行完成！")
}`

	// 解析工作流
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

		// 根据步骤名称返回不同的结果
		switch in.StepName {
		case "step1":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"err": nil,
				},
			}, nil
		case "step2":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"cost": 1000,
					"err":  nil,
				},
			}, nil
		case "step3":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":     "EMP001",
					"username":   in.RealInput["username"],
					"department": in.RealInput["department"],
					"err":        nil,
				},
			}, nil
		case "step4":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"interviewTime":     "2024-01-15 14:00",
					"interviewer":       "李面试官",
					"interviewLocation": "会议室A",
					"err":               nil,
				},
			}, nil
		case "step5", "step6":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"err": nil,
				},
			}, nil
		case "step7":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"candidateId": "CAND001",
					"err":         nil,
				},
			}, nil
		case "step8":
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"buildId":     "BUILD001",
					"buildStatus": "success",
					"err":         nil,
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

	// 验证执行结果
	t.Logf("工作流执行完成，FlowID: %s", result.FlowID)
	t.Logf("步骤数量: %d", len(result.Steps))
	t.Logf("主函数语句数量: %d", len(result.MainFunc.Statements))
	t.Logf("变量数量: %d", len(result.Variables))

	// 验证所有步骤都执行完成
	for i, stmt := range result.MainFunc.Statements {
		if stmt.Type == "function-call" {
			if stmt.Status != "completed" {
				t.Errorf("步骤 %d 执行状态不正确: 期望 completed, 实际 %s", i, stmt.Status)
			}
		}
	}

	// 验证关键变量
	expectedVars := []string{"工号", "用户名", "部门", "面试时间", "面试官名称", "面试地点", "候选人ID", "构建ID", "构建状态"}
	for _, varName := range expectedVars {
		if _, exists := result.Variables[varName]; !exists {
			t.Errorf("缺少变量: %s", varName)
		}
	}
}

// TestIntegrationWithPersistence 测试与持久化的集成
func TestIntegrationWithPersistence(t *testing.T) {
	// 工作流代码
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱") -> (workId: string "工号", username: string "用户名", email: string "邮箱", err: error "是否失败");
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", err: error "是否失败");
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");

func main() {
    sys.Println("开始执行集成测试工作流...")
    
    // 创建用户
    工号, 用户名, 邮箱, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
    
    // 发送通知
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}"
    step3Err := step3(邮箱, 邮件主题, 邮件内容)
    if step3Err != nil {
        step3.Printf("发送邮件失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 邮件发送成功")
    
    sys.Println("🎉 集成测试工作流执行完成！")
}`

	// 解析工作流
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

		// 模拟成功执行
		return &ExecutorOut{
			Success: true,
			WantOutput: map[string]interface{}{
				"workId":        "EMP001",
				"username":      in.RealInput["username"],
				"email":         in.RealInput["email"],
				"interviewTime": "2024-01-15 14:00",
				"interviewer":   "李面试官",
				"err":           nil,
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

	// 验证工作流执行完成
	if result.FlowID == "" {
		t.Error("FlowID不应该为空")
	}

	t.Logf("✅ 集成测试通过，工作流执行完成")
}

// TestIntegrationWithMetadata 测试与元数据的集成
func TestIntegrationWithMetadata(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行带元数据的工作流...")
    
    // 带元数据的函数调用
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000, priority:"high", debug:true, log_level:"debug", ai_model:"gpt-4"}
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
}`

	// 解析工作流
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

		// 验证日志级别
		if in.Options.LogLevel != "debug" {
			t.Errorf("日志级别不匹配: 期望 debug, 实际 %s", in.Options.LogLevel)
		}

		// 验证AI模型
		if in.Options.AIModel != "gpt-4" {
			t.Errorf("AI模型不匹配: 期望 gpt-4, 实际 %s", in.Options.AIModel)
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

	t.Logf("✅ 元数据集成测试通过")
}

// TestIntegrationErrorRecovery 测试错误恢复
func TestIntegrationErrorRecovery(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", err: error "是否失败");

func main() {
    sys.Println("开始执行错误恢复测试工作流...")
    
    // 第一个步骤成功
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    // 第二个步骤失败
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
}`

	// 解析工作流
	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置回调函数，模拟第一个步骤成功，第二个步骤失败
	stepCount := 0
	executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
		stepCount++
		t.Logf("执行步骤 %d: %s", stepCount, in.StepName)

		if stepCount == 1 {
			// 第一个步骤成功
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":   "EMP001",
					"username": in.RealInput["username"],
					"err":      nil,
				},
			}, nil
		} else {
			// 第二个步骤失败，返回错误
			return nil, fmt.Errorf("模拟业务错误")
		}
	}

	executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
		t.Logf("工作流状态更新: FlowID=%s", current.FlowID)
		return nil
	}

	// 执行工作流
	ctx := context.Background()
	err := executor.Start(ctx, result)
	if err != nil {
		t.Logf("工作流执行失败（预期）: %v", err)
	}

	// 验证错误处理
	// 第一个步骤是 step1 调用
	if result.MainFunc.Statements[1].Status != "completed" {
		t.Errorf("第一个步骤状态不正确: 期望 completed, 实际 %s", result.MainFunc.Statements[1].Status)
	}

	// 第二个步骤是 step2 调用
	if result.MainFunc.Statements[3].Status != "failed" {
		t.Errorf("第二个步骤状态不正确: 期望 failed, 实际 %s", result.MainFunc.Statements[3].Status)
	}

	t.Logf("✅ 错误恢复测试通过")
}
