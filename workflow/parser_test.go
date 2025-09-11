package workflow

import (
	"testing"
	"time"
)

// TestSimpleParser_StaticWorkflow 测试静态工作流解析
func TestSimpleParser_StaticWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "项目名称": "my-project",
    "环境": "production",
    "版本": "v1.0.0",
}

step1 = beiluo.test1.devops.git_push[用例001] -> (err: error "是否失败");
step2 = beiluo.test1.devops.deploy_test[用例002] -> (cost: int "成本", err: error "是否失败");
step3 = beiluo.test1.devops.deploy_prod[用例003] -> (msg: string "消息", err: error "是否失败");
step4 = beiluo.test1.notify.send_notification[用例004] -> (err: error "是否失败");

func main() {
    //desc: 开始执行发布流程
    sys.Println("开始执行发布流程...")
    
    //desc: 推送代码到远程仓库
    err := step1()
    if err != nil {
        step1.Printf("推送代码失败: %v", err)
        return
    }
    step1.Printf("✅ 代码推送成功")
    
    //desc: 部署到测试环境
    err = step2()
    if err != nil {
        step2.Printf("发布测试环境失败: %v", err)
        return
    }
    step2.Printf("✅ 测试环境发布成功")
    
    //desc: 部署到生产环境
    err = step3()
    if err != nil {
        step3.Printf("发布线上环境失败: %v", err)
        return
    }
    step3.Printf("✅ 线上环境发布成功")
    
    //desc: 发送部署完成通知
    err = step4()
    if err != nil {
        step4.Printf("发送通知失败: %v", err)
        return
    }
    step4.Printf("✅ 通知发送成功")
    
    sys.Println("🎉 发布流程执行完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证输入变量
	expectedInputVars := map[string]interface{}{
		"项目名称": "my-project",
		"环境":   "production",
		"版本":   "v1.0.0",
	}
	for key, expectedValue := range expectedInputVars {
		if actualValue, exists := result.InputVars[key]; !exists {
			t.Errorf("缺少输入变量: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("输入变量 %s 值不匹配: 期望 %v, 实际 %v", key, expectedValue, actualValue)
		}
	}

	// 验证步骤定义
	if len(result.Steps) != 4 {
		t.Fatalf("步骤数量不匹配: 期望 4, 实际 %d", len(result.Steps))
	}

	// 验证第一个步骤
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤名称不匹配: 期望 step1, 实际 %s", step1.Name)
	}
	if step1.Function != "beiluo.test1.devops.git_push" {
		t.Errorf("函数名称不匹配: 期望 beiluo.test1.devops.git_push, 实际 %s", step1.Function)
	}
	if !step1.IsStatic {
		t.Error("步骤应该是静态工作流")
	}
	if step1.CaseID != "用例001" {
		t.Errorf("用例ID不匹配: 期望 用例001, 实际 %s", step1.CaseID)
	}

	// 验证主函数语句
	if len(result.MainFunc.Statements) < 10 {
		t.Errorf("主函数语句数量不足: 期望至少 10, 实际 %d", len(result.MainFunc.Statements))
	}

	// 验证第一个语句是print类型
	firstStmt := result.MainFunc.Statements[0]
	if firstStmt.Type != "print" {
		t.Errorf("第一个语句类型不匹配: 期望 print, 实际 %s", firstStmt.Type)
	}
	if firstStmt.Desc != "开始执行发布流程" {
		t.Errorf("第一个语句描述不匹配: 期望 '开始执行发布流程', 实际 '%s'", firstStmt.Desc)
	}
}

// TestSimpleParser_DynamicWorkflow 测试动态工作流解析
func TestSimpleParser_DynamicWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", err: error "是否失败");
step3 = beiluo.test1.crm.crm_interview_notify(interviewer: string "面试官名称", message: string "通知信息") -> (err: error "是否失败");

func main() {
    //desc: 开始执行动态工作流
    sys.Println("开始执行动态工作流...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    //desc: 安排面试时间，联系面试官
    面试时间, 面试官名称, step2Err := step2(用户名)
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s", 面试时间)
    
    //desc: 准备通知信息，使用模板变量
    通知信息 := "你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注"
    
    //desc: 发送面试通知给面试官
    step3Err := step3(面试官名称, 通知信息)
    if step3Err != nil {
        step3.Printf("发送通知失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 通知发送成功")
    
    sys.Println("🎉 动态工作流执行完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证步骤定义
	if len(result.Steps) != 3 {
		t.Fatalf("步骤数量不匹配: 期望 3, 实际 %d", len(result.Steps))
	}

	// 验证第一个步骤（动态工作流）
	step1 := result.Steps[0]
	if step1.Name != "step1" {
		t.Errorf("步骤名称不匹配: 期望 step1, 实际 %s", step1.Name)
	}
	if step1.Function != "beiluo.test1.devops.devops_script_create" {
		t.Errorf("函数名称不匹配: 期望 beiluo.test1.devops.devops_script_create, 实际 %s", step1.Function)
	}
	if step1.IsStatic {
		t.Error("步骤应该是动态工作流")
	}
	if step1.CaseID != "" {
		t.Errorf("动态工作流不应该有用例ID: 实际 %s", step1.CaseID)
	}

	// 验证输入参数
	if len(step1.InputParams) != 2 {
		t.Errorf("输入参数数量不匹配: 期望 2, 实际 %d", len(step1.InputParams))
	}

	// 验证第一个输入参数
	param1 := step1.InputParams[0]
	if param1.Name != "username" {
		t.Errorf("第一个参数名称不匹配: 期望 username, 实际 %s", param1.Name)
	}
	if param1.Type != "string" {
		t.Errorf("第一个参数类型不匹配: 期望 string, 实际 %s", param1.Type)
	}
	if param1.Desc != "用户名" {
		t.Errorf("第一个参数描述不匹配: 期望 '用户名', 实际 '%s'", param1.Desc)
	}

	// 验证输出参数
	if len(step1.OutputParams) != 3 {
		t.Errorf("输出参数数量不匹配: 期望 3, 实际 %d", len(step1.OutputParams))
	}

	// 验证第一个输出参数
	output1 := step1.OutputParams[0]
	if output1.Name != "workId" {
		t.Errorf("第一个输出参数名称不匹配: 期望 workId, 实际 %s", output1.Name)
	}
	if output1.Type != "string" {
		t.Errorf("第一个输出参数类型不匹配: 期望 string, 实际 %s", output1.Type)
	}
	if output1.Desc != "工号" {
		t.Errorf("第一个输出参数描述不匹配: 期望 '工号', 实际 '%s'", output1.Desc)
	}

	// 验证主函数中的function-call语句
	functionCallStmt := result.MainFunc.Statements[1]
	if functionCallStmt.Type != "function-call" {
		t.Errorf("第二个语句类型不匹配: 期望 function-call, 实际 %s", functionCallStmt.Type)
	}
	if functionCallStmt.Function != "step1" {
		t.Errorf("函数调用名称不匹配: 期望 step1, 实际 %s", functionCallStmt.Function)
	}

	// 验证参数映射
	if len(functionCallStmt.Args) != 2 {
		t.Errorf("函数调用参数数量不匹配: 期望 2, 实际 %d", len(functionCallStmt.Args))
	}

	// 验证返回值映射
	if len(functionCallStmt.Returns) != 3 {
		t.Errorf("函数调用返回值数量不匹配: 期望 3, 实际 %d", len(functionCallStmt.Returns))
	}
}

// TestSimpleParser_MetadataSupport 测试元数据支持
func TestSimpleParser_MetadataSupport(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    //desc: 开始执行带元数据的工作流
    sys.Println("开始执行带元数据的工作流...")
    
    //desc: 创建用户账号，带重试和超时配置
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry: 3, timeout: 5000, priority: "high", debug: true}
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    sys.Println("🎉 带元数据的工作流执行完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证元数据解析
	functionCallStmt := result.MainFunc.Statements[1]
	if functionCallStmt.Metadata == nil {
		t.Fatal("元数据不应该为nil")
	}

	expectedMetadata := map[string]interface{}{
		"retry":    3,
		"timeout":  5000,
		"priority": "high",
		"debug":    true,
	}

	for key, expectedValue := range expectedMetadata {
		if actualValue, exists := functionCallStmt.Metadata[key]; !exists {
			t.Errorf("缺少元数据: %s", key)
		} else if actualValue != expectedValue {
			t.Errorf("元数据 %s 值不匹配: 期望 %v, 实际 %v", key, expectedValue, actualValue)
		}
	}
}

// TestSimpleParser_ErrorHandling 测试错误处理
func TestSimpleParser_ErrorHandling(t *testing.T) {
	// 测试空代码
	parser := NewSimpleParser()
	result := parser.ParseWorkflow("")
	if result.Success {
		t.Error("空代码应该解析失败")
	}

	// 测试无效语法
	invalidCode := `step1 = invalid syntax`
	result = parser.ParseWorkflow(invalidCode)
	if result.Success {
		t.Error("无效语法应该解析失败")
	}

	// 测试缺少main函数
	noMainCode := `var input = map[string]interface{}{"test": "value"}
step1 = beiluo.test1.test.test_func() -> (err: error "是否失败");`
	result = parser.ParseWorkflow(noMainCode)
	if result.Success {
		t.Error("缺少main函数应该解析失败")
	}
}

// TestSimpleParser_ComplexWorkflow 测试复杂工作流
func TestSimpleParser_ComplexWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
    "部门": "技术部",
    "职位": "高级工程师",
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号", email: string "邮箱", department: string "部门") -> (workId: string "工号", username: string "用户名", department: string "部门", err: error "是否失败");
step2 = beiluo.test1.crm.crm_interview_schedule(username: string "用户名", department: string "部门", position: string "职位") -> (interviewTime: string "面试时间", interviewer: string "面试官名称", interviewLocation: string "面试地点", err: error "是否失败");
step3 = beiluo.test1.notification.send_email(email: string "邮箱", subject: string "主题", content: string "内容") -> (err: error "是否失败");
step4 = beiluo.test1.notification.send_sms(phone: int "手机号", content: string "内容") -> (err: error "是否失败");
step5 = beiluo.test1.crm.crm_create_candidate(workId: string "工号", username: string "用户名", department: string "部门", position: string "职位") -> (candidateId: string "候选人ID", err: error "是否失败");

func main() {
    //desc: 开始用户注册和面试安排流程
    sys.Println("开始用户注册和面试安排流程...")
    
    //desc: 创建用户账号，获取工号
    工号, 用户名, 部门, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"], input["部门"])
    if step1Err != nil {
        step1.Printf("创建用户失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 用户创建成功，工号: %s", 工号)
    
    //desc: 安排面试时间，联系面试官
    面试时间, 面试官名称, 面试地点, step2Err := step2(用户名, 部门, input["职位"])
    if step2Err != nil {
        step2.Printf("安排面试失败: %v", step2Err)
        return
    }
    step2.Printf("✅ 面试安排成功，时间: %s, 地点: %s", 面试时间, 面试地点)
    
    //desc: 发送邮件通知
    邮件主题 := "面试安排通知"
    邮件内容 := "您已成功安排面试，时间: {{面试时间}}, 地点: {{面试地点}}"
    step3Err := step3(input["邮箱"], 邮件主题, 邮件内容)
    if step3Err != nil {
        step3.Printf("发送邮件失败: %v", step3Err)
        return
    }
    step3.Printf("✅ 邮件发送成功")
    
    //desc: 发送短信通知
    短信内容 := "面试安排: {{面试时间}} {{面试地点}}"
    step4Err := step4(input["手机号"], 短信内容)
    if step4Err != nil {
        step4.Printf("发送短信失败: %v", step4Err)
        return
    }
    step4.Printf("✅ 短信发送成功")
    
    //desc: 创建候选人记录
    候选人ID, step5Err := step5(工号, 用户名, 部门, input["职位"])
    if step5Err != nil {
        step5.Printf("创建候选人记录失败: %v", step5Err)
        return
    }
    step5.Printf("✅ 候选人记录创建成功，ID: %s", 候选人ID)
    
    sys.Println("🎉 用户注册和面试安排流程完成！")
}`

	parser := NewSimpleParser()
	result := parser.ParseWorkflow(code)

	if !result.Success {
		t.Fatalf("解析失败: %s", result.Error)
	}

	// 验证步骤数量
	if len(result.Steps) != 5 {
		t.Fatalf("步骤数量不匹配: 期望 5, 实际 %d", len(result.Steps))
	}

	// 验证主函数语句数量
	if len(result.MainFunc.Statements) < 15 {
		t.Errorf("主函数语句数量不足: 期望至少 15, 实际 %d", len(result.MainFunc.Statements))
	}

	// 验证变量映射
	expectedVars := []string{"工号", "用户名", "部门", "面试时间", "面试官名称", "面试地点", "候选人ID", "step1Err", "step2Err", "step3Err", "step4Err", "step5Err"}
	for _, varName := range expectedVars {
		if _, exists := result.Variables[varName]; !exists {
			t.Errorf("缺少变量: %s", varName)
		}
	}
}

// TestSimpleParser_Performance 测试解析性能
func TestSimpleParser_Performance(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(username: string "用户名", phone: int "手机号") -> (workId: string "工号", username: string "用户名", err: error "是否失败");

func main() {
    sys.Println("开始执行工作流...")
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"])
    if step1Err != nil {
        step1.Printf("执行失败: %v", step1Err)
        return
    }
    step1.Printf("✅ 执行成功，工号: %s", 工号)
    sys.Println("工作流执行完成！")
}`

	parser := NewSimpleParser()

	// 测试解析时间
	start := time.Now()
	for i := 0; i < 1000; i++ {
		result := parser.ParseWorkflow(code)
		if !result.Success {
			t.Fatalf("解析失败: %s", result.Error)
		}
	}
	duration := time.Since(start)

	t.Logf("解析1000次耗时: %v", duration)
	t.Logf("平均每次解析耗时: %v", duration/1000)

	// 性能要求：每次解析应该在1ms以内
	if duration/1000 > time.Millisecond {
		t.Errorf("解析性能不达标: 平均耗时 %v > 1ms", duration/1000)
	}
}
