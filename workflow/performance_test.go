package workflow

import (
	"context"
	"testing"
	"time"
)

// TestParserPerformance 测试解析器性能
func TestParserPerformance(t *testing.T) {
	// 简单工作流
	simpleCode := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 执行用户创建
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        //desc: 处理创建失败情况
        return
    }
    //desc: 输出创建成功信息
}`

	// 中等复杂度工作流
	mediumCode := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
    "职位": "高级工程师",
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
//desc: 安排面试时间
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门", position: string "职位") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", interviewLocation: string "面试地点", err: error "是否失败");
//desc: 发送邮件通知
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");
//desc: 发送短信通知
step4 = beiluo.test1.notification.send_sms(phone: int "手机号", content: string "内容") -> (err: error "是否失败");
//desc: 创建候选人记录
step5 = beiluo.test1.crm.crm_create_candidate(workId: string "工号", username: string "用户名", department: string "部门", position: string "职位") -> (candidateId: string "候选人ID", err: error "是否失败");

func main() {
    
    //desc: 执行用户创建
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step1Err != nil {
        //desc: 处理用户创建失败
        return
    }
    //desc: 输出用户创建成功信息
    
    //desc: 执行面试安排
    面试时间, 面试官名称, 面试地点, step2Err := step2(用户名, 部门, input["职位"])
    if step2Err != nil {
        //desc: 处理面试安排失败
        return
    }
    //desc: 输出面试安排成功信息
    
    //desc: 准备邮件内容
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}, 地点: {{面试地点}}"
    //desc: 发送邮件通知
    step3Err := step3(input["邮箱"], 邮件主题, 邮件内容)
    if step3Err != nil {
        //desc: 处理邮件发送失败
        return
    }
    //desc: 输出邮件发送成功信息
    
    //desc: 准备短信内容
    短信内容 := "面试安排: {{面试时间}} {{面试地点}}"
    //desc: 发送短信通知
    step4Err := step4(input["手机号"], 短信内容)
    if step4Err != nil {
        //desc: 处理短信发送失败
        return
    }
    //desc: 输出短信发送成功信息
    
    //desc: 创建候选人记录
    候选人ID, step5Err := step5(工号, 用户名, 部门, input["职位"])
    if step5Err != nil {
        //desc: 处理候选人记录创建失败
        return
    }
    //desc: 输出候选人记录创建成功信息
    
}`

	// 复杂工作流
	complexCode := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
    "职位": "高级工程师",
    "项目名称": "AI工作流引擎",
    "版本": "v1.0.0",
    "环境": "production",
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
//desc: 安排面试时间
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门", position: string "职位") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", interviewLocation: string "面试地点", err: error "是否失败");
//desc: 发送邮件通知
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");
//desc: 发送短信通知
step4 = beiluo.test1.notification.send_sms(phone: int "手机号", content: string "内容") -> (err: error "是否失败");
//desc: 创建候选人记录
step5 = beiluo.test1.crm.crm_create_candidate(workId: string "工号", username: string "用户名", department: string "部门", position: string "职位") -> (candidateId: string "候选人ID", err: error "是否失败");
//desc: 推送代码到仓库
step6 = beiluo.test1.devops.git_push[用例001] -> (err: error "是否失败");
//desc: 构建项目
step7 = beiluo.test1.devops.build_project(projectName: string "项目名称", version: string "版本", environment: string "环境") -> (buildId: string "构建ID", buildStatus: string "构建状态", err: error "是否失败");
//desc: 部署服务
step8 = beiluo.test1.devops.deploy_service(serviceName: string "服务名称", version: string "版本", environment: string "环境") -> (deploymentId: string "部署ID", deploymentStatus: string "部署状态", err: error "是否失败");

func main() {
    
    //desc: 执行用户创建
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step1Err != nil {
        //desc: 处理用户创建失败
        return
    }
    //desc: 输出用户创建成功信息
    
    //desc: 执行面试安排
    面试时间, 面试官名称, 面试地点, step2Err := step2(用户名, 部门, input["职位"])
    if step2Err != nil {
        //desc: 处理面试安排失败
        return
    }
    //desc: 输出面试安排成功信息
    
    //desc: 准备邮件内容
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}, 地点: {{面试地点}}"
    //desc: 发送邮件通知
    step3Err := step3(input["邮箱"], 邮件主题, 邮件内容)
    if step3Err != nil {
        //desc: 处理邮件发送失败
        return
    }
    //desc: 输出邮件发送成功信息
    
    //desc: 准备短信内容
    短信内容 := "面试安排: {{面试时间}} {{面试地点}}"
    //desc: 发送短信通知
    step4Err := step4(input["手机号"], 短信内容)
    if step4Err != nil {
        //desc: 处理短信发送失败
        return
    }
    //desc: 输出短信发送成功信息
    
    //desc: 创建候选人记录
    候选人ID, step5Err := step5(工号, 用户名, 部门, input["职位"])
    if step5Err != nil {
        //desc: 处理候选人记录创建失败
        return
    }
    //desc: 输出候选人记录创建成功信息
    
    //desc: 执行代码推送
    step6Err := step6()
    if step6Err != nil {
        //desc: 处理代码推送失败
        return
    }
    //desc: 输出代码推送成功信息
    
    //desc: 执行项目构建
    构建ID, 构建状态, step7Err := step7(input["项目名称"], input["版本"], input["环境"])
    if step7Err != nil {
        //desc: 处理项目构建失败
        return
    }
    //desc: 输出项目构建成功信息
    
    //desc: 执行服务部署
    部署ID, 部署状态, step8Err := step8(input["项目名称"], input["版本"], input["环境"])
    if step8Err != nil {
        //desc: 处理服务部署失败
        return
    }
    //desc: 输出服务部署成功信息
    
}`

	parser := NewSimpleParser()

	// 测试简单工作流性能
	t.Run("Simple", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < 1000; i++ {
			result := parser.ParseWorkflow(simpleCode)
			if !result.Success {
				t.Fatalf("解析失败: %s", result.Error)
			}
		}
		duration := time.Since(start)
		t.Logf("简单工作流解析1000次耗时: %v", duration)
		t.Logf("平均每次解析耗时: %v", duration/1000)

		if duration/1000 > time.Millisecond {
			t.Errorf("简单工作流解析性能不达标: 平均耗时 %v > 1ms", duration/1000)
		}
	})

	// 测试中等复杂度工作流性能
	t.Run("Medium", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < 500; i++ {
			result := parser.ParseWorkflow(mediumCode)
			if !result.Success {
				t.Fatalf("解析失败: %s", result.Error)
			}
		}
		duration := time.Since(start)
		t.Logf("中等复杂度工作流解析500次耗时: %v", duration)
		t.Logf("平均每次解析耗时: %v", duration/500)

		if duration/500 > 2*time.Millisecond {
			t.Errorf("中等复杂度工作流解析性能不达标: 平均耗时 %v > 2ms", duration/500)
		}
	})

	// 测试复杂工作流性能
	t.Run("Complex", func(t *testing.T) {
		start := time.Now()
		for i := 0; i < 100; i++ {
			result := parser.ParseWorkflow(complexCode)
			if !result.Success {
				t.Fatalf("解析失败: %s", result.Error)
			}
		}
		duration := time.Since(start)
		t.Logf("复杂工作流解析100次耗时: %v", duration)
		t.Logf("平均每次解析耗时: %v", duration/100)

		if duration/100 > 5*time.Millisecond {
			t.Errorf("复杂工作流解析性能不达标: 平均耗时 %v > 5ms", duration/100)
		}
	})
}

// TestExecutorPerformance 测试执行器性能
func TestExecutorPerformance(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 执行用户创建
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        //desc: 处理创建失败情况
        return
    }
    //desc: 输出创建成功信息
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 创建执行器
	executor := NewExecutor()

	// 设置快速回调函数
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

	// 测试执行性能
	start := time.Now()
	for i := 0; i < 100; i++ {
		// 为每次执行创建新的结果副本
		resultCopy := *result
		resultCopy.FlowID = "flow_" + string(rune('0'+i%10))

		ctx := context.Background()
		err := executor.Start(ctx, &resultCopy)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}
	}
	duration := time.Since(start)
	t.Logf("执行100次工作流耗时: %v", duration)
	t.Logf("平均每次执行耗时: %v", duration/100)

	if duration/100 > 10*time.Millisecond {
		t.Errorf("执行性能不达标: 平均耗时 %v > 10ms", duration/100)
	}
}

// TestMemoryUsage 测试内存使用
func TestMemoryUsage(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
//desc: 安排面试时间
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", err: error "是否失败");
//desc: 发送邮件通知
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");
//desc: 发送短信通知
step4 = beiluo.test1.notification.send_sms(phone: int "手机号", content: string "内容") -> (err: error "是否失败");

func main() {
    
    //desc: 执行用户创建
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step1Err != nil {
        //desc: 处理用户创建失败
        return
    }
    //desc: 输出用户创建成功信息
    
    //desc: 执行面试安排
    面试时间, 面试官名称, step2Err := step2(用户名, 部门)
    if step2Err != nil {
        //desc: 处理面试安排失败
        return
    }
    //desc: 输出面试安排成功信息
    
    //desc: 准备邮件内容
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}"
    //desc: 发送邮件通知
    step3Err := step3(input["邮箱"], 邮件主题, 邮件内容)
    if step3Err != nil {
        //desc: 处理邮件发送失败
        return
    }
    //desc: 输出邮件发送成功信息
    
    //desc: 准备短信内容
    短信内容 := "面试安排: {{面试时间}}"
    //desc: 发送短信通知
    step4Err := step4(input["手机号"], 短信内容)
    if step4Err != nil {
        //desc: 处理短信发送失败
        return
    }
    //desc: 输出短信发送成功信息
    
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)
	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 统计内存使用情况
	t.Logf("工作流步骤数量: %d", len(result.Steps))
	t.Logf("主函数语句数量: %d", len(result.MainFunc.Statements))
	t.Logf("变量映射数量: %d", len(result.Variables))

	// 统计参数数量
	totalParams := 0
	totalReturns := 0
	for _, step := range result.Steps {
		totalParams += len(step.InputParams)
		totalReturns += len(step.OutputParams)
	}

	t.Logf("总参数数量: %d", totalParams)
	t.Logf("总返回值数量: %d", totalReturns)

	// 验证内存使用合理
	if len(result.Steps) > 10 {
		t.Errorf("步骤数量过多: %d", len(result.Steps))
	}
	if len(result.MainFunc.Statements) > 50 {
		t.Errorf("主函数语句数量过多: %d", len(result.MainFunc.Statements))
	}
	if len(result.Variables) > 100 {
		t.Errorf("变量数量过多: %d", len(result.Variables))
	}
}

// TestConcurrentPerformance 测试并发性能
func TestConcurrentPerformance(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 执行用户创建
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        //desc: 处理创建失败情况
        return
    }
    //desc: 输出创建成功信息
}`

	parser := NewSimpleParser()

	// 并发执行测试
	numGoroutines := 10
	numExecutions := 100
	done := make(chan bool, numGoroutines)

	start := time.Now()

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			defer func() { done <- true }()

			executor := NewExecutor()
			executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
				return &ExecutorOut{
					Success: true,
					WantOutput: map[string]interface{}{
						"workId":   "EMP" + string(rune('0'+id)),
						"username": in.RealInput["username"],
						"err":      nil,
					},
				}, nil
			}
			executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
				// 不执行任何操作
				return nil
			}

			for j := 0; j < numExecutions/numGoroutines; j++ {
				result := parser.ParseWorkflow(code)
				if !result.Success {
					t.Errorf("解析失败: %s", result.Error)
					return
				}
				result.FlowID = "flow_" + string(rune('0'+id)) + "_" + string(rune('0'+j))

				ctx := context.Background()
				err := executor.Start(ctx, result)
				if err != nil {
					t.Errorf("执行失败: %v", err)
					return
				}
			}
		}(i)
	}

	// 等待所有goroutine完成
	for i := 0; i < numGoroutines; i++ {
		<-done
	}

	duration := time.Since(start)
	t.Logf("并发执行 %d 次工作流耗时: %v", numExecutions, duration)
	t.Logf("平均每次执行耗时: %v", duration/time.Duration(numExecutions))

	if duration/time.Duration(numExecutions) > 20*time.Millisecond {
		t.Errorf("并发执行性能不达标: 平均耗时 %v > 20ms", duration/time.Duration(numExecutions))
	}
}

// TestStressTest 测试压力测试
func TestStressTest(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

//desc: 创建用户账号
step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 执行用户创建
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        //desc: 处理创建失败情况
        return
    }
    //desc: 输出创建成功信息
}`

	parser := NewSimpleParser()

	// 压力测试：连续解析和执行
	numIterations := 1000
	start := time.Now()

	for i := 0; i < numIterations; i++ {
		// 解析
		result := parser.ParseWorkflow(code)
		if !result.Success {
			t.Fatalf("解析失败: %s", result.Error)
		}

		// 创建执行器
		executor := NewExecutor()
		executor.OnFunctionCall = func(ctx context.Context, step SimpleStep, in *ExecutorIn) (*ExecutorOut, error) {
			return &ExecutorOut{
				Success: true,
				WantOutput: map[string]interface{}{
					"workId":   "EMP" + string(rune('0'+i%10)),
					"username": in.RealInput["username"],
					"err":      nil,
				},
			}, nil
		}
		executor.OnWorkFlowUpdate = func(ctx context.Context, current *SimpleParseResult) error {
			// 不执行任何操作
			return nil
		}

		// 执行
		ctx := context.Background()
		err := executor.Start(ctx, result)
		if err != nil {
			t.Fatalf("执行失败: %v", err)
		}
	}

	duration := time.Since(start)
	t.Logf("压力测试 %d 次迭代耗时: %v", numIterations, duration)
	t.Logf("平均每次迭代耗时: %v", duration/time.Duration(numIterations))

	if duration/time.Duration(numIterations) > 50*time.Millisecond {
		t.Errorf("压力测试性能不达标: 平均耗时 %v > 50ms", duration/time.Duration(numIterations))
	}
}
