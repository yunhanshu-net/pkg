package workflow

import (
	"testing"
)

// 测试简单工作流执行
func TestWorkflowExecutor_SimpleWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号) -> (string 工号, string 用户名, err 是否失败);

func main() {
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"]){retry:3, timeout:5000}
    if step1Err != nil {
        fmt.Printf("创建用户失败: %v\n", step1Err)
        return
    }
    fmt.Printf("用户创建成功，工号: %s\n", 工号)
}`

	executor := NewWorkflowExecutor()
	result := executor.ExecuteWorkflow(code)

	if !result.Success {
		t.Fatalf("工作流执行失败: %s", result.Error)
	}

	if len(result.Steps) != 1 {
		t.Errorf("期望执行1个步骤，实际执行%d个步骤", len(result.Steps))
	}

	// 验证步骤信息
	step := result.Steps[0]
	if step.StepName != "step1" {
		t.Errorf("期望步骤名称为step1，实际为%s", step.StepName)
	}

	if step.Function != "step1" {
		t.Errorf("期望函数名为step1，实际为%s", step.Function)
	}

	if !step.Success {
		t.Error("期望步骤执行成功")
	}

	t.Logf("执行成功，总耗时: %dms", result.Duration)
}

// 测试复杂工作流执行
func TestWorkflowExecutor_ComplexWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "用户名": "张三",
    "手机号": 13800138000,
    "邮箱": "zhangsan@example.com",
}

step1 = beiluo.test1.devops.devops_script_create(string 用户名, int 手机号, string 邮箱) -> (string 工号, string 用户名, err 是否失败);
step2 = beiluo.test1.crm.crm_interview_schedule(string 用户名) -> (string 面试时间, string 面试官名称, err 是否失败);
step3 = beiluo.test1.notification.send_email(string 邮箱, string 内容) -> (err 是否失败);

func main() {
    fmt.Println("开始用户注册和面试安排流程...")
    
    // 创建用户
    工号, 用户名, step1Err := step1(input["用户名"], input["手机号"], input["邮箱"]){retry:3, timeout:5000, priority:"high"}
    if step1Err != nil {
        fmt.Printf("创建用户失败: %v\n", step1Err)
        return
    }
    
    // 安排面试
    面试时间, 面试官名称, step2Err := step2(用户名){retry:2, timeout:3000, priority:"normal"}
    if step2Err != nil {
        fmt.Printf("安排面试失败: %v\n", step2Err)
        return
    }
    
    // 发送通知
    通知内容 := "你收到了:{{用户名}},时间：{{面试时间}}的面试安排，请关注"
    step3Err := step3(input["邮箱"], 通知内容){retry:1, timeout:2000, priority:"low"}
    if step3Err != nil {
        fmt.Printf("发送通知失败: %v\n", step3Err)
        return
    }
    
    fmt.Printf("流程完成，工号: %s，面试时间: %s\n", 工号, 面试时间)
}`

	executor := NewWorkflowExecutor()
	result := executor.ExecuteWorkflow(code)

	if !result.Success {
		t.Fatalf("工作流执行失败: %s", result.Error)
	}

	if len(result.Steps) != 3 {
		t.Errorf("期望执行3个步骤，实际执行%d个步骤", len(result.Steps))
	}

	// 验证步骤信息
	expectedSteps := []string{"step1", "step2", "step3"}
	for i, expectedStep := range expectedSteps {
		if i >= len(result.Steps) {
			t.Errorf("步骤%d不存在", i+1)
			continue
		}

		step := result.Steps[i]
		if step.StepName != expectedStep {
			t.Errorf("步骤%d期望名称为%s，实际为%s", i+1, expectedStep, step.StepName)
		}

		if !step.Success {
			t.Errorf("步骤%d执行失败", i+1)
		}
	}

	t.Logf("执行成功，总耗时: %dms", result.Duration)
}

// 测试带条件的工作流执行
func TestWorkflowExecutor_ConditionalWorkflow(t *testing.T) {
	code := `var input = map[string]interface{}{
    "订单号": "ORD-001",
    "金额": 1000,
}

step1 = beiluo.test1.order.validate_order(string 订单号, int 金额) -> (bool 验证结果, err 是否失败);
step2 = beiluo.test1.order.process_payment(string 订单号, int 金额) -> (string 支付流水号, err 是否失败);
step3 = beiluo.test1.order.send_notification(string 订单号) -> (err 是否失败);

func main() {
    fmt.Println("开始订单处理流程...")
    
    // 验证订单
    验证结果, step1Err := step1(input["订单号"], input["金额"]){retry:2, timeout:3000}
    if step1Err != nil {
        fmt.Printf("订单验证失败: %v\n", step1Err)
        return
    }
    
    if 验证结果 {
        fmt.Println("订单验证通过，开始处理支付...")
        
        // 处理支付
        支付流水号, step2Err := step2(input["订单号"], input["金额"]){retry:3, timeout:5000, priority:"high"}
        if step2Err != nil {
            fmt.Printf("支付处理失败: %v\n", step2Err)
            return
        }
        
        // 发送通知
        step3Err := step3(input["订单号"]){retry:1, timeout:2000}
        if step3Err != nil {
            fmt.Printf("发送通知失败: %v\n", step3Err)
            return
        }
        
        fmt.Printf("订单处理完成，支付流水号: %s\n", 支付流水号)
    } else {
        fmt.Println("订单验证失败，流程结束")
    }
}`

	executor := NewWorkflowExecutor()
	result := executor.ExecuteWorkflow(code)

	if !result.Success {
		t.Fatalf("工作流执行失败: %s", result.Error)
	}

	// 验证执行了所有步骤
	if len(result.Steps) != 3 {
		t.Errorf("期望执行3个步骤，实际执行%d个步骤", len(result.Steps))
	}

	t.Logf("执行成功，总耗时: %dms", result.Duration)
}

// 测试执行引擎性能
func TestWorkflowExecutor_Performance(t *testing.T) {
	code := `var input = map[string]interface{}{
    "项目名称": "test-project",
}

step1 = beiluo.test1.devops.git_push(string 项目名称) -> (err 是否失败);
step2 = beiluo.test1.devops.build_image(string 项目名称) -> (string 镜像ID, err 是否失败);
step3 = beiluo.test1.devops.deploy_service(string 项目名称, string 镜像ID) -> (string 服务地址, err 是否失败);

func main() {
    fmt.Println("开始CI/CD流程...")
    
    step1Err := step1(input["项目名称"]){retry:2, timeout:3000}
    if step1Err != nil {
        fmt.Printf("代码推送失败: %v\n", step1Err)
        return
    }
    
    镜像ID, step2Err := step2(input["项目名称"]){retry:3, timeout:10000, priority:"high"}
    if step2Err != nil {
        fmt.Printf("镜像构建失败: %v\n", step2Err)
        return
    }
    
    服务地址, step3Err := step3(input["项目名称"], 镜像ID){retry:2, timeout:5000}
    if step3Err != nil {
        fmt.Printf("服务部署失败: %v\n", step3Err)
        return
    }
    
    fmt.Printf("CI/CD流程完成，服务地址: %s\n", 服务地址)
}`

	executor := NewWorkflowExecutor()

	// 执行多次测试性能
	for i := 0; i < 5; i++ {
		result := executor.ExecuteWorkflow(code)
		if !result.Success {
			t.Fatalf("第%d次执行失败: %s", i+1, result.Error)
		}
		t.Logf("第%d次执行耗时: %dms", i+1, result.Duration)
	}
}
